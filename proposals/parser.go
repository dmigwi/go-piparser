// Copyright 2019 Migwi Ndung'u.
// License that can be found in the LICENSE file.

package proposals

import (
	"net/http"
	"strings"

	"github.com/dmigwi/go-piparser/v1/proposals/gitapi"
	"github.com/dmigwi/go-piparser/v1/proposals/gitcmd"
	"github.com/dmigwi/go-piparser/v1/types"
)

// Initializes and sets the one JournalActionFormat variables. JournalActionFormat
// is a regex expression that helps eliminate unwanted parts of the vote information.
func init() {
	types.SetJournalActionFormat()
}

type ExplorerDataSource interface {
	Proposal(proposalToken string) (items []*types.History, err error)
}

func NewAPIExplorer(accessToken, repoOwner, repoName string,
	newInstance ...*http.Client) (ExplorerDataSource, error) {

	// Trim trailing and leading whitespaces
	repoName = strings.TrimSpace(repoName)
	repoOwner = strings.TrimSpace(repoOwner)
	accessToken = strings.TrimSpace(accessToken)

	// Set defaults if empty values were passed
	ValidateRepoProperties(&repoOwner, &repoName)

	var parser *gitapi.Parser
	if len(newInstance) == 0 {
		parser = gitapi.NewParser(repoOwner, repoName)
	} else {
		parser = gitapi.NewParser(repoOwner, repoName, newInstance[0])
	}

	err := parser.SetAccessToken(accessToken)
	if err != nil {
		return nil, err
	}

	return parser, err
}

func NewCMDExplorer(repoOwner, repoName, rootCloneDir string) (
	ExplorerDataSource, error) {

	// Trim trailing and leading whitespaces
	repoName = strings.TrimSpace(repoName)
	repoOwner = strings.TrimSpace(repoOwner)
	rootCloneDir = strings.TrimSpace(rootCloneDir)

	// Set defaults if empty values were passed
	ValidateRepoProperties(&repoOwner, &repoName)

	parser, err := gitcmd.NewParser(repoName, repoOwner, rootCloneDir)
	if err != nil {
		return nil, err
	}

	return parser, nil
}

// ValidateRepoProperties sets the default repo name and repo user if empty
// value were passed.
func ValidateRepoProperties(repoOwner, repoName *string) {
	if *repoName == "" {
		*repoName = types.DefaultRepo
	}

	if *repoOwner == "" {
		*repoOwner = types.DefaultRepoOwner
	}
}
