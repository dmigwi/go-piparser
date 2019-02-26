package piclient

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	defaultRepo      = "mainnet"
	defaultRepoOwner = "decred-proposals"
	rawContentURL    = "https://api.github.com/repos/%s/%s/commits/%s"
	rawSHAListURL    = "https://api.github.com/repos/%s/%s/commits?path=%s"
)

func NewHTTPClient() *http.Client {
	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}
	return &http.Client{Transport: tr}
}

func GetRequestHandler(client *http.Client, URLPath string) ([]byte, error) {
	resp, err := client.Get(URLPath)
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
	if repoName == "" {
		repoName = defaultRepo
	}

	if repoOwner == "" {
		repoOwner = defaultRepoOwner
	}

	if contentSHA == "" {
		return "", fmt.Errorf("missing content SHA")
	}

	return fmt.Sprintf(rawContentURL, repoOwner, repoName, contentSHA), nil
}

func SHAListURL(proposalToken, repoOwner, repoName string) (string, error) {
	if repoName == "" {
		repoName = defaultRepo
	}

	if repoOwner == "" {
		repoOwner = defaultRepoOwner
	}

	if proposalToken == "" {
		return "", fmt.Errorf("missing proposal token")
	}

	return fmt.Sprintf(rawSHAListURL, repoOwner, repoName, proposalToken), nil
}
