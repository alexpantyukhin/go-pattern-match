package match

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMatch_SimpleTypeInt(t *testing.T) {
	isMatched, _ := Match(42).
		When(42, true).
		Result()

	assert.True(t, isMatched)
}

func TestMatch_AnyPattern(t *testing.T) {
	_, res := Match(5).
		When(1, 1).
		When(2, 2).
		When(ANY, 10).
		Result()

	assert.Equal(t, 10, res)
}

func fib(n int) int {
	_, res := Match(n).
		When(1, 1).
		When(2, 1).
		When(ANY, func() int { return fib(n-1) + fib(n-2) }).
		Result()

	return res.(int)
}

func TestMatch_Fibonacci(t *testing.T) {
	assert.Equal(t, 21, fib(8))
}

func TestMatch_MatchWithoutResult(t *testing.T) {
	Match(10).
		When(10, func() { assert.True(t, true) }).
		When(ANY, func() { assert.False(t, true) }).
		Result()
}

func TestMatch_SimpleTypeIntWithFunc(t *testing.T) {
	_, res := Match(42).
		When(42, func() interface{} { return 84 }).
		Result()

	assert.Equal(t, 84, res)
}

func TestMatch_SimpleTypeIntMatchWithDifferentType(t *testing.T) {
	isMatched, _ := Match(42).
		When(int64(42), true).
		Result()

	assert.False(t, isMatched)
}

func TestMatch_SliceWithHead(t *testing.T) {
	isMatched, _ := Match([]interface{}{1, 2, 3}).
		When([]interface{}{HEAD, 2, 3}, true).
		Result()

	assert.True(t, isMatched)
}

func TestMatch_SliceWithAny(t *testing.T) {
	isMatched, _ := Match([]interface{}{1, 2, 3}).
		When([]interface{}{1, ANY, 3}, true).
		Result()

	assert.True(t, isMatched)
}

func TestMatch_SliceWithValueIntType(t *testing.T) {
	isMatched, _ := Match([]int{1, 2, 3}).
		When([]interface{}{HEAD, 2, 3}, true).
		Result()

	assert.True(t, isMatched)
}

func TestMatch_SlicePanicsWhenHeadNotFirst(t *testing.T) {
	mr := Match([]int{1, 2, 3}).
		When([]interface{}{1, HEAD, 3}, true)

	assert.Panics(t, func() { mr.Result() })
}

func TestMatch_SlicePanicsWhenTailNotLast(t *testing.T) {
	mr := Match([]int{1, 2, 3}).
		When([]interface{}{1, TAIL, 3}, true)

	assert.Panics(t, func() { mr.Result() })
}

func TestMatch_SliceHeadNotMatched(t *testing.T) {
	isMatched, _ := Match([]int{}).
		When([]interface{}{HEAD}, true).
		Result()

	assert.False(t, isMatched)
}

func TestMatch_SliceWithHeadMoreThanOneElement(t *testing.T) {
	isMatched, _ := Match([]interface{}{1, 2, 3}).
		When([]interface{}{HEAD, 3}, true).
		Result()

	assert.True(t, isMatched)
}

func TestMatch_SliceWithOneOf(t *testing.T) {
	isMatched, _ := Match([]interface{}{1, 2, 3}).
		When([]interface{}{1, OneOf(1, 2, 3), 3}, true).
		Result()

	assert.True(t, isMatched)
}

func TestMatch_SliceWithOneOfDoesntMatch(t *testing.T) {
	isMatched, _ := Match([]interface{}{1, 2, 3}).
		When([]interface{}{1, OneOf(4, 5, 6), 3}, true).
		Result()

	assert.False(t, isMatched)
}

func TestMatch_ArrayWithAny(t *testing.T) {
	isMatched, _ := Match([3]int{1, 2, 3}).
		When([]interface{}{1, ANY, 3}, true).
		Result()

	assert.True(t, isMatched)
}

func TestMatch_SliceWithMatchedItems(t *testing.T) {
	isMatched, res := Match([]interface{}{1, 2, 3, 4, 5}).
		When([]interface{}{HEAD, 3, TAIL}, func(head MatchItem, tail MatchItem) [][]interface{} {
			return [][]interface{}{head.valueAsSlice, tail.valueAsSlice}
		}).
		Result()

	convertedRes := res.([][]interface{})

	assert := assert.New(t)

	assert.True(isMatched)
	assert.Equal(1, convertedRes[0][0].(int))
	assert.Equal(2, convertedRes[0][1].(int))

	assert.Equal(4, convertedRes[1][0].(int))
	assert.Equal(5, convertedRes[1][1].(int))
}

func TestMatch_SliceNotMatchWithHeadAndWrongPatternLater(t *testing.T) {
	isMatched, _ := Match([]interface{}{1, 2, 3, 4, 5}).
		When([]interface{}{HEAD, 10, 11}, true).
		Result()

	assert.False(t, isMatched)
}

func TestMatch_SliceWithMatchedItemsWithAny(t *testing.T) {
	isMatched, res := Match([]interface{}{1, 2, 3, 4, 5}).
		When([]interface{}{HEAD, 2, ANY, 4, TAIL}, func(head MatchItem, any MatchItem, tail MatchItem) [][]interface{} {
			return [][]interface{}{head.valueAsSlice, {any.value}, tail.valueAsSlice}
		}).
		Result()

	convertedRes := res.([][]interface{})

	assert := assert.New(t)

	assert.True(isMatched)
	assert.Equal(1, convertedRes[0][0].(int))
	assert.Equal(3, convertedRes[1][0].(int))
	assert.Equal(5, convertedRes[2][0].(int))
}

func TestMatch_SliceWithMatchedItemsWithLessPatternsParamteresInAction(t *testing.T) {
	isMatched, res := Match([]interface{}{1, 2, 3, 4, 5}).
		When([]interface{}{HEAD, 2, ANY, 4, TAIL}, func(head MatchItem) [][]interface{} {
			return [][]interface{}{head.valueAsSlice}
		}).
		Result()

	convertedRes := res.([][]interface{})

	assert := assert.New(t)

	assert.True(isMatched)
	assert.Equal(1, len(convertedRes))
	assert.Equal(1, convertedRes[0][0].(int))
}

func TestMatch_SliceWithMatchedItemsWithGreaterPatternsParamteresInAction(t *testing.T) {
	isMatched, res := Match([]interface{}{1, 2, 3, 4, 5}).
		When([]interface{}{HEAD, 5}, func(head MatchItem, any MatchItem) [][]interface{} {
			return [][]interface{}{head.valueAsSlice}
		}).
		Result()

	convertedRes := res.([][]interface{})

	assert := assert.New(t)

	assert.True(isMatched)
	assert.Equal(1, len(convertedRes))
	assert.Equal(1, convertedRes[0][0].(int))
}

func TestMatch_SlicePatternAndValueAreEmpty(t *testing.T) {
	isMatched, _ := Match([]interface{}{}).
		When([]interface{}{}, true).
		Result()

	assert.True(t, isMatched)
}

func TestMatch_SlicePatternEmptyAndValueNotEmpty(t *testing.T) {
	isMatched, _ := Match([]interface{}{1, 2, 3, 4, 5}).
		When([]interface{}{}, true).
		Result()

	assert.False(t, isMatched)
}

func TestMatch_SlicePatternNotEmptyAndValueEmpty(t *testing.T) {
	isMatched, _ := Match([]interface{}{}).
		When([]interface{}{1, 2, 3, 4, 5}, true).
		Result()

	assert.False(t, isMatched)
}

func TestMatch_Map(t *testing.T) {
	isMatched, _ := Match(map[string]int{
		"rsc": 3711,
		"r":   2138,
		"gri": 1908,
		"adg": 912,
	}).
		When(map[string]int{
			"rsc": 3711,
			"r":   2138,
			"gri": 1908,
			"adg": 912,
		}, true).
		Result()

	assert.True(t, isMatched)
}

func TestMatch_MapPatternWithAny(t *testing.T) {
	isMatched, _ := Match(map[string]int{
		"rsc": 3711,
		"r":   2138,
		"gri": 1908,
		"adg": 912,
	}).
		When(map[string]interface{}{
			"rsc": 3711,
			"r":   2138,
			"gri": 1908,
			"adg": ANY,
		}, true).
		Result()

	assert.True(t, isMatched)
}

func TestMatch_MapPatternWithOneOf(t *testing.T) {
	isMatched, _ := Match(map[string]int{
		"rsc": 3711,
		"r":   2138,
		"gri": 1908,
		"adg": 912,
	}).
		When(map[string]interface{}{
			"rsc": 3711,
			"r":   2138,
			"gri": 1908,
			"adg": OneOf(111, 912),
		}, true).
		Result()

	assert.True(t, isMatched)
}

func TestMatch_MapPatternWithOneOfNotMatch(t *testing.T) {
	isMatched, _ := Match(map[string]int{
		"rsc": 3711,
		"r":   2138,
		"gri": 1908,
		"adg": 912,
	}).
		When(map[string]interface{}{
			"rsc": 3711,
			"r":   2138,
			"gri": 1908,
			"adg": OneOf(111, 913),
		}, true).
		Result()

	assert.False(t, isMatched)
}

func TestMatch_MapPatternDifferentValue(t *testing.T) {
	isMatched, _ := Match(map[string]int{
		"rsc": 3711,
		"r":   2138,
		"gri": 1908,
		"adg": 912,
	}).
		When(map[string]interface{}{
			"rsc": 3711,
			"r":   2138,
			"gri": 1908,
			"adg": 1,
		}, true).
		Result()

	assert.False(t, isMatched)
}

func TestMatch_MapWithInnerSlice(t *testing.T) {
	isMatched, _ := Match(map[string]interface{}{
		"rsc": 3711,
		"r":   2138,
		"gri": 1908,
		"adg": []int{1, 2, 3},
	}).
		When(map[string]interface{}{
			"rsc": 3711,
			"r":   2138,
			"gri": 1908,
			"adg": []interface{}{HEAD, 2, TAIL},
		}, true).
		Result()

	assert.True(t, isMatched)
}

func TestMatch_String(t *testing.T) {
	isMatched, _ := Match("gophergopher").
		When("gophergopher", true).
		Result()

	assert.True(t, isMatched)
}

func TestMatch_Regexp(t *testing.T) {
	isMatched, _ := Match("gophergopher").
		When(regexp.MustCompile("(gopher){2}"), true).
		Result()

	assert.True(t, isMatched)
}

func TestMatch_RegisterPattern(t *testing.T) {
	myMagicChecker := func(pattern interface{}, value interface{}) bool {

		if pattern == 12345 {
			return true
		}

		return false
	}

	RegisterMatcher(myMagicChecker)
	isMatched, _ := Match(1000).
		When(12345, true).
		When(1000, false).
		Result()

	assert.True(t, isMatched)
}

type TestStruct struct {
	value int
}

func TestMatch_SimpleStructMatch(t *testing.T) {
	val := TestStruct{1}

	isMatched, _ := Match(val).
		When(TestStruct{1}, 1).
		Result()

	assert.True(t, isMatched)
}

func TestMatch_SimpleStructNotMatch(t *testing.T) {
	val := TestStruct{1}

	isMatched, _ := Match(val).
		When(TestStruct{2}, 1).
		Result()

	assert.False(t, isMatched)
}

func TestMatch_StructTypeMatch(t *testing.T) {
	val := TestStruct{1}

	isMatched, _ := Match(val).
		When(func(TestStruct) {}, 1).
		Result()

	assert.True(t, isMatched)
}

func TestMatch_StructFuncMatch(t *testing.T) {
	val := TestStruct{1}

	isMatched, _ := Match(val).
		When(func(ts TestStruct) bool { return ts.value == 1 }, 1).
		Result()

	assert.True(t, isMatched)
}

func TestMatch_StructFuncNotMatch(t *testing.T) {
	val := TestStruct{1}

	isMatched, _ := Match(val).
		When(func(ts TestStruct) bool { return ts.value == 2 }, 1).
		Result()

	assert.False(t, isMatched)
}

type AnotherTestStruct struct {
	value int
}

func TestMatch_StructDifferentTypeNotMatch(t *testing.T) {
	val := TestStruct{1}

	isMatched, _ := Match(val).
		When(func(AnotherTestStruct) {}, 1).
		Result()

	assert.False(t, isMatched)
}
