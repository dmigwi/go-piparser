// Copyright 2019 Migwi Ndung'u.
// License that can be found in the LICENSE file.

// Package types defines global data structures (used in gitcmd and gitapi) and
// regular expressions that helps unmarshal commit history string into History
// struct that can be shared with the outside world.
package types

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/decred/politeia/politeiad/backend/gitbe"
)

const (
	// DefaultRepo is the default github repository name where Politea Votes
	// are stored.
	DefaultRepo = "mainnet"

	// DefaultRepoOwner is the owner of the default github repository where
	// Politeia votes are stored.
	DefaultRepoOwner = "decred-proposals"

	// DefaultVotesCommitMsg defines the message of the commits that holds
	// the votes data for various proposal token(s).
	DefaultVotesCommitMsg = "Flush vote journals"

	// cmdDateFormat defines the date format returned by github via git cmd data
	// source.
	cmdDateFormat = "Mon Jan 2 15:04:05 2006 -0700"
)

var journalActionFormat, proposalToken string

// Confirm that Votes implements the unmarshalling interface.
var _ json.Unmarshaler = (*Votes)(nil)

// History defines the standard single commit history contents to be shared
// with the outside world.
type History struct {
	Author    string
	CommitSHA string
	Date      time.Time
	VotesInfo Votes
}

// Votes defines a slice type of all votes cast data.
type Votes []CastVoteData

// CastVoteData defines the struct of a cast vote and the reciept response.
type CastVoteData struct {
	*PiVote `json:"castvote"`
	// Receipt string `json:"receipt"`
}

// PiVote defines the finer details about a vote.
type PiVote struct {
	// Token     string  `json:"token"`
	Ticket  string  `json:"ticket"`
	VoteBit bitCast `json:"votebit"`
	// Signature string  `json:"signature"`
}

// bitCast defines the votebit cast.
type bitCast string

// String is the default string for bitCast.
func (b bitCast) String() string {
	var data = map[bitCast]string{
		"1": "No",
		"2": "Yes",
	}
	if vote, ok := data[b]; ok {
		return vote
	}
	return "Unknown"
}

// UnmarshalJSON defines the global unmarshaller for Votes in package gitapi and
// gitcmd. The votes unmarshalling happens for all tokens in the current commit
// message string unless the specific proposal token is set.
func (v *Votes) UnmarshalJSON(b []byte) error {
	// create a custom unmarshalling type to avoid being trapped in endless loop.
	type votes2 Votes
	var v2 votes2

	err := json.Unmarshal(b, &v2)
	if err != nil {
		return err
	}

	*v = Votes(v2)

	return nil
}

// CustomUnmashaller is the default unmarshaller for the History. The string
// argument passed here is not a valid json string.
func (h *History) CustomUnmashaller(str string) error {
	if isMatched := IsMatching(str, VotesJSONSignature()); !isMatched {
		// Required string payload could not be matched.
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

	date, err := RetrieveCMDDate(str)
	if err != nil {
		return err // Missing Date
	}

	str = RetrieveAllPatchSelection(str)

	str = ReplaceJournalSelection(str, "")

	// Add the square brackets to complete the JSON string array format.
	str = "[" + str + "]"

	var v Votes

	if err = json.Unmarshal([]byte(str), &v); err != nil {
		return fmt.Errorf("Unmarshalling Votes failed: %v", err)
	}

	// Do not store any empty votes data.
	if len(v) == 0 {
		return nil
	}

	*h = History{
		Author:    author,
		CommitSHA: commit,
		Date:      date,
		VotesInfo: v,
	}

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

// GetProposalToken returns the current proposal token value.
func GetProposalToken() string {
	return proposalToken
}

// ClearProposalToken deletes the outdated proposal token value.
func ClearProposalToken() {
	proposalToken = ""
}

// SetJournalActionFormat sets journal (struct with the version and the journal
// action) format to use for the regexp.
func SetJournalActionFormat() {
	f, err := json.Marshal(gitbe.JournalAction{
		Version: `[[:digit:]]*`,
		Action:  "(add)?(del)?(addlike)?",
	})
	if err != nil {
		panic(err)
	}

	journalActionFormat = string(f)
}
