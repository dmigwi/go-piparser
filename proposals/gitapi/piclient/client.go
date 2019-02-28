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
	rawCommitURL = "https://api.github.com/repos/%s/%s/commits/%s"

	// rawContentsURL is a URL format that returns contents of a github directory
	// or a file. Its general format is 'GET /repos/{owner}/{repo}/contents/{path}'.
	rawContentsURL = "https://api.github.com/repos/%s/%s/contents/%s"

	// rawSHAListURL is a URL format that returns a list of commit SHA strings
	// associated with the provided github file path. page and per page query values
	// are used to manage pagination. Its general format is
	// 'GET /repos/{owner}/{repo}?path={file_path}&page={page_number}&per_page={page_size}'
	rawSHAListURL = "https://api.github.com/repos/%s/%s/commits?path=%s/%s/plugins/decred/ballot.journal&page=%d&per_page=%d"
)

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

// GetRequestHandler accepts a http client, github access token and a URL path
// needed to make a http GET request. It makes the request and returns a
// byte slice read from the body. A github access token is required to avoid the
// app from being rate limited by github API endpoints. For unauthenticated requests
// a limit of 60 requests per hour is set but for authenticated requests more than
// 5000 requests can be made per hour. https://developer.github.com/v3/#rate-limiting
func GetRequestHandler(client *http.Client, accessToken, URLPath string) ([]byte, error) {
	if URLPath == "" {
		return nil, fmt.Errorf("empty URL path found")
	}

	// Obtain a GET request.
	req, err := http.NewRequest("GET", URLPath, nil)
	if err != nil || req == nil {
		return nil, fmt.Errorf("http.NewRequest error : %v", err)
	}

	// If the access token is not empty, set it to the request else return an error.
	// The app could work without the access token but 60req/hr is not enough to
	// enable the app make any meaningful number of queries before the rate
	// limit is hit.
	if accessToken != "" {
		req.Header.Set("Authorization", "token "+accessToken)
	} else {
		return nil, fmt.Errorf("empty github access token found (https://developer.github.com/v3/#rate-limiting)")
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

// CommitURL constructs the complete commit query github API endpoint.
func CommitURL(commitSHA, repoOwner, repoName string) (string, error) {
	if commitSHA == "" {
		return "", fmt.Errorf("missing commit SHA")
	}
	return fmt.Sprintf(rawCommitURL, repoOwner, repoName, commitSHA), nil
}

// SHAListURL constructs the complete SHA list query github API endpoint.
func SHAListURL(proposalToken, repoOwner, repoName, votesDirName string,
	page, pageSize int) (string, error) {
	if proposalToken == "" {
		return "", fmt.Errorf("missing proposal token")
	}
	return fmt.Sprintf(rawSHAListURL, repoOwner, repoName, proposalToken,
		votesDirName, page, pageSize), nil
}

// ContentsURL constructs the complete directory of file content query github
// API endpoint.
func ContentsURL(proposalToken, repoOwner, repoName string) (string, error) {
	if proposalToken == "" {
		return "", fmt.Errorf("missing proposal token")
	}
	return fmt.Sprintf(rawContentsURL, repoOwner, repoName, proposalToken), nil
}
