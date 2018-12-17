package match

import (
	"reflect"
	"regexp"
)

type MatchKey int

type matchItem struct {
	pattern interface{}
	action  func() interface{}
}

// PatternChecker is func for checking pattern
type PatternChecker func(pattern interface{}, value interface{}) bool

var (
	registeredMatchers []PatternChecker
)

const (
	ANY  MatchKey = 0
	HEAD MatchKey = 1
	TAIL MatchKey = 2
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

// RegisterMatcher register custim pattern
func RegisterMatcher(pattern PatternChecker) {
	registeredMatchers = append(registeredMatchers, pattern)
}

// Result returns the result value of matching process.
func (matcher *Matcher) Result() interface{} {
	simpleTypes := []reflect.Kind{reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16,
		reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16,
		reflect.Uint32, reflect.Uint64, reflect.Uintptr, reflect.Float32, reflect.Float64,
		reflect.Complex64, reflect.Complex128,
	}

	valueKind := reflect.TypeOf(matcher.value).Kind()
	valueIsSimpleType := containsKind(simpleTypes, valueKind)

	for _, mi := range matcher.matchItems {
		for _, registerMatcher := range registeredMatchers {
			if registerMatcher(mi.pattern, matcher.value) {
				return mi.action()
			}
		}

		if (valueIsSimpleType) && matcher.value == mi.pattern {
			return mi.action()
		}

		miKind := reflect.TypeOf(mi.pattern).Kind()

		if valueKind == reflect.Slice &&
			miKind == reflect.Slice &&
			matchSlice(mi.pattern, matcher.value) {

			return mi.action()
		}

		if valueKind == reflect.Map &&
			miKind == reflect.Map &&
			matchMap(mi.pattern, matcher.value) {

			return mi.action()
		}

		if valueKind == reflect.String {
			if miKind == reflect.String {
				if mi.pattern == matcher.value {
					return mi.action()
				}
			}

			reg, ok := mi.pattern.(*regexp.Regexp)
			if ok {
				if matchRegexp(reg, matcher.value) {
					return mi.action()
				}
			}
		}
	}

	return nil
}

// todo: implement
func matchSlice(pattern interface{}, value interface{}) bool {
	patternSlice := reflect.ValueOf(pattern)
	valueSlice := reflect.ValueOf(value)

	patternSliceLength := patternSlice.Len()
	valueSliceLength := valueSlice.Len()

	if patternSliceLength == 0 || valueSliceLength == 0 {
		if patternSliceLength == valueSliceLength {
			return true
		}
		return false
	}

	patternSliceMaxIndex := patternSliceLength - 1
	valueSliceMaxIndex := valueSliceLength - 1

	for i := 0; i < max(patternSliceLength, valueSliceLength); i++ {
		currPatternIndex := min(i, patternSliceMaxIndex)
		currValueIndex := min(i, valueSliceMaxIndex)

		currPattern := patternSlice.Index(currPatternIndex).Interface()
		currValue := valueSlice.Index(currValueIndex).Interface()

		if currPattern == HEAD {
			if i != 0 {
				panic("HEAD can only be in first position of a pattern.")
			} else {
				if i > valueSliceMaxIndex {
					return false
				}
			}
		} else if currPattern == TAIL {
			if patternSliceMaxIndex > i {
				panic("TAIL must me in last position of the pattern.")
			} else {
				break
			}
		} else {
			if currPattern != currValue && currPattern != ANY {
				return false
			}
		}
	}

	return true
}

func matchMap(pattern interface{}, value interface{}) bool {
	patternMap := reflect.ValueOf(pattern)
	valueMap := reflect.ValueOf(value)

	stillUsablePatternKeys := patternMap.MapKeys()
	stillUsableValueKeys := valueMap.MapKeys()

	for _, pKey := range patternMap.MapKeys() {
		if !containsValue(stillUsablePatternKeys, pKey) {
			continue
		}
		pVal := patternMap.MapIndex(pKey)
		matchedLeftAndRight := false

		for _, vKey := range valueMap.MapKeys() {
			if !containsValue(stillUsableValueKeys, vKey) {
				continue
			}

			if !containsValue(stillUsablePatternKeys, pKey) {
				continue
			}

			vVal := valueMap.MapIndex(vKey)
			keyMatched := pKey.Interface() == vKey.Interface()
			if keyMatched {
				valueMatched := pVal.Interface() == vVal.Interface() || pVal.Interface() == ANY
				if valueMatched {
					matchedLeftAndRight = true
					removeValue(stillUsablePatternKeys, pKey)
					removeValue(stillUsableValueKeys, vKey)
				}
			}
		}

		if !matchedLeftAndRight {
			return false
		}
	}

	return true
}

func matchRegexp(regexp *regexp.Regexp, value interface{}) bool {
	valueStr := value.(string)

	return regexp.MatchString(valueStr)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a < b {
		return b
	}
	return a
}

func removeValue(vals []reflect.Value, val reflect.Value) []reflect.Value {
	indexOf := -1
	for index, v := range vals {
		if val.Interface() == v.Interface() {
			indexOf = index
			break
		}
	}

	vals[indexOf] = vals[len(vals)-1]
	vals = vals[:len(vals)-1]

	return vals
}

func containsValue(vals []reflect.Value, val reflect.Value) bool {
	for _, v := range vals {
		if val.Interface() == v.Interface() {
			return true
		}
	}
	return false
}

func containsKind(vals []reflect.Kind, val reflect.Kind) bool {
	for _, v := range vals {
		if val == v {
			return true
		}
	}
	return false
}
