package types

import (
	"os"
	"strconv"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	SetJournalActionFormat()

	os.Exit(m.Run())
}

// testData defines the general struct with possible inputs and outputs that can
// be used to test the various regex expressions.
type testData struct {
	src, repl, output string
	isError           bool
	timestamp         time.Time
}

func TestReplaceLineEndingChars(t *testing.T) {
	td := []testData{
		{
			src:    "Hello, World",
			repl:   "xx",
			output: "Hello, World",
		},
		{
			src: ` Flush vote journals.
			`,
			repl: "xx",
			output: " Flush vote journals.xx			",
		},
		{
			src: `Politeia 
			is 
			the 
			Decred 
			proposal 
			system. `,
			repl: "xx",
			output: `Politeia xx			is xx			the xx			Decred xx			proposal xx			system. `,
		},
		{
			src: `line 
			
			
			
			break`,
			repl: "xx",
			output: "line xx			xx			xx			xx			break",
		},
	}

	for i, val := range td {
		t.Run("Test_#"+strconv.Itoa(i), func(t *testing.T) {
			result := ReplaceLineEndingChars(val.src, val.repl)
			if val.output != result {
				t.Fatalf("expected the returned string to be equal to '%s' but was '%s'",
					val.output, result)
			}
		})
	}
}

func TestRetrieveCMDAuthor(t *testing.T) {
	td := []testData{
		{
			src:     "Politeia is the Decred proposal system. ",
			output:  "",
			isError: true,
		},
		{
			src:     "Author:",
			output:  "",
			isError: true,
		},
		{
			src: `Author: Politeia <noreply@decred.org>
			`,
			output: "Politeia <noreply@decred.org>",
		},
		{
			src: `Author: Politeia <noreply@decred.org> 
			Author: Politeia <migwindungu0@gmail.com>`,
			output: "Politeia <noreply@decred.org> ",
		},
		{
			src:    `Author: Politeia <noreply@decred.org> Author: Politeia <migwindungu0@gmail.com>`,
			output: "Politeia <noreply@decred.org> Author: Politeia <migwindungu0@gmail.com>",
		},
	}

	for i, val := range td {
		t.Run("Test_#"+strconv.Itoa(i), func(t *testing.T) {
			result, err := RetrieveCMDAuthor(val.src)
			if err == nil && val.isError {
				t.Fatalf("expected to find an error but none was returned")
			}

			if !val.isError && err != nil {
				t.Fatalf("expected no error but '%v' was returned ", err)
			}

			if result != val.output {
				t.Fatalf("expected the returned Author to be equal to '%s' but found '%s'",
					val.output, result)
			}
		})
	}
}

func TestRetrieveCMDDate(t *testing.T) {
	date, _ := time.Parse(cmdDateFormat, "Thu Feb 28 15:35:56 2019 -0600")
	defaultTime := time.Time{}

	td := []testData{
		{
			src:       "Today is a good day.",
			timestamp: defaultTime,
			isError:   true,
		},
		{
			src:       "Date:",
			timestamp: defaultTime,
			isError:   true,
		},
		{
			src: `Date: Thu Feb 28 15:35:56 2019 -0600
			`,
			timestamp: date,
		},
		{
			src: `Date: Thu Feb 28 15:35:56 2019 -0600
			Author: Politeia <migwindungu0@gmail.com>`,
			timestamp: date,
		},
		{
			src:       `Date: Thu Feb 28 15:35:56 2019 -0600 Author: Politeia <migwindungu0@gmail.com> Date: Fri Mar 1 01:27:13 2019 +0300`,
			timestamp: defaultTime,
			isError:   true,
		},
	}

	for i, val := range td {
		t.Run("Test_#"+strconv.Itoa(i), func(t *testing.T) {
			result, err := RetrieveCMDDate(val.src)
			if err == nil && val.isError {
				t.Fatalf("expected to find an error but none was returned")
			}

			if !val.isError && err != nil {
				t.Fatalf("expected no error but '%v' was returned ", err)
			}

			if result.UnixNano() != val.timestamp.UnixNano() {
				t.Fatalf("expected the returned Date to be equal to '%s' but found '%s'",
					val.timestamp, result)
			}
		})
	}
}

func TestRetrieveCMDCommit(t *testing.T) {
	td := []testData{
		{
			src:     "commit:",
			output:  "",
			isError: true,
		},
		{
			src:    " 855cad7c76087645a8f3c3525bb79513e35fe4ac",
			output: "855cad7c76087645a8f3c3525bb79513e35fe4ac",
		},
		{
			src:    "855cad7c76087645a8f3c3525bb79513e35fe4ac",
			output: "855cad7c76087645a8f3c3525bb79513e35fe4ac",
		},
		{
			src: `commit: 855cad7c76087645a8f3c3525bb79513e35fe4ac
			`,
			output: "855cad7c76087645a8f3c3525bb79513e35fe4ac",
		},
		{
			src: `commit: e79d60ca76187b35e723fe9f09fba6169e7de300 
			Author: Politeia <migwindungu0@gmail.com>`,
			output: "e79d60ca76187b35e723fe9f09fba6169e7de300 ",
		},
	}

	for i, val := range td {
		t.Run("Test_#"+strconv.Itoa(i), func(t *testing.T) {
			result, err := RetrieveCMDCommit(val.src)
			if err == nil && val.isError {
				t.Fatalf("expected to find an error but none was returned")
			}

			if !val.isError && err != nil {
				t.Fatalf("expected no error but '%v' was returned ", err)
			}

			if result != val.output {
				t.Fatalf("expected the returned commit SHA to be equal to '%s' but found '%s'",
					val.output, result)
			}
		})
	}
}

func TestReplaceJournalSelection(t *testing.T) {
	td := []testData{
		{
			src:    "Hello, World",
			repl:   "xx",
			output: "Hello, World",
		},
		{
			src:    `{"version":"1","action":"add"}{"castvote":{"token":"27f87171d98b7923a1bd2bee6affed929fa2d2a6e178b5c80a9971a92a5c7f50",`,
			repl:   "xx",
			output: `{"version":"1","action":"add"}{"castvote":{"token":"27f87171d98b7923a1bd2bee6affed929fa2d2a6e178b5c80a9971a92a5c7f50",`,
		},
		{
			src:    `+{"version":"1","action":"addlike"}{"castvote":{"token":"27f87171d98b7923a1bd2bee6affed929fa2d2a6e178b5c80a9971a92a5c7f50","`,
			repl:   "xx",
			output: `xx{"castvote":{"token":"27f87171d98b7923a1bd2bee6affed929fa2d2a6e178b5c80a9971a92a5c7f50","`,
		},
		{
			src:    `+{"version":"1","action":"del"}{"data":"blahblah"}+{"version":"1","action":"add"}`,
			repl:   "xx",
			output: `xx{"data":"blahblah"}xx`,
		},
		{
			src: `+{"version":"1","action":"add"}{"data":"blahblah"}
			+{"version":"1","action":"addlike"}`,
			repl: "xx",
			output: `xx{"data":"blahblah"}
			xx`,
		},
	}

	for i, val := range td {
		t.Run("Test_#"+strconv.Itoa(i), func(t *testing.T) {
			result := ReplaceJournalSelection(val.src, val.repl)
			if val.output != result {
				t.Fatalf("expected the returned string to be equal to '%s' but was '%s'",
					val.output, result)
			}
		})
	}
}

func TestRetrieveAllPatchSelection(t *testing.T) {
	td := []testData{
		{
			src:    `{"version":"1","action":"add"}{"castvote":{"token":"27f87171d98b7923a1bd2bee6affed929fa2d2a6e178b5c80a9971a92a5c7f50"}`,
			output: "", // match not found
		},
		{
			src:    `+{"version":"1","action":"add"}{"castvote":{"token":"27f87171d98b7923a1bd2bee6affed929fa2d2a6e178b5c80a9971a92a5c7f50"}`,
			output: `+{"version":"1","action":"add"}{"castvote":{"token":"27f87171d98b7923a1bd2bee6affed929fa2d2a6e178b5c80a9971a92a5c7f50"}`,
		},
		{
			src: `+{"version":"1","action":"add"}{"castvote":{"token":"27f87171d98b7923a1bd2bee6affed929fa2d2a6e178b5c80a9971a92a5c7f50"}
			`,
			output: `+{"version":"1","action":"add"}{"castvote":{"token":"27f87171d98b7923a1bd2bee6affed929fa2d2a6e178b5c80a9971a92a5c7f50"}`,
		},
		{
			src: `{"version":"1","action":"add"}{"castvote":{"token":"27f8","ticket":"dab3d276","votebit":"2","signature":"20942f2a"},"receipt":"547416f8f"}
			{"version":"1","action":"add"}{"castvote":{"token":"27f8","ticket":"8863036","votebit":"2","signature":"2066ca72"},"receipt":"5b15f2b3c2"}
			{"version":"1","action":"add"}{"castvote":{"token":"27f8","ticket":"f837268","votebit":"2","signature":"20d3d731"},"receipt":"ea6fb3c02f"}
		   +{"version":"1","action":"addlike"}{"castvote":{"token":"27f8","ticket":"14af967","votebit":"2","signature":"20e1e03d"},"receipt":"9a65775122"}
		   +{"version":"1","action":"del"}{"castvote":{"token":"27f8","ticket":"b4f7f4e","votebit":"2","signature":"1fd2e407"},"receipt":"02da25c951"}
		   +{"version":"1","action":"del"}{"castvote":{"token":"27f8","ticket":"c6d2633","votebit":"2","signature":"1ffec594"},"receipt":"0045001444"}
		   
		   `,
			output: `+{"version":"1","action":"addlike"}{"castvote":{"token":"27f8","ticket":"14af967","votebit":"2","signature":"20e1e03d"},"receipt":"9a65775122"},` +
				`+{"version":"1","action":"del"}{"castvote":{"token":"27f8","ticket":"b4f7f4e","votebit":"2","signature":"1fd2e407"},"receipt":"02da25c951"},` +
				`+{"version":"1","action":"del"}{"castvote":{"token":"27f8","ticket":"c6d2633","votebit":"2","signature":"1ffec594"},"receipt":"0045001444"}`,
		},
	}
	for i, val := range td {
		t.Run("Test_#"+strconv.Itoa(i), func(t *testing.T) {
			result := RetrieveAllPatchSelection(val.src)
			if val.output != result {
				t.Fatalf("expected the returned string to be equal to '%s' but was '%s'",
					val.output, result)
			}
		})
	}
}

func TestIsMatching(t *testing.T) {
	type testData struct {
		src, regex string
		isFound    bool
	}

	td := []testData{
		{src: `
		`, regex: "\n", isFound: true}, // matching newline character.
		{src: `
		`, regex: "\\n", isFound: true}, //matching escaped newline character.
		{src: `\n`, regex: "\\n", isFound: false},
		{src: `    Flush vote journals.
		`, regex: defaultVotesCommitMsg, isFound: true},
		{src: "b/27f87171d98b7923a1bd2bee6af/3/plugins/decred/ballot.journal",
			regex: "27f87171d98b7923a1bd2bee6af", isFound: true},
	}

	for i, val := range td {
		t.Run("Test_#"+strconv.Itoa(i), func(t *testing.T) {
			result := IsMatching(val.src, val.regex)
			if result != val.isFound {
				t.Fatalf("expected the matching src to the regex to be %v but found %v",
					val.isFound, result)
			}
		})
	}
}
