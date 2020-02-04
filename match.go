package match

import (
	"reflect"
	"regexp"
)

type matchKey int

type matchItem struct {
	pattern interface{}
	action  interface{}
}

// PatternChecker is func for checking pattern.
type PatternChecker func(pattern interface{}, value interface{}) bool

var (
	registeredMatchers []PatternChecker
)

const (
	// ANY is the pattern which allows any value.
	ANY matchKey = 0
	// HEAD is the pattern for start element of silce.
	HEAD matchKey = 1
	// TAIL is the pattern for end element(s) of slice.
	TAIL matchKey = 2
)

// MatchItem defines a matched item value.
type MatchItem struct {
	value        interface{}
	valueAsSlice []interface{}
}

type oneOfContainer struct {
	items []interface{}
}

// OneOf defines the pattern where at least one item matches.
func OneOf(items ...interface{}) oneOfContainer {
	return oneOfContainer{items}
}

// Matcher struct
type Matcher struct {
	value      interface{}
	matchItems []matchItem
}

// Match function takes a value for matching and
func Match(val interface{}) *Matcher {
	matchItems := []matchItem{}
	return &Matcher{val, matchItems}
}

// When function adds new pattern for checking matching.
// If pattern matched with value the func will be called.
func (matcher *Matcher) When(val interface{}, fun interface{}) *Matcher {
	newMatchItem := matchItem{val, fun}
	matcher.matchItems = append(matcher.matchItems, newMatchItem)

	return matcher
}

// RegisterMatcher register custom pattern.
func RegisterMatcher(pattern PatternChecker) {
	registeredMatchers = append(registeredMatchers, pattern)
}

// Result returns the result value of matching process.
func (matcher *Matcher) Result() (bool, interface{}) {
	for _, mi := range matcher.matchItems {
		matchedItems, matched := matchValue(mi.pattern, matcher.value)
		if matched {
			miActionType := reflect.TypeOf(mi.action)
			if miActionType.Kind() == reflect.Func {
				numberOfArgs := miActionType.NumIn()
				lenMatchedItems := len(matchedItems)
				if numberOfArgs > lenMatchedItems {
					for i := lenMatchedItems; i < numberOfArgs; i++ {
						matchedItems = append(matchedItems, MatchItem{value: nil})
					}
				} else if lenMatchedItems > numberOfArgs {
					matchedItems = matchedItems[:numberOfArgs]
				}

				var params []reflect.Value
				for i := 0; i < len(matchedItems); i++ {
					params = append(params, reflect.ValueOf(matchedItems[i]))
				}

				funcRes := reflect.ValueOf(mi.action).Call(params)
				if (len(funcRes)) > 0 {
					return true, funcRes[0].Interface()
				}

				return true, nil
			}

			return true, mi.action
		}
	}

	return false, nil
}

func matchValue(pattern interface{}, value interface{}) ([]MatchItem, bool) {
	simpleTypes := []reflect.Kind{reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16,
		reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16,
		reflect.Uint32, reflect.Uint64, reflect.Uintptr, reflect.Float32, reflect.Float64,
		reflect.Complex64, reflect.Complex128,
	}

	if pattern == ANY {
		return nil, true
	}

	valueKind := reflect.TypeOf(value).Kind()
	valueIsSimpleType := containsKind(simpleTypes, valueKind)

	for _, registerMatcher := range registeredMatchers {
		if registerMatcher(pattern, value) {
			return nil, true
		}
	}

	if (valueIsSimpleType) && value == pattern {
		return nil, true
	}

	patternType := reflect.TypeOf(pattern)
	patternKind := patternType.Kind()

	if (valueKind == reflect.Slice || valueKind == reflect.Array) &&
		patternKind == reflect.Slice {

		matchedItems, isMatched := matchSlice(pattern, value)
		if isMatched {
			return matchedItems, isMatched
		}
	}

	if patternKind == reflect.Func && patternType.NumIn() == 1 &&
		matchStruct(patternType.In(0), value) {
		return nil, true
	}

	if valueKind == reflect.Map &&
		patternKind == reflect.Map &&
		matchMap(pattern, value) {

		return nil, true
	}

	if valueKind == reflect.String {
		if patternKind == reflect.String {
			if pattern == value {
				return nil, true
			}
		}

		reg, ok := pattern.(*regexp.Regexp)
		if ok {
			if matchRegexp(reg, value) {
				return nil, true
			}
		}
	}

	if valueKind == reflect.Struct {
		if patternKind == reflect.Struct {
			if value == pattern {
				return nil, true
			}
		}
	}

	return nil, false
}

func matchSlice(pattern interface{}, value interface{}) ([]MatchItem, bool) {
	patternSlice := reflect.ValueOf(pattern)
	patternSliceLen := patternSlice.Len()

	valueSlice := reflect.ValueOf(value)
	valueSliceLen := valueSlice.Len()

	if patternSliceLen > 0 && patternSlice.Index(0).Interface() == HEAD {
		if valueSliceLen == 0 {
			return nil, false
		}

		patternSliceVal := patternSlice.Slice(1, patternSliceLen)
		patternSliceLen = patternSliceVal.Len()
		patternSliceInterface := patternSliceVal.Interface()

		for i := 0; i < valueSliceLen-patternSliceLen+1; i++ {
			matchedItems, isMatched := matchSubSlice(patternSliceInterface, valueSlice.Slice(i, valueSliceLen).Interface())
			resMatchedItems := append([]MatchItem{{valueAsSlice: sliceValueToSliceOfInterfaces(valueSlice.Slice(0, i))}}, matchedItems...)
			if isMatched {
				return resMatchedItems, true
			}
		}

		return nil, false
	}

	return matchSubSlice(pattern, value)
}

func matchSubSlice(pattern interface{}, value interface{}) ([]MatchItem, bool) {
	patternSlice := reflect.ValueOf(pattern)
	valueSlice := reflect.ValueOf(value)

	patternSliceLength := patternSlice.Len()
	valueSliceLength := valueSlice.Len()

	matchedItems := make([]MatchItem, 0)

	if patternSliceLength == 0 || valueSliceLength == 0 {
		if patternSliceLength == valueSliceLength {
			return nil, true
		}
		return nil, false
	}

	patternSliceMaxIndex := patternSliceLength - 1
	valueSliceMaxIndex := valueSliceLength - 1
	oneOfContainerType := reflect.TypeOf(oneOfContainer{})

	for i := 0; i < max(patternSliceLength, valueSliceLength); i++ {
		currPatternIndex := min(i, patternSliceMaxIndex)
		currValueIndex := min(i, valueSliceMaxIndex)

		currPattern := patternSlice.Index(currPatternIndex).Interface()
		currValue := valueSlice.Index(currValueIndex).Interface()

		if currPattern == HEAD {
			panic("HEAD can only be in first position of a pattern.")
		} else if currPattern == TAIL {
			if patternSliceMaxIndex > i {
				panic("TAIL must me in last position of the pattern.")
			} else {
				matchedItems = append(matchedItems, MatchItem{valueAsSlice: sliceValueToSliceOfInterfaces(valueSlice.Slice(i, valueSliceMaxIndex+1))})
				break
			}
		} else if reflect.TypeOf(currPattern).AssignableTo(oneOfContainerType) {
			if !oneOfContainerPatternMatch(currPattern, currValue) {
				return matchedItems, false
			}
		} else if currPattern == ANY {
			matchedItems = append(matchedItems, MatchItem{value: currValue})
			continue
		} else {
			isMatched := matchValueBool(currPattern, currValue)

			if !isMatched {
				return matchedItems, false
			}
		}
	}

	return matchedItems, true
}

func sliceValueToSliceOfInterfaces(val reflect.Value) []interface{} {
	var res []interface{}
	for i := 0; i < val.Len(); i++ {
		res = append(res, val.Index(i).Interface())
	}

	return res
}

func matchStruct(patternType reflect.Type, value interface{}) bool {
	if patternType.AssignableTo(reflect.TypeOf(value)) {
		return true
	}

	return false
}

func matchMap(pattern interface{}, value interface{}) bool {
	patternMap := reflect.ValueOf(pattern)
	valueMap := reflect.ValueOf(value)

	stillUsablePatternKeys := patternMap.MapKeys()
	stillUsableValueKeys := valueMap.MapKeys()
	oneOfContainerType := reflect.TypeOf(oneOfContainer{})

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
				pValInterface := pVal.Interface()
				vValInterface := vVal.Interface()
				valueMatched := pValInterface == ANY || matchValueBool(pValInterface, vValInterface) ||
					(reflect.TypeOf(pValInterface).AssignableTo(oneOfContainerType) && oneOfContainerPatternMatch(pValInterface, vValInterface))
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

func matchValueBool(pattern interface{}, value interface{}) bool {
	_, res := matchValue(pattern, value)
	return res
}

func oneOfContainerPatternMatch(oneOfPattern interface{}, value interface{}) bool {
	oneOfContainerPatternInstance := oneOfPattern.(oneOfContainer)
	matched := false
	for _, item := range oneOfContainerPatternInstance.items {
		if matchValueBool(item, value) {
			matched = true
			break
		}
	}

	return matched
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
	valInterface := val.Interface()
	for _, v := range vals {
		if valInterface == v.Interface() {
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
