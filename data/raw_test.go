package data

import (
	"reflect"
	"strings"
	"testing"

	"github.com/dmigwi/go-piparser/v1/types"
)

// TestUnmarshallingHistory uses the data.RawGitCommit stored in data/raw.go file
// to test if the parser tool can unmarshall the input data correctly into a
// History struct that can be shared with the outside world. The returned result
// of the unmarshalling test compared with data.VotesData stored in
// data/processed.go.
func TestUnmarshallingHistory(t *testing.T) {
	// currentProposalToken sample data being unmarshalled belong to this proposal
	// token.
	currentProposalToken := "27f87171d98b7923a1bd2bee6affed929fa2d2a6e178b5c80a9971a92a5c7f50"

	types.SetProposalToken(currentProposalToken)
	types.SetJournalActionFormat()

	var hist []*types.History
	commits := strings.Split(RawGitCommit, "commit")

	for _, c := range commits {
		// strings.Split returns some split strings as empty or with just
		// whitespaces and other special charactes. This happens when
		// the seperating argument is the first in the source string or is
		// surrounded by whitespaces and other special characters.
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

	if !reflect.DeepEqual(hist, VotesData) {
		t.Fatalf("expected the returned history to be equal to data.VotesData but it wasn't")
	}
}
