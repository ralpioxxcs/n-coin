package main

import (
	"github.com/ralpioxxcs/n-coin/cli"
	"github.com/ralpioxxcs/n-coin/db"
)

func main() {
	defer db.Close() // execute when main() exited
	cli.Start()
}
