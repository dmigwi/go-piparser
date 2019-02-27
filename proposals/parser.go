package proposals

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/decred/politeia/politeiad/backend/gitbe"
	"github.com/dmigwi/go-piparser/v1/piclient"
	"github.com/dmigwi/go-piparser/v1/types"
)

var defaultPageSize = 20

const (
	defaultRepo      = "mainnet"
	defaultRepoOwner = "decred-proposals"
)

// Initializes and sets the one JournalActionFormat variables. JournalActionFormat
// is a regex expression that helps eliminate unwanted parts of the vote information.
func init() {
	f, err := json.Marshal(gitbe.JournalAction{
		Version: `[[:digit:]]*`,
		Action:  "[add]*[del]*[addlike]*",
	})
	if err != nil {
		panic(err)
	}

	format := string(f)
	types.SetJournalActionFormat(format)
}

// Parser holds the http client instance, repo Owner and repo Name.
type Parser struct {
	repoOwner string
	repoName  string
	client    *http.Client
}

// NewParser returns a Parser instance with an initialized client, repoName and
// repoOwner. It allows a http client with custom configuration to be set. It
// defaults to the piclient.NewHTTPClient() with general settings if no http client
// was provided or a nil client was passed.
func NewParser(repoOwner, repoName string, newInstance ...*http.Client) *Parser {
	p := new(Parser)

	if len(newInstance) > 0 && newInstance[0] != nil {
		p.client = newInstance[0]
	} else {
		p.client = piclient.NewHTTPClient()
	}

	if repoName == "" {
		repoName = defaultRepo
	}

	if repoOwner == "" {
		repoOwner = defaultRepoOwner
	}

	p.repoName = repoName
	p.repoOwner = repoOwner

	return p
}

// SetProposalToken set the specific proposal token whose vote details only should
// be unmarshalled from the current list of commits.
func (p *Parser) SetProposalToken(token string) error {
	if token == "" {
		return fmt.Errorf("empty token hash string found")
	}
	types.SetProposalToken(token)
	return nil
}

// Proposal returns the commit history data associated with the provided token.
func (p *Parser) Proposal(token string) (items []*types.History, err error) {
	defer types.ClearProposalToken()

	var page int
	var data []*types.History

	items = make([]*types.History, 0)

	votesDir, err := p.retrieveVotesDirName(token)
	if err != nil {
		return nil, err
	}

	for {
		page++
		log.Printf("Handling batch %d of max %d commits...", page, defaultPageSize)

		data, err = p.proposal(token, votesDir, page, defaultPageSize)
		if err != nil {
			return
		}

		items = append(items, data...)

		if len(data) < defaultPageSize {
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

// retrieveSHAList returns a list of commit SHA with the provided token string
// as a file path.
func (p *Parser) retrieveSHAList(token, piVotesDirName string, page, pageSize int) (*types.HistorySHAs, error) {
	// Constructs the full commit SHA list url.
	URLPath, err := piclient.SHAListURL(token, p.repoOwner, p.repoName, piVotesDirName, page, pageSize)
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

	if len(data) == 0 {
		fmt.Println(" >>>> <<<<<< ", string(list))
	}

	return &data, nil
}

// retrieveCommit returns the unmarshalled commit data identified by the provided
// commit SHA if it exists.
func (p *Parser) retrieveCommit(commitSHA string) (*types.History, error) {
	// Constructs full commit content url path.
	URLPath, err := piclient.ContentURL(commitSHA, p.repoOwner, p.repoName)
	if err != nil {
		return nil, fmt.Errorf("creating the ContentURL failed: %v", err)
	}

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
	URLPath, err := piclient.PropDirContentsURL(token, p.repoOwner, p.repoName)
	if err != nil {
		return "", fmt.Errorf("creating the PropDirContentsURL failed: %v", err)
	}

	// Fetch the list of directories in the proposal token string directory.
	data, err := piclient.GetRequestHandler(p.client, URLPath)
	if err != nil {
		return "", err
	}

	// Unmarshals the returned directories data.
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
