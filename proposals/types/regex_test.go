package types

import (
	"strconv"
	"testing"
	"time"
)

// testData defines the general struct with possible inputs and outputs that can
// be used to test the various regex expressions.
type testData struct {
	src, repl, output string
	isError           bool
	timestamp         time.Time
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
	date, _ := time.Parse(CmdDateFormat, "Thu Feb 28 15:35:56 2019 -0600")
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

	// set current proposal token to 27f87171d98b7923a1bd2bee6affed929fa2d2a6e178b5c80a9971a92a5c7f50
	SetProposalToken("27f87171d98b7923a1bd2bee6affed929fa2d2a6e178b5c80a9971a92a5c7f50")

	td := []testData{
		{src: `
		`, regex: "\n", isFound: true}, // matching newline character.
		{src: `
		`, regex: "\\n", isFound: true}, //matching escaped newline character.
		{src: `\n`, regex: "\\n", isFound: false},
		{src: `    Flush vote journals.
		`, regex: DefaultVotesCommitMsg, isFound: true},
		{src: "b/27f87171d98b7923a1bd2bee6af/3/plugins/decred/ballot.journal",
			regex: "27f87171d98b7923a1bd2bee6af", isFound: true},
		{src: `{"version":"1","action":"add"}{"castvote":{"token":"27f87171d98b7923a1bd2bee6affed929fa2d2a6e178b5c80a9971a92a5c7f50",.
		`, regex: VotesJSONSignature(), isFound: true}, // Match the VotesJSONSignature() regex
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

// TestVotesJSONSignature tests the strings that can be matched with
// VotesJSONSignature output.
func TestVotesJSONSignature(t *testing.T) {
	testString1 := `{"version":"1","action":"add"}{"castvote":{"token":"a3def199af812b796887f4eae22e11e45f112b50c2e17252c60ed190933ec14f","ticket":"03d4f5888a0a7bf983852b379de539acf8eff272534cf2be6846ac55eaae878b","votebit":"1","signature":"1f06c29926a871a501f91fd0bca0b68b2d12226c582f0277b4be59eb48454b8e894824c4a02ec312b87245d285a99f835492dd766bfd34d9d32222a6f03c60a413"},"receipt":"7e0f760157cf8d3cb7bfe76e4c76aaf41a6571dc4a9519d603be30986fb36028203cf21c9e81e2819adaa3660b4195a0868daf068c5a39f7949f822b53977f05"}`

	testString2 := `{"version":"1","action":"add"}{"castvote":{"token":"27f87171d98b7923a1bd2bee6affed929fa2d2a6e178b5c80a9971a92a5c7f50","ticket":"03d4f5888a0a7bf983852b379de539acf8eff272534cf2be6846ac55eaae878b","votebit":"1","signature":"1f06c29926a871a501f91fd0bca0b68b2d12226c582f0277b4be59eb48454b8e894824c4a02ec312b87245d285a99f835492dd766bfd34d9d32222a6f03c60a413"},"receipt":"7e0f760157cf8d3cb7bfe76e4c76aaf41a6571dc4a9519d603be30986fb36028203cf21c9e81e2819adaa3660b4195a0868daf068c5a39f7949f822b53977f05"}`

	// invalid proposal vote data.
	testString3 := `{"version":"1","action":"add"}{"castvote":{"token":"27f87171d98b7923a1bd2bee6affed929fa2d26e178b5c80a9971a92a5c7f50","ticket":"03d4f5888a0a7bf983852b379de539acf8eff272534cf2be6846ac55eaae878b","votebit":"1","signature":"1f06c29926a871a501f91fd0bca0b68b2d12226c582f0277b4be59eb48454b8e894824c4a02ec312b87245d285a99f835492dd766bfd34d9d32222a6f03c60a413"},"receipt":"7e0f760157cf8d3cb7bfe76e4c76aaf41a6571dc4a9519d603be30986fb36028203cf21c9e81e2819adaa3660b4195a0868daf068c5a39f7949f822b53977f05"}`

	// invalid proposal vote data.
	testString4 := `{"version":"1","action":"add"}{"castvote":{"token":"27f87171d98b7923a1bd2bee6affed929 fa2d26e178b5c80a9971a92a5c7f50","ticket":"03d4f5888a0a7bf983852b379de539acf8eff272534cf2be6846ac55eaae878b","votebit":"1","signature":"1f06c29926a871a501f91fd0bca0b68b2d12226c582f0277b4be59eb48454b8e894824c4a02ec312b87245d285a99f835492dd766bfd34d9d32222a6f03c60a413"},"receipt":"7e0f760157cf8d3cb7bfe76e4c76aaf41a6571dc4a9519d603be30986fb36028203cf21c9e81e2819adaa3660b4195a0868daf068c5a39f7949f822b53977f05"}`

	defaultToken := "27f87171d98b7923a1bd2bee6affed929fa2d2a6e178b5c80a9971a92a5c7f50"

	type testData struct {
		src, token string
		isMatched  bool
	}

	td := []testData{
		{src: testString1, token: "", isMatched: true}, // Any token matched
		{src: testString1, token: defaultToken, isMatched: false},
		{src: testString2, token: "", isMatched: true},           // Any token matched
		{src: testString2, token: defaultToken, isMatched: true}, // only defaultToken matched
		{src: testString3, token: "", isMatched: false},
		{src: testString3, token: defaultToken, isMatched: false},
		{src: testString4, token: "", isMatched: false},
		{src: testString4, token: defaultToken, isMatched: false},
	}

	for i, val := range td {
		t.Run("Test_#"+strconv.Itoa(i), func(t *testing.T) {
			if len(val.token) > 0 {
				// set proposal token if it exists.
				SetProposalToken(val.token)
			} else {
				// drop previous set proposal value.
				ClearProposalToken()
			}

			result := IsMatching(val.src, VotesJSONSignature())
			if result != val.isMatched {
				t.Fatalf("expected the matching src to the regex to be %v but found %v",
					val.isMatched, result)
			}
		})
	}
}

func TestRetrieveProposalToken(t *testing.T) {
	testString1 := `{"version":"1","action":"add"}{"castvote":{"token":"a3def199af812b796887f4eae22e11e45f112b50c2e17252c60ed190933ec14f","ticket":"03d4f5888a0a7bf983852b379de539acf8eff272534cf2be6846ac55eaae878b","votebit":"1","signature":"1f06c29926a871a501f91fd0bca0b68b2d12226c582f0277b4be59eb48454b8e894824c4a02ec312b87245d285a99f835492dd766bfd34d9d32222a6f03c60a413"},"receipt":"7e0f760157cf8d3cb7bfe76e4c76aaf41a6571dc4a9519d603be30986fb36028203cf21c9e81e2819adaa3660b4195a0868daf068c5a39f7949f822b53977f05"}`

	testString2 := `{"version":"1","action":"add"}{"castvote":{"token":"27f87171d98b7923a1bd2bee6affed929fa2d2a6e178b5c80a9971a92a5c7f50","ticket":"03d4f5888a0a7bf983852b379de539acf8eff272534cf2be6846ac55eaae878b","votebit":"1","signature":"1f06c29926a871a501f91fd0bca0b68b2d12226c582f0277b4be59eb48454b8e894824c4a02ec312b87245d285a99f835492dd766bfd34d9d32222a6f03c60a413"},"receipt":"7e0f760157cf8d3cb7bfe76e4c76aaf41a6571dc4a9519d603be30986fb36028203cf21c9e81e2819adaa3660b4195a0868daf068c5a39f7949f822b53977f05"}`

	// invalid proposal vote data.
	testString3 := `{"version":"1","action":"add"}{"castvote":{"token":"27f87171d98b7923a1bd2bee6affed929fa2d26e178b5c80a9971a92a5c7f50","ticket":"03d4f5888a0a7bf983852b379de539acf8eff272534cf2be6846ac55eaae878b","votebit":"1","signature":"1f06c29926a871a501f91fd0bca0b68b2d12226c582f0277b4be59eb48454b8e894824c4a02ec312b87245d285a99f835492dd766bfd34d9d32222a6f03c60a413"},"receipt":"7e0f760157cf8d3cb7bfe76e4c76aaf41a6571dc4a9519d603be30986fb36028203cf21c9e81e2819adaa3660b4195a0868daf068c5a39f7949f822b53977f05"}`

	// invalid proposal vote data.
	testString4 := `{"version":"1","action":"add"}{"castvote":{"token":"27f87171d98b7923a1bd2bee6affed929 fa2d26e178b5c80a9971a92a5c7f50","ticket":"03d4f5888a0a7bf983852b379de539acf8eff272534cf2be6846ac55eaae878b","votebit":"1","signature":"1f06c29926a871a501f91fd0bca0b68b2d12226c582f0277b4be59eb48454b8e894824c4a02ec312b87245d285a99f835492dd766bfd34d9d32222a6f03c60a413"},"receipt":"7e0f760157cf8d3cb7bfe76e4c76aaf41a6571dc4a9519d603be30986fb36028203cf21c9e81e2819adaa3660b4195a0868daf068c5a39f7949f822b53977f05"}`

	token1 := "27f87171d98b7923a1bd2bee6affed929fa2d2a6e178b5c80a9971a92a5c7f50"
	token2 := "a3def199af812b796887f4eae22e11e45f112b50c2e17252c60ed190933ec14f"

	type testData struct {
		src, token string
		isError    bool
	}

	td := []testData{
		{src: testString1, token: token2, isError: false},
		{src: testString2, token: token1, isError: false},
		{src: testString3, token: "", isError: true},
		{src: testString4, token: "", isError: true},
	}

	for i, val := range td {
		t.Run("Test_#"+strconv.Itoa(i), func(t *testing.T) {
			resp, err := RetrieveProposalToken(val.src)
			if err != nil && !val.isError {
				t.Fatalf("expected no error but found: %v", err)
			}

			if err == nil && val.isError {
				t.Fatal("expected an error but found none")
			}

			if val.token != resp {
				t.Fatalf("expected the returned value (%s) to be equal to (%s) but it wasn't", resp, val.token)
			}
		})
	}
}
