// Copyright 2019 Migwi Ndung'u.
// License that can be found in the LICENSE file.

// Package types defines the data structures and regular expressions that helps
// to unmarshal commits history string into data to be shared with the outside world.
package types

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

const (
	// DefaultRepo is the default github repository name where Politea Votes
	// are stored.
	DefaultRepo = "mainnet"

	// DefaultRepoOwner is the owner of the default github repository where
	// Politeia votes are stored.
	DefaultRepoOwner = "decred-proposals"

	// DefaultVotesCommitMsg defines the message of the commits that holds
	// the votes data for the various proposal token(s).
	DefaultVotesCommitMsg = "Flush vote journals"

	// CmdDateFormat defines the date format of the time returned by git commandline
	// interface. Time format is known as RFC2822.
	CmdDateFormat = "Mon Jan 2 15:04:05 2006 -0700"

	// journalActionFormat is the format of the journal action struct appended
	// to all votes. Its a struct with the version and the journal action.
	journalActionFormat = `{"version":"[[:digit:]]*","action":"(add)?(del)?(addlike)?"}`

	// RemoteURL uses the https protocol instead of git or ssh protocol. git
	// protocol may be faster but requires a dedicated port (9418) to be
	// open always. ssh requires authentication which is clearly not necessary
	// in this case scenario. Find out more on the access protcols here:
	// https://git-scm.com/book/en/v2/Git-on-the-Server-The-Protocols
	RemoteURL = "https://github.com/%s/%s.git"

	// DirPrefix defines the temporary folder prefix.
	DirPrefix = "go-piparser"

	// CloneRepoAlias defines the clone repository alias used by default instead
	// of the actual repo name.
	CloneRepoAlias = "prop-repo"

	// RemoteURLRef references the remote URL used to clone the repository.
	// 'origin' is the default set.
	RemoteURLRef = "origin"
)

var proposalToken string

// DataSource defines the required methods needed to query the proposals data.
type DataSource interface {
	SetUpEnv() error
	PullData(proposalToken string, since ...time.Time) ([]*History, error)
	FetchProporties() (owner, name, cloneDir string)
}

// semVer defines the semantic version structure.
type semVer struct {
	Major int
	Minor int
	Patch int
}

// minGitVersion defines the minimum supported git version. This version
// allows git log -p --reverse to be run. Here are the release notes:
// https://github.com/git/git/blob/53f9a3e157dbbc901a02ac2c73346d375e24978c/Documentation/RelNotes/1.5.1.txt
var minGitVersion = semVer{1, 5, 1}

// ErrGitVersion is the default error returned if an invalid git version was found.
var ErrGitVersion = errors.New("invalid git version found. A minimum of v" +
	minGitVersion.String() + " was expected")

// String() is the default stringer for the semVer data type.
func (s semVer) String() string {
	return fmt.Sprintf("%d.%d.%d", s.Major, s.Minor, s.Patch)
}

// Confirm that Votes implements the unmarshalling interface.
var _ json.Unmarshaler = (*Votes)(nil)

// History defines the standard single commit history contents to be shared
// with the outside world.
type History struct {
	Author    string
	CommitSHA string
	Date      time.Time
	Patch     []*File
}

// File defines the votes cast for a single token in a commit. A commit can
// votes cast for several commits joined together.
type File struct {
	Token     string
	VotesInfo Votes
}

// Votes defines a slice type of votes cast data.
type Votes []CastVoteData

// CastVoteData defines the struct of a cast vote and the reciept response.
type CastVoteData struct {
	*PiVote `json:"castvote"`
	// Receipt string `json:"receipt"`
}

// PiVote defines the ticket hash and vote bit type details about a vote.
type PiVote struct {
	// Token     string  `json:"token"`
	Ticket  string  `json:"ticket"`
	VoteBit bitCast `json:"votebit"`
	// Signature string  `json:"signature"`
}

// bitCast defines the votebit cast.
type bitCast string

// UnmarshalJSON defines the bitCast unmarshaller that sets the vote id of the
// vote bit cast.
func (b *bitCast) UnmarshalJSON(d []byte) error {
	// bitCast data mapping.
	var data = map[bitCast]string{
		`"1"`: "No",
		`"2"`: "Yes",
		`"3"`: "Unknown", // invalid entry. Wherever its found something went wrong.
	}
	vote, ok := data[bitCast(d)]
	if !ok {
		vote = "Unknown"
	}

	*b = bitCast(vote)
	return nil
}

// ToBitcast casts the vote Id bitCast type.
func ToBitcast(in string) bitCast {
	return bitCast(in)
}

// UnmarshalJSON defines the default unmarshaller for Votes.
func (v *Votes) UnmarshalJSON(b []byte) error {
	// create a custom unmarshalling type to avoid being trapped in
	// an endless loop.
	type votes2 Votes
	var v2 votes2

	err := json.Unmarshal(b, &v2)
	if err != nil {
		return err
	}

	*v = Votes(v2)

	return nil
}

// CustomUnmashaller unmarshals the string argument passed. Its not in a JSON
// format. History unmarshalling happens ONLY for the set proposal token and
// for all proposal tokens available if otherwise (not set).
func CustomUnmashaller(h *History, str string, since ...time.Time) error {
	// If no votes data detected, ignore the current str payload.
	if isMatched := IsMatching(str, DefaultVotesCommitMsg); !isMatched {
		return nil
	}

	date, err := RetrieveCMDDate(str)
	if err != nil {
		return err // Missing Date
	}

	if len(since) > 0 && date.Equal(since[0]) {
		// If this date matches the date in the record being unmarshalled
		// then it already existed earlier on thus ignore it.
		return nil
	}

	commit, err := RetrieveCMDCommit(str)
	if err != nil {
		return err // Missing commit SHA
	}

	author, err := RetrieveCMDAuthor(str)
	if err != nil {
		return err // Missing Author
	}

	var changes []*File

	for _, filePatch := range SplitCommitDiff(str) {

		// If the proposal token has been set, check if this payload has the required
		// proposal token data. If it exists proceed otherwise ignore it.
		if isMatched := IsMatching(filePatch, VotesJSONSignature()); !isMatched {
			continue
		}

		proposalToken, err := RetrieveProposalToken(filePatch)
		if err != nil {
			return err // Missing proposal token
		}

		filePatch = RetrieveAllPatchSelection(filePatch)

		filePatch = ReplaceJournalSelection(filePatch, "")

		// Drop any special characters left.
		filePatch = ReplaceAny(filePatch, `\s`, "")

		// Add the square brackets and commas to complete the JSON string array format.
		filePatch = "[" + ReplaceAny(filePatch, "}{", "},{") + "]"

		var v Votes

		if err = json.Unmarshal([]byte(filePatch), &v); err != nil {
			panic(err)
			return fmt.Errorf("Unmarshalling File failed: %v %s", err, filePatch)
		}

		// If votes data was found, append it the File patch data else ignore it.
		if len(v) > 0 {
			changes = append(changes, &File{proposalToken, v})
		}
	}

	if len(changes) == 0 {
		return nil
	}

	h.Author = author
	h.CommitSHA = commit
	h.Date = date
	h.Patch = changes

	return nil
}

// SetProposalToken sets the current proposal token string whose data is being
// unmarshalled.
func SetProposalToken(token string) error {
	if len(token) == 0 {
		return fmt.Errorf("empty token hash string found")
	}

	proposalToken = token
	return nil
}

// GetProposalToken returns the current proposal token value set.
func GetProposalToken() string {
	return proposalToken
}

// ClearProposalToken deletes the current proposal token value.
func ClearProposalToken() {
	proposalToken = ""
}
