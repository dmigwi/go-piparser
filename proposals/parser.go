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
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/dmigwi/go-piparser/proposals/gitlib"
	"github.com/dmigwi/go-piparser/proposals/gittool"
	"github.com/dmigwi/go-piparser/proposals/types"
)

// updateFlag defines the flag to be set when the client want to receive updates
const updateFlag = int32(2)

// Parser holds the clone directory, repo owner and repo name. This data is
// used to query politeia data via git command line tool.
type Parser struct {
	mtx         sync.RWMutex
	source      types.DataSource
	triggerFlag int32
}

// triggerChan is a channel used to notify the client if updates are available.
var triggerChan chan struct{}

// NewParser returns a Parser instance with repoName, cloneDir and repoOwner
// set. If the repoName and repoOwner provided are empty, the defaults are set.
// If the cloneDir is not provided or an invalid path is provided, a dir in the
// tmp folder is created and set. It also sets ups the environment by cloning
// the repo if it doesn't exist or fetches the latest updates if it does. It
// initiates an asynchronous fetch of hourly politiea updates and there after
// triggers the client to fetch the new updates via a signal channel if the
// trigger flag was set and the channel isn't blocked. isTool is an optional
// argument that if set to true gittool package is set as the data source
// otherwise it defaults to gitlib package as the data source.
func NewParser(repoOwner, repo, rootCloneDir string, isTool ...bool) (*Parser, error) {
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
		rootCloneDir, err = ioutil.TempDir(os.TempDir(), types.DirPrefix)
		if err != nil || rootCloneDir == "" {
			return nil, fmt.Errorf("failed to create a temp cloning dir: %v", err)
		}
	}

	p := new(Parser)
	if len(isTool) > 0 && isTool[0] {
		p.source = gittool.NewDataSource(repoOwner, repo, rootCloneDir)
	} else {
		p.source = gitlib.NewDataSource(repoOwner, repo, rootCloneDir)
	}

	// If tests are running do not proceed further to clone the test git repos.
	if flag.Lookup("test.v") != nil {
		return p, nil
	}

	// For the first time, initiate git update outside the goroutine and on
	// consecutive times at intervals of 1hr fetch the updates in a goroutine.
	if err := p.source.SetUpEnv(); err != nil {
		return nil, fmt.Errorf("updateEnv failed: %v", err)
	}

	// This git updates fetch is made asynchronous.
	go func() {
		// Initiate a repo update at intervals of 1h. Politeia updates are made hourly.
		// https://docs.decred.org/advanced/navigating-politeia-data/#voting-and-comment-data
		timer := time.NewTicker(1 * time.Hour)
		for range timer.C {
			p.mtx.Lock()

			if err := p.source.SetUpEnv(); err != nil {
				log.Fatalf("updateEnv failed: %v", err)
				return
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

			p.mtx.Unlock()
		}
	}()

	return p, nil
}

// UpdateSignal sends a read only signal channel used to inform the client that
// some updates exists.
func (p *Parser) UpdateSignal() <-chan struct{} {
	p.mtx.Lock()
	defer p.mtx.Unlock()

	atomic.StoreInt32(&(p.triggerFlag), updateFlag)

	if triggerChan == nil {
		triggerChan = make(chan struct{})
	}

	return triggerChan
}

// TriggerUpdates allows the user to have a way to trigger updates retrieval
// from github should they choose not to wait for the hourly updates or
// are confident that new updates exists but the default update may take a while.
// This is as a fail safe method to trigger updates but its usage
// should be limited to very neccessary instances to avoid blocking the default
// hourly updates retrieval system.
func (p *Parser) TriggerUpdates() error {
	return p.source.SetUpEnv()
}

// ProposalHistory returns the all the commits history data associated with the
// provided proposal token. This method is thread-safe.
func (p *Parser) ProposalHistory(proposalToken string) ([]*types.History, error) {
	p.mtx.Lock()
	defer p.mtx.Unlock()

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
	p.mtx.Lock()
	defer p.mtx.Unlock()

	if err := types.SetProposalToken(proposalToken); err != nil {
		// error returned, indicates that the proposal token was empty.
		return nil, err
	}

	return p.proposal(proposalToken, since)
}

// ProposalsHistory returns all the commits history data for the current proposal
// tokens available. This method is thread-safe.
func (p *Parser) ProposalsHistory() ([]*types.History, error) {
	p.mtx.Lock()
	defer p.mtx.Unlock()

	return p.proposal("")
}

// ProposalsHistorySince returns all the commits history updates for the current
// proposal tokens available since the provided date. This method is thread-safe.
func (p *Parser) ProposalsHistorySince(since time.Time) ([]*types.History, error) {
	p.mtx.Lock()
	defer p.mtx.Unlock()

	return p.proposal("", since)
}

// proposal queries and parses the provided proposal token(s) data from the
// cloned repository using the installed git command line interface tool. If
// the optional since time argument is provided, only the proposal(s) history
// returned was created after the since time.
func (p *Parser) proposal(proposalToken string, since ...time.Time) ([]*types.History, error) {
	defer types.ClearProposalToken()

	return p.source.PullData(proposalToken)
}
