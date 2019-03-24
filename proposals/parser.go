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
	"time"

	"github.com/dmigwi/go-piparser/v1/types"
)

const (
	// gitCmd defines the prefix string for all commands issued to the git
	// command line interface.
	gitCmd = "git"

	// listCommitsArg defines the git command line argument that lists the repo
	// commit history in a reverse chronoligical order. The oldest commit is
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
)

// Parser holds the clone directory, repo owner and repo name. This data is
// used to query politeia data via git command line tool.
type Parser struct {
	sync.RWMutex
	cloneDir  string
	repoName  string
	repoOwner string
}

// init sets the JournalActionFormat regex expression that helps eliminate
// unwanted journal action votes data pushed to github.
func init() {
	types.SetJournalActionFormat()
}

// NewExplorer returns a Parser instance with repoName, cloneDir and repoOwner
// set. If the repoName and repoOwner provided are empty, the defaults are set.
// If the cloneDir is not provided or an invalid path is provided, a dir in the
// tmp folder is created and set. It also sets ups the environment by cloning
// the repo if it doesn't exist or fetches the latest updates if it does. It
// initiates an asynchronous fetch of hourly politiea updates and there after
// triggers the client to fetch the new updates via the notificationHander function.
func NewExplorer(repoOwner, repo, rootCloneDir string, notificationHander func()) (*Parser, error) {
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
		timer := time.NewTimer(1 * time.Hour)
		for range timer.C {
			if err := p.updateEnv(); err != nil {
				log.Fatalf("updateEnv failed: %v", err)
				return
			}

			// notificationHander is invoked to trigger the process of
			// retrieving the newly fetched updates.
			notificationHander()
		}
	}()

	return p, nil
}

// Proposal returns the all the commits history data associated with the provided
// proposal token. This method is thread-safe.
func (p *Parser) Proposal(proposalToken string) ([]*types.History, error) {
	p.Lock()
	defer p.Unlock()

	if err := types.SetProposalToken(proposalToken); err != nil {
		// error returned, indicates that the proposal token was empty.
		return nil, err
	}

	return p.proposal(proposalToken)
}

// ProposalUpdate returns the commits history data associated with the provided
// proposal token and was made after the since argument time provided. This
// method is thread-safe.
func (p *Parser) ProposalUpdate(proposalToken string, since time.Time) ([]*types.History, error) {
	p.Lock()
	defer p.Unlock()

	if err := types.SetProposalToken(proposalToken); err != nil {
		// error returned, indicates that the proposal token was empty.
		return nil, err
	}

	return p.proposal(proposalToken, since)
}

// Proposals returns all the commits history data for the current proposal tokens
// available. This method is thread-safe.
func (p *Parser) Proposals() ([]*types.History, error) {
	p.Lock()
	defer p.Unlock()

	return p.proposal("")
}

// ProposalsUpdate returns all the commits history updates for the current
// proposal tokens available since the provided date. This method is thread-safe.
func (p *Parser) ProposalsUpdate(since time.Time) ([]*types.History, error) {
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
	args := []string{listCommitsArg, reverseOrder, commitPatchArg, proposalToken}

	// Append the time limiting arguments.
	if len(since) > 0 && since[0] != t {
		args = append(args, []string{sinceArg, since[0].Format(types.CmdDateFormat)}...)
	}

	// Fetch the data via git cmd.
	patchData, err := p.readCommandOutput(gitCmd, args...)
	if err != nil {
		return nil, fmt.Errorf("fetching proposal(s) history failed: %v", err)
	}

	data := strings.Split(patchData, "commit")

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
		if err = types.CustomUnmashaller(&h, entry, since...); err != nil {
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
			return nil
		}

		// git pull command failed.
		// Drop the old repo and proceed to clone the repo a fresh.
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
	dir := filepath.Join(p.cloneDir, p.repoName)
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
