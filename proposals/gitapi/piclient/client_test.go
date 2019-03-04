package piclient

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
)

var testServer *httptest.Server

// TestMain creates a test httptest server.
func TestMain(m *testing.M) {
	testServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "OK")
	}))

	responseCode := m.Run()
	testServer.Close()

	os.Exit(responseCode)
}

// TestSetAccessToken tests the setting of the access token.
func TestSetAccessToken(t *testing.T) {
	errMsg := "empty github access token found (https://developer.github.com/v3/#rate-limiting)"
	t.Run("Test_set_empty_token", func(t *testing.T) {
		err := SetAccessToken("")
		if err != nil && err.Error() != errMsg {
			t.Fatalf("expected to find error message '%s' but found '%v'", errMsg, err)
		}
	})

	t.Run("Test_set_non_empty_token", func(t *testing.T) {
		err := SetAccessToken("sample-token")
		if err != nil {
			t.Fatalf("expect to find no err but found: %v", err)
		}
	})
}

// TestGetRequestHandler tests validation of the URLPath and the accessToken.
// It also tests the actual API request call.
func TestGetRequestHandler(t *testing.T) {
	type testData struct {
		client  *http.Client
		URLPath string

		Results []byte
		errMsg  string
	}

	client := testServer.Client()
	testURL := testServer.URL

	td := []testData{
		{client, "", []byte(""), "empty URL path found"},
		{client, testURL, []byte("OK"), ""},
	}

	for i, val := range td {
		t.Run("Test_#"+strconv.Itoa(i), func(t *testing.T) {
			resp, err := GetRequestHandler(val.client, val.URLPath)
			if err != nil && err.Error() != val.errMsg {
				t.Fatalf("expected to find error '%s' but found '%v", val.errMsg, err)
			}

			if err == nil && val.errMsg != "" {
				t.Fatalf("expected no error but found '%v'", val.errMsg)
			}

			if !bytes.Equal(resp, val.Results) {
				t.Fatalf("expected the returned result to be '%s' but found '%s",
					string(val.Results), string(resp))
			}
		})
	}
}

// TestCommitURL tests the functionality of CommitURL. The repoOwner and repoName
// are preset before the invoking of CommitURL.
func TestCommitURL(t *testing.T) {
	repoOwner := "dmigwi"
	repoName := "mainnet"
	baseURL := "https://api.github.com"
	expected1 := "https://api.github.com/repos/dmigwi/mainnet/commits"
	commitSHA := "eced3135d573509e4460af56d148f177498be122"
	parentCommitSHA := "97b8cba954735ec428f140b8d2a4bfa8"
	expected2 := "https://api.github.com/repos/dmigwi/mainnet/commits/" +
		"eced3135d573509e4460af56d148f177498be122#diff-97b8cba954735ec428f140b8d2a4bfa8"

	type testData struct {
		commit, parentCommit, URL string
	}

	td := []testData{
		{"", "", expected1},
		{commitSHA, "", expected1},
		{"", parentCommitSHA, expected1},
		{commitSHA, parentCommitSHA, expected2},
	}

	for i, val := range td {
		t.Run("Test_#"+strconv.Itoa(i), func(t *testing.T) {
			resp := CommitURL(baseURL, val.commit, val.parentCommit, repoOwner, repoName)
			if resp != val.URL {
				t.Fatalf("expected the returned result to be '%s' but found '%s' ",
					val.URL, resp)
			}
		})
	}
}

// TestSHAListURL tests the functionality of SHAListURL. The repoOwner, repoName,
// votesDirName, page and pageSizes values are preset before invoking SHAListURL.
func TestSHAListURL(t *testing.T) {
	repoOwner := "dmigwi"
	repoName := "mainnet"
	votesDirName := "3"
	baseURL := "https://api.github.com"
	page := 1
	pageSize := 20
	proposalToken := ""

	t.Run("Test_with_empty_proposal_token", func(t *testing.T) {
		resp, err := SHAListURL(baseURL, proposalToken, repoOwner, repoName, votesDirName, page, pageSize)
		if err == nil || err.Error() != "missing proposal token" {
			t.Fatalf("expected to find 'missing proposal token' error but found: %v", err)
		}

		if resp != "" {
			t.Fatalf("expected to find an empty an empty URL path but found '%s' ", resp)
		}
	})

	proposalToken = "27f87171d98b7923a1bd2bee6affed929fa2d2a6e178b5c80a9971a92a5c7f50"
	expected := "https://api.github.com/repos/dmigwi/mainnet/commits?path=" +
		"27f87171d98b7923a1bd2bee6affed929fa2d2a6e178b5c80a9971a92a5c7f50/3/plugins/decred/ballot.journal&page=1&per_page=20"

	t.Run("Test_with_non_empty_token", func(t *testing.T) {
		resp, err := SHAListURL(baseURL, proposalToken, repoOwner, repoName, votesDirName, page, pageSize)
		if err != nil {
			t.Fatalf("expected to find no error but found: %v", err)
		}

		if resp != expected {
			t.Fatalf("expected the returned result to be '%s' but found '%s' ", expected, resp)
		}
	})
}

// TestContentsURL tests the functionality of ContentsURL. The repoOwner and repoName
// are preset before the invoking of ContentsURL.
func TestContentsURL(t *testing.T) {
	repoOwner := "dmigwi"
	repoName := "mainnet"
	baseURL := "https://api.github.com"
	proposalToken := ""

	t.Run("Test_with_empty_proposal_token", func(t *testing.T) {
		resp, err := ContentsURL(baseURL, proposalToken, repoOwner, repoName)
		if err == nil || err.Error() != "missing proposal token" {
			t.Fatalf("expected to find 'missing proposal token' error but found: %v", err)
		}

		if resp != "" {
			t.Fatalf("expected to find an empty an empty URL path but found '%s' ", resp)
		}
	})

	proposalToken = "27f87171d98b7923a1bd2bee6affed929fa2d2a6e178b5c80a9971a92a5c7f50"
	expected := "https://api.github.com/repos/dmigwi/mainnet/contents/" +
		"27f87171d98b7923a1bd2bee6affed929fa2d2a6e178b5c80a9971a92a5c7f50"

	t.Run("Test_with_non_empty_token", func(t *testing.T) {
		resp, err := ContentsURL(baseURL, proposalToken, repoOwner, repoName)
		if err != nil {
			t.Fatalf("expected to find no error but found: %v", err)
		}

		if resp != expected {
			t.Fatalf("expected the returned result to be '%s' but found '%s' ", expected, resp)
		}
	})
}
