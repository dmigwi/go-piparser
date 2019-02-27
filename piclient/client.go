package piclient

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	rawContentURL = "https://api.github.com/repos/%s/%s/commits/%s"
	rawSHAListURL = "https://api.github.com/repos/%s/%s/commits?path=%s/%s/plugins/decred/ballot.journal&page=%d&per_page=%d"

	dirContentsURL = "https://api.github.com/repos/%s/%s/contents/%s"
)

func NewHTTPClient() *http.Client {
	tr := &http.Transport{
		MaxIdleConns:       5,
		IdleConnTimeout:    10 * time.Second,
		DisableCompression: true,
	}
	return &http.Client{Transport: tr}
}

func GetRequestHandler(client *http.Client, URLPath string) ([]byte, error) {
	req, err := http.NewRequest("GET", URLPath, nil)
	if err != nil || req == nil {
		return nil, fmt.Errorf("http.NewRequest error : %v", err)
	}

	req.Header.Set("Authorization", "token 0a4384d36f93f6db4ae227fb186ea9e06ef6d771")

	resp, err := client.Do(req)
	if err != nil || resp == nil {
		return nil, fmt.Errorf("client.Get error : %v", err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ioutil.ReadAll error : %v", err)
	}

	return body, nil
}

func ContentURL(contentSHA, repoOwner, repoName string) (string, error) {
	if contentSHA == "" {
		return "", fmt.Errorf("missing content SHA")
	}
	return fmt.Sprintf(rawContentURL, repoOwner, repoName, contentSHA), nil
}

func SHAListURL(proposalToken, repoOwner, repoName, votesDirName string,
	page, pageSize int) (string, error) {
	if proposalToken == "" {
		return "", fmt.Errorf("missing proposal token")
	}
	return fmt.Sprintf(rawSHAListURL, repoOwner, repoName, proposalToken,
		votesDirName, page, pageSize), nil
}

func PropDirContentsURL(proposalToken, repoOwner, repoName string) (string, error) {
	if proposalToken == "" {
		return "", fmt.Errorf("missing proposal token")
	}
	return fmt.Sprintf(dirContentsURL, repoOwner, repoName, proposalToken), nil
}
