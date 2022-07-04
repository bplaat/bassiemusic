package main

import (
	"flag"
	"os"
	"fmt"
)

func startRemove() {
	command := flag.NewFlagSet("remove", flag.ExitOnError)
	query := os.Args[2]
	command.Parse(os.Args[3:])

	fmt.Println("Removing " + query)
	// TODO
}
