package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/ralpioxxcs/nocoin/explorer"
	"github.com/ralpioxxcs/nocoin/rest"
)

func usage() {
	fmt.Printf("Welcome to 노마드 코인\n")
	fmt.Printf("Plase use the following flags:\n\n")
	fmt.Printf("-port=4000:		Set the PORT of the server\n")
	fmt.Printf("-mode=rest:		Choose between 'html','rest' and 'both' (html, rest)'\n")

	os.Exit(2)
}

func Start() {
	if len(os.Args) == 1 {
		usage()
	}

	port := flag.Int("port", 4000, "Set port of the server")
	mode := flag.String("mode", "rest", "Choose between 'html' and 'rest'")

	flag.Parse()

	switch *mode {
	case "rest":
		// start rest api
		rest.Start(*port)
	case "html":
		// start html explorer
		explorer.Start(*port)
	case "both":
		// start both of html explorer, rest api
		go rest.Start(*port)
		explorer.Start(*port + 1)
	default:
		usage()
	}

	fmt.Println(*port, *mode)

}
