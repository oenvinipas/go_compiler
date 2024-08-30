package main

import (
	"fmt"
	"strconv"
)

var builtins = map[string]func([]value, map[string]any) any{}

func copyContext (in map[string]any) map[string]any {
	out := map[string]any{}
	for key, val := range in {
		out[key] = val
	}

	return out
}

func initializeBuiltins() {
	builtins["if"] = func(args []value, ctx map[string]any) any {
		condition := astWalk2(args[0], ctx)
		then := args[1]
		_else := args[2]

		if condition.(bool) {
			return astWalk2(then, ctx)
		}

		return astWalk2(_else, ctx)
	}
	builtins["+"] = func(args []value, ctx map[string]any) any {
		var i int64
		for _, arg := range args {
			i += astWalk2(arg, ctx).(int64)
		}

		return i
	}
	builtins["-"] = func(args []value, ctx map[string]any) any {
		i := astWalk2(args[0], ctx).(int64)
		for _, arg := range args[1:] {
			i -= astWalk2(arg, ctx).(int64)
		}

		return i
	}
	builtins["begin"] = func(args []value, ctx map[string]any) any {
		var last any
		for _, arg := range args {
			last = astWalk2(arg, ctx)
		}

		return last
	}
	builtins["func"] = func(args []value, ctx map[string]any) any {
		functionName := (*args[0].literal).value
		params := *args[1].list

		body := *args[2].list

		ctx[functionName] = func(args []any, ctx map[string]any) any {
			childCtx := copyContext(ctx)
			if len(params) != len(args) {
				panic(fmt.Sprintf("Expected %d args to `%s`, got %d", len(params), functionName, len(args)))
			}
			for i, param := range params {
				childCtx[(*param.literal).value] = args[i]
			}
			return astWalk(body, childCtx)
		}
		return ctx[functionName]
	}
}

/*
Example of evaluation:
( + 13 ( - 12 1 ) )
	+ 		13
    	-
	12		1
*/

/* 

	Example file:
	( + 12 2 )
	( - 11 1 )
	
	That file is transformed into:
	(begin
		( + 12 2 )
		( - 11 1 )
	)
*/
func astWalk(ast []value, ctx map[string]any) any {

	// Example `if` `+`
	functionName := (*ast[0].literal).value

	if builtinFunction, ok := builtins[functionName]; ok {
		return builtinFunction(ast[1:], ctx)
	}

	

	userDefinedFunction := ctx[functionName].(func([]any, map[string]any) any)
	var args []any
	for _, unevaluatedArg := range ast[1:] {
		args = append(args, astWalk2(unevaluatedArg, ctx))
	}
	return userDefinedFunction(args, ctx)
}

func astWalk2(v value, ctx map[string]any) any {
	if v.kind == literalValue {
		token := *v.literal
		switch token.kind {
			// `12` `1`
			case integerToken:
				i, err := strconv.ParseInt(token.value, 10, 64)
				if err != nil {
					fmt.Println("Expected an integer, got: " + token.value)
					panic(err)
				}

				return i

			// `+` `-`
			case identifierToken:
				return ctx[token.value]
		}
	}

	return astWalk(*v.list, ctx)
}