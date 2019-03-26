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
			output: `+{"version":"1","action":"add"}{"castvote":{"token":"27f87171d98b7923a1bd2bee6affed929fa2d2a6e178b5c80a9971a92a5c7f50"}
`,
		},
		{
			src: `{"version":"1","action":"add"}{"castvote":{"token":"27f8","ticket":"dab3d276","votebit":"2","signature":"20942f2a"},"receipt":"547416f8f"}
			{"version":"1","action":"add"}{"castvote":{"token":"27f8","ticket":"8863036","votebit":"2","signature":"2066ca72"},"receipt":"5b15f2b3c2"}
			{"version":"1","action":"add"}{"castvote":{"token":"27f8","ticket":"f837268","votebit":"2","signature":"20d3d731"},"receipt":"ea6fb3c02f"}
			+{"version":"1","action":"addlike"}{"castvote":{"token":"27f8","ticket":"14af967","votebit":"2","signature":"20e1e03d"},"receipt":"9a65775122"}
			+{"version":"1","action":"del"}{"castvote":{"token":"27f8","ticket":"b4f7f4e","votebit":"2","signature":"1fd2e407"},"receipt":"02da25c951"}
			+{"version":"1","action":"del"}{"castvote":{"token":"27f8","ticket":"c6d2633","votebit":"2","signature":"1ffec594"},"receipt":"0045001444"}
		   
		   `,
			output: `+{"version":"1","action":"addlike"}{"castvote":{"token":"27f8","ticket":"14af967","votebit":"2","signature":"20e1e03d"},"receipt":"9a65775122"}
			+{"version":"1","action":"del"}{"castvote":{"token":"27f8","ticket":"b4f7f4e","votebit":"2","signature":"1fd2e407"},"receipt":"02da25c951"}
			+{"version":"1","action":"del"}{"castvote":{"token":"27f8","ticket":"c6d2633","votebit":"2","signature":"1ffec594"},"receipt":"0045001444"}
`,
		},
		{
			src: ` /5431da8ff4eda8cdbf8f4f2e08566ffa573464b97ef6d6bae78e749f27800d3a/3/plugins/decred/ballot.journal b/5431da8ff4eda8cdbf8f4f2e08566ffa573464b97ef6d6bae78e749f27800d3a/3/plugins/decred/ballot.journal
			index 63e27702..8cd5006f 100644
			--- a/5431da8ff4eda8cdbf8f4f2e08566ffa573464b97ef6d6bae78e749f27800d3a/3/plugins/decred/ballot.journal
			+++ b/5431da8ff4eda8cdbf8f4f2e08566ffa573464b97ef6d6bae78e749f27800d3a/3/plugins/decred/ballot.journal
			@@ -2612,3 +2612,29 @@
			 {"version":"1","action":"add"}{"castvote":{"token":"5431da8ff4eda8cdbf8f4f2e08566ffa573464b97ef6d6bae78e749f27800d3a","ticket":"d5e69b8a51faa13f600c001c9b7a409b38d002140a39b09557420982dbd98192","votebit":"2","signature":"200fc10d9c9eae81f33832c935b1b4b7c7ad752d5764
			baecdf15752ca9bae68aeb1475df57ffd07628ae41c5d2e2ad754930e7e6ccbb7376e6b496df094075e1e9"},"receipt":"904457b9e903a2b603fe6bfd12b19e26a4fb652220c6dcef62043031f85ed70d5c5032872ecfbd05dee2094a8b1935c738bea3310f11ea7f60cc7572b1e7f404"}
			 {"version":"1","action":"add"}{"castvote":{"token":"5431da8ff4eda8cdbf8f4f2e08566ffa573464b97ef6d6bae78e749f27800d3a","ticket":"35e4b55f35a535061031ca5f9277f6d561ba08f7f716f302f856a52816eb7bdb","votebit":"2","signature":"1f0ca8ef3b83f8329eb0d64e9626847421b63139217d
			e25f210cf5f3b3b35f42980c096cce332424fb429d0fe13baa638280eaa1235778029d2f6149d9b329370b"},"receipt":"b7504e3fb0dc84b3f40b0efd693a85258333070fd8739450443e72fa8baabfdcefa823b3644eaa929ad2b9b0fcb1e26761b82975807a6cdc7497d796fb43f004"}
			 {"version":"1","action":"add"}{"castvote":{"token":"5431da8ff4eda8cdbf8f4f2e08566ffa573464b97ef6d6bae78e749f27800d3a","ticket":"fa03104fa79583cf76289eded12d99075612d6cc6d05be279c32cb1290f96753","votebit":"2","signature":"1f87d1b48635e6ab019d0d8c548fa76ff7dae48b4a24
			92445575f88dc8d4b9eb0e7cdb4d517d92868bf772d13d891dc5afbf47469923477038e08a4a3b2dbf39ea"},"receipt":"ba72f63b8aebbe2066f445c69b3090ebab695468635f985fde4268248dbccacaa522350651a8f6ce32c14267e3a9924262398217f9aa01573e8f866e705ba903"}
			+{"version":"1","action":"add"}{"castvote":{"token":"5431da8ff4eda8cdbf8f4f2e08566ffa573464b97ef6d6bae78e749f27800d3a","ticket":"03cca8c7d0d8d6f8904e8535bed958063a45fd0b0e2a336492b1518d543366fc","votebit":"2","signature":"209ffbb7f5593cf67bbce7c29c8282de160d4dd2ef84
			bd62f9bfa69e671614db5414b1b80806dda66a5134c964f80b53983be568ef1c2e940f8a6e8202b5065f42"},"receipt":"6786d64ed1c2b06d053f51b14a175a6457a8a87a2bbc920c56ffbdf8
			12a1905f3bcb091910c1730187dcbf5acb2a25b01a2641ef29fba4236e2bb9e64e56b40f"}
			+{"version":"1","action":"add"}{"castvote":{"token":"5431da8ff4eda8cdbf8f4f2e08566ffa573464b97ef6d6bae78e749f27800d3a","ticket":"4fe96f731451a49d944a4d42c259ba1f13ac64019fe5929a1bd28ff4f192d249","votebit":"2","signature":"1f1c4ae1e165c30f0dea6daf2f7e89ac58ba08ae6bf1
			4835bb138baa0bba1d19181aea0b03ee3b7621163a9a4c772d798719e51bfed9f36702abb3ce0ad32b01d6"},"receipt":"2dde8844cef5c4149edc5c81a3cae4d8fca8a9d2f00fc8ff0f9ac7fbcd1afe08ce7a3f0de5103d99aee0fca406c59c5741e116bec4e983de335d110267a57c00"}
			+{"version":"1","action":"add"}{"castvote":{"token":"5431da8ff4eda8cdbf8f4f2e08566ffa573464b97ef6d6bae78e749f27800d3a","ticket":"81db496d21a2719e685f53f1d2916a065773ffa50741ca65c1e4ca1914a974ac","votebit":"2","signature":"1f510e55a571c8aa4f98a21708fcffc2ab35ac577f8b
			e85ef3141e2272709e01da4ff2be5e8b52c20a0712ddfae621c87e330f806ebc250c2ea3671306c879f1a8"},"receipt":"899a314801c2fac8bb8b1c4ad7d75bea6c63e52ca2ff5e3bcc532a21e75d11e378a618ff61475f3707620fd2defc45556ddbe1b194e182f5bffc23fca6bde206"}
			+{"version":"1","action":"add"}{"castvote":{"token":"5431da8ff4eda8cdbf8f4f2e08566ffa573464b97ef6d6bae78e749f27800d3a","ticket":"96bee2741fee0052b2f166dc27ffa985cee0c3201954695ba566712afa441d7e","votebit":"2","signature":"1fd6b346faa901bea78ed93fd822e74cdb6e4443db05
			094198e248e0945cc6fcfc20b9e747c69c775844fa9ddfeb7003e1369eb8e5bced4be4143c833545da4260"},"receipt":"da898c111cc7525f230f62a5cb333dd6ffd8860d4936539bb94388c9895a23c8f2475e18b16e32c23711a7e45e4281c503c262bf71a5e04d7e6617a3df682805"}
			`,
			output: `+{"version":"1","action":"add"}{"castvote":{"token":"5431da8ff4eda8cdbf8f4f2e08566ffa573464b97ef6d6bae78e749f27800d3a","ticket":"03cca8c7d0d8d6f8904e8535bed958063a45fd0b0e2a336492b1518d543366fc","votebit":"2","signature":"209ffbb7f5593cf67bbce7c29c8282de160d4dd2ef84
			bd62f9bfa69e671614db5414b1b80806dda66a5134c964f80b53983be568ef1c2e940f8a6e8202b5065f42"},"receipt":"6786d64ed1c2b06d053f51b14a175a6457a8a87a2bbc920c56ffbdf8
			12a1905f3bcb091910c1730187dcbf5acb2a25b01a2641ef29fba4236e2bb9e64e56b40f"}
			+{"version":"1","action":"add"}{"castvote":{"token":"5431da8ff4eda8cdbf8f4f2e08566ffa573464b97ef6d6bae78e749f27800d3a","ticket":"4fe96f731451a49d944a4d42c259ba1f13ac64019fe5929a1bd28ff4f192d249","votebit":"2","signature":"1f1c4ae1e165c30f0dea6daf2f7e89ac58ba08ae6bf1
			4835bb138baa0bba1d19181aea0b03ee3b7621163a9a4c772d798719e51bfed9f36702abb3ce0ad32b01d6"},"receipt":"2dde8844cef5c4149edc5c81a3cae4d8fca8a9d2f00fc8ff0f9ac7fbcd1afe08ce7a3f0de5103d99aee0fca406c59c5741e116bec4e983de335d110267a57c00"}
			+{"version":"1","action":"add"}{"castvote":{"token":"5431da8ff4eda8cdbf8f4f2e08566ffa573464b97ef6d6bae78e749f27800d3a","ticket":"81db496d21a2719e685f53f1d2916a065773ffa50741ca65c1e4ca1914a974ac","votebit":"2","signature":"1f510e55a571c8aa4f98a21708fcffc2ab35ac577f8b
			e85ef3141e2272709e01da4ff2be5e8b52c20a0712ddfae621c87e330f806ebc250c2ea3671306c879f1a8"},"receipt":"899a314801c2fac8bb8b1c4ad7d75bea6c63e52ca2ff5e3bcc532a21e75d11e378a618ff61475f3707620fd2defc45556ddbe1b194e182f5bffc23fca6bde206"}
			+{"version":"1","action":"add"}{"castvote":{"token":"5431da8ff4eda8cdbf8f4f2e08566ffa573464b97ef6d6bae78e749f27800d3a","ticket":"96bee2741fee0052b2f166dc27ffa985cee0c3201954695ba566712afa441d7e","votebit":"2","signature":"1fd6b346faa901bea78ed93fd822e74cdb6e4443db05
			094198e248e0945cc6fcfc20b9e747c69c775844fa9ddfeb7003e1369eb8e5bced4be4143c833545da4260"},"receipt":"da898c111cc7525f230f62a5cb333dd6ffd8860d4936539bb94388c9895a23c8f2475e18b16e32c23711a7e45e4281c503c262bf71a5e04d7e6617a3df682805"}
`,
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
