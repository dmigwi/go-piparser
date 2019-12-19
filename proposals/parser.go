// Copyright 2019 Migwi Ndung'u.
// License that can be found in the LICENSE file.

// Package proposals holds the various methods and functions that facilitate
// access to Politeia votes data cloned from github and accessed using
// the git commandline interface. Pre-Installation of the git cmd tool is a
// requirement for effective functionality use.
package proposals

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/dmigwi/go-piparser/proposals/types"
)

const (
	// gitCmd defines the prefix string for all commands issued to the git
	// command line interface.
	gitCmd = "git"

	// listCommitsArg defines the git command line argument that lists the repo
	// commit history in a reverse chronological order. The oldest commit is
	// listed as the last.
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
	// remote repository set. git pull = git fetch + git merge.
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

	// sinceArg with syntax "--since <date>" returns commits more recent than a
	// specific date.
	sinceArg = "--since"

	// By default, the commits are shown in reverse chronological order. Using
	// the reverse order argument ensures that all commits returned are listed
	// in chronological order.
	reverseOrder = "--reverse"

	// cloneRepoAlias defines the clone repository alias used by default instead
	// of the actual repo name.
	cloneRepoAlias = "prop-repo"

	// updateFlag defines the flag to be set when the client want to receive updates
	updateFlag = int32(2)

	// trackedRemoteURL is part of the arguments that help retrieve the clone url
	// of the tracked repository.
	trackedRemoteURL = "get-url"

	// remoteDef defines the remote argument
	remoteDef = "remote"
)

// Parser holds the clone directory, repo owner and repo name. This data is
// used to query politeia data via git command line tool.
type Parser struct {
	sync.RWMutex
	cloneDir    string
	repoName    string
	repoOwner   string
	triggerFlag int32
}

// triggerChan is a channel used to notify the client if updates are available.
var triggerChan chan struct{}

// NewParser returns a Parser instance with repoName, cloneDir and repoOwner
// set. If the repoName and repoOwner provided are empty, the defaults are set.
// If the cloneDir is not provided or an invalid path is provided, a dir in the
// tmp folder is created and set. It also sets ups the environment by cloning
// the repo if it doesn't exist or fetches the latest updates if it does. It
// initiates an asynchronous fetch of hourly politeia updates and there after
// triggers the client to fetch the new updates via a signal channel if the
// trigger flag was set and the channel isn't blocked.
func NewParser(repoOwner, repo, rootCloneDir string) (*Parser, error) {
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

	// If tests are running do not proceed further to clone the test git repos.
	if flag.Lookup("test.v") != nil {
		return p, nil
	}

	// For the first time, initiate git update outside the goroutine and on
	// consecutive times at intervals of 1hr fetch the updates in a goroutine.
	if err := p.updateEnv(); err != nil {
		return nil, fmt.Errorf("updateEnv failed: %v", err)
	}

	// This git updates fetch is made asynchronous.
	go func() {
		// Initiate a repo update at intervals of 1h. Politeia updates are made hourly.
		// https://docs.decred.org/advanced/navigating-politeia-data/#voting-and-comment-data
		timer := time.NewTicker(time.Hour)
		for range timer.C {
			if err := p.updateEnv(); err != nil {
				log.Printf("updateEnv failed: %v", err)
				continue
			}

			// If UpdateSignal() was invoked, the trigger flag must have
			// been set, an indication that the client wants to fetch updates
			// after the tool retrieves them.
			if atomic.LoadInt32(&(p.triggerFlag)) == updateFlag {
				// Attempt to send updates signal if the channel isn't blocked
				// otherwise ignore it till next interval.
				select {
				case triggerChan <- struct{}{}:
				default:
				}
			}
		}
	}()

	return p, nil
}

// UpdateSignal sends a read only signal channel used to inform the client that
// some updates exists.
func (p *Parser) UpdateSignal() <-chan struct{} {
	p.Lock()
	defer p.Unlock()

	atomic.StoreInt32(&(p.triggerFlag), updateFlag)

	if triggerChan == nil {
		triggerChan = make(chan struct{})
	}

	return triggerChan
}

// TriggerUpdates allows the user to have a way to trigger updates retrieval
// from github should they choose not to wait for the hourly updates or are
// confident that new updates exists but the default update may take a while.
// This is as a fail safe method to trigger updates but its usage should be
// limited to very necessary instances to avoid blocking the default hourly
// updates retrieval system.
func (p *Parser) TriggerUpdates() error {
	return p.updateEnv()
}

// ProposalHistory returns the all the commits history data associated with the
// provided proposal token. This method is thread-safe.
func (p *Parser) ProposalHistory(proposalToken string) ([]*types.History, error) {
	p.Lock()
	defer p.Unlock()

	if err := types.SetProposalToken(proposalToken); err != nil {
		// error returned, indicates that the proposal token was empty.
		return nil, err
	}

	return p.proposal(proposalToken)
}

// ProposalHistorySince returns the commits history data associated with the
// provided proposal token and was made after the since argument time provided.
// This method is thread-safe.
func (p *Parser) ProposalHistorySince(proposalToken string, since time.Time) ([]*types.History, error) {
	p.Lock()
	defer p.Unlock()

	if err := types.SetProposalToken(proposalToken); err != nil {
		// error returned, indicates that the proposal token was empty.
		return nil, err
	}

	return p.proposal(proposalToken, since)
}

// ProposalsHistory returns all the commits history data for the current proposal
// tokens available. This method is thread-safe.
func (p *Parser) ProposalsHistory() ([]*types.History, error) {
	p.Lock()
	defer p.Unlock()

	return p.proposal("")
}

// ProposalsHistorySince returns all the commits history updates for the current
// proposal tokens available since the provided date. This method is thread-safe.
func (p *Parser) ProposalsHistorySince(since time.Time) ([]*types.History, error) {
	p.Lock()
	defer p.Unlock()

	return p.proposal("", since)
}

// proposal queries and parses the provided proposal token(s) data from the
// cloned repository using the installed git command line interface tool. If
// the optional since time argument is provided, only the proposal(s) history
// returned was created after the since time.
func (p *Parser) proposal(proposalToken string,
	since ...time.Time) (items []*types.History, err error) {

	defer types.ClearProposalToken()

	var t time.Time
	args := []string{listCommitsArg, reverseOrder, commitPatchArg}

	// Append the proposal token limiting argument if it exists.
	if proposalToken != "" {
		args = append(args, proposalToken)
	}

	// Append the time limiting argument if it exists.
	if len(since) > 0 && since[0] != t {
		args = append(args, []string{sinceArg,
			since[0].Format(types.CmdDateFormat)}...)
	}

	// Fetch the data via git cmd.
	patchData, err := p.readCommandOutput(gitCmd, args...)
	if err != nil {
		return nil, fmt.Errorf("fetching proposal(s) history failed: %v", err)
	}

	data := strings.Split(patchData, "commit")

	for _, entry := range data {
		// strings.Split returns some split strings as empty or with just
		// whitespaces and other special characters. This happens when the
		// seperating argument is the first in the source string or is
		// surrounded by whitespaces and other special characters.
		if len(strings.TrimSpace(entry)) == 0 {
			continue
		}

		var h types.History

		// entry string is not a valid JSON string format thus the use of a
		// customized unmarshaller.
		if err = types.CustomUnmashaller(&h, entry, since...); err != nil {
			return nil, fmt.Errorf("CustomUnmashaller failed: %v", err)
		}

		// Do not store any empty history data.
		if len(h.Patch) == 0 || h.Author == "" || h.CommitSHA == "" {
			continue
		}

		items = append(items, &h)
	}

	return
}

// updateEnv pulls changes from github if they exists or otherwise it clones the
// repository. It also ensures that a working git commandline tool is installed
// in the underlying platform and has the minimum version required. If the
// required repo was cloned earlier, only the latest changes are pulled
// otherwise a fresh clone is made. Should an error occurs while pulling
// updates, the old repo is dropped and a fresh clone made.
func (p *Parser) updateEnv() error {
	p.Lock()
	defer p.Unlock()

	// check if git exists by checking the git installation version.
	versionStr, err := p.readCommandOutput(gitCmd, versionArg)
	if err != nil {
		return fmt.Errorf("checking git version(%s %s) failed: %v",
			gitCmd, versionArg, err)
	}

	// Check if a valid git version exists. A minimum of v1.5.1 is required.
	if err = types.IsGitVersionSupported(versionStr); err != nil {
		// invalid/unsupported git version must have been found.
		return err
	}

	// full clone directory: includes the expected repository name.
	workingDir := filepath.Join(p.cloneDir, cloneRepoAlias)
	_, err = os.Stat(workingDir)

	completeRemoteURL := fmt.Sprintf(remoteURL, p.repoOwner, p.repoName)

	switch {
	case !os.IsNotExist(err):
		// The working directory was found thus check if the tracked repo is the
		// same as the required one.
		trackedRepo, err := p.readCommandOutput(gitCmd, remoteDef, trackedRemoteURL, remoteURLRef)

		// If the required tracked repo was found initiate the updates fetch process
		if err == nil && types.IsMatching(trackedRepo, completeRemoteURL) {
			if err = p.execCommand(gitCmd, pullChangesArg, remoteURLRef); err == nil {
				return nil
			}
		}

		// git pull command failed or the required tracked repo wasn't found.
		// Drop the old repo and proceed to clone the repo a fresh.
		if err = os.RemoveAll(workingDir); err != nil {
			return err
		}

		fallthrough

	default:
		// The required working directory could not be found or the repo update
		// process failed. Clone the remote repository into the clone directory.
		err = p.execCommand(gitCmd, cloneArg, completeRemoteURL, cloneRepoAlias)
		if err != nil {
			return fmt.Errorf("failed to clone %s : %v", completeRemoteURL, err)
		}
	}

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
		return "", formatError(cmdName, args, err)
	}

	return string(stdOutput), nil
}

// execCommand executes commands that do not return necessary std output messages.
func (p *Parser) execCommand(cmdName string, args ...string) error {
	cmd, err := p.processCommand(cmdName, args...)
	if err != nil {
		return err
	}

	err = cmd.Run()

	return formatError(cmdName, args, err)
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
	dir := filepath.Join(p.cloneDir, cloneRepoAlias)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		dir = p.cloneDir
	}
	return dir
}

// formatError replaces "exit status 128" with "git args... failed to execute."
func formatError(cmd string, args []string, err error) error {
	if err == nil {
		return err
	}
	str := fmt.Sprintf("%s command with %v failed to execute", cmd, args)
	return fmt.Errorf(types.ReplaceAny(err.Error(), "exit status 128", str))
}
