package main

import (
	"html/template"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/dmigwi/go-piparser/v1/proposals"
)

var tmpl *template.Template
var parser *proposals.Parser

// chartData defines the charts data to be used by charts js.
type chartData struct {
	Yes  int
	No   int
	Date time.Time
}

// Load the template files first.
func init() {
	tmpl = template.Must(template.ParseFiles("index.tmpl"))
}

// handleProposal handles /{proposal-token} route.
func handleProposal(w http.ResponseWriter, r *http.Request) {
	proposalToken := r.URL.

	data, err := parser.Proposal(proposalToken)
	if err != nil {
		log.Fatalf("unexpected error occured: %v", err)
	}

	var graph []chartData
	var yes, no int

	for _, val := range data {
		for _, vote := range val.VotesInfo {
			switch vote.VoteBit {
			case "No":
				no++
			case "Yes":
				yes++
			default:
				log.Fatalf("Invalid vote bit found: %v", vote.VoteBit)
			}
		}

		graph = append(graph, chartData{Yes: yes, No: no, Date: val.Date})
	}

	payload := struct {
		Data  []chartData
		Token string
	}{
		graph,
		proposalToken,
	}

	err = tmpl.Execute(w, payload)
	if err != nil {
		log.Fatalf("error found: %v", err)
	}
}

func main() {
	log.Println("Please Wait... Setting up the environment")

	var err error
	cloneDir := "/home/dmigwi/playground"

	parser, err = proposals.NewExplorer("", "", cloneDir)
	if err != nil {
		log.Fatalf("unexpected error occured: %v", err)
		return
	}

	http.HandleFunc("/{token:[A-z0-9]{64}}", handleProposal)

	wg := new(sync.WaitGroup)

	wg.Add(1)

	go func() {
		log.Println("Serving on 127.0.0.1:8080")

		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Fatalf("error occured: %v ", err)

			wg.Done()
		}
	}()

	wg.Wait()
}
