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
	// for _, token := range tokens {
	// 	fmt.Println(token.value)
	// }

	// parsing next -> taking flat structure and turning into tree
	// ( + 13 (- 12 1) )
	//     +
	//  13    -
	// 		12  1

	ast, _ := parse(tokens, 0)
	fmt.Print(ast.pretty())
}