package main

import (
	"log"
	"time"

	"github.com/dmigwi/go-piparser/v1/proposals"
)

func main() {
	log.Println("Please Wait... ")
	t := time.Now()
	parser, err := proposals.NewCMDExplorer("", "", "/home/dmigwi/playground")
	if err != nil {
		log.Fatalf("unexpected error occured: %v", err)
		return
	}

	data, err := parser.Proposal("27f87171d98b7923a1bd2bee6affed929fa2d2a6e178b5c80a9971a92a5c7f50")
	if err != nil {
		log.Fatalf("unexpected error occured: %v", err)
		return
	}

	log.Println("Data >>> ", len(data))

	log.Println(" >>> Took :", time.Since(t))
}
