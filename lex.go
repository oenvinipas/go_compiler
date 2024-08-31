package main

import (
	"fmt"
	"os"
	"unicode"
)

type lexingContext struct {
	source []rune
	sourceFileName string
}

type tokenKind uint

const (
	// e.g. "(", ")"
	syntaxToken tokenKind = iota
	// e.g. "+", "define"
	identifierToken
	// e.g. "1", "12"
	integerToken
)

type token struct {
	value string
	kind  tokenKind
	location int
	lc lexingContext
}

func (t token) debug(description string) {
	// 1. Grab the entire line from the source code where the token is at
	// 2. Print the entire line
	// 3. Print a marker to the column where the token is at
	// 4. Print the error/debug description

	var tokenLine []rune
	var tokenLineNumber int
	var tokenColumn int
	var inTokenLine bool
	var i int

	for i < len(t.lc.source) {
		r := t.lc.source[i]

		if i < t.location {
			tokenColumn++
		}

		tokenLine = append(tokenLine, r)

		if r == '\n' {
			tokenLineNumber++
			// Got to the end of the line that the token is in.
			if inTokenLine {
				// Now outside the loop, `tokenLine`
				// will contain the entire source code
				// line where the token was. And
				// `tokenColumn` will be the column
				// number of the token.
				break
			}

			tokenColumn = 1
			tokenLine = nil
		}

		if i == t.location {
			inTokenLine = true
		}

		i++
	}

	fmt.Printf("%s [at line %d, column %d in file %s]\n", description, tokenLineNumber, tokenColumn, t.lc.sourceFileName)
	fmt.Println(string(tokenLine))

	// WILL NOT IF THERE ARE TABS OR OTHER WEIRD CHARACTERS
	for tokenColumn >= 1 {
		fmt.Printf(" ")
		tokenColumn--
	}
	fmt.Println("^ near here")
}

// "skip" the whitespace
func eatWhitespace(source []rune, cursor int) int {
	for cursor < len(source) {
		if unicode.IsSpace(source[cursor]) {
			cursor++
			continue
		}
		break
	}

	return cursor
}

func (lc lexingContext) lexSyntaxToken(cursor int) (int, *token) {
	if lc.source[cursor] == '(' || lc.source[cursor] == ')' {
		return cursor + 1, &token {
			value: string([]rune {lc.source[cursor]}),
			kind: syntaxToken,
			location: cursor,
			lc: lc,
		}
	}

	return cursor, nil
}


// "lexIntegerToken("foo 123", 4) => "123"
// "lexIntegerToken("foo 12 3", 4) => "12"
// "lexIntegerToken("foo 12a 3", 4) <- ignore (keeping it simple)
func (lc lexingContext) lexIntegerToken(cursor int) (int, *token) {
	// position of the first integer
	originalCursor := cursor

	var value []rune
	for cursor < len(lc.source) {
		r := lc.source[cursor]
		if r >= '0' && r <= '9' {
			value = append(value, r)
			cursor++
			continue
		}
		break
	}

	if len(value) == 0 {
		return originalCursor, nil
	}

	return cursor, &token {
		value: string(value),
		kind: integerToken,
		location: originalCursor,
		lc: lc,
	}
}

// lexIdentifierToken("123 ab +", 4) => "ab"
// lexIdentifierToken("123 ab123 +", 4) => "ab123"
func (lc lexingContext) lexIdentifierToken(cursor int) (int, *token) {
	// position of the first identifier
	originalCursor := cursor
	
	var value []rune

	for cursor < len(lc.source) {
		r := lc.source[cursor]
		if !(unicode.IsSpace(r) || r == ')') {
			value = append(value, r)
			cursor++
			continue
		}
		break
	}

	if len(value) == 0 {
		return originalCursor, nil
	}

	return cursor, &token {
		value: string(value),
		kind: identifierToken,
		location: originalCursor,
		lc: lc,
	}
}

// example: " (+ 13 2  )"
// (" (+ 13 2  )") => ["(", "+", "13", "2", ")"]
func (lc lexingContext) lex() []token {
	// runes represent the smallest valid individual thing in unicode
	var tokens []token
	var t *token
	// keeps tracks what we're looking at in the source rune array
	cursor := 0
	for cursor < len(lc.source) {
		// eat whitespace
		cursor = eatWhitespace(lc.source, cursor)
		if cursor == len(lc.source) {
			break
		}

		// check for syntaxToken: if so, 'continue'
		cursor, t = lc.lexSyntaxToken(cursor)
		if t != nil {
			tokens = append(tokens, *t)
			continue
		}

		// check for integerToken: if so, 'continue'
		cursor, t = lc.lexIntegerToken(cursor)
		if t != nil {
			tokens = append(tokens, *t)
			continue
		}

		// check for identifierToken: if so, 'continue'
		cursor, t = lc.lexIdentifierToken(cursor)
		if t != nil {
			tokens = append(tokens, *t)
			continue
		}

		panic("Could not lex")
	}

	return tokens
}

func newLexingContext(file string) lexingContext {
	program, err := os.ReadFile(os.Args[1])
	
	if err != nil {
		panic(err)
	}

	return lexingContext {
		source: []rune(string(program)),
		sourceFileName: file,
	}
}