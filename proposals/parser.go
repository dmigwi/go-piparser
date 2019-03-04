// Copyright 2019 Migwi Ndung'u.
// License that can be found in the LICENSE file.

// Package proposals holds the various methods and functions that facilitate access
// to Politeia votes data that are cloned from github and accessed using the git
// commandline interface.
package proposals

import (
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

	// versionArg list the version of the active git install. Its primarily used
	// to check if git is installed.
	versionArg = "--version"

	// pullChangesArg is an argument that helps pull latest changes from the
	// remote repository set.
	pullChangesArg = "pull"

	// remoteURLRef references the remote URL used to clone the repository.
	// 'origin' is the default set.
	remoteURLRef = "origin"

	// remoteURL set uses https protocol instead of git or ssh protocol. git
	// protocol may be faster but it requires a dedicated port (9418) to be
	// open always. ssh requires authentication which is clearly not necessary
	// in this case scenario. Find more on the access protcols here:
	// https://git-scm.com/book/en/v2/Git-on-the-Server-The-Protocols
	remoteURL = "https://github.com/%s/%s.git"
)

// Parser holds the clone Directory repo Owner, repo Name and last Update time.
// This data is used to politeia votes data via git commandline tool..
type Parser struct {
	cloneDir   string
	repoName   string
	repoOwner  string
	lastUpdate time.Time
}

// Initializes and sets the one JournalActionFormat variables. JournalActionFormat
// is a regex expression that helps eliminate unwanted parts of the vote information.
func init() {
	types.SetJournalActionFormat()
}

// NewExplorer returns a Parser instance with a repoName, cloneDir and repoOwner
// set. If the repoName and repoOwner provided are empty the defaults are set.
// If the cloneDir is not provided a dir in the tmp folder is created and set.
func NewExplorer(repoOwner, repo, rootCloneDir string) (*Parser, error) {
	// Trim trailing and leading whitespaces
	repo = strings.TrimSpace(repo)
	repoOwner = strings.TrimSpace(repoOwner)
	rootCloneDir = strings.TrimSpace(rootCloneDir)

	// Set the default repo name if an empty values was passed.
	if repo == "" {
		repo = types.DefaultRepo
	}

	// Set the default repo owner if an empty values was passed.
	if repoOwner == "" {
		repoOwner = types.DefaultRepoOwner
	}

	// If no directory was provided or the provided directory does not exist
	// create a temporary folder.
	var err error
	if _, err = os.Stat(rootCloneDir); os.IsNotExist(err) {
		rootCloneDir, err = ioutil.TempDir("temp", "go-piparser")
		if err != nil || rootCloneDir == "" {
			return nil, fmt.Errorf("failed to create a temp cloning dir: %v", err)
		}
	}

	p := &Parser{
		repoName:  repo,
		repoOwner: repoOwner,
		cloneDir:  rootCloneDir,
	}
	return p, nil
}

// Proposal returns the commit history data associated with the provided token.
func (p *Parser) Proposal(proposalToken string) (items []*types.History, err error) {
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

		// entry string is not a valid JSON string format.
		v2 := &types.History{}
		err = v2.CustomUnmashaller(entry)
		if err != nil {
			return nil, fmt.Errorf("History.CustomUnmashaller failed: %v", err)
		}

		h := types.History(*v2)

		// Do not store any empty votes data.
		if len(h.VotesInfo) == 0 || h.Author == "" || h.CommitSHA == "" {
			continue
		}

		items = append(items, &h)
	}

	return
}

// proposal queries the provided proposal data from cloned repository using the
// installed git command line interface. The single string of commit messages
// is split into a slice of individual commit messages and returned.
func (p *Parser) proposal(proposalToken string) ([]string, error) {
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
func (p *Parser) updateEnv() error {
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
	}

	p.lastUpdate = time.Now()

	return nil
}

// readCommandOutput executes and reads the stdout message of the run command.
func (p *Parser) readCommandOutput(cmdName string, args ...string) (string, error) {
	cmd, err := p.processCommand(cmdName, args...)
	if err != nil {
		return "", err
	}

	stdOutput, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	return string(stdOutput), nil
}

// execCommand executes commands that do not return any necessary stdout message.
func (p *Parser) execCommand(cmdName string, args ...string) error {
	cmd, err := p.processCommand(cmdName, args...)
	if err != nil {
		return err
	}

	return cmd.Run()
}

// processCommand checks if an empty command prefix was provided. It also sets the
// the working directory.
func (p *Parser) processCommand(cmdName string, args ...string) (*exec.Cmd, error) {
	if cmdName == "" {
		return nil, fmt.Errorf("missing command")
	}
	cmd := exec.Command(cmdName, args...)

	// set the working directory.
	cmd.Dir = p.workingDir()

	return cmd, nil
}

// workingDir return (cloneDir + repoName) directory path if the target repo
// exists otherwise returns cloneDir as the working directory.
func (p *Parser) workingDir() string {
	dir := filepath.Join(p.cloneDir, p.repoName)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		dir = p.cloneDir
	}
	return dir
}
