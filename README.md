# go-piparser
[![Build Status](https://travis-ci.org/dmigwi/go-piparser.svg?branch=master)](https://travis-ci.org/dmigwi/go-piparser)
[![ISC License](https://img.shields.io/badge/license-ISC-blue.svg)](http://copyfree.org)
![GitHub top language](https://img.shields.io/github/languages/top/dmigwi/go-piparser.svg?color=green&style=popout)

go-piparser is tool that parses [Politeia](https://proposals.decred.org) proposals votes data stored in github.
It adds a timestamp field obtained from the commit history. The tool makes use of the git commandline interface to clone 
and query the politeia votes data. Github repository updates are fetched at intervals of 1hr after setting up the
environment. [doc](https://docs.decred.org/advanced/navigating-politeia-data/#voting-and-comment-data)


## Table of Contents

- [Requirement](#requirement)
- [Installation](#installation)
- [Import go-piparser](#import-go-piparser)
- [Initialize the Explorer](#initialize-the-explorer)
- [Fetch the Proposal's Votes](#fetch-the-proposal's-votes)
- [Full Sample Program](#full-sample-program)
- [Test Client](#test-client)


## Requirement

- git -  The tool requires a functional git commandline installation.
To install git visit [here](https://git-scm.com/book/en/v2/Getting-Started-Installing-Git)

## Installation

```bash
    go get -u github.com/dmigwi/go-piparser/v1/proposals
```

## Import go-piparser

```go
    import "github.com/dmigwi/go-piparser/v1/proposals"
```

## Initialize the Explorer

```go
    repoOwner := ""
    repoName  := ""
    cloneDir  := "/path/to/root/clone/directory"

    parser, err := proposals.NewExplorer(repoOwner, repoName, cloneDir)
    if err != nil {
		log.Fatalf("unexpected error occured: %v", err)
    }
```

- `repoOwner` - defines the owner of the repository where the Politeia votes are to be queries from. If not set, it defaults to `decred-proposals`
- `repoName` - defines the name of the repository holding the Politeia votes. If not set, it defaults to `mainnet`.
- `cloneDir` - defines the directory where the said repository will be cloned into. If not set, a tmp folder is created and set.

## Fetch the Proposal's Votes

```go
    // Decred Integration for IDAX Exchange Proposal token. 
    // https://proposals.decred.org/proposals/60adb9c0946482492889e85e9bce05c309665b3438dd85cb1a837df31fbf57fb
    proposalToken := "60adb9c0946482492889e85e9bce05c309665b3438dd85cb1a837df31fbf57fb"

    // Proposal returns votes data only associated with the set proposal token. 
    data, err := parser.Proposal(proposalToken)
	if err != nil {
		log.Fatalf("unexpected error occured: %v", err)
    }
    
    ...
```

## Full Sample Program

```go 
    package main

    import (
        "log"

        "github.com/dmigwi/go-piparser/v1/proposals"
    )

    func main() {
        cloneDir := "/path/to/root/clone/directory"

        // Create a new Explorer
        parser, err := proposals.NewExplorer("", "", cloneDir)
        if err != nil {
            log.Fatalf("unexpected error occured: %v", err)
        }

        // Set the Proposal token
        proposalToken := "60adb9c0946482492889e85e9bce05c309665b3438dd85cb1a837df31fbf57fb"

        // Retrieve the proposal token's votes data.
        data, err := parser.Proposal(proposalToken)
        if err != nil {
            log.Fatalf("unexpected error occured: %v", err)
        }


        ...
    }

```

## Test Client

Find a complete test go-piparser implementation [here](https://github.com/dmigwi/go-piparser/tree/master/testutil)