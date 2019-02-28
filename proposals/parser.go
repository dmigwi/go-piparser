package proposals

import (
	"encoding/json"

	"github.com/decred/politeia/politeiad/backend/gitbe"
	"github.com/dmigwi/go-piparser/v1/types"
)

// Initializes and sets the one JournalActionFormat variables. JournalActionFormat
// is a regex expression that helps eliminate unwanted parts of the vote information.
func init() {
	f, err := json.Marshal(gitbe.JournalAction{
		Version: `[[:digit:]]*`,
		Action:  "[add]*[del]*[addlike]*",
	})
	if err != nil {
		panic(err)
	}

	format := string(f)
	types.SetJournalActionFormat(format)
}
