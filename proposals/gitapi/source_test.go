package gitapi

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
)

var testClient *http.Client
var testURL string

// TestMain sets up s tests server used to mock the http requests.
func TestMain(m *testing.M) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload string
		switch r.URL.Path {
		case "/repos/dmigwi/mainnet/commits":
			payload = ""
		case "/repos/dmigwi/mainnet/commits/d15fa2317cc5f37508a470413e09f11e2a5faa9c":
			payload = ""
		case "/repos/dmigwi/mainnet/contents/27f87171d98b7923a1bd2bee6affed929fa2d2a6e178b5c80a9971a92a5c7f50":
			payload = ""
		case "/repos/dmigwi/mainnet/commits?path=27f87171d98b7923a1bd2bee6affed929fa2d2a6e178b5c80a9971a92a5c7f50/3/plugins/decred/ballot.journal&page=1&per_page=2":
			payload = ""
		default:
			payload = ""
		}

		w.WriteHeader(http.StatusOK)
		io.WriteString(w, payload)
	}),
	)

	testClient = server.Client()
	testURL = server.URL

	returnCode := m.Run()

	// close the server
	server.Close()
	// exit with the returned code
	os.Exit(returnCode)
}

// TestNewParser tests assignment of the default values.
func TestNewParser(t *testing.T) {
	type testData struct {
		repoName   string
		repoOwner  string
		testClient *http.Client
	}

	rName := "samaki"
	rOwner := "dmigwi"

	td := []testData{
		{"", "", nil},
		{rName, "", nil},
		{"", rOwner, nil},
		{rName, rOwner, nil},
		{"", rOwner, testClient},
		{rName, "", testClient},
		{rName, rOwner, testClient},
	}

	for i, val := range td {
		expectedRName, expectedROwner := rName, rOwner

		t.Run("Test_#"+strconv.Itoa(i), func(t *testing.T) {
			testParser := NewParser(val.repoOwner, val.repoName, val.testClient)

			if val.repoName == "" {
				expectedRName = defaultRepo
			}
			if testParser.repoName != expectedRName {
				t.Fatalf("expected the repo name to be %s but found %s", expectedRName, testParser.repoName)
			}

			if val.repoOwner == "" {
				expectedROwner = defaultRepoOwner
			}
			if testParser.repoOwner != expectedROwner {
				t.Fatalf("expected the repo owner to be %s but found %s ", expectedROwner, testParser.repoOwner)
			}

			if testParser.client == nil {
				t.Fatal("expected the http client not to be nil but it was")
			}

			if val.testClient != nil && val.testClient != testClient {
				t.Fatal("expected the client set to be equal the testClient but i wasn't")
			}
		})
	}
}
