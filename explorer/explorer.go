package explorer

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/ralpioxxcs/n-coin/blockchain"
)

const (
	templateDir = "explorer/templates/"
)

var templates *template.Template

type homeData struct {
	PageTitle string
	Blocks    []*blockchain.Block
}

func home(rw http.ResponseWriter, r *http.Request) {
	data := homeData{"Home", nil}
	templates.ExecuteTemplate(rw, "home", data)
}

func add(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		templates.ExecuteTemplate(rw, "add", "Add")
	case "POST":
		r.ParseForm()
		blockchain.Blockchain().AddBlock()

		// redirect to home
		http.Redirect(rw, r, "/home", http.StatusPermanentRedirect)

	}
	//data := homeData{"Home", blockchain.GetBlockchain().AllBlocks()}
	//templates.ExecuteTemplate(rw, "add", data)
}

func Start(port int) {
	handler := http.NewServeMux()

	// template.Must check if error occured
	// template.ParseGlob load template files
	templates = template.Must(template.ParseGlob(templateDir + "pages/*.html"))
	templates = template.Must(templates.ParseGlob(templateDir + "partials/*.html"))

	handler.HandleFunc("/", home)
	handler.HandleFunc("/add", add)

	fmt.Printf("Listening on http://localhost:%d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), handler))
}
