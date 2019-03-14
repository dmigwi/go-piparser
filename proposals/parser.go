// Copyright 2019 Migwi Ndung'u.
// License that can be found in the LICENSE file.

// Package proposals holds the various methods and functions that facilitate
// access to Politeia votes data cloned from github and accessed using
// the git commandline interface. Pre-Installation of the git cmd tool is a
// requirement for effective functionality use.
package proposals

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/dmigwi/go-piparser/v1/types"
)

const (
	// gitCmd defines the prefix string for all commands issued to the git
	// commandline interface.
	gitCmd = "git"

	// listCommitsArg defines the git commandline argument that lists the repo
	// commit history in a chronoligical order. The oldest commit is listed as
	// the last.
	listCommitsArg = "log"

	// commitPatchArg is an optional argument that is added to show the commit
	// history with a patch (changes made) field included.
	commitPatchArg = "-p"

	// cloneArg is the argument added between the git prefix command and the
	// repository download URL. It is used to download the repository into the
	// underlying platform in the set working directory.
	cloneArg = "clone"

	// versionArg lists the version of the active git installation. It is
	// primarily used to check if git is installed on the underlying platform.
	versionArg = "--version"

	// pullChangesArg is an argument that helps pull latest changes from the
	// remote repository set.
	pullChangesArg = "pull"

	// remoteURLRef references the remote URL used to clone the repository.
	// 'origin' is the default set.
	remoteURLRef = "origin"

	// remoteURL uses the https protocol instead of git or ssh protocol. git
	// protocol may be faster but requires a dedicated port (9418) to be
	// open always. ssh requires authentication which is clearly not necessary
	// in this case scenario. Find out more on the access protcols here:
	// https://git-scm.com/book/en/v2/Git-on-the-Server-The-Protocols
	remoteURL = "https://github.com/%s/%s.git"

	// DirPrefix defines the temporary folder prefix.
	DirPrefix = "go-piparser"
)

// Parser holds the clone directory, repo owner, repo name and last Update time.
// This data is used to query politeia votes data via git commandline tool.
type Parser struct {
	sync.RWMutex
	cloneDir   string
	repoName   string
	repoOwner  string
	lastUpdate time.Time
}

// init sets the JournalActionFormat regex expression that helps eliminate
// unwanted journal action votes data pushed to github.
func init() {
	types.SetJournalActionFormat()
}

// NewExplorer returns a Parser instance with a repoName, cloneDir and repoOwner
// set. If the repoName and repoOwner provided are empty, the defaults are set.
// If the cloneDir is not provided or an invalid path is provided, a dir in the
// tmp folder is created and set.
func NewExplorer(repoOwner, repo, rootCloneDir string) (*Parser, error) {
	// Trim trailing and leading whitespaces
	repo = strings.TrimSpace(repo)
	repoOwner = strings.TrimSpace(repoOwner)
	rootCloneDir = strings.TrimSpace(rootCloneDir)

	// Set the default repo name if an empty value was passed.
	if repo == "" {
		repo = types.DefaultRepo
	}

	// Set the default repo owner if an empty value was passed.
	if repoOwner == "" {
		repoOwner = types.DefaultRepoOwner
	}

	// If no directory was provided or the provided directory does not exist
	// create a temp folder.
	var err error
	if _, err = os.Stat(rootCloneDir); os.IsNotExist(err) {
		rootCloneDir, err = ioutil.TempDir(os.TempDir(), DirPrefix)
		if err != nil || rootCloneDir == "" {
			return nil, fmt.Errorf("failed to create a temp cloning dir: %v", err)
		}
	}

	p := &Parser{
		repoName:  repo,
		repoOwner: repoOwner,
		cloneDir:  rootCloneDir,
	}

	// For the first time, initiate git update outside the goroutine and on
	// consecutive times at intervals of 1hr fetch the updates in a goroutine.
	if err := p.updateEnv(); err != nil {
		return nil, fmt.Errorf("updateEnv failed: %v", err)
	}

	// This git updates fetch is run in goroutine but when the actual update is
	// being made, fetching data via methods provided here is temporary halted.
	go func() {
		// Initiate a repo update at intervals of 1h. Politeia updates are made hourly.
		// https://docs.decred.org/advanced/navigating-politeia-data/#voting-and-comment-data
		timer := time.NewTimer(1 * time.Hour)
		for range timer.C {
			if err := p.updateEnv(); err != nil {
				log.Fatalf("updateEnv failed: %v", err)
				return
			}
		}
	}()

	return p, nil
}

// Proposal returns the all the commits history data associated with the provided
// proposal token. This method is thread-safe.
func (p *Parser) Proposal(proposalToken string) (items []*types.History, err error) {
	p.Lock()
	defer p.Unlock()

	if err = types.SetProposalToken(proposalToken); err != nil {
		// error returned, indicates that the proposal token was empty.
		return nil, err
	}

	defer types.ClearProposalToken()

	data, err := p.proposal(proposalToken)
	if err != nil {
		return nil, fmt.Errorf("fetching proposal failed: %v", err)
	}

	for _, entry := range data {
		// strings.Split returns some split strings as empty or with just
		// whitespaces and other special charactes. This happens when
		// the seperating argument is the first in the source string or is
		// surrounded by whitespaces and other special characters.
		if len(strings.TrimSpace(entry)) == 0 {
			continue
		}

		var h types.History

		// entry string is not a valid JSON string format thus the use of a
		// customized unmarshaller.
		if err = types.CustomUnmashaller(&h, entry); err != nil {
			return nil, fmt.Errorf("CustomUnmashaller failed: %v", err)
		}

		// Do not store any empty votes data.
		if len(h.VotesInfo) == 0 || h.Author == "" || h.CommitSHA == "" {
			continue
		}

		items = append(items, &h)
	}

	return
}

// proposal queries the provided proposal token's data from the cloned
// repository using the installed git commandline interface. The single string
// of commits message is split into a slice of individual commit messages and
// returned.
func (p *Parser) proposal(proposalToken string) ([]string, error) {
	patchData, err := p.readCommandOutput(gitCmd, listCommitsArg,
		commitPatchArg, proposalToken)
	if err != nil {
		return nil, fmt.Errorf("fetching proposal history failed: %v", err)
	}

	return strings.Split(patchData, "commit"), nil
}

// updateEnv ensures that a working git commandline tool is installed in the
// underlying platform. It also checks if the required repo was cloned earlier.
// If the repo was cloned earlier, the latest changes are pulled and if an error
// occurs while pulling updates, the old repo version is dropped and a fresh
// clone is made.
func (p *Parser) updateEnv() error {
	p.Lock()
	defer p.Unlock()

	// check if git exists by checking the git installation version.
	err := p.execCommand(gitCmd, versionArg)
	if err != nil {
		return fmt.Errorf("checking git version(%s %s) failed: %v",
			gitCmd, versionArg, err)
	}

	// full clone directory: includes the expected repository name.
	workingDir := filepath.Join(p.cloneDir, p.repoName)
	_, err = os.Stat(workingDir)

	switch {
	case !os.IsNotExist(err):
		// The working directory was found thus initiate the updates fetch process.
		if err := p.execCommand(gitCmd, pullChangesArg, remoteURLRef); err == nil {
			p.lastUpdate = time.Now()
			return nil
		}

		// git pull changes command failed.
		// Drop the old repo and proceed with a full fresh repo clone.
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

// readCommandOutput reads the std output messages of the run command.
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

// execCommand executes commands that do not return necessary std output messages.
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
