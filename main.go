package main

import (
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/ralpioxxcs/nocoin/blockchain"
)

const (
	port        = ":4001"
	templateDir = "templates/"
)

var templates *template.Template

type homeData struct {
	PageTitle string
	Blocks    []*blockchain.Block
}

func home(rw http.ResponseWriter, r *http.Request) {
	data := homeData{"Home", blockchain.GetBlockchain().AllBlocks()}
	templates.ExecuteTemplate(rw, "home", data)
}

func add(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		templates.ExecuteTemplate(rw, "add", nil)
	case "POST":
		r.ParseForm()
		data := r.Form.Get("blockData")
		blockchain.GetBlockchain().AddBlock(data)

	}
	//data := homeData{"Home", blockchain.GetBlockchain().AllBlocks()}
	//templates.ExecuteTemplate(rw, "add", data)
}

func main() {
	// update
	templates = template.Must(template.ParseGlob(templateDir + "pages/*.html"))
	templates = template.Must(templates.ParseGlob(templateDir + "partials/*.html"))

	http.HandleFunc("/", home)
	http.HandleFunc("/add", add)

	fmt.Printf("Listening on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))

	//chain := blockchain.GetBlockchain()
	//chain.AddBlock("2nd Block")
	//chain.AddBlock("3rd Block")
	//chain.AddBlock("4th Block")

	//for _, block := range chain.AllBlocks() {
	//    fmt.Printf("Data : %s\n", block.Data)
	//    fmt.Printf("Hash : %s\n", block.Hash)
	//    fmt.Printf("Prev Hash : %s\n", block.PrevHash)
	//}
}
