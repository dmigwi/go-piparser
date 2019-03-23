// Copyright 2019 Migwi Ndung'u.
// License that can be found in the LICENSE file.

package types

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

// All regular expressions used in this tool are documented and implemented here.
// More on the regular expressions syntax used can be found here,
// https://github.com/google/re2/wiki/Syntax.

// PiRegExp helps defines the various regex expression supported. It also helps
// to easily compile them.
type PiRegExp string

var (
	// cmdAuthorSelection matches a text line that starts with 'Author' and ends
	// with line ending character(s) or its the actual end of the line.
	cmdAuthorSelection PiRegExp = `Author[:\s]*(.*)`

	// cmdCommitSelection matches a text line that starts with 'commit' or a
	// white space character and ends with line ending character(s) or its the
	// actual end of the line. The commit SHA part will always be the start of
	// the commit message after the whole git cmd history string is split into
	// individual messages.
	cmdCommitSelection PiRegExp = `[(^ )commit]*[:\s]*(.*)`

	// cmdDateSelection matches a text line that starts with 'Date' and ends with
	// line ending character(s) or its the actual end of the line.
	cmdDateSelection PiRegExp = `Date[:\s]*(.*)`

	// journalSelection matches the vote journal text line that takes the format,
	// +{"version":"\d","action":"(add|del|addlike)"} e.g +{"version":"1","action":"add"}
	// This journal section is appended to every individual vote cast result.
	journalSelection = func() PiRegExp {
		return PiRegExp(`[+]` + journalActionFormat)
	}

	// patchSelection matches the flushed votes changes pushed for the current
	// commit SHA. It matches all valid votes for the current time and commit.
	// It starts where the journalSelection on a text line matches and ends where
	// the line ending characters at matched at on the same text line.
	patchSelection = func() PiRegExp {
		return PiRegExp(journalSelection() + `.*`)
	}

	// anyTokenSelection matches any proposal token. A proposal token is
	// defined by 64 alphanumeric characters which can be upper case or lower
	// case of any letter exclusive of punctuations and white space characters.
	anyTokenSelection PiRegExp = `[A-z0-9]{64}`
)

// VotesJSONSignature defines a part of the json string signature that matches
// the commit patch string required. The matched commit patch string contains
// the needed votes data.
func VotesJSONSignature() string {
	if proposalToken == "" {
		return fmt.Sprintf(`{"castvote":{"token":"%s",`, anyTokenSelection)
	}
	return fmt.Sprintf(`{"castvote":{"token":"%s",`, proposalToken)
}

// exp compiles the PiRegExp regex expression type.
func (e PiRegExp) exp() *regexp.Regexp { return regexp.MustCompile(string(e)) }

// RetrieveCMDAuthor uses cmdAuthorSelection regex expression to retrieve the
// Author value in the provided parent string.
func RetrieveCMDAuthor(parent string) (string, error) {
	data := cmdAuthorSelection.exp().FindStringSubmatch(parent)
	if len(data) > 1 && data[1] != "" {
		return data[1], nil
	}
	return "", fmt.Errorf("missing Author from the parsed string")
}

// RetrieveCMDDate uses cmdDateSelection regex expression to retrieve the Date
// value in the provided parent string. The fetched date string is converted into
// a time.Time objected using "Mon Jan 2 15:04:05 2006 -0700" date format.
func RetrieveCMDDate(parent string) (time.Time, error) {
	data := cmdDateSelection.exp().FindStringSubmatch(parent)
	if len(data) > 1 && data[1] != "" {
		return time.Parse(CmdDateFormat, data[1])
	}
	return time.Time{}, fmt.Errorf("missing Date from the parsed string")
}

// RetrieveCMDCommit uses cmdCommitSelection to retrieve the SHA commit value
// from the provided parent string.
func RetrieveCMDCommit(parent string) (string, error) {
	data := cmdCommitSelection.exp().FindStringSubmatch(parent)
	if len(data) > 1 && data[1] != "" {
		return data[1], nil
	}
	return "", fmt.Errorf("missing commit from the parsed string")
}

// ReplaceJournalSelection uses journalSelection regex expression to replace the
// journal action in the provided parent string using the provided replacement.
func ReplaceJournalSelection(parent, with string) string {
	return journalSelection().exp().ReplaceAllLiteralString(parent, with)
}

// RetrieveAllPatchSelection uses patchSelection regex expression to fetch all
// individual matching lines from the provided parent string. All the individually
// matched strings are commas separated to form one single string.
func RetrieveAllPatchSelection(parent string) string {
	matches := patchSelection().exp().FindAllString(parent, -1)
	return strings.Join(matches, ",")
}

// IsMatching returns boolean true if the matchRegex can be matched in the parent
// string.
func IsMatching(parent, matchRegex string) bool {
	isMatched, err := regexp.MatchString(matchRegex, parent)
	if !isMatched || err != nil {
		return false
	}
	return true
}
