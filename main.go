package main

import (
	"log"
	"time"

	"github.com/dmigwi/go-piparser/v1/proposals"
)

func main() {
	log.Println("Please Wait... ")
	t := time.Now()
	parser := proposals.NewParser("", "")

	proposalToken := "27f87171d98b7923a1bd2bee6affed929fa2d2a6e178b5c80a9971a92a5c7f50"

	if err := parser.SetProposalToken(proposalToken); err != nil {
		log.Fatalf("parser.SetProposalToken unexpected error occured: %v", err)
	}

	data, err := parser.Proposal(proposalToken)
	if err != nil {
		log.Fatalf("parser.Proposal unexpected error occured: %v", err)
		return
	}

	// s, _ := json.Marshal(data)
	log.Println("Data >>> ", len(data))

	log.Println(" >>> Took :", time.Since(t))
}
