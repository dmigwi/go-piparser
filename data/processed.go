// Copyright 2019 Migwi Ndung'u.
// License that can be found in the LICENSE file.

package data

import (
	"time"

	"github.com/dmigwi/go-piparser/proposals/types"
)

// go-piparser parses the politeia's github stored data to primarily fetch the
// votes data. The votes data returned contain a time.Time field obtained from
// the individual commit timestamp. The go-piparser data output can be used to
// do lots of things which include but not limited to ploting charts.
// For charts: See the testutil.

var t, _ = time.Parse(types.CmdDateFormat, "Mon Nov 5 17:58:13 2018 +0000")

// SingleTokenVotesData defines the unmarshalled votes data result returned
// after parsing the the raw commmits message string in data/raw.go to
// unmarshal data for a single proposal token.
var SingleTokenVotesData = []*types.History{
	&types.History{
		Author:    "",
		CommitSHA: "",
		Date:      time.Time{},
		Token:     "",
		VotesInfo: nil,
	},
	&types.History{
		Author:    "Politeia <noreply@decred.org>",
		CommitSHA: "62f715e00c50e7c506acc4b6e33eb86d02bab6d1",
		Date:      t,
		Token:     "27f87171d98b7923a1bd2bee6affed929fa2d2a6e178b5c80a9971a92a5c7f50",
		VotesInfo: []types.CastVoteData{
			{
				PiVote: &types.PiVote{
					Ticket:  "1e4e075ef0346cbb07a42f9a15a1960939e8ee052a6c95fd276fa507fb9f89f7",
					VoteBit: "Yes",
				},
			},
			{
				PiVote: &types.PiVote{
					Ticket:  "c311797d4e2faf9d5e800ba0192061249ff578a041d972d81010b80f4e139fa5",
					VoteBit: "No",
				},
			},
			{
				PiVote: &types.PiVote{
					Ticket:  "90a4b53b5280cf621e06b94d106dd02c934846776f83ecbdd6c8374eb073deae",
					VoteBit: "Yes",
				},
			},
			{
				PiVote: &types.PiVote{
					Ticket:  "e272d314b1f6a15c4480145ab286a54bb9b6735718b776755fea7c77eba030b8",
					VoteBit: "Yes",
				},
			},
			{
				PiVote: &types.PiVote{
					Ticket:  "551988264328bd8ed75f87276ec4a94b3961bf0fe3698b9d976b3cd28b18d31d",
					VoteBit: "Yes",
				},
			},
			{
				PiVote: &types.PiVote{
					Ticket:  "caefdc114219ca2725618c06b87af5bf1ce67d18bc9a06718738b0acf08da57b",
					VoteBit: "No",
				},
			},
			{
				PiVote: &types.PiVote{
					Ticket:  "caf9aed8253bf7c03424d35b39550b7a4394149cfa4425155722ca995ef1a2fc",
					VoteBit: "Yes",
				},
			},
			{
				PiVote: &types.PiVote{
					Ticket:  "2c6c484f87c19df267e4316122dae5450120e892f04de81f8c0672ee41e2d94f",
					VoteBit: "No",
				},
			},
			{
				PiVote: &types.PiVote{
					Ticket:  "374d89180bbe0b11d22f1001c3933c766d1f1c2896e7b85dd4515bffc390ccdd",
					VoteBit: "No",
				},
			},
			{
				PiVote: &types.PiVote{
					Ticket:  "bd24343290a608dba0cdf103bc1390ce1fb863669a0eef5ae73e1765e841401d",
					VoteBit: "No",
				},
			},
			{
				PiVote: &types.PiVote{
					Ticket:  "af5501345b32e149b3da9710acb1887210d5efcd4fb4be3b711f62c69e4db95a",
					VoteBit: "No",
				},
			},
			{
				PiVote: &types.PiVote{
					Ticket:  "32d2b6259e7a33c4c0d472db77d7c69eb8f5e7deaa550b47cd8ac8f134f50755",
					VoteBit: "Yes",
				},
			},
		},
	},
	&types.History{
		Author:    "",
		CommitSHA: "",
		Date:      time.Time{},
		Token:     "",
		VotesInfo: nil,
	},
}

var t2, _ = time.Parse(types.CmdDateFormat, "Tue Mar 5 01:58:01 2019 +0000")

// AllTokensVotesData defines the unmarshalled votes data result returned
// after parsing the the raw commmits message string in data/raw.go to
// unmarshal data for a all proposal tokens.
var AllTokensVotesData = []*types.History{
	&types.History{
		Author:    "",
		CommitSHA: "",
		Date:      time.Time{},
		Token:     "",
		VotesInfo: nil,
	},
	&types.History{
		Author:    "Politeia <noreply@decred.org>",
		CommitSHA: "62f715e00c50e7c506acc4b6e33eb86d02bab6d1",
		Date:      t,
		Token:     "27f87171d98b7923a1bd2bee6affed929fa2d2a6e178b5c80a9971a92a5c7f50",
		VotesInfo: []types.CastVoteData{
			{
				PiVote: &types.PiVote{
					Ticket:  "1e4e075ef0346cbb07a42f9a15a1960939e8ee052a6c95fd276fa507fb9f89f7",
					VoteBit: "Yes",
				},
			},
			{
				PiVote: &types.PiVote{
					Ticket:  "c311797d4e2faf9d5e800ba0192061249ff578a041d972d81010b80f4e139fa5",
					VoteBit: "No",
				},
			},
			{
				PiVote: &types.PiVote{
					Ticket:  "90a4b53b5280cf621e06b94d106dd02c934846776f83ecbdd6c8374eb073deae",
					VoteBit: "Yes",
				},
			},
			{
				PiVote: &types.PiVote{
					Ticket:  "e272d314b1f6a15c4480145ab286a54bb9b6735718b776755fea7c77eba030b8",
					VoteBit: "Yes",
				},
			},
			{
				PiVote: &types.PiVote{
					Ticket:  "551988264328bd8ed75f87276ec4a94b3961bf0fe3698b9d976b3cd28b18d31d",
					VoteBit: "Yes",
				},
			},
			{
				PiVote: &types.PiVote{
					Ticket:  "caefdc114219ca2725618c06b87af5bf1ce67d18bc9a06718738b0acf08da57b",
					VoteBit: "No",
				},
			},
			{
				PiVote: &types.PiVote{
					Ticket:  "caf9aed8253bf7c03424d35b39550b7a4394149cfa4425155722ca995ef1a2fc",
					VoteBit: "Yes",
				},
			},
			{
				PiVote: &types.PiVote{
					Ticket:  "2c6c484f87c19df267e4316122dae5450120e892f04de81f8c0672ee41e2d94f",
					VoteBit: "No",
				},
			},
			{
				PiVote: &types.PiVote{
					Ticket:  "374d89180bbe0b11d22f1001c3933c766d1f1c2896e7b85dd4515bffc390ccdd",
					VoteBit: "No",
				},
			},
			{
				PiVote: &types.PiVote{
					Ticket:  "bd24343290a608dba0cdf103bc1390ce1fb863669a0eef5ae73e1765e841401d",
					VoteBit: "No",
				},
			},
			{
				PiVote: &types.PiVote{
					Ticket:  "af5501345b32e149b3da9710acb1887210d5efcd4fb4be3b711f62c69e4db95a",
					VoteBit: "No",
				},
			},
			{
				PiVote: &types.PiVote{
					Ticket:  "32d2b6259e7a33c4c0d472db77d7c69eb8f5e7deaa550b47cd8ac8f134f50755",
					VoteBit: "Yes",
				},
			},
		},
	},
	&types.History{
		Author:    "Politeia <noreply@decred.org>",
		CommitSHA: "4913ebaef7eac7f70913f285d49de03f5ed08e87",
		Date:      t2,
		Token:     "a3def199af812b796887f4eae22e11e45f112b50c2e17252c60ed190933ec14f",
		VotesInfo: []types.CastVoteData{
			{
				PiVote: &types.PiVote{
					Ticket:  "03d4f5888a0a7bf983852b379de539acf8eff272534cf2be6846ac55eaae878b",
					VoteBit: "No",
				},
			},
			{
				PiVote: &types.PiVote{
					Ticket:  "dca8cf91d55a5f1b00979723cdb7ceb66bc83234f1851328232b77c3d0062ec2",
					VoteBit: "No",
				},
			},
			{
				PiVote: &types.PiVote{
					Ticket:  "1e3836b86ff7a5809fe834c0f03f8f04c54ff06afff3ba8a3620c17434b94d86",
					VoteBit: "No",
				},
			},
			{
				PiVote: &types.PiVote{
					Ticket:  "b254bb2f69d1335009c9c64f7f80b36a1f30714cab99a97e6011dfa03fd623a3",
					VoteBit: "No",
				},
			},
			{
				PiVote: &types.PiVote{
					Ticket:  "879c3994bd3b69bc334ba584a7cbf2a0449a9841435f9dca1b4bd0a1496b7007",
					VoteBit: "No",
				},
			},
			{
				PiVote: &types.PiVote{
					Ticket:  "339f78bd215672003b23909d45a4489b97d52c454425501449a4ac51f59ca029",
					VoteBit: "No",
				},
			},
			{
				PiVote: &types.PiVote{
					Ticket:  "2052550015c6efbc67a71294f02f089900d3bb9dfb07b623da3a5797d75b1816",
					VoteBit: "No",
				},
			},
			{
				PiVote: &types.PiVote{
					Ticket:  "3e7d140a43defca57354436e8a7829d22854d8b5d6a7dac7cc4acb419eeb979d",
					VoteBit: "No",
				},
			},
			{
				PiVote: &types.PiVote{
					Ticket:  "7836719f9829af92cbf39f40096d90df30f529ea70ac47faeed0ad018770ac13",
					VoteBit: "No",
				},
			},
			{
				PiVote: &types.PiVote{
					Ticket:  "0713aae531246030672ca417c4f89b4006a258282e78b7441337f5c0c5dbfb0c",
					VoteBit: "No",
				},
			},
		},
	},
}
