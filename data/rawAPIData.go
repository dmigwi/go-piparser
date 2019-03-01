package gitapi

var listCommitsData = `[
	{
	  "sha": "8ed7ae063166d6ea0b7c2e8fb89daae95b384163",
	  "node_id": "MDY6Q29tbWl0MTcyODA0MzIyOjhlZDdhZTA2MzE2NmQ2ZWEwYjdjMmU4ZmI4OWRhYWU5NWIzODQxNjM=",
	  "commit": {
		"author": {
		  "name": "Politeia",
		  "email": "noreply@decred.org",
		  "date": "2019-02-26T22:58:00Z"
		},
		"committer": {
		  "name": "Politeia",
		  "email": "noreply@decred.org",
		  "date": "2019-02-26T22:58:00Z"
		},
		"message": "Flush comment journals.\n\n60adb9c0946482492889e85e9bce05c309665b3438dd85cb1a837df31fbf57fb\na3def199af812b796887f4eae22e11e45f112b50c2e17252c60ed190933ec14f\naea224a561cfed183f514a9ac700d68ba8a6c71dfbee71208fb9bff5fffab51d\nc84a76685e4437a15760033725044a15ad832f68f9d123eb837337060a09f86e\nfb8e6ca361c807168ea0bd6ddbfb7e05896b78f2576daf92f07315e6f8b5cd83",
		"tree": {
		  "sha": "edcdab581e190975dad668554484e35811046257",
		  "url": "https://api.github.com/repos/dmigwi/mainnet/git/trees/edcdab581e190975dad668554484e35811046257"
		},
		"url": "https://api.github.com/repos/dmigwi/mainnet/git/commits/8ed7ae063166d6ea0b7c2e8fb89daae95b384163",
		"comment_count": 0,
		"verification": {
		  "verified": false,
		  "reason": "unsigned",
		  "signature": null,
		  "payload": null
		}
	  },
	  "url": "https://api.github.com/repos/dmigwi/mainnet/commits/8ed7ae063166d6ea0b7c2e8fb89daae95b384163",
	  "html_url": "https://github.com/dmigwi/mainnet/commit/8ed7ae063166d6ea0b7c2e8fb89daae95b384163",
	  "comments_url": "https://api.github.com/repos/dmigwi/mainnet/commits/8ed7ae063166d6ea0b7c2e8fb89daae95b384163/comments",
	  "author": null,
	  "committer": null,
	  "parents": [
		{
		  "sha": "f6a97ec7c23468676b2b69abc7d275a85126ed19",
		  "url": "https://api.github.com/repos/dmigwi/mainnet/commits/f6a97ec7c23468676b2b69abc7d275a85126ed19",
		  "html_url": "https://github.com/dmigwi/mainnet/commit/f6a97ec7c23468676b2b69abc7d275a85126ed19"
		}
	  ]
	},
	{
	  "sha": "f6a97ec7c23468676b2b69abc7d275a85126ed19",
	  "node_id": "MDY6Q29tbWl0MTcyODA0MzIyOmY2YTk3ZWM3YzIzNDY4Njc2YjJiNjlhYmM3ZDI3NWE4NTEyNmVkMTk=",
	  "commit": {
		"author": {
		  "name": "Politeia",
		  "email": "noreply@decred.org",
		  "date": "2019-02-26T22:36:25Z"
		},
		"committer": {
		  "name": "Politeia",
		  "email": "noreply@decred.org",
		  "date": "2019-02-26T22:36:25Z"
		},
		"message": "Anchor confirmation bdf1271f1343753c0f3590b11fecf1603f3b864b6d1c9aaf84eec8d17a86d9d2\n\nbdf1271f1343753c0f3590b11fecf1603f3b864b6d1c9aaf84eec8d17a86d9d2 anchored in TX 8b5cc74cc607538dcf0bd5b536c33fc8fda53d2e88e3c829edbe9a552bba86fb",
		"tree": {
		  "sha": "918f8d34addc12e665329dbcabb3ae8963aaaa6f",
		  "url": "https://api.github.com/repos/dmigwi/mainnet/git/trees/918f8d34addc12e665329dbcabb3ae8963aaaa6f"
		},
		"url": "https://api.github.com/repos/dmigwi/mainnet/git/commits/f6a97ec7c23468676b2b69abc7d275a85126ed19",
		"comment_count": 0,
		"verification": {
		  "verified": false,
		  "reason": "unsigned",
		  "signature": null,
		  "payload": null
		}
	  },
	  "url": "https://api.github.com/repos/dmigwi/mainnet/commits/f6a97ec7c23468676b2b69abc7d275a85126ed19",
	  "html_url": "https://github.com/dmigwi/mainnet/commit/f6a97ec7c23468676b2b69abc7d275a85126ed19",
	  "comments_url": "https://api.github.com/repos/dmigwi/mainnet/commits/f6a97ec7c23468676b2b69abc7d275a85126ed19/comments",
	  "author": null,
	  "committer": null,
	  "parents": [
		{
		  "sha": "5b4ce6d2c2704255ac19b2a085928279c2476b0c",
		  "url": "https://api.github.com/repos/dmigwi/mainnet/commits/5b4ce6d2c2704255ac19b2a085928279c2476b0c",
		  "html_url": "https://github.com/dmigwi/mainnet/commit/5b4ce6d2c2704255ac19b2a085928279c2476b0c"
		}
	  ]
	}
  ]`

var pathContents = `[
	{
	  "name": "1",
	  "path": "27f87171d98b7923a1bd2bee6affed929fa2d2a6e178b5c80a9971a92a5c7f50/1",
	  "sha": "fdda7add76ddd743395aedbb3b4792f349ad0ed7",
	  "size": 0,
	  "url": "https://api.github.com/repos/dmigwi/mainnet/contents/27f87171d98b7923a1bd2bee6affed929fa2d2a6e178b5c80a9971a92a5c7f50/1?ref=master",
	  "html_url": "https://github.com/dmigwi/mainnet/tree/master/27f87171d98b7923a1bd2bee6affed929fa2d2a6e178b5c80a9971a92a5c7f50/1",
	  "git_url": "https://api.github.com/repos/dmigwi/mainnet/git/trees/fdda7add76ddd743395aedbb3b4792f349ad0ed7",
	  "download_url": null,
	  "type": "dir",
	  "_links": {
		"self": "https://api.github.com/repos/dmigwi/mainnet/contents/27f87171d98b7923a1bd2bee6affed929fa2d2a6e178b5c80a9971a92a5c7f50/1?ref=master",
		"git": "https://api.github.com/repos/dmigwi/mainnet/git/trees/fdda7add76ddd743395aedbb3b4792f349ad0ed7",
		"html": "https://github.com/dmigwi/mainnet/tree/master/27f87171d98b7923a1bd2bee6affed929fa2d2a6e178b5c80a9971a92a5c7f50/1"
	  }
	},
	{
	  "name": "2",
	  "path": "27f87171d98b7923a1bd2bee6affed929fa2d2a6e178b5c80a9971a92a5c7f50/2",
	  "sha": "29316ee748643fe53bd90bf3c49bbd7ecb4201b9",
	  "size": 0,
	  "url": "https://api.github.com/repos/dmigwi/mainnet/contents/27f87171d98b7923a1bd2bee6affed929fa2d2a6e178b5c80a9971a92a5c7f50/2?ref=master",
	  "html_url": "https://github.com/dmigwi/mainnet/tree/master/27f87171d98b7923a1bd2bee6affed929fa2d2a6e178b5c80a9971a92a5c7f50/2",
	  "git_url": "https://api.github.com/repos/dmigwi/mainnet/git/trees/29316ee748643fe53bd90bf3c49bbd7ecb4201b9",
	  "download_url": null,
	  "type": "dir",
	  "_links": {
		"self": "https://api.github.com/repos/dmigwi/mainnet/contents/27f87171d98b7923a1bd2bee6affed929fa2d2a6e178b5c80a9971a92a5c7f50/2?ref=master",
		"git": "https://api.github.com/repos/dmigwi/mainnet/git/trees/29316ee748643fe53bd90bf3c49bbd7ecb4201b9",
		"html": "https://github.com/dmigwi/mainnet/tree/master/27f87171d98b7923a1bd2bee6affed929fa2d2a6e178b5c80a9971a92a5c7f50/2"
	  }
	},
	{
	  "name": "3",
	  "path": "27f87171d98b7923a1bd2bee6affed929fa2d2a6e178b5c80a9971a92a5c7f50/3",
	  "sha": "9277aab787a22bd0e195e14ebfc969ad2dc716f7",
	  "size": 0,
	  "url": "https://api.github.com/repos/dmigwi/mainnet/contents/27f87171d98b7923a1bd2bee6affed929fa2d2a6e178b5c80a9971a92a5c7f50/3?ref=master",
	  "html_url": "https://github.com/dmigwi/mainnet/tree/master/27f87171d98b7923a1bd2bee6affed929fa2d2a6e178b5c80a9971a92a5c7f50/3",
	  "git_url": "https://api.github.com/repos/dmigwi/mainnet/git/trees/9277aab787a22bd0e195e14ebfc969ad2dc716f7",
	  "download_url": null,
	  "type": "dir",
	  "_links": {
		"self": "https://api.github.com/repos/dmigwi/mainnet/contents/27f87171d98b7923a1bd2bee6affed929fa2d2a6e178b5c80a9971a92a5c7f50/3?ref=master",
		"git": "https://api.github.com/repos/dmigwi/mainnet/git/trees/9277aab787a22bd0e195e14ebfc969ad2dc716f7",
		"html": "https://github.com/dmigwi/mainnet/tree/master/27f87171d98b7923a1bd2bee6affed929fa2d2a6e178b5c80a9971a92a5c7f50/3"
	  }
	}
  ]`

var commitHistory = `{
	"sha": "d15fa2317cc5f37508a470413e09f11e2a5faa9c",
	"node_id": "MDY6Q29tbWl0MTcyODA0MzIyOmQxNWZhMjMxN2NjNWYzNzUwOGE0NzA0MTNlMDlmMTFlMmE1ZmFhOWM=",
	"commit": {
	  "author": {
		"name": "Politeia",
		"email": "noreply@decred.org",
		"date": "2018-11-01T04:58:07Z"
	  },
	  "committer": {
		"name": "Politeia",
		"email": "noreply@decred.org",
		"date": "2018-11-01T04:58:07Z"
	  },
	  "message": "Flush vote journals.\n\n27f87171d98b7923a1bd2bee6affed929fa2d2a6e178b5c80a9971a92a5c7f50\nbc8776180b5ea8f5d19e7d08e9fcc35f0d1e3d16974963e3e5ded65139e7b092",
	  "tree": {
		"sha": "c1e157da2c5e1fb61246002d85a2f99d91cda237",
		"url": "https://api.github.com/repos/dmigwi/mainnet/git/trees/c1e157da2c5e1fb61246002d85a2f99d91cda237"
	  },
	  "url": "https://api.github.com/repos/dmigwi/mainnet/git/commits/d15fa2317cc5f37508a470413e09f11e2a5faa9c",
	  "comment_count": 0,
	  "verification": {
		"verified": false,
		"reason": "unsigned",
		"signature": null,
		"payload": null
	  }
	},
	"url": "https://api.github.com/repos/dmigwi/mainnet/commits/d15fa2317cc5f37508a470413e09f11e2a5faa9c",
	"html_url": "https://github.com/dmigwi/mainnet/commit/d15fa2317cc5f37508a470413e09f11e2a5faa9c",
	"comments_url": "https://api.github.com/repos/dmigwi/mainnet/commits/d15fa2317cc5f37508a470413e09f11e2a5faa9c/comments",
	"author": null,
	"committer": null,
	"parents": [
	  {
		"sha": "dee00630b85022aaace5ece09089e3cf8a9f2d19",
		"url": "https://api.github.com/repos/dmigwi/mainnet/commits/dee00630b85022aaace5ece09089e3cf8a9f2d19",
		"html_url": "https://github.com/dmigwi/mainnet/commit/dee00630b85022aaace5ece09089e3cf8a9f2d19"
	  }
	],
	"stats": {
	  "total": 83,
	  "additions": 83,
	  "deletions": 0
	},
	"files": [
	  {
		"sha": "e36a48088b8669964f853d78276a4d7277d688dd",
		"filename": "27f87171d98b7923a1bd2bee6affed929fa2d2a6e178b5c80a9971a92a5c7f50/3/plugins/decred/ballot.journal",
		"status": "modified",
		"additions": 65,
		"deletions": 0,
		"changes": 65,
		"blob_url": "https://github.com/dmigwi/mainnet/blob/d15fa2317cc5f37508a470413e09f11e2a5faa9c/27f87171d98b7923a1bd2bee6affed929fa2d2a6e178b5c80a9971a92a5c7f50/3/plugins/decred/ballot.journal",
		"raw_url": "https://github.com/dmigwi/mainnet/raw/d15fa2317cc5f37508a470413e09f11e2a5faa9c/27f87171d98b7923a1bd2bee6affed929fa2d2a6e178b5c80a9971a92a5c7f50/3/plugins/decred/ballot.journal",
		"contents_url": "https://api.github.com/repos/dmigwi/mainnet/contents/27f87171d98b7923a1bd2bee6affed929fa2d2a6e178b5c80a9971a92a5c7f50/3/plugins/decred/ballot.journal?ref=d15fa2317cc5f37508a470413e09f11e2a5faa9c",
		"patch": "@@ -9415,3 +9415,68 @@\n {\"version\":\"1\",\"action\":\"add\"}{\"castvote\":{\"token\":\"27f87171d98b7923a1bd2bee6affed929fa2d2a6e178b5c80a9971a92a5c7f50\",\"ticket\":\"bb93071d144f751966a75ed8049601ac8df934e2a74373d2ecdb60ebf6a6a94b\",\"votebit\":\"1\",\"signature\":\"1fd0bcad3cdb378fe9b388def5bdd7d71832d31c16be6359bb4d31b3230d1e8f525c08e0f72ce30e406b8185275b3254bc541b7fc1bc202c007865cff973896163\"},\"receipt\":\"53ebb75bf794c6b6e03c38c9b69fec180d311c260b79470249d02510ce32f24ce1f6d0bf4ff6ce1538be4a880b961f7a30261a9b1a6438b219712514be428005\"}\n {\"version\":\"1\",\"action\":\"add\"}{\"castvote\":{\"token\":\"27f87171d98b7923a1bd2bee6affed929fa2d2a6e178b5c80a9971a92a5c7f50\",\"ticket\":\"ef10f3aa28bb510acc4390a54f7baf411a31b2ce74b439f0cac7ddcc848b7e60\",\"votebit\":\"2\",\"signature\":\"20db884804b1a2d47c73bc7e5b690947fb5bb2c95c39f95a69caa9f730d023ffcf0e68f86dc720a90df7025eb7f67a6f49e3c62574b754406c92e8d2b4831e170c\"},\"receipt\":\"9040532f5a09f3b390aff2daba97179fddd430456a072d3992e6d988077f878ba11feb97acc170bcfbe69850b839debd9e7aaf4e64783bfa5c1ee84d2821ac0a\"}\n {\"version\":\"1\",\"action\":\"add\"}{\"castvote\":{\"token\":\"27f87171d98b7923a1bd2bee6affed929fa2d2a6e178b5c80a9971a92a5c7f50\",\"ticket\":\"1a6202c5448f71890e6bd0a49692ccd3307c783c7f801712ab66f922927fa24c\",\"votebit\":\"1\",\"signature\":\"1f4c2076b6a912e7d9536d32dfc60e672bf2278578ddf3decec32470e2139a472100834031bc0fc6dc05e9a2584da5ad8c1de3dd1c8ca1fd9e038f63f8f85fe96c\"},\"receipt\":\"f4921f87ea828d2f5ee135d9f4477712ea3f939ef37304c03f523129028e588876b34a5a319973fc24bc9162356b395ba2b78abfc90e78c80218ffaba59bc10c\"}\n+{\"version\":\"1\",\"action\":\"add\"}{\"castvote\":{\"token\":\"27f87171d98b7923a1bd2bee6affed929fa2d2a6e178b5c80a9971a92a5c7f50\",\"ticket\":\"ef262b7ec242febb2bd266033d7fc962926a687d5d2b7dbcae33a1c32e86126e\",\"votebit\":\"2\",\"signature\":\"1f8bea3a3a73d90f0094263d516f53e1e7eb145b4f7b29ef527b3a7b2a689c8fd544835c2ee363dbde4e09562c71346500859b4f4c9a94d6ac8713910dd11ebf5a\"},\"receipt\":\"1344b114d666e357bf388d2effebe58ad3ce0cfe66a96e1405e1406aa9dc09ff0e5fd6c2088316e568f0083245c0b1338a207f518f2d713f4e4b0b69bb21e908\"}\n+{\"version\":\"1\",\"action\":\"add\"}{\"castvote\":{\"token\":\"27f87171d98b7923a1bd2bee6affed929fa2d2a6e178b5c80a9971a92a5c7f50\",\"ticket\":\"efc84f62b764b0288de5475fbc0991dc077bec06d35ddfe98fdae96089b5d651\",\"votebit\":\"2\",\"signature\":\"207fcc0bf3b9a2913e6f7dd570dfc32ef27793b96550a8fe24a09e85e5236310cb3f74ea556dc8d43f19a31ac74bda13d6436ff199a5449ee6e40a169fbbba000f\"},\"receipt\":\"5b7b1512367cc1ef528a7a58aad4096932acf680f80749deeb48e4771dbbbddc36f9ed6c4a427ff35519251720b291c2f1c863627d53d52fbe2e5284579d790d\"}\n+{\"version\":\"1\",\"action\":\"add\"}{\"castvote\":{\"token\":\"27f87171d98b7923a1bd2bee6affed929fa2d2a6e178b5c80a9971a92a5c7f50\",\"ticket\":\"2469fc75e79d1ffce699b39fa678c56ce2d91747fd4ac013e5db9f3ce0c85a1f\",\"votebit\":\"1\",\"signature\":\"20a7dcd2b87c9e19ebfb06f64efe4fca28ae1531520fdac7d4a3ccdeaadee95255024d39dfc54ea30fb2078b6b019ff3a4f7a3afd2e1c75dbaca709365219fc43d\"},\"receipt\":\"df503a7dfaaacee20e1a5f1003f99d651a996d459f205628650c1c718d1c853805f96f3f7787213845528bd8bb2dc8726c9a5ff6f1a81a8db62aaf19f82fda07\"}\n+{\"version\":\"1\",\"action\":\"add\"}{\"castvote\":{\"token\":\"27f87171d98b7923a1bd2bee6affed929fa2d2a6e178b5c80a9971a92a5c7f50\",\"ticket\":\"000c9a7e12cef0f86a73c0671df521bfeccbf5b13362355d236286c87840626b\",\"votebit\":\"1\",\"signature\":\"203b71ccf041c9e6de5f63a35ced923eceb824220f9e1125c1a5dae5ad3f531d9b47513a3726448761f5a2dd18d3e91d16063cbed6d081d462898a527861da4b54\"},\"receipt\":\"96d62bfd1c8c283d6d14f5cd04342549c3e8d74bdb2f8460ef9db4a5d7b05fb2abb81ce0bdf9f0339f434af555679a8a47a92dd686dbc7142ac5ae2a9f1d7109\"}\n+{\"version\":\"1\",\"action\":\"add\"}{\"castvote\":{\"token\":\"27f87171d98b7923a1bd2bee6affed929fa2d2a6e178b5c80a9971a92a5c7f50\",\"ticket\":\"f00ab37a8c2ba58c9898bea55b7534be79dda8972940e46c56bd6cbaeb23c4c2\",\"votebit\":\"2\",\"signature\":\"20d5172b9d9b83e5750f78b194b133435c9019ed92dba3faa1843332e549d3addd2ea6bb4d4f4284ab0b7427307ec3e8d5fb2d1eabcc41d4815e309ff91a76920b\"},\"receipt\":\"a2ba3a04424edb262ceaef2b8ace5e6a80b3ffd3cac8ed5e0fd728bf9a7ea71e5aab60891ba4179652bf8d9b2b363a98d0a4b9a7995762c92e97fa426a08370a\"}\n+{\"version\":\"1\",\"action\":\"add\"}{\"castvote\":{\"token\":\"27f87171d98b7923a1bd2bee6affed929fa2d2a6e178b5c80a9971a92a5c7f50\",\"ticket\":\"a81a32d173dafabbfd3bc09682a594981ddee1a62beb3d35928f4183443cc95e\",\"votebit\":\"1\",\"signature\":\"20546d2c8aaf894c5f0cef43bc29c855a8c6c89532d996753534307a46838d98145f9aa26b58823dd84c916f96169c67e7df7bcccd839126ce70ad3677200145fd\"},\"receipt\":\"2d9b08a2bc7fc1773767cc30dc621b015cb0c666bced2f4df86a838081bad4217d3d11619e21afeb3b0b369eb9d02f6bd444f797bdff699e352dab9f98593d0a\"}\n+{\"version\":\"1\",\"action\":\"add\"}{\"castvote\":{\"token\":\"27f87171d98b7923a1bd2bee6affed929fa2d2a6e178b5c80a9971a92a5c7f50\",\"ticket\":\"deec0c1dd84e74b77e6484a0d9b5f68cd191a035d95656e3338ba1ca41658b80\",\"votebit\":\"1\",\"signature\":\"2033d0dfb742a03484c17bfb35795146f1b30ccb76d7a07e9525ee4ffec511d3706003d224d12c25b351deabf116ce85bfe94292999a3646ea0ba1b8e7f5d943d8\"},\"receipt\":\"b5de799a85b7ff0d5fc6433fd09e032d78530cf17e1ae7a88888963b7c475afdd689b396f1824f0a689884b38c49a74d70f1bf170ab633b6154bf16e3cbdd201\"}\n+{\"version\":\"1\",\"action\":\"add\"}{\"castvote\":{\"token\":\"27f87171d98b7923a1bd2bee6affed929fa2d2a6e178b5c80a9971a92a5c7f50\",\"ticket\":\"016ac5c073a8b5ca008de7d41c5fb200c3d87c239958f7500321dd13e52da298\",\"votebit\":\"1\",\"signature\":\"1fe69de7fb63279f3c54b0954389960193e0f1fca3c9c614056f0d08a8d29780983498a3b30c0f79abfdbe91acf2649b6e48aae3a16dbacfee2b90ef34d79d1650\"},\"receipt\":\"e4e937315a272090cbe9878c16b3a89d3ea4967c26b14c31ff6e0e1175d4a7c381c7f3ba31bfed9524b21a3db1b04648358950711238497174c2b26b237c8e07\"}\n+{\"version\":\"1\",\"action\":\"add\"}{\"castvote\":{\"token\":\"27f87171d98b7923a1bd2bee6affed929fa2d2a6e178b5c80a9971a92a5c7f50\",\"ticket\":\"f19905e6d782d27912779752ba6ea960c2324050903a620d156093d056a3db0f\",\"votebit\":\"2\",\"signature\":\"1f953b5d2eda7d9a8d2d5801484e30f1efa8c3bd9e54b2477b9743ee3e6abe84d16098d175d0984f4f5ef3c07440981cbbd646a0ff5d2e14cb656028d0122794fc\"},\"receipt\":\"412741608fb7ee0ac1e887ce7d2b689edcdde65739bca38afd2026e6bcd4c6a525a6bccb8c7f9d7b9e628ae0b568ecc069da5d40e6a832ec3e5d1de1fdcb6c08\"}\n+{\"version\":\"1\",\"action\":\"add\"}{\"castvote\":{\"token\":\"27f87171d98b7923a1bd2bee6affed929fa2d2a6e178b5c80a9971a92a5c7f50\",\"signature\":\"2097c7f5f47b5bf29cd640ff2e2924f52a20829246166574055c34dcc7b2d5038a0eecae85131e08f5d6c51320495ef41bb55cfaf223675a50338112528aab5b9f\"},\"receipt\":\"b94fb6eb28589e98f6c858287906a62b04d5c5bcd803ffd75983b7129a0af63995dc183fecb381f5311b421d53da66a5ded1ed80b9781fda52129d0e31040901\"}\n+{\"version\":\"1\",\"action\":\"add\"}{\"castvote\":{\"token\":\"27f87171d98b7923a1bd2bee6affed929fa2d2a6e178b5c80a9971a92a5c7f50\",\"ticket\":\"02b471689e31586ee4c189a7b1122383868d5ffa0d7556142d6e9501da8f0fe4\",\"votebit\":\"1\",\"signature\":\"20157d81854ac952581b77ab7930300fb1346de4c0a58bb5c7068edf7f3ee8e1a63e640133f47a0fbbdc0d0e8baabab2c629c3cbe31589cdab9e713f566b979582\"},\"receipt\":\"348c7818d885191ce289549a6eafb58c4eb2e6167ef115c8c962528ede1869e19f5d9853b5b2d0d96708c2d2516ec680a190ad20e757e65e10e2a301a59e480c\"}\n+{\"version\":\"1\",\"action\":\"add\"}{\"castvote\":{\"token\":\"27f87171d98b7923a1bd2bee6affed929fa2d2a6e178b5c80a9971a92a5c7f50\",\"ticket\":\"381ede1199381e90337347e61c9ed065ee152da4c36574a34c236e2245505fc6\",\"votebit\":\"1\",\"signature\":\"1f13d929e6cbf0809897c8748cf7bd9286fdea2e4d6c7f078ad782c56333dc53196dc59cdbf87f4130d702ada45bc1dff72c0082a13b4a8a82332a6278cb5dc232\"},\"receipt\":\"8fd8e4dad6d918ca1fc905980a87605f9641b86974b50001cb2bf9cf95895d1737ccbd9e6b05ca67ff709f90aceebe280fe4e1487b27f10b2133b1d78a3b070d\"}\n+{\"version\":\"1\",\"action\":\"add\"}{\"castvote\":{\"token\":\"27f87171d98b7923a1bd2bee6affed929fa2d2a6e178b5c80a9971a92a5c7f50\",\"ticket\":\"f296078bc7ab77b14458250f6cfa56c7b4f992be509f15d6e687c82ece8bb858\",\"votebit\":\"2\",\"signature\":\"1f4cde5d24c2c3a2059ad3a55168b43c76d4b8ec02747ca65e1680ed7db3e1dba74cf5eb92f9c3b49c27f0157eef57dc296777508f4681f0a7d0987efbd93fd045\"},\"receipt\":\"29e5ed9cf0673587ba5fe296fe75db109b3f15b57095c32ac1be0063dfdcad38d6dbfe46251016cc1622448deb4a7342a88be5c6c8cc7cb88c3dab54a039bb0a\"}\n+{\"version\":\"1\",\"action\":\"add\"}{\"castvote\":{\"token\":\"27f87171d98b7923a1bd2bee6affed929fa2d2a6e178b5c80a9971a92a5c7f50\",\"ticket\":\"35534b9f0243cfadc82f7a05fa41ea778cf2108474fe2e533644734ca7b279e3\",\"votebit\":\"1\",\"signature\":\"20b96a983af0c45f8610361fe7a4e6baefdc4e9e769a6bd8763056f4b728c0358514101539fd9f40d15c364abe6ab628eb787e9979f302a7ccbca647e44da38013\"},\"receipt\":\"dee5467e3269d2bba78a1b511efa86e547e6e5d052f2802c9c0c104a5e22873fd909aa793b49609f05acac3864b9f44b2168dc9f48d37e7040815f09d9970e0d\"}
	  },
	  {
		"sha": "2d80361044eb32f3b416284155e3258b6aaabfad",
		"filename": "bc8776180b5ea8f5d19e7d08e9fcc35f0d1e3d16974963e3e5ded65139e7b092/6/plugins/decred/ballot.journal",
		"status": "modified",
		"additions": 18,
		"deletions": 0,
		"changes": 18,
		"blob_url": "https://github.com/dmigwi/mainnet/blob/d15fa2317cc5f37508a470413e09f11e2a5faa9c/bc8776180b5ea8f5d19e7d08e9fcc35f0d1e3d16974963e3e5ded65139e7b092/6/plugins/decred/ballot.journal",
		"raw_url": "https://github.com/dmigwi/mainnet/raw/d15fa2317cc5f37508a470413e09f11e2a5faa9c/bc8776180b5ea8f5d19e7d08e9fcc35f0d1e3d16974963e3e5ded65139e7b092/6/plugins/decred/ballot.journal",
		"contents_url": "https://api.github.com/repos/dmigwi/mainnet/contents/bc8776180b5ea8f5d19e7d08e9fcc35f0d1e3d16974963e3e5ded65139e7b092/6/plugins/decred/ballot.journal?ref=d15fa2317cc5f37508a470413e09f11e2a5faa9c",
		"patch": "@@ -9042,3 +9042,21 @@\n {\"version\":\"1\",\"action\":\"add\"}{\"castvote\":{\"token\":\"bc8776180b5ea8f5d19e7d08e9fcc35f0d1e3d16974963e3e5ded65139e7b092\",\"ticket\":\"9544954e46b6c1219efea0f7b1ae4e9937484443beff73457d9285a609fd85b0\",\"votebit\":\"1\",\"signature\":\"2040a247b8a7283ad599bf25f44e8f5f255610fa1e1cfc69b7c99795a7589991870551c3e17a65c86f62bcdbdfd987b3996c1c8efc6a59120a46c41412a4bbf6d5\"},\"receipt\":\"93f8fdf375a05706b13eef21563362fd07f7f8883043561f58ae0bb27c2008d4b58017b9674d593d7edfa839bed9f96873e8329f55f14fddbaa2afca0fec520f\"}\n {\"version\":\"1\",\"action\":\"add\"}{\"castvote\":{\"token\":\"bc8776180b5ea8f5d19e7d08e9fcc35f0d1e3d16974963e3e5ded65139e7b092\",\"ticket\":\"cb48bc60d8c2e0924f1467a6d4aab8a73b25a1a2047530e37a35652ed8176242\",\"votebit\":\"1\",\"signature\":\"201e22163837c01d4f7030c29e901a5b57927be88344a522810dfae5ec04e2cd9912eb2ee45c57967a65e59bd6ad7b374de93b3f5994e550c9054aba9202234eb7\"},\"receipt\":\"b5d11fd79b168e5bba1d2959d42129c3c38880bdf8269e19b965f684e7a0862de9b50406e71e6445f6419a910ef90e30e2785e57ad51959aaad20fd7becece06\"}\n {\"version\":\"1\",\"action\":\"add\"}{\"castvote\":{\"token\":\"bc8776180b5ea8f5d19e7d08e9fcc35f0d1e3d16974963e3e5ded65139e7b092\",\"ticket\":\"ebb0625d1cedfa68fdbc184ea482980c888a5d8bf81ff5b9dfbd28bbfc6d71e7\",\"votebit\":\"1\",\"signature\":\"2064e53976814bf4bdef2ee0ef2e8a3ee52cac08a9eaa33346a0affac557ad5ee720325facc6829e93b6cba2120096961385f48d3c64dc15e5018fecf034be115e\"},\"receipt\":\"a2c3de6830e6783e726214511ddbb4c39bd1b840cb591bdf3bad32eea2d5db3e4f129414902c794dc8effa1f705e2bf8163d6a5b307654f71d9c1b6112a9ec09\"}\n+{\"version\":\"1\",\"action\":\"add\"}{\"castvote\":{\"token\":\"bc8776180b5ea8f5d19e7d08e9fcc35f0d1e3d16974963e3e5ded65139e7b092\",\"ticket\":\"000c9a7e12cef0f86a73c0671df521bfeccbf5b13362355d236286c87840626b\",\"votebit\":\"2\",\"signature\":\"20fe8b4d8b1097b7c6831a90c772a3fbbd8192eb59c6f3b71f6c2efb2b1785f691734a62bd678046ef6cda1049fc2c48bc92f095e2505191621f60507ba1c5c420\"},\"receipt\":\"97df31bc1767b821c2f37aa2109c50993c8ffd71e8f252ecda59455942dabd0831876c2dfc9935389349d8aab3f6f7d86169b2a61fc0c7457470353811101e0d\"}\n+{\"version\":\"1\",\"action\":\"add\"}{\"castvote\":{\"token\":\"bc8776180b5ea8f5d19e7d08e9fcc35f0d1e3d16974963e3e5ded65139e7b092\",\"ticket\":\"006298836801cab607dc06f40f91122d3e2af97a337a75494f432a5923fe0701\",\"votebit\":\"2\",\"signature\":\"1f5d551e4bec3d34369929a347a2b76a9d2f97374d72644b1c9bde54f35f6dfea615a486ae4b2e01cf416453afd88fd4ee95d676362935f3d26d81e09e74867e96\"},\"receipt\":\"8a0610d5ca9fc0056c4e0efed682483be9f9ff0b2c0cb0161dc7e962d940ed0224c6714d417e4e93d99e23c924570eb739fc771d6131886d1751ef7dca92520d\"}\n+{\"version\":\"1\",\"action\":\"add\"}{\"castvote\":{\"token\":\"bc8776180b5ea8f5d19e7d08e9fcc35f0d1e3d16974963e3e5ded65139e7b092\",\"ticket\":\"0066e6b5e4657cd80f5cfd72092921ea5a392d86072910b9214eebfcc8f7a35b\",\"votebit\":\"2\",\"signature\":\"1fd5ec85c8eb6d77b65e74295ddbf12bb95e6bf58b0cc34ce425b1839f9e3e4baf35c89e249552dbb0df8001d671db195d1b3f2f7ebabac7ac12d377b4ca466b12\"},\"receipt\":\"541c94bd33e95632442ab59b1c5f025742179b96702eaadde3c887e8015c51b6236df22d86720319a055ab9a21b93c30f2905b5ef452683ff47c87bb74f43e0e\"}\n+{\"version\":\"1\",\"action\":\"add\"}{\"castvote\":{\"token\":\"bc8776180b5ea8f5d19e7d08e9fcc35f0d1e3d16974963e3e5ded65139e7b092\",\"ticket\":\"00ad9a497d96c5c0c55e7ee840950ab7fd2a5b26a54cc69544e3ffe3e80b13d9\",\"votebit\":\"2\",\"signature\":\"1f122a5744efde63b3d37db2f66e5efd25b431bc9bd3216b712e9f83dfc664acc82da3bb18d930d3a55a63e662b02bc5ae03425b1331a6077fd51baf72d9ec2746\"},\"receipt\":\"a69af47b735cfb7e54c3a6ceb31c6ab48cb52b23c0dd8a829bdad2672466d80ca537f3141a9ef03f93f1eca7209b92ad58fd2cba0a145786c795a0a796944902\"}\n+{\"version\":\"1\",\"action\":\"add\"}{\"castvote\":{\"token\":\"bc8776180b5ea8f5d19e7d08e9fcc35f0d1e3d16974963e3e5ded65139e7b092\",\"ticket\":\"00f6b17c7514ec3e1697cd0aed340a8799c799d9855b2899bab09fdfa811f612\",\"votebit\":\"2\",\"signature\":\"1fd2a059ab08ea24f40acf812bf7a86e4ae505fb04a1724bb3288331acff94b9f76aac4b2e3f77d31fe809206eb3d68fc5251a79c16478def999a7336e17272e7a\"},\"receipt\":\"f9e405239e8a39575ec0ffd16ee18155f252c95199a6da9b75e983426c847b11024740938b2a4110c99382f5cbd6a0541cd40367c241a6bf2401667e23db7e06\"}\n+{\"version\":\"1\",\"action\":\"add\"}{\"castvote\":{\"token\":\"bc8776180b5ea8f5d19e7d08e9fcc35f0d1e3d16974963e3e5ded65139e7b092\",\"ticket\":\"01698537f1648f4c7c8861a6cfb33d4b01a6a35aef74a48aa4d11a5c2d6c3946\",\"votebit\":\"2\",\"signature\":\"208b956adf363600942bc1f9a3adda6f53c12b0ce63262a1ae507ffe159f3be1f57e2de3ad409bfe491af36500bf09c55baa10ab248faac01e372dc125321748e2\"},\"receipt\":\"d09a0a3cb47f51cbcaa25db38c60ddc64d379ebd882fc656b585bc67025904cd7f761e77dfc13d94c5a820e6c1bc1b9609fa53e9fa4c894d0ed8dc067ca18f03\"}\n+{\"version\":\"1\",\"action\":\"add\"}{\"castvote\":{\"token\":\"bc8776180b5ea8f5d19e7d08e9fcc35f0d1e3d16974963e3e5ded65139e7b092\",\"ticket\":\"016ac5c073a8b5ca008de7d41c5fb200c3d87c239958f7500321dd13e52da298\",\"votebit\":\"2\",\"signature\":\"1f7599fdb5fd0647fd7ab00d864e3872fb099cce36e9476c074810e3fe09c72fed22faabba6915ea0f1215605347098c58c5a7cf091337e03b5c69309f6c719d6c\"},\"receipt\":\"a1c2dd0a26dc564837f7df960a9bf4bbcfcba35f46fa222fa15b24b0a54ff3ede29ff369200726baa224b6b0a76919c8106b65db4793f99f79cf5c0361c4e50a\"}\n+{\"version\":\"1\",\"action\":\"add\"}{\"castvote\":{\"token\":\"bc8776180b5ea8f5d19e7d08e9fcc35f0d1e3d16974963e3e5ded65139e7b092\",\"ticket\":\"0176545607a0b9a14caa0c1050acaed5352c301deb54d5fa95ee5b6d65a72b0d\",\"votebit\":\"2\",\"signature\":\"2060934560fe1910f2a1b3fabe557a9a740aaa43ac6d686bb78404ff6845adeb3d4ce423493b0e17760a322437ffc53b9375709f0c2c899319cd73d6ea8e46123b\"},\"receipt\":\"a60f56ae93fb59a678cc6627476869ccb8d4cd756858e27285a2ca01dfa3ec0afdb9a04d6f4643066d1e1637abce75d8b59bc668709f58bce2966f4372fd400e\"}\n+{\"version\":\"1\",\"action\":\"add\"}{\"castvote\":{\"token\":\"bc8776180b5ea8f5d19e7d08e9fcc35f0d1e3d16974963e3e5ded65139e7b092\",\"ticket\":\"0178070aef01eaa74aef102e08dd2c30b4a98404c98fa9346542f4618bc8457d\",\"votebit\":\"2\",\"signature\":\"1fe110b7343b1883e3e4833c432214f7833345a9e74ab9202fdcae80e995aa930847d6e38f75000d771da9012adc21c7914c72317f568db95ea823f47c7ca2c3e6\"},\"receipt\":\"f4535972403c92533eff66ed24ed4423925d723fc4ca8b15faf07515250bb8bfd5db0bdd2e5de39c6989538dcaa0b0ebb0777a721a8c6f76ff3308fa448a2100\"}
	  }
	]
  }`
