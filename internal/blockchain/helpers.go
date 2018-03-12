package blockchain

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

//helper to get response of nodeChain
func getNodeChain(node string) *Blockchain {
	var bc *Blockchain
	url := fmt.Sprintf("%v/chain", node)
	req, _ := http.NewRequest("GET", url, nil)

	// get []bytes & unmarshal into *Blockchain
	data := executeRequest(req)

	err := json.Unmarshal(data, &bc)
	if err != nil {
		fmt.Println(err)
	}
	return bc
}

//executeRequest wraps an HTTP client and returns request response as []bytes.
func executeRequest(r *http.Request) []byte {
	c := &http.Client{Timeout: time.Second * 60}
	resp, err := c.Do(r)
	if err != nil {
		fmt.Println(err)
	}

	// Read the response body
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	return data
}
