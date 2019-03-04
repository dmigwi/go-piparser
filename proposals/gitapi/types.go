package gitapi

import (
	"encoding/json"
	"fmt"
	"strconv"

	pitypes "github.com/dmigwi/go-piparser/v1/types"
)

// HistorySHAs holds a slice of the commit history SHA token strings.
type historySHAs []commitSHA

// commitSHA holds the specific commit unique SHA string value.
type commitSHA struct {
	SHA string `json:"sha"`
}

// rawHistory defines the commit full information about a commit
type rawHistory struct {
	SHA    string      `json:"sha"`
	Commit rawCommit   `json:"commit"`
	Files  []votesData `json:"files"`
	*errResponse
}

// rawCommit defines information about the committer and the commit message used.
type rawCommit struct {
	Committer commitInfo `json:"committer"`
}

// CommitInfo defines information about the committer
type commitInfo struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Date  string `json:"date"`
}

// votesData defines the changes made in the commit content.
type votesData struct {
	Data pitypes.Votes `json:"patch"`
}

// listDirs defines the directories in the name of the repository directories.
// The directories in github hold data about the proposals token and thier metadata.
type listDirs struct {
	Name string `json:"name"`
}

// UnmarshalJSON is the default unmarshaller for HistorySHA.
func (h *historySHAs) UnmarshalJSON(b []byte) error {
	// Match the DefaultVotesCommitMsg string
	isMatched := pitypes.IsMatching(string(b), pitypes.DefaultVotesCommitMsg)
	if !isMatched {
		return nil
	}

	pToken := pitypes.GetProposalToken()

	// Match the proposalToken string
	if pToken != "" {
		if isMatched := pitypes.IsMatching(string(b), pToken); !isMatched {
			return fmt.Errorf("missing proposal token %s", pToken)
		}
	}

	// Create a custom unmarshalling type to avoid being trapped in endless loop.
	type history historySHAs
	var h2 history

	if err := json.Unmarshal(b, &h2); err != nil {
		return err
	}

	*h = historySHAs(h2)
	return nil
}

// UnmarshalJSON is the default unmarshaller for votesData struct.
func (v *votesData) UnmarshalJSON(b []byte) error {
	str := string(b)
	if isMatched := pitypes.IsMatching(str, pitypes.VotesJSONSignature()); !isMatched {
		// Required string payload could not be matched.
		return nil
	}

	str = pitypes.ReplaceAddnDelMetrics(str, "")

	str = pitypes.RetrieveAllPatchSelection(str)

	str = pitypes.ReplaceJournalSelection(str, "")

	// Drops github added newline special characters.
	str, _ = strconv.Unquote(str)

	// Add the square brackets to complete the JSON string array format.
	str = "[" + str + "]"

	// create a custom unmarshalling type to avoid being trapped in an endless loop.
	type votes2 votesData
	var v2 votes2

	if err := json.Unmarshal([]byte(str), &v2); err != nil {
		return err
	}

	*v = votesData(v2)

	return nil
}

// errResponse is the alternative response returned if the default one
// wasn't successful. It majorly contains error message details.
type errResponse struct {
	Message string `json:"message"`
	URL     string `json:"documentation_url"`
}

// constructAPIErr attempts to construct an error message from the data that can
// be marshalled into errResponse. Attempts to construct the error message are
// initiated if something went wrong in the unmarshalling of the expected API
// endpoint response. It returns a more descriptive from the github API.
func constructAPIErr(b []byte) error {
	var resp errResponse
	err := json.Unmarshal(b, &resp)
	if err != nil {
		// No suitable errResponse struct data found.
		return nil
	}
	return fmt.Errorf("%s (%s)", resp.Message, resp.URL)
}
