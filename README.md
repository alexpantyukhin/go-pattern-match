# Go pattern matching
[![Mentioned in Awesome Go](https://awesome.re/mentioned-badge.svg)](https://github.com/avelino/awesome-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/alexpantyukhin/go-pattern-match)](https://goreportcard.com/report/github.com/alexpantyukhin/go-pattern-match)
[![Build Status](https://travis-ci.org/alexpantyukhin/go-pattern-match.svg?branch=master
)](https://travis-ci.org/alexpantyukhin/go-pattern-match)
[![codecov](https://codecov.io/gh/alexpantyukhin/go-pattern-match/branch/master/graph/badge.svg)](https://codecov.io/gh/alexpantyukhin/go-pattern-match)
[![GoDoc](https://godoc.org/alexpantyukhin/go-pattern-match?status.svg)](https://godoc.org/github.com/alexpantyukhin/go-pattern-match)
[![LICENSE](https://img.shields.io/github/license/alexpantyukhin/go-pattern-match.svg)](https://github.com/alexpantyukhin/go-pattern-match/blob/master/LICENSE)

It's just another implementation of pattern matching in Go. I have been inspired by [Python pattern matching](https://github.com/santinic/pampy), that's why I wanted to try writing something similar in Go :)
For now the following matching are implemented :
   - [x] Simple types (like int, int64, float, float64, bool..).
   - [x] Struct type.
   - [x] Slices (with HEAD, TAIL, OneOf patterns).
   - [x] Dictionary (with ANY, OneOf pattern).
   - [x] Regexp.
   - [x] Additional custom matching (ability to add special matching for some, structs for example).
   
# Usages

## Fibonacci example:

```go
func fib(n int) int {
	_, res := match.Match(n).
		When(1, 1).
		When(2, 1).
		When(match.ANY, func() int { return fib(n-1) + fib(n-2) }).
		Result()

	return res.(int)
}
```

## Simple types:

```go
isMatched, mr := match.Match(42).
                When(42, 10).
                Result()
// isMatched - true, mr - 10
```

## With Structs:
- Simple check value by type
```go
val := TestStruct{1}

isMatched, _ := Match(val).
    When(func(TestStruct) {},  1).
    Result()
```

- Check value by type and condition
```go
val := TestStruct{1}

isMatched, _ := Match(val).
	When(func(ts TestStruct) bool { return ts.value == 42 }, 1).
	When(func(ts AnotherStruct) bool { return ts.stringValue == "hello" }, 2).
	Result()
```

## With Maps:
```go
isMatched, mr := match.Match(map[string]int{
                	"rsc": 3711,
                	"r":   2138,
                	"gri": 1908,
                	"adg": 912,
                }).
        	    When(map[string]interface{}{
                	"rsc": 3711,
                	"r":   2138,
                	"gri": 1908,
                	"adg": match.ANY,
            	}, true).
            	Result()
```

## With Slices:
```go
isMatched, mr := match.Match([]int{1, 2, 3, 4, 5, 6}).
            	When([]interface{}{match.HEAD, 3, match.OneOf(3, 4), 5, 6}, 125).
            	Result()
```

## With regexps:
```go
isMatched, mr := match.Match("gophergopher").
            	When("gophergopher", func() interface{} { return true }).
            	Result()
```

## Without result:
```go
func main() {
	Match(val).
	When(42, func() { fmt.Println("You found the answer to life, universe and everything!") }).
	When(ANY, func() { fmt.Println("No.. It's not an answer.") }).
	Result()
}
```

# Installation
Just `go get` this repository in the following way:

```
go get github.com/alexpantyukhin/go-pattern-match
```

# Full example
```go
package main

import (
    "fmt"
    "github.com/alexpantyukhin/go-pattern-match"
)

func main() {
    isMatched, mr := match.Match([]int{1, 2, 3}).
        When(42, false).
        When([]interface{}{match.HEAD, 2, 3}, true).
        Result()


    if isMatched {
        fmt.Println(mr)
    }
}
```
