package types

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	//"github.com/decred/politeia/politeiad/backend/gitbe"
)

type HistoryList struct {
	SHA        string `json:"sha"`
	Commit     Commit `json:"commit"`
	APIURLPath string `json:"url"`
}

type CommitInfo struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Date  string `json:"date"`
}

type Commit struct {
	Committer CommitInfo `json:"committer"`
	Message   string     `json:"message"`
}

type Content struct {
	FileSHA    string `json:"sha"`
	FileName   string `json:"filename"`
	Status     string `json:"status"`
	Additions  int32  `json:"additions"`
	Deletions  int32  `json:"deletions"`
	Change     int32  `json:"changes"`
	BlobURL    string `json:"blob_url"`
	RawURL     string `json:"raw_url"`
	ContentURL string `json:"contents_url"`
	Data       *Votes `json:"patch"`
}

type History struct {
	SHA        string    `json:"sha"`
	Commit     Commit    `json:"commit"`
	APIURLPath string    `json:"url"`
	Files      []Content `json:"files"`
}

type PiVote struct {
	Token     string `json:"token"`
	Ticket    string `json:"ticket"`
	VoteBit   string `json:"votebit"`
	Signature string `json:"signature"`
}
type Votes []CastVoteData

type CastVoteData struct {
	Vote    *PiVote `json:"castvote"`
	Receipt string  `json:"receipt"`
}

var journalActionFormat string
var ProposalToken string

func (v *Votes) UnmarshalJSON(b []byte) error {
	str := string(b)
	isMatched, err := regexp.MatchString(ProposalToken, str)
	if !isMatched || err != nil {
		return err
	}

	// Delete the special characters indicating addition and deletion metrics.
	r := regexp.MustCompile(`(@{2}[\s\S]*@{2})`)
	str = r.ReplaceAllLiteralString(str, "")
	str, _ = strconv.Unquote(str)

	str = "[" + str + "]"

	journalActionFormat = `([[][\s]*{"version":"\d","action":"[add]*[del]*"})`
	r = regexp.MustCompile(string(journalActionFormat))
	str = r.ReplaceAllLiteralString(str, "[")

	journalActionFormat = `(}[\s+]*{"version":"\d","action":"[add]*[del]*"})`
	r = regexp.MustCompile(string(journalActionFormat))
	str = r.ReplaceAllLiteralString(str, "},")

	type votes2 Votes
	var v2 votes2

	err = json.Unmarshal([]byte(str), &v2)
	if err != nil {
		return err
	}

	fmt.Println(" >>>> ", len(v2))

	*v = Votes(v2)

	return nil
}
