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

type Parser struct {
	client    *http.Client
	repoOwner string
	repoName  string
}

func NewParser(repoOwner, repoName string, newInstance ...*http.Client) *Parser {
	p := new(Parser)

	if len(newInstance) > 0 && newInstance[0] != nil {
		p.client = newInstance[0]
	} else {
		p.client = piclient.NewHTTPClient()
	}

	p.repoName = repoName
	p.repoOwner = repoOwner

	return p
}

func (p *Parser) Proposal(token string) ([]*types.History, error) {
	commitsSHA, err := p.retrieveSHAList(token)
	if err != nil {
		return nil, err
	}

	types.SetProposalToken(token)

	items := make([]*types.History, 0)
	for _, hash := range *commitsSHA {
		elem, err := p.retrieveCommit(hash.SHA)
		if err != nil {
			return nil, err
		}

		s, _ := json.Marshal(elem)
		log.Println(" >>> ", string(s))

		if elem != nil {
			items = append(items, elem)
		}
	}

	return items, nil
}

func (p *Parser) retrieveSHAList(token string) (*types.HistorySHAs, error) {
	URLPath, err := piclient.SHAListURL(token, p.repoOwner, p.repoName)
	if err != nil {
		return nil, fmt.Errorf("creating the SHAListURL failed: %v", err)
	}

	list, err := piclient.GetRequestHandler(p.client, URLPath)
	if err != nil {
		return nil, fmt.Errorf("piclient.GetRequestHandler failed: %v", err)
	}

	var data *types.HistorySHAs
	err = json.Unmarshal(list, &data)
	if err != nil {
		return nil, fmt.Errorf("types.HistorySHAs unmarshalling failed: %v", err)
	}

	return data, nil
}

func (p *Parser) retrieveCommit(commitSHA string) (*types.History, error) {
	URLPath, err := piclient.ContentURL(commitSHA, p.repoOwner, p.repoName)
	if err != nil {
		return nil, fmt.Errorf("creating the ContentURL failed: %v", err)
	}

	content, err := piclient.GetRequestHandler(p.client, URLPath)
	if err != nil {
		return nil, fmt.Errorf("piclient.GetRequestHandler failed: %v", err)
	}

	var data *types.History
	err = json.Unmarshal(content, &data)
	if err != nil {
		return nil, fmt.Errorf("types.History unmarshalling failed: %v", err)
	}

	return data, nil
}
