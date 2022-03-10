package main

import (
	"fmt"
	"time"

	"github.com/ralpioxxcs/n-coin/cli"
	"github.com/ralpioxxcs/n-coin/db"
)

func send(c chan<- int) /** send only **/ {
	for i := range [10]int{} {
		// time.Sleep(100 * time.Millisecond)
		fmt.Printf(">> sending %d <<\n", i)
		c <- i // will block until receive somewhere
		fmt.Printf(">> sent%d <<\n", i)
	}
	close(c)
}

func receive(c <-chan int) /** receive only **/ {
	for {
		time.Sleep(1 * time.Second)
		a, ok := <-c
		if !ok {
			fmt.Println("we're Done!")
			break
		}
		fmt.Printf("|| received %d ||\n", a)
	}
}

func main() {
	defer db.Close()
	db.InitDB()
	cli.Start()

	// //c := make(chan int)    // unbufferd channel
	// c := make(chan int, 10) // bufferd channel (queue)
	// go send(c)
	// receive(c)
}
