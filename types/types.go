// Copyright 2019 Migwi Ndung'u.
// License that can be found in the LICENSE file.

// Package types defines the data structures and regular expressions that helps
// to unmarshal commits history string into data to be shared with the outside world.
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
	// the votes data for the various proposal token(s).
	DefaultVotesCommitMsg = "Flush vote journals"

	// CmdDateFormat defines the date format of the time returned by git commandline
	// interface. Time format is known as RFC2822.
	CmdDateFormat = "Mon Jan 2 15:04:05 2006 -0700"
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

// CustomUnmashaller unmarshals the string argument passed. Its not in JSON
// format. History unmarshalling happens ONLY for the set proposal token and
// for none if otherwise (not set).
func CustomUnmashaller(h *History, str string, since ...time.Time) error {
	if isMatched := IsMatching(str, VotesJSONSignature()); !isMatched {
		// Required string payload could not be matched.
		return nil
	}

	date, err := RetrieveCMDDate(str)
	if err != nil {
		return err // Missing Date
	}

	if len(since) > 0 && since[0] == date {
		// It this date matches then the record being marshalled then it already
		// exists thus ignore it.
		return nil
	}

	proposalToken, err := RetrieveProposalToken(str)
	if err != nil {
		return err // Missing proposal token
	}

	commit, err := RetrieveCMDCommit(str)
	if err != nil {
		return err // Missing commit SHA
	}

	author, err := RetrieveCMDAuthor(str)
	if err != nil {
		return err // Missing Author
	}

	str = RetrieveAllPatchSelection(str)

	str = ReplaceJournalSelection(str, "")
	// Add the square brackets to complete the JSON string array format.
	str = "[" + str + "]"

	var v Votes

	if err = json.Unmarshal([]byte(str), &v); err != nil {
		return fmt.Errorf("Unmarshalling Votes failed: %v", err)
	}

	h.Author = author
	h.CommitSHA = commit
	h.Date = date
	h.Token = proposalToken
	h.VotesInfo = v

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

// SetJournalActionFormat sets journal (struct with the version and the journal
// action) format to use in the regexp.
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
