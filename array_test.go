package validator_test

import (
	"testing"

	validator "github.com/Jamess-Lucass/validator-go"
	"github.com/Jamess-Lucass/validator-go/rule"
	"github.com/stretchr/testify/assert"
)

func TestArray_WithMessage_AfterEachValue(t *testing.T) {
	type matrix struct {
		Row [2]int `json:"row"`
	}

	m := matrix{Row: [2]int{0, 5}}
	v, _ := validator.New(&m)
	validator.Array(v, &m.Row).EachValue(rule.Int().Positive()).WithMessage("bad cell")
	result := v.Validate()
	assert.Len(t, result.Errors, 1)
	assert.Equal(t, "row[0]", result.Errors[0].Field)
	assert.Equal(t, "bad cell", result.Errors[0].Message)

	assert.Panics(t, func() { validator.Array(v, &m.Row).EachValue(nil) })
}

func TestArray_EachValue(t *testing.T) {
	type Matrix struct {
		Row [3]int `json:"row"`
	}
	m := Matrix{Row: [3]int{1, -2, 3}}
	v, _ := validator.New(&m)
	validator.Array(v, &m.Row).EachValue(rule.Int().Positive())
	res := v.Validate()
	assert.False(t, res.IsValid())
	assert.Equal(t, "row[1]", res.Errors[0].Field)
	assert.Equal(t, "must be positive", res.Errors[0].Message)
}

func TestArray_Length(t *testing.T) {
	type S struct {
		Tags [2]string `json:"tags"`
	}
	s := S{Tags: [2]string{"a", "b"}}
	v, _ := validator.New(&s)
	validator.Array(v, &s.Tags).Min(3)
	res := v.Validate()
	assert.False(t, res.IsValid())
	assert.Equal(t, "tags", res.Errors[0].Field)
	assert.Equal(t, "must have at least 3 items", res.Errors[0].Message)
}

func TestArray_CountRules(t *testing.T) {
	type S struct {
		Pair  [2]string `json:"pair"`
		Items []int     `json:"items"`
	}
	s := S{Pair: [2]string{"a", "b"}}
	v, _ := validator.New(&s)
	validator.Array(v, &s.Items).NotEmpty()
	validator.Array(v, &s.Pair).Max(1).WithMessage("too many")
	validator.Array(v, &s.Pair).Length(3, 4)
	res := v.Validate()
	assert.Len(t, res.Errors, 3)
	assert.Equal(t, "must not be empty", res.Errors[0].Message)
	assert.Equal(t, "too many", res.Errors[1].Message)
	assert.Equal(t, "must have between 3 and 4 items", res.Errors[2].Message)
}

func TestArray_NilPointerSkipped(t *testing.T) {
	type S struct {
		Row *[2]int `json:"row"`
	}
	s := S{Row: nil}

	v, _ := validator.New(&s)
	validator.Array(v, s.Row).Min(1)
	assert.True(t, v.Validate().IsValid())

	v, _ = validator.New(&s)
	validator.Array(v, s.Row).NotNil().WithName("row")
	res := v.Validate()
	assert.False(t, res.IsValid())
	assert.Equal(t, "row", res.Errors[0].Field)
	assert.Equal(t, "must not be nil", res.Errors[0].Message)
}

func TestArray_PassedByValue_NoPanic(t *testing.T) {
	// Forgetting the & validates a copy instead of panicking; the name falls
	// back to "unknown" because a value has no address to resolve.
	type S struct {
		Row [2]int `json:"row"`
	}
	s := S{Row: [2]int{1, 2}}
	v, _ := validator.New(&s)
	validator.Array(v, s.Row).Min(3)
	res := v.Validate()
	assert.False(t, res.IsValid())
	assert.Equal(t, "unknown", res.Errors[0].Field)
	assert.Equal(t, "must have at least 3 items", res.Errors[0].Message)
}

func TestSliceOf_AcceptsArrayValue(t *testing.T) {
	// In map/standalone mode, an array value is accepted, not just a slice.
	m := map[string]any{"nums": [3]int{1, 2, 3}}
	v := validator.NewObject(m)
	v.Field("nums", validator.SliceOf(rule.Int().Positive()).Min(2))
	assert.True(t, v.Validate().IsValid())

	m["nums"] = [3]int{1, -2, 3}
	v = validator.NewObject(m)
	v.Field("nums", validator.SliceOf(rule.Int().Positive()))
	res := v.Validate()
	assert.False(t, res.IsValid())
	assert.Equal(t, "nums[1]", res.Errors[0].Field)
}

func TestArray_NotNil_NilSlice(t *testing.T) {
	// Array accepts a pointer to a slice; a nil slice fails NotNil, matching
	// SliceRule rather than silently passing.
	type S struct {
		Items []int `json:"items"`
	}
	s := S{}
	v, _ := validator.New(&s)
	validator.Array(v, &s.Items).NotNil()
	res := v.Validate()
	assert.False(t, res.IsValid())
	assert.Equal(t, "items", res.Errors[0].Field)
	assert.Equal(t, "must not be nil", res.Errors[0].Message)
}
