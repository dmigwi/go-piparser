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

// SingleTokenVotesData defines how the unmarshalled history data returned
// after parsing the the raw commmits message string in data/raw.go to
// unmarshal data for a single proposal token should look like. History objects
// with empty fields are emitted from the final data returned.
var SingleTokenVotesData = []*types.History{
	&types.History{
		Author:    "",
		CommitSHA: "",
		Date:      time.Time{},
		Patch:     nil,
	},
	&types.History{
		Author:    "Politeia <noreply@decred.org>",
		CommitSHA: "62f715e00c50e7c506acc4b6e33eb86d02bab6d1",
		Date:      t,
		Patch: []*types.File{
			{
				Token: "27f87171d98b7923a1bd2bee6affed929fa2d2a6e178b5c80a9971a92a5c7f50",
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
		},
	},
	&types.History{
		Author:    "",
		CommitSHA: "",
		Date:      time.Time{},
		Patch:     nil,
	},
	&types.History{
		Author:    "",
		CommitSHA: "",
		Date:      time.Time{},
		Patch:     nil,
	},
}

var t2, _ = time.Parse(types.CmdDateFormat, "Tue Mar 5 01:58:01 2019 +0000")
var t3, _ = time.Parse(types.CmdDateFormat, "Wed Mar 6 12:58:01 2019 +0000")

// AllTokensVotesData defines  how the unmarshalled history data returned
// after parsing the the raw commmits message string in data/raw.go to
// unmarshal data for a all proposal tokens should look like. History objects
// with empty fields are emitted from the final data returned.
var AllTokensVotesData = []*types.History{
	&types.History{
		Author:    "",
		CommitSHA: "",
		Date:      time.Time{},
		Patch:     nil,
	},
	&types.History{
		Author:    "Politeia <noreply@decred.org>",
		CommitSHA: "62f715e00c50e7c506acc4b6e33eb86d02bab6d1",
		Date:      t,
		Patch: []*types.File{
			{
				Token: "27f87171d98b7923a1bd2bee6affed929fa2d2a6e178b5c80a9971a92a5c7f50",
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
		},
	},
	&types.History{
		Author:    "Politeia <noreply@decred.org>",
		CommitSHA: "4913ebaef7eac7f70913f285d49de03f5ed08e87",
		Date:      t2,
		Patch: []*types.File{
			{
				Token: "a3def199af812b796887f4eae22e11e45f112b50c2e17252c60ed190933ec14f",
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
		},
	},
	&types.History{
		Author:    "Politeia <noreply@decred.org>",
		CommitSHA: "1d6edd806dd8bf043cdbd343c9d7d8e5dcc90b4f",
		Date:      t3,
		Patch: []*types.File{
			{
				Token: "5431da8ff4eda8cdbf8f4f2e08566ffa573464b97ef6d6bae78e749f27800d3a",
				VotesInfo: []types.CastVoteData{
					{
						PiVote: &types.PiVote{
							Ticket:  "03cca8c7d0d8d6f8904e8535bed958063a45fd0b0e2a336492b1518d543366fc",
							VoteBit: "Yes",
						},
					},
					{
						PiVote: &types.PiVote{
							Ticket:  "4fe96f731451a49d944a4d42c259ba1f13ac64019fe5929a1bd28ff4f192d249",
							VoteBit: "Yes",
						},
					},
					{
						PiVote: &types.PiVote{
							Ticket:  "81db496d21a2719e685f53f1d2916a065773ffa50741ca65c1e4ca1914a974ac",
							VoteBit: "Yes",
						},
					},
					{
						PiVote: &types.PiVote{
							Ticket:  "96bee2741fee0052b2f166dc27ffa985cee0c3201954695ba566712afa441d7e",
							VoteBit: "Yes",
						},
					},
					{
						PiVote: &types.PiVote{
							Ticket:  "ca409f7aed1fb83e4b84705c96d93810654c2985e8abefd6a441c2eaf68b4a91",
							VoteBit: "Yes",
						},
					},
					{
						PiVote: &types.PiVote{
							Ticket:  "f21106045ac661e0e1f6f50330a64efd29c079339b3128ac64529f28e04ac794",
							VoteBit: "Yes",
						},
					},
					{
						PiVote: &types.PiVote{
							Ticket:  "f8db8898b25420370734963399511b7ff94621c1b1ed3911c01cc3b5ba1b06a5",
							VoteBit: "Yes",
						},
					},
					{
						PiVote: &types.PiVote{
							Ticket:  "fe8260ff855253ad29e1a31b77ca68b10ee6825b81b5503cc897ac910a1467ef",
							VoteBit: "Yes",
						},
					},
					{
						PiVote: &types.PiVote{
							Ticket:  "3f9bc2620457d17b8f5524ee0a879c468f482ee0fe47141f66ed9e8155d53979",
							VoteBit: "Yes",
						},
					},
					{
						PiVote: &types.PiVote{
							Ticket:  "63cba705b84825fc88357fcbe5dec7025bd7054d2d0815109452fd72f5542f97",
							VoteBit: "Yes",
						},
					},
					{
						PiVote: &types.PiVote{
							Ticket:  "783d463dbf710a4a1a98a30c5fe6eeaf30c6a0cffbd51a00f18cf7919e73beeb",
							VoteBit: "Yes",
						},
					},
					{
						PiVote: &types.PiVote{
							Ticket:  "7a1ef899f0e1cc3b021d2d2bcf85dc7157cc22b1b6d856092bbede28d6af437b",
							VoteBit: "Yes",
						},
					},
					{
						PiVote: &types.PiVote{
							Ticket:  "7c6d6f99097e5f0bed5a2871da45c64e1bf364e8451538e4c2ae73d3661b7329",
							VoteBit: "Yes",
						},
					},
					{
						PiVote: &types.PiVote{
							Ticket:  "f76fee67d96c8b23ba41760fb23c0a0f0b1d135fe545a839b0fd73f42420be9a",
							VoteBit: "Yes",
						},
					},
				},
			},
			{
				Token: "60adb9c0946482492889e85e9bce05c309665b3438dd85cb1a837df31fbf57fb",
				VotesInfo: []types.CastVoteData{
					{
						PiVote: &types.PiVote{
							Ticket:  "28d890c4801a8691c9ba9594a9aa2abb167321ef3bcf1f331b6ec863553f8b51",
							VoteBit: "No",
						},
					},
					{
						PiVote: &types.PiVote{
							Ticket:  "3b6f7a70321d463a7d1921eedea8c189d6c80e4a489844a2a386c56e9a63cb08",
							VoteBit: "No",
						},
					},
					{
						PiVote: &types.PiVote{
							Ticket:  "81d40bb7f3fb98869fc6086e1adf00f8afe3362b6838f1b3446226b267410d31",
							VoteBit: "No",
						},
					},
					{
						PiVote: &types.PiVote{
							Ticket:  "9cb65e4580c73a0e50276c53d807a9c8929de7b8283aae4afa6c5d72ba14411a",
							VoteBit: "No",
						},
					},
				},
			},
			{
				Token: "a3def199af812b796887f4eae22e11e45f112b50c2e17252c60ed190933ec14f",
				VotesInfo: []types.CastVoteData{
					{
						PiVote: &types.PiVote{
							Ticket:  "347642ffc492a484aa223b06b6420ebbe71f132af19021e8c3f42701dd0c63fa",
							VoteBit: "No",
						},
					},
					{
						PiVote: &types.PiVote{
							Ticket:  "f299d3e5700300491eb91d45065b8e8635bcc38118a1887f88b9c70b1dbf9aff",
							VoteBit: "No",
						},
					},
					{
						PiVote: &types.PiVote{
							Ticket:  "25f1d9344d3c49e6185f0050b5fc862852a2b0a04f03f10acec6a8e044c310c3",
							VoteBit: "Yes",
						},
					},
					{
						PiVote: &types.PiVote{
							Ticket:  "3db37e214787bff8f298c084772a65cb26b279f0a2d964d83edb7521d704e47c",
							VoteBit: "No",
						},
					},
					{
						PiVote: &types.PiVote{
							Ticket:  "93e43cfad33cd1a635704bf8e0e11c5b825b276c10b2e883e143e56dd8e22ab6",
							VoteBit: "No",
						},
					},
					{
						PiVote: &types.PiVote{
							Ticket:  "51bc6206766ae0f7913b228575a81624733ad4057ab6b966f249c95e4ba1cf94",
							VoteBit: "No",
						},
					},
					{
						PiVote: &types.PiVote{
							Ticket:  "980a0588cc6cb908048b72a437b56f4a76411a7593a25c4abe1879ca063158b8",
							VoteBit: "No",
						},
					},
					{
						PiVote: &types.PiVote{
							Ticket:  "7f9efcb4d9ee8214918186a6054ffa362bb55de039a0631275d086dbfedb70ce",
							VoteBit: "No",
						},
					},
					{
						PiVote: &types.PiVote{
							Ticket:  "1b76219ca6e124185f37756ca2b4aacd32129db6cc7f2467776c481a4a153118",
							VoteBit: "No",
						},
					},
					{
						PiVote: &types.PiVote{
							Ticket:  "f6c97ffaf9964d2dfc821be3cace223bc99b739a8cf7d9a206c6307a49464edd",
							VoteBit: "No",
						},
					},
				},
			},
		},
	},
}
