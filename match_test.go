package match

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMatch_SimpleTypeInt(t *testing.T) {
	isMatched, _ := Match(42).
		When(42, func() interface{} { return true }).
		Result()

	assert.Equal(t, true, isMatched)
}

func TestMatch_SimpleTypeIntMatchWithDifferentType(t *testing.T) {
	isMatched, _ := Match(42).
		When(int64(42), func() interface{} { return true }).
		Result()

	assert.Equal(t, false, isMatched)
}

func TestMatch_SliceWithHead(t *testing.T) {
	isMatched, _ := Match([]interface{}{1, 2, 3}).
		When([]interface{}{HEAD, 2, 3}, func() interface{} { return true }).
		Result()

	assert.Equal(t, true, isMatched)
}

func TestMatch_SliceWithAny(t *testing.T) {
	isMatched, _ := Match([]interface{}{1, 2, 3}).
		When([]interface{}{1, ANY, 3}, func() interface{} { return true }).
		Result()

	assert.Equal(t, true, isMatched)
}

func TestMatch_SliceWithValueIntType(t *testing.T) {
	isMatched, _ := Match([]int{1, 2, 3}).
		When([]interface{}{HEAD, 2, 3}, func() interface{} { return true }).
		Result()

	assert.Equal(t, true, isMatched)
}

func TestMatch_SlicePanicsWhenHeadNotFirst(t *testing.T) {
	mr := Match([]int{1, 2, 3}).
		When([]interface{}{1, HEAD, 3}, func() interface{} { return true })

	assert.Panics(t, func() { mr.Result() })
}

func TestMatch_SlicePanicsWhenTailNotLast(t *testing.T) {
	mr := Match([]int{1, 2, 3}).
		When([]interface{}{1, TAIL, 3}, func() interface{} { return true })

	assert.Panics(t, func() { mr.Result() })
}

func TestMatch_SliceHeadNotMatched(t *testing.T) {
	isMatched, _ := Match([]int{}).
		When([]interface{}{HEAD}, func() interface{} { return true }).
		Result()

	assert.Equal(t, false, isMatched)
}

func TestMatch_SliceWithHeadMoreThanOneElement(t *testing.T) {
	isMatched, _ := Match([]interface{}{1, 2, 3}).
		When([]interface{}{HEAD, 3}, func() interface{} { return true }).
		Result()

	assert.Equal(t, true, isMatched)
}

func TestMatch_SliceWithOneOf(t *testing.T) {
	isMatched, _ := Match([]interface{}{1, 2, 3}).
		When([]interface{}{1, OneOf(1, 2, 3), 3}, func() interface{} { return true }).
		Result()

	assert.Equal(t, true, isMatched)
}

func TestMatch_SliceWithOneOfDoesntMatch(t *testing.T) {
	isMatched, _ := Match([]interface{}{1, 2, 3}).
		When([]interface{}{1, OneOf(4, 5, 6), 3}, func() interface{} { return true }).
		Result()

	assert.Equal(t, false, isMatched)
}

func TestMatch_ArrayWithAny(t *testing.T) {
	isMatched, _ := Match([3]int{1, 2, 3}).
		When([]interface{}{1, ANY, 3}, func() interface{} { return true }).
		Result()

	assert.Equal(t, true, isMatched)
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
		}, func() interface{} { return true }).
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
		}, func() interface{} { return true }).
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
		}, func() interface{} { return true }).
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
		}, func() interface{} { return true }).
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
		}, func() interface{} { return true }).
		Result()

	assert.Equal(t, false, isMatched)
}

func Match_MapWithInnerSlice(t *testing.T) {
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
		}, func() interface{} { return true }).
		Result()

	assert.Equal(t, true, isMatched)
}

func TestMatch_String(t *testing.T) {
	isMatched, _ := Match("gophergopher").
		When("gophergopher", func() interface{} { return true }).
		Result()

	assert.Equal(t, true, isMatched)
}

func TestMatch_Regexp(t *testing.T) {
	isMatched, _ := Match("gophergopher").
		When(regexp.MustCompile("(gopher){2}"), func() interface{} { return true }).
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
		When(12345, func() interface{} { return true }).
		When(1000, func() interface{} { return false }).
		Result()

	assert.Equal(t, true, isMatched)
}

type TestStruct struct {
	value int
}

func TestMatch_SimpleStructMathc(t *testing.T) {
	val := TestStruct{1}

	isMatched, _ := Match(val).
		When(TestStruct{1}, func() interface{} { return 1 }).
		Result()

	assert.Equal(t, true, isMatched)
}

func TestMatch_SimpleStructNotMatch(t *testing.T) {
	val := TestStruct{1}

	isMatched, _ := Match(val).
		When(TestStruct{2}, func() interface{} { return 1 }).
		Result()

	assert.Equal(t, false, isMatched)
}

func TestMatch_StructTypeMatch(t *testing.T) {
	val := TestStruct{1}

	isMatched, _ := Match(val).
		When(func(TestStruct) {}, func() interface{} { return 1 }).
		Result()

	assert.Equal(t, true, isMatched)
}

type AnotherTestStruct struct {
	value int
}

func TestMatch_StructDifferentTypeNotMatch(t *testing.T) {
	val := TestStruct{1}

	isMatched, _ := Match(val).
		When(func(AnotherTestStruct) {}, func() interface{} { return 1 }).
		Result()

	assert.Equal(t, false, isMatched)
}
