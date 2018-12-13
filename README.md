# Pattern matching
It's just another approach for using pattern matching in Go.
For now the following matching are implemented :
   - [x] Simple types (like int, int64, float, float64, bool..).
   - [x] Slices (with HEAD, TAIL patterns).
   - [x] Dictionary (with ANY pattern).
   - [ ] Regexp.
   - [ ] Adding custom matching (ability to add special matching for some structs for example)

# Installation
Just `go get` this repository with the following way:

```
go get github.com/alexpantyukhin/go-pattern-match
```

# Usage
```go
package main

import (
    "fmt"
    "github.com/alexpantyukhin/go-pattern-match"
)

func main() {
    mr := match.Match([]int{1, 2, 3}).
        When(42, func() interface{} { return false } ).
        When([]interface{}{match.HEAD, 2, 3}, func() interface{} { return true }).
        Result()

    fmt.Println(mr)
}
```
