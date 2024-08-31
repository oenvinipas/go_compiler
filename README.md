# go_interpreter

Learning Go by creating an interpreter for a Scheme-like language

Requirements:

* Go 1.18+

To Run:

```shell
$ go mod tidy
$ go test
$ go build
$ cat examples/fib.scm
$ (func fib (a)
      (if (< a 2)
   a
   (+ (fib (- a 1)) (fib (- a 2)))))

(fib 11)
$ ./basicinterpreter examples/fib.scm
Result: 89
```
