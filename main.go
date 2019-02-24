package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const gitURL = "https://api.github.com/repos/decred-proposals/mainnet/commits?path=27f87171d98b7923a1bd2bee6affed929fa2d2a6e178b5c80a9971a92a5c7f50/3/plugins/decred/ballot.journal"

func main() {
	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}
	client := &http.Client{Transport: tr}

	resp, err := client.Get(gitURL)
	if err != nil {
		fmt.Printf("unexpected error occured: %v", err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	fmt.Println(" >>>>> ", string(body))
}
