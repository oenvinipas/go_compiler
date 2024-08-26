package main

import (
	"fmt"
	"os"
)

func main() {
	// accept input
	program, err := os.ReadFile(os.Args[1])
	
	if err != nil {
		panic(err)
	}

	tokens := lex([]rune(string(program)))
	fmt.Println(tokens)
}