# Go pattern matching
===============

.. image:: https://travis-ci.org/alexpantyukhin/go-pattern-match.svg?branch=master
    :target: https://travis-ci.org/alexpantyukhin/go-pattern-match

It's just another approach for using pattern matching in Go. I have been inspired by [python pattern matching](https://github.com/santinic/pampy) that's why I wanted to make an attempt to rewrite something similar in Go :)
For now the following matching are implemented :
   - [x] Simple types (like int, int64, float, float64, bool..).
   - [x] Slices (with HEAD, TAIL patterns).
   - [x] Dictionary (with ANY pattern).
   - [x] Regexp.
   - [x] Adding custom matching (ability to add special matching for some structs for example)
   
# Usages
It's possible to try use matching Simple types:

```go
isMatched, mr := match.Match(42).
                When(42, func() interface{} { return true }).
                Result()
```

With Maps:
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
            	}, func() interface{} { return true }).
            	Result()
```

With Slices:
```go
isMatched, mr := match.Match([]int{1, 2, 3}).
            	When([]interface{}{match.HEAD, 2, 3}, func() interface{} { return true }).
            	Result()
```

With regexps:
```go
isMatched, mr := match.Match("gophergopher").
            	When("gophergopher", func() interface{} { return true }).
            	Result()
```

# Plans:
 - [ ] I would like to implement a recursive pattern matching (for matching inner elements of objects)
 - [ ] Possibly to have matching without result.

# Installation
Just `go get` this repository with the following way:

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
        When(42, func() interface{} { return false } ).
        When([]interface{}{match.HEAD, 2, 3}, func() interface{} { return true }).
        Result()


    if isMatched {
        fmt.Println(mr)
    }
}
```
