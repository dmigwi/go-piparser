// Copyright 2019 Migwi Ndung'u.
// License that can be found in the LICENSE file.

package proposals

import (
	"net/http"

	"github.com/dmigwi/go-piparser/v1/proposals/gitapi"
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

func NewCMDExplorer() {

}
