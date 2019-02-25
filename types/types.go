package types

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
	Data       string `json:"patch"`
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
	Receipt   string `json:"receipt"`
}

type CastVoteData struct {
	Vote PiVote `json:"castvote"`
}
