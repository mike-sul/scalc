package main

import (
	"fmt"
	"github.com/mike-sul/scalc/pkg/scalc"
	"io"
	"os"
)

func PrintUsage() {
	fmt.Println("Usage: scalc <expression>")
	fmt.Println("expression := “[“ operator set_1 set_2 set_3 … set_n “]”")
	fmt.Println("set := file | expression")
	fmt.Println("operator := “SUM” | “INT” | “DIF”")
	fmt.Println("Example: scalc [ SUM [ DIF a.txt b.txt c.txt ] [ INT b.txt c.txt ] ]")
}

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Expression hasn't been specified")
		PrintUsage()
		os.Exit(1)
	}

	resSet, pos, err := scalc.ParseExpression(os.Args, 1)

	if err != nil || pos != len(os.Args)-1 {
		fmt.Printf("Failed to parse an input expression, error position: %d: %s\n", pos, err.Error())
		PrintUsage()
		os.Exit(1)
	}

	var val  int
	for val, err = resSet.Next(); err == nil; val, err = resSet.Next() {
		fmt.Println(val)
	}

	if err != nil && err != io.EOF {
		fmt.Printf("Error while processing an input expression: %s\n", err.Error())
		os.Exit(1)
	}

	os.Exit(0)
}
