package main

import (
	"html/template"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/dmigwi/go-piparser/v1/proposals"
	"github.com/gorilla/mux"
)

var tmpl *template.Template
var parser *proposals.Parser

// chartData defines the charts data to be used by charts js.
type chartData struct {
	Yes  []int       // Yes
	No   []int       // No
	Date []time.Time // Date
}

// Load the template files first.
func init() {
	tmpl = template.Must(template.ParseFiles("index.tmpl"))
}

// handleProposal handles /{proposal-token} route.
func handleProposal(w http.ResponseWriter, r *http.Request) {
	proposalToken := mux.Vars(r)["token"]

	log.Printf("Retrieving details for %s ...\n", proposalToken)
	data, err := parser.Proposal(proposalToken)
	if err != nil {
		log.Fatalf("unexpected error occured: %v", err)
	}

	log.Printf("Processing charts data for %s ...\n", proposalToken)
	var graph1, graph2 chartData
	var yes, no int

	for _, val := range data {
		var y, n int
		for _, vote := range val.VotesInfo {
			switch vote.VoteBit {
			case "No":
				no++
				n++
			case "Yes":
				yes++
				y++
			default:
				log.Fatalf("Invalid vote bit found: %v", vote.VoteBit)
			}
		}

		graph1.Yes = append(graph1.Yes, yes)
		graph1.No = append(graph1.No, no)
		graph1.Date = append(graph1.Date, val.Date)

		graph2.Yes = append(graph2.Yes, y)
		graph2.No = append(graph2.No, n)
		graph2.Date = append(graph2.Date, val.Date)
	}

	payload := struct {
		Cummulative chartData
		VotesByTime chartData
		Token       string
	}{
		graph1,
		graph2,
		proposalToken,
	}

	c := len(graph1.Yes)
	if c > 1 {
		var yes, no, total int

		c--
		yes, no = graph1.Yes[c], graph1.No[c]
		total = yes + no

		log.Printf("Found Yes: %d No: %d and Total %d \n", yes, no, total)
	} else {
		log.Fatalf(" **Found no valid chart data** ")
	}

	err = tmpl.Execute(w, payload)
	if err != nil {
		log.Fatalf("error found: %v", err)
	}

	log.Println("Done.")
}

func main() {
	log.Println("Please Wait... Setting up the environment")

	var err error
	cloneDir := "~/playground"

	parser, err = proposals.NewExplorer("", "", cloneDir)
	if err != nil {
		log.Fatalf("unexpected error occured: %v", err)
		return
	}

	rtr := mux.NewRouter()
	rtr.HandleFunc("/{token:[A-z0-9]{64}}", handleProposal).Methods("GET")

	http.Handle("/", rtr)

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
