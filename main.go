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

	resp, err := client.Get(historyContent)
	if err != nil {
		log.Fatal("unexpected error occured: %v", err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("unexpected error occured: %v", err)
	}

	var data types.History
	err = json.Unmarshal(body, &data)

	fmt.Println(" >>>>> ", len(data.Files[0].Data))
}
