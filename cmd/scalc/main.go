package main

import (
	"fmt"
	"github.com/mike-sul/scalc/pkg/scalc"
	"log"
	"os"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Expression hasn't been specified")
		os.Exit(1)
	}
	ds, pos, err := scalc.ParseExpression(os.Args, 1)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	if pos != len(os.Args)-1 {
		fmt.Println("Failed to parse an input expression")
		os.Exit(1)
	}

	scalc.DumpDataStrem(ds)
	os.Exit(0)
}
