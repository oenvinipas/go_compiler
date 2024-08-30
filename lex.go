package main

import "unicode"

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
}

// func (t token) debug(source []rune) {
	
// }

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

func lexSyntaxToken(source []rune, cursor int) (int, *token) {
	if source[cursor] == '(' || source[cursor] == ')' {
		return cursor + 1, &token {
			value: string([]rune {source[cursor]}),
			kind: syntaxToken,
			location: cursor,
		}
	}

	return cursor, nil
}


// "lexIntegerToken("foo 123", 4) => "123"
// "lexIntegerToken("foo 12 3", 4) => "12"
// "lexIntegerToken("foo 12a 3", 4) <- ignore (keeping it simple)
func lexIntegerToken(source []rune, cursor int) (int, *token) {
	// position of the first integer
	originalCursor := cursor

	var value []rune
	for cursor < len(source) {
		r := source[cursor]
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
	}
}

// lexIdentifierToken("123 ab +", 4) => "ab"
// lexIdentifierToken("123 ab123 +", 4) => "ab123"
func lexIdentifierToken(source []rune, cursor int) (int, *token) {
	// position of the first identifier
	originalCursor := cursor
	
	var value []rune

	for cursor < len(source) {
		r := source[cursor]
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
	}

}

// example: " (+ 13 2  )"
// (" (+ 13 2  )") => ["(", "+", "13", "2", ")"]
func lex(source []rune) []token {
	// runes represent the smallest valid individual thing in unicode
	var tokens []token
	var t *token
	// keeps tracks what we're looking at in the source rune array
	cursor := 0
	for cursor < len(source) {
		// eat whitespace
		cursor = eatWhitespace(source, cursor)
		if cursor == len(source) {
			break
		}

		// check for syntaxToken: if so, 'continue'
		cursor, t = lexSyntaxToken(source, cursor)
		if t != nil {
			tokens = append(tokens, *t)
			continue
		}

		// check for integerToken: if so, 'continue'
		cursor, t = lexIntegerToken(source, cursor)
		if t != nil {
			tokens = append(tokens, *t)
			continue
		}

		// check for identifierToken: if so, 'continue'
		cursor, t = lexIdentifierToken(source, cursor)
		if t != nil {
			tokens = append(tokens, *t)
			continue
		}

		//Lexed nothing
		// fmt.Println(tokens[len(tokens) - 1].debug()) // line of code
		panic("Could not lex")
	}

	return tokens
}