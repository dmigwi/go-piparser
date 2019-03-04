// Copyright 2019 Migwi Ndung'u.
// License that can be found in the LICENSE file.

// Package gitcmd holds the various methods and functions that facilitate access
// to Politeia votes data that are cloned from github and accessed using the git
// commandline interface.
package gitcmd

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/dmigwi/go-piparser/v1/types"
)

const (
	// cmdTimeout defines the length in time, where the outcome of a command is
	// expected to be returned before the current context times out.
	cmdTimeout = 5 * time.Second

	// gitCmd is the prefix of all command issued to the git commandline interface.
	gitCmd = "git"

	// listCommitsArg defines the git commandline argument that lists the repo
	// commit history in a chronoligical order. The oldest commit is listed as
	// the last one.
	listCommitsArg = "log"

	// commitPatchArg is an optional argument added to show the commit history
	// with the patch(changes made) field.
	commitPatchArg = "-p"

	// cloneArg is the argument added between the git prefix command and the
	// repository download URL.
	cloneArg = "clone"

	// remoteAddArg is a combination of two arguments that are used to set a
	// reference to the remote repository.
	remoteAddArg = "remote add"

	// versionArg list the version of the active git install. Its primarily used
	// to check if git is installed.
	versionArg = "--version"

	// pullChangesArg is an argument that helps pull latest changes from the
	// remote repository set.
	pullChangesArg = "pull"

	// remoteURLRef references the remote URL used to clone the repository.
	remoteURLRef = "origin"

	// remoteURL set uses https protocol instead of git or ssh protocol. git
	// protocol may be faster but it requires a dedicated port (9418) to be
	// open always. ssh requires authentication which is clearly not necessary
	// in this case scenario. Find more on the access protcols here:
	// https://git-scm.com/book/en/v2/Git-on-the-Server-The-Protocols
	remoteURL = "https://github.com/%s/%s.git"
)

type CMDParser struct {
	repoName   string
	repoOwner  string
	cloneDir   string
	lastUpdate time.Time
}

func NewParser(repo, repoOwner, rootCloneDir string) (*CMDParser, error) {
	// If no directory was provided or the provided directory does not exist
	// create a temporary folder.
	var err error
	if _, err = os.Stat(rootCloneDir); os.IsNotExist(err) {
		rootCloneDir, err = ioutil.TempDir("temp", "go-piparser")
		if err != nil || rootCloneDir == "" {
			return nil, fmt.Errorf("failed to create a temp cloning dir: %v", err)
		}
	}

	p := &CMDParser{
		repoName:  repo,
		repoOwner: repoOwner,
		cloneDir:  rootCloneDir,
	}
	return p, nil
}

func (p *CMDParser) Proposal(proposalToken string) (items []*types.History, err error) {
	if err = types.SetProposalToken(proposalToken); err != nil {
		// error returned indicate that the proposal token was empty.
		return nil, err
	}

	defer types.ClearProposalToken()

	// Initiate a repo update if the last time repo updates were fetched is more
	// than an hour ago since Politeia updates are made hourly.
	// https://docs.decred.org/advanced/navigating-politeia-data/#voting-and-comment-data
	if time.Since(p.lastUpdate) > 1*time.Hour {
		if err := p.updateEnv(); err != nil {
			return nil, fmt.Errorf("updateEnv failed: %v", err)
		}
	}

	data, err := p.proposal(proposalToken)
	if err != nil {
		return nil, fmt.Errorf("fetching proposal failed: %v", err)
	}

	for _, entry := range data {
		if len(entry) == 0 {
			continue
		}

		v2 := &copyHistory{}
		err = v2.customUnmashaller(entry)
		if err != nil {
			return nil, fmt.Errorf("copyHistory.customUnmashaller failed: %v", err)
		}

		h := types.History(*v2)
		items = append(items, &h)
	}

	return
}

func (p *CMDParser) proposal(proposalToken string) ([]string, error) {
	patchData, err := p.readCommandOutput(gitCmd, listCommitsArg,
		commitPatchArg, proposalToken)
	if err != nil {
		return nil, fmt.Errorf("fetching proposal history failed: %v", err)
	}

	return strings.Split(patchData, "commit"), nil
}

// updateEnv ensures that a working git command line tool is installed in the
// underlying platform. It also checks if the required repo was cloned earlier.
// If the repo was cloned earlier the latest changes are pulled. It an error
// occurs while pull updates, the old repo version is dropped and a fresh repos
// is cloned.
func (p *CMDParser) updateEnv() error {
	// check if git exists but checking the git version.
	err := p.execCommand(gitCmd, versionArg)
	if err != nil {
		return fmt.Errorf("checking git version(%s %s) failed: %v",
			gitCmd, versionArg, err)
	}

	// full clone directory: includes the expected cloned repository.
	workingDir := filepath.Join(p.cloneDir, p.repoName)
	_, err = os.Stat(workingDir)

	switch {
	case !os.IsNotExist(err):
		// The working directory was found thus initiate the updates fetch process.
		if err := p.execCommand(gitCmd, pullChangesArg, remoteURLRef); err == nil {
			p.lastUpdate = time.Now()
			return nil
		}

		// git pull changes command failed. Trigger a full fresh repo clone and
		// drop the old error.
		err = nil

		// drop the old working directory version.
		if err = os.RemoveAll(workingDir); err != nil {
			return err
		}

		fallthrough

	default:
		// The required working directory could not be found or the repo update
		// process failed.
		completeRemoteURL := fmt.Sprintf(remoteURL, p.repoOwner, p.repoName)

		// Clone the remote repository into the clone directory.
		if err := p.execCommand(gitCmd, cloneArg, completeRemoteURL); err != nil {
			return fmt.Errorf("failed to clone %s : %v", completeRemoteURL, err)
		}

		// Set the remote url to enable pulling changes added in the future
		if err := p.execCommand(gitCmd, remoteAddArg, remoteURLRef, completeRemoteURL); err != nil {
			return fmt.Errorf("failed to set %s as remote url ref to %s: %v",
				remoteURLRef, completeRemoteURL, err)
		}
	}

	p.lastUpdate = time.Now()
	return nil
}

func (p *CMDParser) readCommandOutput(cmdName string, args ...string) (string, error) {
	cmd, cancel, err := p.processCommand(cmdName, args...)

	defer cancel()
	if err != nil {
		return "", err
	}

	stdOutput, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	return string(stdOutput), nil
}

func (p *CMDParser) execCommand(cmdName string, args ...string) error {
	cmd, cancel, err := p.processCommand(cmdName, args...)

	defer cancel()
	if err != nil {
		return err
	}

	return cmd.Run()
}

func (p *CMDParser) processCommand(cmdName string, args ...string) (*exec.Cmd,
	context.CancelFunc, error) {
	if cmdName == "" {
		return nil, nil, fmt.Errorf("missing command")
	}

	ctx, cancel := context.WithTimeout(context.Background(), cmdTimeout)
	cmd := exec.CommandContext(ctx, cmdName, args...)

	// set the working directory.
	cmd.Dir = p.workingDir()

	return cmd, cancel, nil
}

// workingDir return (cloneDir + repoName) directory path if the target repo
// exists otherwise returns cloneDir as the working directory.
func (p *CMDParser) workingDir() string {
	dir := filepath.Join(p.cloneDir, p.repoName)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		dir = p.cloneDir
	}
	return dir
}
