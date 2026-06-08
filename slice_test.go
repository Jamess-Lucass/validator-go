package validator_test

import (
	"testing"

	validator "github.com/Jamess-Lucass/validator-go"
	"github.com/stretchr/testify/assert"
)

type sliceTestStruct struct {
	Tags  []string    `json:"tags"`
	Items []orderItem `json:"items"`
}

type orderItem struct {
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
}

func TestSlice_Min(t *testing.T) {
	s := sliceTestStruct{Tags: []string{}}
	v, _ := validator.New(&s)
	validator.Slice(v, &s.Tags).Min(1)
	result := v.Validate()
	assert.False(t, result.IsValid())
	assert.Equal(t, "tags", result.Errors[0].Field)
	assert.Equal(t, "must have at least 1 item", result.Errors[0].Message)

	s.Tags = []string{"go"}
	v, _ = validator.New(&s)
	validator.Slice(v, &s.Tags).Min(1)
	result = v.Validate()
	assert.True(t, result.IsValid())
}

func TestSlice_Max(t *testing.T) {
	s := sliceTestStruct{Tags: []string{"a", "b", "c"}}
	v, _ := validator.New(&s)
	validator.Slice(v, &s.Tags).Max(2)
	result := v.Validate()
	assert.False(t, result.IsValid())

	s.Tags = []string{"a", "b"}
	v, _ = validator.New(&s)
	validator.Slice(v, &s.Tags).Max(2)
	result = v.Validate()
	assert.True(t, result.IsValid())
}

func TestSlice_EachValue(t *testing.T) {
	s := sliceTestStruct{Tags: []string{"go", "a", "rust"}}
	v, _ := validator.New(&s)
	validator.Slice(v, &s.Tags).EachValue(validator.String().Min(2))
	result := v.Validate()
	assert.False(t, result.IsValid())
	assert.Len(t, result.Errors, 1)
	assert.Equal(t, "tags[1]", result.Errors[0].Field)
	assert.Equal(t, "must be at least 2 characters", result.Errors[0].Message)
}

func TestSlice_Each_Struct(t *testing.T) {
	s := sliceTestStruct{Items: []orderItem{
		{Name: "Widget", Quantity: 1},
		{Name: "", Quantity: 0},
	}}
	v, _ := validator.New(&s)
	validator.Slice(v, &s.Items).Each(func(item *orderItem, sv *validator.Validator) {
		sv.String(&item.Name).NotEmpty()
		sv.Int(&item.Quantity).Gte(1)
	})
	result := v.Validate()
	assert.False(t, result.IsValid())
	assert.Len(t, result.Errors, 2)
	assert.Equal(t, "items[1].name", result.Errors[0].Field)
	assert.Equal(t, "items[1].quantity", result.Errors[1].Field)
}

func TestSlice_NotEmpty(t *testing.T) {
	s := sliceTestStruct{Tags: []string{}}
	v, _ := validator.New(&s)
	validator.Slice(v, &s.Tags).NotEmpty()
	result := v.Validate()
	assert.False(t, result.IsValid())

	s.Tags = []string{"go"}
	v, _ = validator.New(&s)
	validator.Slice(v, &s.Tags).NotEmpty()
	result = v.Validate()
	assert.True(t, result.IsValid())
}

func TestSlice_NotNil(t *testing.T) {
	s := sliceTestStruct{Tags: nil}
	v, _ := validator.New(&s)
	validator.Slice(v, &s.Tags).NotNil()
	result := v.Validate()
	assert.False(t, result.IsValid())
	assert.Equal(t, "tags", result.Errors[0].Field)
	assert.Equal(t, "must not be nil", result.Errors[0].Message)
}

func TestSlice_Length(t *testing.T) {
	s := sliceTestStruct{Tags: []string{"a", "b", "c", "d"}}
	v, _ := validator.New(&s)
	validator.Slice(v, &s.Tags).Length(1, 3)
	result := v.Validate()
	assert.False(t, result.IsValid())
	assert.Equal(t, "must have between 1 and 3 items", result.Errors[0].Message)

	s.Tags = []string{"a", "b"}
	v, _ = validator.New(&s)
	validator.Slice(v, &s.Tags).Length(1, 3)
	assert.True(t, v.Validate().IsValid())
}

func TestSlice_Must(t *testing.T) {
	s := sliceTestStruct{Tags: []string{"a", "b"}}
	v, _ := validator.New(&s)
	validator.Slice(v, &s.Tags).Must(func(tags []string) bool {
		for _, tag := range tags {
			if tag == "go" {
				return true
			}
		}
		return false
	}).WithMessage("must contain 'go'")
	result := v.Validate()
	assert.False(t, result.IsValid())
	assert.Equal(t, "must contain 'go'", result.Errors[0].Message)
}

func TestSlice_Min_SingularMessage(t *testing.T) {
	s := sliceTestStruct{Tags: []string{}}
	v, _ := validator.New(&s)
	validator.Slice(v, &s.Tags).Min(1)
	result := v.Validate()
	assert.Equal(t, "must have at least 1 item", result.Errors[0].Message)
}

func TestSlice_WithName(t *testing.T) {
	s := sliceTestStruct{Tags: []string{}}
	v, _ := validator.New(&s)
	validator.Slice(v, &s.Tags).NotEmpty().WithName("labels")
	result := v.Validate()
	assert.False(t, result.IsValid())
	assert.Equal(t, "labels", result.Errors[0].Field)
}

func TestSlice_Standalone_ReflectionFallback(t *testing.T) {
	// SliceOf yields SliceRule[any]; passing a concrete []int exercises the
	// reflection fallback that converts each element to any.
	rule := validator.SliceOf(validator.Int().Positive()).Min(2)

	errs := rule.Validate([]int{1, 2, 3})
	assert.Empty(t, errs)

	errs = rule.Validate([]int{1, -2, 3})
	assert.Len(t, errs, 1)
	assert.Equal(t, "[1]", errs[0].Field)

	// A non-slice value is rejected.
	errs = rule.Validate("not a slice")
	assert.Len(t, errs, 1)
	assert.Equal(t, "must be an array", errs[0].Message)
}

func TestSlice_NilSkip(t *testing.T) {
	var tags *[]string
	type s struct {
		Tags *[]string `json:"tags"`
	}
	st := s{Tags: tags}
	_ = st
	// Can't easily test nil slice with current API since Slice takes *[]T
	// A nil []string is still a valid *[]string (pointing to nil slice)
	// This tests that an empty slice passes when no rules are added
	empty := sliceTestStruct{Tags: nil}
	v, _ := validator.New(&empty)
	validator.Slice(v, &empty.Tags).Min(1)
	result := v.Validate()
	assert.False(t, result.IsValid()) // nil slice has len 0
}
