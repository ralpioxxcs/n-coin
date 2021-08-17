package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/ralpioxxcs/nocoin/blockchain"
	"github.com/ralpioxxcs/nocoin/utils"
)

const port string = ":4000"

type URL string

// MarshalText interface
func (u URL) MarshalText() ([]byte, error) {
	url := fmt.Sprintf("http://localhost%s%s", port, u)
	return []byte(url), nil
}

type URLDescription struct {
	URL         URL    `json:"url"`
	Method      string `json:"method"`
	Description string `json:"description"`
	Payload     string `json:"payload,omitempty"`
}

type AddBlockBody struct {
	Message string
}

func documentation(rw http.ResponseWriter, r *http.Request) {
	data := []URLDescription{
		{
			URL:         URL("/"),
			Method:      "GET",
			Description: "See Documentation",
		},
		{
			URL:         URL("/blocks/{id}"),
			Method:      "POST",
			Description: "See A Block",
		},
	}
	fmt.Println(data)

	// Adding json header, browser can to parse json
	rw.Header().Add("Content-Type", "application/json")

	json.NewEncoder(rw).Encode(data)
	// Equivalent to this code {{
	//b, err := json.Marshal(data)
	//utils.HandleErr(err)
	//fmt.Fprintf(rw, "%s", b)
	// }}
}

func blocks(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		rw.Header().Add("Content-Type", "application/json")
		json.NewEncoder(rw).Encode(blockchain.GetBlockchain().AllBlocks())
	case "POST":
		// decode body & binding struct
		var addBlockBody AddBlockBody
		utils.HandleErr(json.NewDecoder(r.Body).Decode(&addBlockBody))
		blockchain.GetBlockchain().AddBlock(addBlockBody.Message)
		rw.WriteHeader(http.StatusCreated)

	}
}

func main() {

	//explorer.Start()
	http.HandleFunc("/", documentation)
	http.HandleFunc("/blocks", blocks)

	fmt.Printf("Li.stening on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
