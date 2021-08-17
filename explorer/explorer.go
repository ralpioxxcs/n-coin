package explorer

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/ralpioxxcs/nocoin/blockchain"
)

const (
	port        = ":4001"
	templateDir = "explorer/templates/"
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
		templates.ExecuteTemplate(rw, "add", "Add")
	case "POST":
		r.ParseForm()
		data := r.Form.Get("blockData")
		blockchain.GetBlockchain().AddBlock(data)

		// redirect to home
		http.Redirect(rw, r, "/home", http.StatusPermanentRedirect)

	}
	//data := homeData{"Home", blockchain.GetBlockchain().AllBlocks()}
	//templates.ExecuteTemplate(rw, "add", data)
}

func Start() {
	// template.Must check if error occured
	// template.ParseGlob load template files
	templates = template.Must(template.ParseGlob(templateDir + "pages/*.html"))
	templates = template.Must(templates.ParseGlob(templateDir + "partials/*.html"))

	http.HandleFunc("/", home)
	http.HandleFunc("/add", add)

	fmt.Printf("Listening on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
