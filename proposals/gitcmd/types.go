package gitcmd

import (
	"encoding/json"
	"fmt"

	pitypes "github.com/dmigwi/go-piparser/v1/types"
)

// copyHistory defines a local type of the global types.History struct. It is
// used to unmarshal the gitcmd commit history payload which unique to gitcmd.
type copyHistory pitypes.History

// customUnmashaller is the default unmarshaller for copyHistory used in package
// gitcmd. The string argument passed here is not a valid json string.
func (h *copyHistory) customUnmashaller(str string) error {
	if isMatched := pitypes.IsMatching(str, pitypes.VotesJSONSignature()); !isMatched {
		// Required string payload could not be matched.
		return nil
	}

	commit, err := pitypes.RetrieveCMDCommit(str)
	if err != nil {
		return err // Missing commit SHA
	}

	author, err := pitypes.RetrieveCMDAuthor(str)
	if err != nil {
		return err // Missing Author
	}

	date, err := pitypes.RetrieveCMDDate(str)
	if err != nil {
		return err // Missing Date
	}

	str = pitypes.RetrieveAllPatchSelection(str)

	str = pitypes.ReplaceJournalSelection(str, "")

	// Add the square brackets to complete the JSON string array format.
	str = "[" + str + "]"

	var v pitypes.Votes

	if err = json.Unmarshal([]byte(str), &v); err != nil {
		return fmt.Errorf("Unmarshalling pitypes.Votes failed: %v", err)
	}

	// Do not store any empty votes data.
	if len(v) == 0 {
		return nil
	}

	*h = copyHistory{
		Author:    author,
		CommitSHA: commit,
		Date:      date,
		VotesInfo: v,
	}

	return nil
}
