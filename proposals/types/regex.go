// Copyright 2019 Migwi Ndung'u.
// License that can be found in the LICENSE file.

package types

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
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
		return PiRegExp(`(` + journalSelection() + `[[:ascii:]]*(}\n?))`)
	}

	// anyTokenSelection matches any proposal token. A proposal token is
	// defined by 64 alphanumeric characters which can be upper case or lower
	// case of any letter, exclusive of punctuations and white space characters.
	anyTokenSelection = `[A-z0-9]{64}`

	// In a git commit history, the changes made per file always start with
	// "diff --git a". commitDiff is therefore used to split the single commit
	// string into file changes in an array. "diff --git a" is documented here:
	// https://github.com/git/git/blob/b58f23b38a9a9f28d751311353819d3cdf6a86da/t/t4000-diff-format.sh#L29-L46
	commitDiff = `diff --git a`

	// gitVersionSelection selects the underlying platform git semantic version.
	gitVersionSelection PiRegExp = "([[:digit:]]+).([[:digit:]]+).([[:digit:]]+)"
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

// RetrieveCMDCommit uses cmdCommitSelection to retrieve the commit SHA value
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
// individual matching lines from the provided parent string.
func RetrieveAllPatchSelection(parent string) string {
	matches := patchSelection().exp().FindAllString(parent, -1)
	return strings.Join(matches, "")
}

// RetrieveProposalToken uses the anyTokenSelection regex to build a regex
// expression used to select the proposal token from the parent string.
func RetrieveProposalToken(parent string) (string, error) {
	regex := fmt.Sprintf(`"token":"(%s)`, anyTokenSelection)
	data := PiRegExp(regex).exp().FindStringSubmatch(parent)
	if len(data) > 1 && data[1] != "" {
		return data[1], nil
	}

	return "", fmt.Errorf("missing token from the parsed string")
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

// ReplaceAny replaces the all occurence of "regex" in string "parent" with
// replacement "with" for all the possible occurences.
func ReplaceAny(parent, regex, with string) string {
	r := regexp.MustCompile(regex)
	return r.ReplaceAllLiteralString(parent, with)
}

// SplitCommitDiff uses the commitDiff separating string to split the string
// into an array.
func SplitCommitDiff(parent string) []string {
	return strings.Split(parent, commitDiff)
}

// IsGitVersionSupported checks if the git version in the parse string is supported.
// An error is returned if otherwise.
func IsGitVersionSupported(parsedStr string) error {
	data := gitVersionSelection.exp().FindAllString(parsedStr, -1)
	data = strings.Split(strings.Join(data, ""), ".")
	if len(data) != 3 {
		return ErrGitVersion
	}

	currentVersion, err := parseVersion(data)
	if err != nil {
		return fmt.Errorf("%v: %s", err, parsedStr)
	}

	// for each semantic version field a max of four digits value are reserved.
	minimumRequired := int64(100050001) // => 1.5.1

	if minimumRequired > currentVersion {
		return ErrGitVersion
	}

	return nil
}

// parseVersion converts the semantic version into an int value that can be compared.
func parseVersion(strList []string) (int64, error) {
	str := ""
	for _, value := range strList {
		str = fmt.Sprintf("%s%04s", str, value)
	}

	result, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return -1, errors.New("error parsing the git version")
	}
	return result, nil
}
