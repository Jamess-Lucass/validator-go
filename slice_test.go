package validator_test

import (
	"testing"

	validator "github.com/Jamess-Lucass/validator-go"
	"github.com/Jamess-Lucass/validator-go/rule"
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
	validator.Slice(v, &s.Tags).EachValue(rule.String().Min(2))
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
		validator.String(sv, &item.Name).NotEmpty()
		validator.Number(sv, &item.Quantity).Gte(1)
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

func TestSlice_NamedSliceType(t *testing.T) {
	type tags []string
	type S struct {
		Tags tags `json:"tags"`
	}
	s := S{Tags: tags{"go", ""}}
	v, _ := validator.New(&s)
	validator.Slice(v, &s.Tags).Min(1).EachValue(rule.String().NotEmpty())
	res := v.Validate()
	assert.Len(t, res.Errors, 1)
	assert.Equal(t, "tags[1]", res.Errors[0].Field)
}

func TestSlice_Each_PrimitiveNaming(t *testing.T) {
	s := sliceTestStruct{Tags: []string{"go", ""}}
	v, _ := validator.New(&s)
	validator.Slice(v, &s.Tags).Each(func(item *string, sv *validator.Validator) {
		validator.String(sv, item).NotEmpty()
	})
	res := v.Validate()
	assert.Len(t, res.Errors, 1)
	// The element itself can't be named by reflection, so the index is used.
	assert.Equal(t, "tags[1]", res.Errors[0].Field)
}

func TestSlice_Standalone_ElementTypeMismatch(t *testing.T) {
	r := &validator.SliceRule[[]int, int]{}
	errs := r.NotEmpty().Validate([]string{"a"})
	assert.Len(t, errs, 1)
	assert.Equal(t, "[0]", errs[0].Field)
	assert.Equal(t, "is of the wrong type", errs[0].Message)
}

func TestSliceOf_TypedNilSlice(t *testing.T) {
	payload := map[string]any{"tags": []string(nil)}
	v := validator.NewObject(payload)
	v.Field("tags", validator.SliceOf(rule.String().NotEmpty()).NotNil())
	res := v.Validate()
	assert.False(t, res.IsValid())
	assert.Equal(t, "tags", res.Errors[0].Field)
	assert.Equal(t, "must not be nil", res.Errors[0].Message)
}

func TestSlice_Standalone_ReflectionFallback(t *testing.T) {
	// SliceOf is built over []any, so a concrete []int goes through the
	// reflection fallback that converts each element.
	r := validator.SliceOf(rule.Int().Positive()).Min(2)

	errs := r.Validate([]int{1, 2, 3})
	assert.Empty(t, errs)

	errs = r.Validate([]int{1, -2, 3})
	assert.Len(t, errs, 1)
	assert.Equal(t, "[1]", errs[0].Field)

	// A non-slice value is rejected.
	errs = r.Validate("not a slice")
	assert.Len(t, errs, 1)
	assert.Equal(t, "must be an array", errs[0].Message)
}

func TestSlice_NilSliceFailsLengthRules(t *testing.T) {
	// A nil slice has length 0, so Min still applies even without NotNil.
	s := sliceTestStruct{Tags: nil}
	v, _ := validator.New(&s)
	validator.Slice(v, &s.Tags).Min(1)
	assert.False(t, v.Validate().IsValid())
}

func TestSlice_StandaloneReportsEveryTypeMismatch(t *testing.T) {
	// A typed slice rule fed a []any whose elements don't all match T reports
	// every mismatch, not just the first.
	type S struct {
		Nums []int `json:"nums"`
	}
	s := S{}
	v, _ := validator.New(&s)
	r := validator.Slice(v, &s.Nums)

	errs := r.Validate([]any{1, "x", 2.5, 3})
	assert.Len(t, errs, 2)
	assert.Equal(t, "[1]", errs[0].Field)
	assert.Equal(t, "[2]", errs[1].Field)
	assert.Equal(t, "is of the wrong type", errs[0].Message)
}

func TestSlice_Each_PointerElements(t *testing.T) {
	type entry struct {
		Name string `json:"name"`
	}
	type holder struct {
		Items []*entry `json:"items"`
	}

	h := holder{Items: []*entry{{Name: ""}}}
	v, _ := validator.New(&h)
	validator.Slice(v, &h.Items).Each(func(item **entry, sv *validator.Validator) {
		validator.String(sv, &(*item).Name).NotEmpty()
	})
	result := v.Validate()
	assert.Len(t, result.Errors, 1)
	assert.Equal(t, "items[0].name", result.Errors[0].Field)
}

func TestSlice_WithMessage_TargetsMostRecent(t *testing.T) {
	type holder struct {
		Tags []string `json:"tags"`
	}

	// After EachValue, WithMessage replaces the element errors' message and
	// leaves the earlier count rule's message alone.
	h := holder{Tags: []string{""}}
	v, _ := validator.New(&h)
	validator.Slice(v, &h.Tags).Min(2).EachValue(rule.String().NotEmpty()).WithMessage("bad tag")
	result := v.Validate()
	assert.Len(t, result.Errors, 2)
	assert.Equal(t, "must have at least 2 items", result.Errors[0].Message)
	assert.Equal(t, "tags[0]", result.Errors[1].Field)
	assert.Equal(t, "bad tag", result.Errors[1].Message)

	// A rule added after EachValue takes WithMessage back.
	v, _ = validator.New(&h)
	validator.Slice(v, &h.Tags).EachValue(rule.String().NotEmpty()).Min(2).WithMessage("too few")
	result = v.Validate()
	assert.Len(t, result.Errors, 2)
	assert.Equal(t, "too few", result.Errors[0].Message)
	assert.Equal(t, "must not be empty", result.Errors[1].Message)
}

func TestSlice_WithMessage_AfterEach(t *testing.T) {
	type entry struct {
		Name string `json:"name"`
	}
	type holder struct {
		Items []entry `json:"items"`
	}

	h := holder{Items: []entry{{Name: ""}}}
	v, _ := validator.New(&h)
	validator.Slice(v, &h.Items).Each(func(item *entry, sv *validator.Validator) {
		validator.String(sv, &item.Name).NotEmpty()
	}).WithMessage("bad item")
	result := v.Validate()
	assert.Len(t, result.Errors, 1)
	assert.Equal(t, "items[0].name", result.Errors[0].Field)
	assert.Equal(t, "bad item", result.Errors[0].Message)
}

func TestSlice_NilErrors_DeclarationOrder(t *testing.T) {
	type holder struct {
		Tags []string `json:"tags"`
	}

	// NotNil declared after NotEmpty reports after it too.
	h := holder{Tags: nil}
	v, _ := validator.New(&h)
	validator.Slice(v, &h.Tags).NotEmpty().NotNil()
	result := v.Validate()
	assert.Len(t, result.Errors, 2)
	assert.Equal(t, "must not be empty", result.Errors[0].Message)
	assert.Equal(t, "must not be nil", result.Errors[1].Message)
}

func TestSlice_EachValueReassign_ClearsMessage(t *testing.T) {
	type holder struct {
		Tags []string `json:"tags"`
	}

	// Replacing the EachValue rule drops the previous WithMessage override.
	h := holder{Tags: []string{""}}
	v, _ := validator.New(&h)
	validator.Slice(v, &h.Tags).EachValue(rule.String().Min(5)).WithMessage("stale").EachValue(rule.String().NotEmpty())
	result := v.Validate()
	assert.Len(t, result.Errors, 1)
	assert.Equal(t, "must not be empty", result.Errors[0].Message)
}

func TestSlice_NilEachInputsPanic(t *testing.T) {
	type holder struct {
		Tags []string `json:"tags"`
	}
	h := holder{}
	v, _ := validator.New(&h)

	assert.Panics(t, func() { validator.Slice(v, &h.Tags).EachValue(nil) })
	// A typed nil inside the Rule interface is just as unusable as plain nil.
	var typedNil *validator.StringRule[string]
	assert.Panics(t, func() { validator.Slice(v, &h.Tags).EachValue(typedNil) })
	assert.Panics(t, func() { validator.Slice(v, &h.Tags).Each(nil) })
	assert.Panics(t, func() { validator.Slice(v, &h.Tags).Must(nil) })
}
