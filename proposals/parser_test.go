package proposals

import (
	"strconv"
	"testing"
)

// TestNewAPIExplorer tests assignment of the default values.
func TestNewAPIExplorer(t *testing.T) {
	type testData struct {
		repoName    string
		repoOwner   string
		accessToken string
		isError     bool
	}

	var (
		rName           = "samaki"
		rOwner          = "dmigwi"
		testAccessToken = "sample"
		defaultErrorMsg = "empty github access token found (https://developer.github.com/v3/#rate-limiting)"
	)

	td := []testData{
		{"", "", "", true},
		{rName, "", "", true},
		{"", rOwner, testAccessToken, false},
		{rName, rOwner, testAccessToken, false},
	}

	// repoName and repoOwner are already set but data is not available to the
	// public because of the interface use.
	for i, val := range td {
		t.Run("Test_#"+strconv.Itoa(i), func(t *testing.T) {
			_, err := NewAPIExplorer(val.accessToken, val.repoOwner, val.repoName)
			var errMsg string
			if err != nil {
				errMsg = err.Error()
			}

			if val.isError && errMsg != defaultErrorMsg {
				t.Fatalf("expected to find '%s' but found '%s'", defaultErrorMsg, errMsg)
			}
		})
	}
}
