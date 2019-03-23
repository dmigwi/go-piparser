package data

import (
	"reflect"
	"strings"
	"testing"

	"github.com/dmigwi/go-piparser/v1/types"
)

// TestUnmarshalSingleTokenHistory uses the data.RawGitCommit stored in data/raw.go
// file to test if the parser tool can unmarshal the input data for the set token
// correctly into data that can be shared with the outside world. The returned
// results of the unmarshalling test is compared with data.SingleTokenVotesData
// stored in data/processed.go.
func TestUnmarshalSingleTokenHistory(t *testing.T) {
	// currentProposalToken is the proposal token whose sample data is being
	// unmarshalled.
	currentProposalToken := "27f87171d98b7923a1bd2bee6affed929fa2d2a6e178b5c80a9971a92a5c7f50"

	types.SetProposalToken(currentProposalToken)
	types.SetJournalActionFormat()

	var hist []*types.History
	commits := strings.Split(RawGitCommit, "commit")

	for _, c := range commits {
		if len(strings.TrimSpace(c)) == 0 {
			continue
		}

		var h types.History
		if err := types.CustomUnmashaller(&h, c); err != nil {
			if err != nil {
				t.Fatalf("expected to find no error but found: %v", err)
				return
			}
		}
		hist = append(hist, &h)
	}

	if !reflect.DeepEqual(hist, SingleTokenVotesData) {
		t.Fatalf("expected the returned history to be equal to data.SingleTokenVotesData but it wasn't")
	}
}

// TestUnmarshalAllTokensHistory uses the data.RawGitCommit stored in data/raw.go
// file to test if the parser tool can unmarshal the input data for all tokens
// correctly into data that can be shared with the outside world. The returned
// results of the unmarshalling test is compared with data.AllTokensVotesData
// stored in data/processed.go.
func TestUnmarshalAllTokensHistory(t *testing.T) {
	types.ClearProposalToken()
	types.SetJournalActionFormat()

	var hist []*types.History
	commits := strings.Split(RawGitCommit, "commit")

	for _, c := range commits {
		if len(strings.TrimSpace(c)) == 0 {
			continue
		}

		var h types.History
		if err := types.CustomUnmashaller(&h, c); err != nil {
			if err != nil {
				t.Fatalf("expected to find no error but found: %v", err)
				return
			}
		}
		hist = append(hist, &h)
	}

	if !reflect.DeepEqual(hist, AllTokensVotesData) {
		t.Fatalf("expected the returned history to be equal to data.AllTokensVotesData but it wasn't")
	}
}
