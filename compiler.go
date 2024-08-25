package main

import (
	"fmt"
	"log"
	"os"
)

type TokenType int

type TokenInfo struct {
	tokType TokenType
	tokStr  string
}

// Token Types
const (
	TT_ILLEGAL TokenType = iota
	TT_EOFILE
	TT_NEW_LINE
	TT_SPACE

	// Types
	TT_IDENT  // main
	TT_INT    // 1234567890
	TT_CHAR   // 'a', 'b', 'c', 'd'
	TT_STRING // "abc"

	// Operations
	TT_ADD    // +
	TT_SUB    // -
	TT_MUL    // *
	TT_DIV    // /
	TT_REM    // %
	TT_AND    // &
	TT_OR     // |
	TT_XOR    // ^
	TT_SHL    // <<
	TT_SHR    // >>
	TT_BAND   // &&
	TT_BOR    // ||
	TT_EQL    // ==
	TT_LSS    // <
	TT_GTR    // >
	TT_ASSIGN // =
	TT_NOT    // !
	TT_NOTEQL // !=
	TT_LEQL   // <=
	TT_GEQL   // >=

	TT_LPAREN    // (
	TT_RPAREN    // )
	TT_LBRACK    // [
	TT_RBRACK    // ]
	TT_LBRACE    // {
	TT_RBRACE    // }
	TT_COMMA     // ,
	TT_PERIOD    // .
	TT_SEMICOLON // ;
	TT_COLON     // :

	// Keywords
	TT_WHILE
	TT_BREAK
	TT_CONTINUE
	TT_IF
	TT_ELSE
	TT_FUNC
	TT_RETURN
	TT_VAR
)

func GenerateToken(src []byte) (TokenInfo, int) {
	var bytesConsumed int = 0

	var curToken TokenInfo
	curToken.tokType = TT_ILLEGAL

	if len(src) == 0 {
		curToken.tokType = TT_EOFILE
		bytesConsumed = 0
		return curToken, bytesConsumed
	} else if src[0] == 0x0a { // 0x0a -> '\n'
		curToken.tokType = TT_NEW_LINE
		bytesConsumed = 1
		return curToken, bytesConsumed
	} else if src[0] == 0x20 { // 0x20 -> ' '
		curToken.tokType = TT_SPACE
		bytesConsumed = 1
		return curToken, bytesConsumed
	}

	var srcStr string

	for idx, arg := range src {
		if arg == 0x0a { // 0x0a -> '\n'
			srcStr = string(src[:idx])
			break
		}
	}

	TokensStrings := map[TokenType]string {
		TT_ADD:  		"+",
		TT_SUB:  		"-",
		TT_MUL:  		"*",
		TT_DIV:  		"/",
		TT_REM:  		"%",
		TT_AND:  		"&",
		TT_OR:  		"|",
		TT_XOR:  		"^",
		TT_SHL:  		"<<",
		TT_SHR:  		">>",
		TT_BAND:  		"&&",
		TT_BOR:  		"||",
		TT_EQL:  		"==",
		TT_LSS:  		"<",
		TT_GTR:  		">",
		TT_ASSIGN:  	"=",
		TT_NOT:  		"!",
		TT_NOTEQL:  	"!=",
		TT_LEQL:  		"<=",
		TT_GEQL:  		">=",

		TT_LPAREN: 		"(",
		TT_RPAREN: 		")",
		TT_LBRACK: 		"[",
		TT_RBRACK: 		"]",
		TT_LBRACE: 		"{",
		TT_RBRACE: 		"}",
		TT_COMMA: 		",",
		TT_PERIOD: 		".",
		TT_SEMICOLON: 	";",
		TT_COLON: 		":",

		TT_WHILE: 		"while",
		TT_BREAK:		"break",
		TT_CONTINUE:	"continue",
		TT_IF:			"if",
		TT_ELSE:		"else",
		TT_FUNC:		"func",
		TT_RETURN:		"return",
		TT_VAR:			"var",
	}

	for tokType, tokStr := range TokensStrings {
		if len(srcStr) >= len(tokStr) && srcStr[:len(tokStr)] == tokStr {
			if curToken.tokType == TT_ILLEGAL || len(curToken.tokStr) < len(tokStr) {
				curToken.tokType = tokType
				curToken.tokStr = tokStr
				bytesConsumed = len(tokStr)
			}
		}
	}

	if curToken.tokType != TT_ILLEGAL {
		return curToken, bytesConsumed
	}

	isDigit := func(c byte) bool {
		// is c between 0 and 9
		return c >= 0x30 && c <= 0x39
	}

	isAlphabet := func(c byte) bool {
		// is c A-Z? is c a-z? is c '_'?
		return (c >= 0x41 && c <= 0x5a) || (c >= 0x61 && c <= 0x7a) || (c == 0x5f)
	}

	var i int = 0

	if isAlphabet(srcStr[i]) {
		curToken.tokType = TT_IDENT
		for (i < len(srcStr)) && (isAlphabet(srcStr[i]) || isDigit(src[i])) {
			i++
		}
		curToken.tokStr = srcStr[:i]
		bytesConsumed = len(srcStr[:i])
	} else if isDigit(srcStr[i]) {
		curToken.tokType = TT_INT
		for (i < len(srcStr)) && (isAlphabet(srcStr[i]) || isDigit(src[i])) {
			i++
		}
		curToken.tokStr = srcStr[:i]
		bytesConsumed = len(srcStr[:i])
	}

	return curToken, bytesConsumed
}

func main() {
	data, err := os.ReadFile(os.Args[1])

	if err != nil {
		log.Fatal(err)
	}

	data = append(data, 0x0a)

	fmt.Println(GenerateToken((data)))
}