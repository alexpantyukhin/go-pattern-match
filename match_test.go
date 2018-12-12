package match

import (
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