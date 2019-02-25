package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/dmigwi/go-piparser/v1/types"
)

const historyList = "https://api.github.com/repos/decred-proposals/mainnet/commits?path=27f87171d98b7923a1bd2bee6affed929fa2d2a6e178b5c80a9971a92a5c7f50/3/plugins/decred/ballot.journal"
const historyContent = "https://api.github.com/repos/decred-proposals/mainnet/commits/11ac7a9b777d94e1e5c14b0abafab41c744b51ee"

func main() {
	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}
	client := &http.Client{Transport: tr}

	t := time.Now()
	resp, err := client.Get(historyContent)
	if err != nil {
		log.Fatalf("client.Get unexpected error occured: %v", err)
		return
	}
	fmt.Println("client.Get >>> Took :", time.Since(t))

	t = time.Now()
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("ioutil.ReadAll unexpected error occured: %v", err)
		return
	}
	fmt.Println("ioutil.ReadAll >>> Took :", time.Since(t))

	types.ProposalToken = "27f87171d98b7923a1bd2bee6affed929fa2d2a6e178b5c80a9971a92a5c7f50"

	t = time.Now()
	var data types.History
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Fatalf("json.Unmarshal unexpected error occured: %v", err)
		return
	}
	fmt.Println("json.Unmarshal >>> Took :", time.Since(t))

	t = time.Now()
	m := *data.Files[0].Data
	// index := len(*data.Files[0].Data) - 1
	s, _ := json.Marshal(m[0])
	fmt.Println("First >>> ", string(s))

	s, _ = json.Marshal(m[len(m)-1])
	fmt.Println("Last >>>> ", string(s))
	fmt.Println(" >>> Took :", time.Since(t))
}
