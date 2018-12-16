package match

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMatch_SimpleTypeInt(t *testing.T) {
	mr := Match(42).
		When(42, func() interface{} { return true }).
		Result()

	assert.Equal(t, true, mr)
}

func TestMatch_SimpleTypeIntMatchWithDifferentType(t *testing.T) {
	mr := Match(42).
		When(int64(42), func() interface{} { return true }).
		Result()

	assert.Equal(t, nil, mr)
}

func TestMatch_SliceWithHead(t *testing.T) {
	mr := Match([]interface{}{1, 2, 3}).
		When([]interface{}{HEAD, 2, 3}, func() interface{} { return true }).
		Result()

	assert.Equal(t, true, mr)
}

func TestMatch_SliceWithAny(t *testing.T) {
	mr := Match([]interface{}{1, 2, 3}).
		When([]interface{}{1, ANY, 3}, func() interface{} { return true }).
		Result()

	assert.Equal(t, true, mr)
}

func TestMatch_SliceWithValueIntType(t *testing.T) {
	mr := Match([]int{1, 2, 3}).
		When([]interface{}{HEAD, 2, 3}, func() interface{} { return true }).
		Result()

	assert.Equal(t, true, mr)
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
	mr := Match([]int{}).
		When([]interface{}{HEAD}, func() interface{} { return true }).
		Result()

	assert.Equal(t, nil, mr)
}

func TestMatch_SliceWithHeadNotMuch(t *testing.T) {
	mr := Match([]interface{}{1, 2, 3}).
		When([]interface{}{HEAD, 3}, func() interface{} { return true }).
		Result()

	assert.Equal(t, nil, mr)
}

func TestMatch_Map(t *testing.T) {
	mr := Match(map[string]int{
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

	assert.Equal(t, true, mr)
}

func TestMatch_MapPatternWithAny(t *testing.T) {
	mr := Match(map[string]int{
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

	assert.Equal(t, true, mr)
}

func TestMatch_MapPatternDifferentValue(t *testing.T) {
	mr := Match(map[string]int{
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

	assert.Equal(t, nil, mr)
}

func TestMatch_String(t *testing.T) {
	mr := Match("gophergopher").
		When("gophergopher", func() interface{} { return true }).
		Result()

	assert.Equal(t, true, mr)
}

func TestMatch_Regexp(t *testing.T) {
	mr := Match("gophergopher").
		When(regexp.MustCompile("(gopher){2}"), func() interface{} { return true }).
		Result()

	assert.Equal(t, true, mr)
}

func TestMatch_RegisterPattern(t *testing.T) {
	myMagicChecker := func(pattern interface{}, value interface{}) bool {

		if pattern == 12345 {
			return true
		}

		return false
	}

	RegisterMatcher(myMagicChecker)
	mr := Match(1000).
		When(12345, func() interface{} { return true }).
		When(1000, func() interface{} { return false }).
		Result()

	assert.Equal(t, true, mr)
}
