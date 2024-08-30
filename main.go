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
	debug := false
	if debug {
		for _, token := range tokens {
			fmt.Println(token.value)
		}
	}

	// parsing next -> taking flat structure and turning into tree
	// ( + 13 (- 12 1) )
	//     +
	//  13    -
	// 		12  1
	var parseIndex int
	var a = ast {
		value {
			kind: literalValue,
			literal: &token {
				value: "begin",
				kind: identifierToken,
			},
		},
	}
	for parseIndex < len(tokens) {
		childAst, nextIndex := parse(tokens, parseIndex)
		a = append(a, value {
				kind: listValue,
				list: &childAst,
			},
		)
		parseIndex = nextIndex
	}

	fmt.Println(a.pretty())

	initializeBuiltins()
	ctx := map[string]any{}
	value := astWalk(a, ctx)
	fmt.Println("Result:", value)
}