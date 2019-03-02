// Package piclient manages how http API endpoints are queried using a http client.
package piclient

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	// rawCommitURL is a URL format that returns a single commit. Its general
	// format is 'GET /repos/{owner}/{repo}/commits/{sha}'.
	rawCommitURL = "/repos/%s/%s/commits/%s"

	// rawCommitsURL is a URL format that returns all the commit. Its general
	// format is 'GET /repos/{owner}/{repo}/commits'.
	rawCommitsURL = "/repos/%s/%s/commits"

	// rawContentsURL is a URL format that returns contents of a github directory
	// or a file. Its general format is 'GET /repos/{owner}/{repo}/contents/{path}'.
	rawContentsURL = "/repos/%s/%s/contents/%s"

	// rawSHAListURL is a URL format that returns a list of commit SHA strings
	// associated with the provided github file path. page and per page query values
	// are used to manage pagination. Its general format is
	// 'GET /repos/{owner}/{repo}?path={file_path}&page={page_number}&per_page={page_size}'
	rawSHAListURL = "/repos/%s/%s/commits?path=%s/%s/plugins/decred/ballot.journal&page=%d&per_page=%d"
)

// This a github Oauth token that is required to avoid the 60req/hour for the
// core github API. With Oauth token set, authenticated requests are made which
// do get very high to no rate limit set at all.
var accessToken string

// NewHTTPClient returns an initialized http client whose timeout has been
// capped at 10 seconds.
func NewHTTPClient() *http.Client {
	tr := &http.Transport{
		MaxIdleConns:        5,
		IdleConnTimeout:     10 * time.Second,
		TLSHandshakeTimeout: 5 * time.Second,
		DisableCompression:  true,
	}
	return &http.Client{
		Transport: tr,
		Timeout:   time.Second * 10,
	}
}

// SetAccessToken sets the github access token, return an error if an empty
// access token is being set.
func SetAccessToken(token string) error {
	if token == "" {
		return fmt.Errorf("empty github access token found " +
			"(https://developer.github.com/v3/#rate-limiting)")
	}

	accessToken = token
	return nil
}

// GetRequestHandler accepts a http client, github access token and a URL path
// needed to make a http GET request. It makes a request and returns a byte slice
// read from the body. A github access token is maybe needed to avoid the app
// from being rate limited by github API endpoints. For unauthenticated requests
// a limit of 60 requests per hour is set but for authenticated requests more than
// 5000 requests can be made per hour. https://developer.github.com/v3/#rate-limiting
func GetRequestHandler(client *http.Client, URLPath string) ([]byte, error) {
	if URLPath == "" {
		return nil, fmt.Errorf("empty URL path found")
	}

	// Obtain a GET request.
	req, err := http.NewRequest("GET", URLPath, nil)
	if err != nil || req == nil {
		return nil, fmt.Errorf("http.NewRequest error : %v", err)
	}

	// If the access token is not empty, set it to the request. The access token
	// is needed to increase rate limit to past 60req/hr.
	if accessToken != "" {
		req.Header.Set("Authorization", "token "+accessToken)
	}

	// Make the actual GET request.
	resp, err := client.Do(req)
	if err != nil || resp == nil {
		return nil, fmt.Errorf("client.Get error : %v", err)
	}

	// Read the request body.
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ioutil.ReadAll error : %v", err)
	}

	return body, nil
}

// CommitURL constructs the complete commit or commits query github API endpoint.
func CommitURL(baseURL, commitSHA, repoOwner, repoName string) string {
	var path string
	if commitSHA == "" {
		path = fmt.Sprintf(rawCommitsURL, repoOwner, repoName)
	} else {
		path = fmt.Sprintf(rawCommitURL, repoOwner, repoName, commitSHA)
	}

	return baseURL + path
}

// SHAListURL constructs the complete SHA list query github API endpoint.
func SHAListURL(baseURL, proposalToken, repoOwner, repoName, votesDirName string,
	page, pageSize int) (string, error) {
	if proposalToken == "" {
		return "", fmt.Errorf("missing proposal token")
	}
	return baseURL + fmt.Sprintf(rawSHAListURL, repoOwner, repoName, proposalToken,
		votesDirName, page, pageSize), nil
}

// ContentsURL constructs the complete directory of file content query github
// API endpoint.
func ContentsURL(baseURL, proposalToken, repoOwner, repoName string) (string, error) {
	if proposalToken == "" {
		return "", fmt.Errorf("missing proposal token")
	}
	return baseURL + fmt.Sprintf(rawContentsURL, repoOwner, repoName, proposalToken), nil
}
