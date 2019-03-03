// Copyright 2019 Migwi Ndung'u.
// License that can be found in the LICENSE file.

package gitapi

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/dmigwi/go-piparser/v1/proposals/gitapi/piclient"
	"github.com/dmigwi/go-piparser/v1/types"
)

// pageSize defines the number of records to be fetched from the github API in
// a single API call.
var pageSize = 20

const (
	defaultAPIURL = "https://api.github.com"

	// Sets that a maximum of 200 records per page can be fetched for processing
	// from the github API endpoints in a single API call.
	maxQueriedPageSize = 200
)

// Parser holds the http client instance, repo Owner and repo Name.
type Parser struct {
	repoOwner  string
	repoName   string
	baseAPIURL string
	client     *http.Client
}

// NewParser returns a Parser instance with an initialized client, repoName and
// repoOwner. It allows a http client with custom configuration to be set. It
// defaults to the piclient.NewHTTPClient() with general settings if no http client
// was provided or a nil client was passed. If the repoName and repoOwner provided
// are empty the defaults are set.
func NewParser(repoOwner, repoName string, newInstance ...*http.Client) *Parser {
	p := &Parser{
		repoName:   repoName,
		repoOwner:  repoOwner,
		baseAPIURL: defaultAPIURL,
	}

	if len(newInstance) > 0 && newInstance[0] != nil {
		// Assign custom http client.
		p.client = newInstance[0]
	} else {
		// Assign default http client.
		p.client = piclient.NewHTTPClient()
	}

	return p
}

// SetProposalToken sets the specific proposal token whose vote details only should
// be unmarshalled from the current list of commits.
func (p *Parser) SetProposalToken(token string) error {
	return types.SetProposalToken(token)
}

// SetAccessToken sets the github access token needed to enable making github
// requests with a higher the rate limit than available for unauthenticated
// requests.
func (p *Parser) SetAccessToken(token string) error {
	return piclient.SetAccessToken(token)
}

// SetPageSize sets maximum number of records that can fetched in a single
// API call.
func (p *Parser) SetPageSize(size int) error {
	if size < 1 || size > maxQueriedPageSize {
		return fmt.Errorf("page size should be between 1 and %v", maxQueriedPageSize)
	}

	pageSize = size
	return nil
}

// Proposal returns the commit history data associated with the provided token.
func (p *Parser) Proposal(proposalToken string) (items []*types.History, err error) {
	defer types.ClearProposalToken()

	var page int
	var data []*types.History

	items = make([]*types.History, 0)

	votesDir, err := p.retrieveVotesDirName(proposalToken)
	if err != nil {
		return nil, err
	}

	for {
		page++
		log.Printf("Handling batch %d of max %d commits...", page, pageSize)

		data, err = p.proposal(proposalToken, votesDir, page, pageSize)
		if err != nil {
			return
		}

		items = append(items, data...)

		if len(data) < pageSize {
			return
		}
	}
}

// proposal returns a proposal whose commit history count and page is defined by
// the pageSize and the page values respectively.
func (p *Parser) proposal(token, piVotesDirName string, page, pageSize int) ([]*types.History, error) {
	commitsSHA, err := p.retrieveSHAList(token, piVotesDirName, page, pageSize)
	if err != nil {
		return nil, err
	}

	items := make([]*types.History, 0)
	for _, hash := range *commitsSHA {
		if hash.SHA == "" {
			continue
		}

		elem, err := p.retrieveCommit(hash.SHA)
		if err != nil {
			return nil, err
		}

		if elem != nil {
			if elem.AltResponse != nil && elem.Message != "" {
				return nil, fmt.Errorf(elem.Message + " " + elem.URL)
			}

			items = append(items, elem)
		}
	}
	return items, nil
}

// retrieveSHAList returns a list of commits SHA associated with the provided
// proposal token string. The token string is used as a file path in the github
// proposals repo.
func (p *Parser) retrieveSHAList(token, piVotesDirName string,
	page, pageSize int) (*types.HistorySHAs, error) {
	// Constructs the full commits SHA list url.
	URLPath, err := piclient.SHAListURL(p.baseAPIURL, token, p.repoOwner, p.repoName,
		piVotesDirName, page, pageSize)
	if err != nil {
		return nil, fmt.Errorf("creating the SHAListURL failed: %v", err)
	}

	// Makes a GET request to the commit SHA list url created.
	list, err := piclient.GetRequestHandler(p.client, URLPath)
	if err != nil {
		return nil, fmt.Errorf("piclient.GetRequestHandler failed: %v", err)
	}

	// Unmarshalls the commit SHA list into types.HistorySHAs.
	var data types.HistorySHAs
	err = json.Unmarshal(list, &data)
	if err != nil {
		return nil, fmt.Errorf("types.HistorySHAs unmarshalling failed: %v", err)
	}

	return &data, nil
}

// retrieveCommit returns the unmarshalled commit data identified by the provided
// commit SHA if it exists.
func (p *Parser) retrieveCommit(commitSHA string) (*types.History, error) {
	if commitSHA == "" {
		return nil, fmt.Errorf("missing commit SHA")
	}

	// Constructs full commit content url path.
	URLPath := piclient.CommitURL(p.baseAPIURL, commitSHA, p.repoOwner, p.repoName)

	// Fetch the commit content GET url path request.
	content, err := piclient.GetRequestHandler(p.client, URLPath)
	if err != nil {
		return nil, fmt.Errorf("piclient.GetRequestHandler failed: %v", err)
	}

	// Unmarshal the commit content url response data.
	var data types.History
	err = json.Unmarshal(content, &data)
	if err != nil {
		return nil, fmt.Errorf("types.History unmarshalling failed: %v", err)
	}

	return &data, nil
}

// retrieveVotesDirName returns the folder name that contain ballot.journal file
// in the git repo.
func (p *Parser) retrieveVotesDirName(token string) (string, error) {
	// Construct the full proposal token string directories content url.
	URLPath, err := piclient.ContentsURL(p.baseAPIURL, token, p.repoOwner, p.repoName)
	if err != nil {
		return "", fmt.Errorf("creating the ContentsURL failed: %v", err)
	}

	// Fetch the list of directories in the proposal token string directory.
	data, err := piclient.GetRequestHandler(p.client, URLPath)
	if err != nil {
		return "", err
	}

	// Unmarshals the returned directories data list.
	dirs := make([]types.GitPropDirectories, 0)
	err = json.Unmarshal(data, &dirs)
	if err != nil {
		return "", err
	}

	// The last directory name in the returned slice is the directory that holds
	// ballot.journal file.
	if len(dirs) > 0 {
		return dirs[len(dirs)-1].Name, nil

	}

	return "", fmt.Errorf("No valid vote journal directory name found")
}
