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

	assert.Equal(t, true, isMatched)
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

	assert.Equal(t, false, isMatched)
}

func TestMatch_SliceWithHead(t *testing.T) {
	isMatched, _ := Match([]interface{}{1, 2, 3}).
		When([]interface{}{HEAD, 2, 3}, true).
		Result()

	assert.Equal(t, true, isMatched)
}

func TestMatch_SliceWithAny(t *testing.T) {
	isMatched, _ := Match([]interface{}{1, 2, 3}).
		When([]interface{}{1, ANY, 3}, true).
		Result()

	assert.Equal(t, true, isMatched)
}

func TestMatch_SliceWithValueIntType(t *testing.T) {
	isMatched, _ := Match([]int{1, 2, 3}).
		When([]interface{}{HEAD, 2, 3}, true).
		Result()

	assert.Equal(t, true, isMatched)
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

	assert.Equal(t, false, isMatched)
}

func TestMatch_SliceWithHeadMoreThanOneElement(t *testing.T) {
	isMatched, _ := Match([]interface{}{1, 2, 3}).
		When([]interface{}{HEAD, 3}, true).
		Result()

	assert.Equal(t, true, isMatched)
}

func TestMatch_SliceWithOneOf(t *testing.T) {
	isMatched, _ := Match([]interface{}{1, 2, 3}).
		When([]interface{}{1, OneOf(1, 2, 3), 3}, true).
		Result()

	assert.Equal(t, true, isMatched)
}

func TestMatch_SliceWithOneOfDoesntMatch(t *testing.T) {
	isMatched, _ := Match([]interface{}{1, 2, 3}).
		When([]interface{}{1, OneOf(4, 5, 6), 3}, true).
		Result()

	assert.Equal(t, false, isMatched)
}

func TestMatch_ArrayWithAny(t *testing.T) {
	isMatched, _ := Match([3]int{1, 2, 3}).
		When([]interface{}{1, ANY, 3}, true).
		Result()

	assert.Equal(t, true, isMatched)
}

func TestMatch_SliceWithMatchedItems(t *testing.T) {
	isMatched, res := Match([]interface{}{1, 2, 3, 4, 5}).
		When([]interface{}{HEAD, 3, TAIL}, func(head MatchItem, tail MatchItem) [][]interface{} {
			return [][]interface{}{head.valueAsSlice, tail.valueAsSlice}
		}).
		Result()

	convertedRes := res.([][]interface{})

	assert.Equal(t, true, isMatched)
	assert.Equal(t, convertedRes[0][0].(int), 1)
	assert.Equal(t, convertedRes[0][1].(int), 2)

	assert.Equal(t, convertedRes[1][0].(int), 4)
	assert.Equal(t, convertedRes[1][1].(int), 5)
}

func TestMatch_SliceWithMatchedItemsWithAny(t *testing.T) {
	isMatched, res := Match([]interface{}{1, 2, 3, 4, 5}).
		When([]interface{}{HEAD, 2, ANY, 4, TAIL}, func(head MatchItem, any MatchItem, tail MatchItem) [][]interface{} {
			return [][]interface{}{head.valueAsSlice, []interface{}{any.value}, tail.valueAsSlice}
		}).
		Result()

	convertedRes := res.([][]interface{})

	assert.Equal(t, true, isMatched)
	assert.Equal(t, convertedRes[0][0].(int), 1)
	assert.Equal(t, convertedRes[1][0].(int), 3)
	assert.Equal(t, convertedRes[2][0].(int), 5)
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

	assert.Equal(t, true, isMatched)
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

	assert.Equal(t, true, isMatched)
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

	assert.Equal(t, true, isMatched)
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

	assert.Equal(t, false, isMatched)
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

	assert.Equal(t, false, isMatched)
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

	assert.Equal(t, true, isMatched)
}

func TestMatch_String(t *testing.T) {
	isMatched, _ := Match("gophergopher").
		When("gophergopher", true).
		Result()

	assert.Equal(t, true, isMatched)
}

func TestMatch_Regexp(t *testing.T) {
	isMatched, _ := Match("gophergopher").
		When(regexp.MustCompile("(gopher){2}"), true).
		Result()

	assert.Equal(t, true, isMatched)
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

	assert.Equal(t, true, isMatched)
}

type TestStruct struct {
	value int
}

func TestMatch_SimpleStructMathc(t *testing.T) {
	val := TestStruct{1}

	isMatched, _ := Match(val).
		When(TestStruct{1}, 1).
		Result()

	assert.Equal(t, true, isMatched)
}

func TestMatch_SimpleStructNotMatch(t *testing.T) {
	val := TestStruct{1}

	isMatched, _ := Match(val).
		When(TestStruct{2}, 1).
		Result()

	assert.Equal(t, false, isMatched)
}

func TestMatch_StructTypeMatch(t *testing.T) {
	val := TestStruct{1}

	isMatched, _ := Match(val).
		When(func(TestStruct) {}, 1).
		Result()

	assert.Equal(t, true, isMatched)
}

type AnotherTestStruct struct {
	value int
}

func TestMatch_StructDifferentTypeNotMatch(t *testing.T) {
	val := TestStruct{1}

	isMatched, _ := Match(val).
		When(func(AnotherTestStruct) {}, 1).
		Result()

	assert.Equal(t, false, isMatched)
}
