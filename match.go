package match

import (
	"reflect"
)

type MatchKey int

type matchItem struct {
	value  interface{}
	action func() interface{}
}

const (
	ANY  MatchKey = 0
	TAIL MatchKey = 1
)

// Matcher struct
type Matcher struct {
	value      interface{}
	matchItems []matchItem
}

// Match func
func Match(val interface{}) *Matcher {
	matchItems := []matchItem{}
	return &Matcher{val, matchItems}
}

// When func
func (matcher *Matcher) When(val interface{}, fun func() interface{}) *Matcher {
	newMatchItem := matchItem{val, fun}
	matcher.matchItems = append(matcher.matchItems, newMatchItem)

	return matcher
}

// Result returns the result value of matching process.
func (matcher *Matcher) Result() interface{} {
	simpleTypes := []reflect.Kind{reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16,
		reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16,
		reflect.Uint32, reflect.Uint64, reflect.Uintptr, reflect.Float32, reflect.Float64,
		reflect.Complex64, reflect.Complex128,
	}

	valueKind := reflect.TypeOf(matcher.value).Kind()
	valueIsSimpleType := contains(simpleTypes, valueKind)

	for _, mi := range matcher.matchItems {
		if (valueIsSimpleType) && matcher.value == mi.value {
			return mi.action()
		}

	}

	return nil
}

// todo: implement
func matchSlice(pattern interface{}, value interface{}) bool {

	return true
}

func contains(vals []reflect.Kind, val reflect.Kind) bool {
	for _, v := range vals {
		if val == v {
			return true
		}
	}
	return false
}
