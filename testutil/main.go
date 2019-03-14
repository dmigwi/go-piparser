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
		log.Printf("unexpected error occured: %v", err)
	}

	log.Printf("Processing charts data for %s ...\n", proposalToken)
	var graph chartData
	var y, n int

	for _, val := range data {
		var yes, no int
		for _, vote := range val.VotesInfo {
			switch vote.VoteBit {
			case "No":
				no++
				n++

			case "Yes":
				yes++
				y++

			default:
				log.Printf("Invalid vote bit found: %v", vote.VoteBit)
			}
		}

		graph.Yes = append(graph.Yes, yes)
		graph.No = append(graph.No, no)
		graph.Date = append(graph.Date, val.Date)
	}

	payload := struct {
		Data   chartData
		Token  string
		ErrMsg string
	}{
		graph,
		proposalToken,
		"",
	}

	if y > 0 {
		log.Printf("Found Yes: %d No: %d and Total %d \n", y, n, n+y)
	}

	if err != nil {
		payload.ErrMsg = err.Error()
	}

	if err := tmpl.Execute(w, payload); err != nil {
		log.Fatalf("error found: %v", err)
	}

	log.Println("Done.")
}

func main() {
	var err error
	cloneDir := "~/playground"

	log.Println("Setting up the environment. Please Wait...")

	parser, err = proposals.NewExplorer("", "", cloneDir)
	if err != nil {
		log.Printf("unexpected error occured: %v", err)
		return
	}

	rtr := mux.NewRouter()
	rtr.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		tmpl.Execute(w, nil)
	})
	rtr.HandleFunc("/{token:[A-z0-9]{64}}", handleProposal).Methods("GET")
	fs := http.FileServer(http.Dir("public"))

	http.Handle("/", rtr)
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	wg := new(sync.WaitGroup)

	wg.Add(1)

	go func() {
		log.Println("Serving on http://127.0.0.1:8080")

		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Fatalf("error occured: %v ", err)

			wg.Done()
		}
	}()

	wg.Wait()
}
