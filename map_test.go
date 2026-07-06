package validator_test

import (
	"math"
	"testing"

	validator "github.com/Jamess-Lucass/validator-go"
	"github.com/Jamess-Lucass/validator-go/rule"
	"github.com/stretchr/testify/assert"
)

type appContainer struct {
	Image    string `json:"image"`
	Replicas int    `json:"replicas"`
}

type deployTestStruct struct {
	Name       string                  `json:"name"`
	Labels     map[string]string       `json:"labels"`
	Containers map[string]appContainer `json:"containers"`
	Weights    map[int]float64         `json:"weights"`
	Extra      *map[string]string      `json:"extra"`
}

func TestMapField_EachValue(t *testing.T) {
	s := deployTestStruct{Labels: map[string]string{"env": "", "team": "core"}}
	v, _ := validator.New(&s)
	validator.Map(v, &s.Labels).EachValue(rule.String().NotEmpty())
	res := v.Validate()
	assert.False(t, res.IsValid())
	assert.Len(t, res.Errors, 1)
	assert.Equal(t, "labels[env]", res.Errors[0].Field)
	assert.Equal(t, "must not be empty", res.Errors[0].Message)

	s.Labels = map[string]string{"env": "prod", "team": "core"}
	v, _ = validator.New(&s)
	validator.Map(v, &s.Labels).EachValue(rule.String().NotEmpty())
	assert.True(t, v.Validate().IsValid())
}

func TestMapField_EachKey(t *testing.T) {
	s := deployTestStruct{Labels: map[string]string{"": "x", "ok": "y"}}
	v, _ := validator.New(&s)
	validator.Map(v, &s.Labels).EachKey(rule.String().NotEmpty())
	res := v.Validate()
	assert.False(t, res.IsValid())
	assert.Len(t, res.Errors, 1)
	assert.Equal(t, "labels[]", res.Errors[0].Field)
	assert.Equal(t, "key must not be empty", res.Errors[0].Message)
}

func TestMapField_Each_StructValues(t *testing.T) {
	s := deployTestStruct{Containers: map[string]appContainer{
		"web": {Image: "", Replicas: 0},
		"db":  {Image: "postgres:16", Replicas: 1},
	}}
	v, _ := validator.New(&s)
	validator.Map(v, &s.Containers).Each(func(name string, c *appContainer, sv *validator.Validator) {
		validator.String(sv, &c.Image).NotEmpty()
		validator.Number(sv, &c.Replicas).Gte(1)
	})
	res := v.Validate()
	assert.False(t, res.IsValid())
	assert.Len(t, res.Errors, 2)
	assert.Equal(t, "containers[web].image", res.Errors[0].Field)
	assert.Equal(t, "containers[web].replicas", res.Errors[1].Field)
}

func TestMapField_Each_PrimitiveValues(t *testing.T) {
	s := deployTestStruct{Labels: map[string]string{"env": ""}}
	v, _ := validator.New(&s)
	validator.Map(v, &s.Labels).Each(func(key string, val *string, sv *validator.Validator) {
		validator.String(sv, val).NotEmpty()
	})
	res := v.Validate()
	assert.False(t, res.IsValid())
	// The value itself can't be named by reflection, so the entry name is used.
	assert.Equal(t, "labels[env]", res.Errors[0].Field)
}

func TestMapField_HasKey(t *testing.T) {
	s := deployTestStruct{Labels: map[string]string{"team": "core"}}
	v, _ := validator.New(&s)
	validator.Map(v, &s.Labels).HasKey("env")
	res := v.Validate()
	assert.False(t, res.IsValid())
	assert.Equal(t, "labels", res.Errors[0].Field)
	assert.Equal(t, `must contain key "env"`, res.Errors[0].Message)

	s.Labels["env"] = "prod"
	v, _ = validator.New(&s)
	validator.Map(v, &s.Labels).HasKey("env")
	assert.True(t, v.Validate().IsValid())
}

func TestMapField_CountRules(t *testing.T) {
	s := deployTestStruct{Labels: map[string]string{"a": "1"}}

	v, _ := validator.New(&s)
	validator.Map(v, &s.Labels).Min(2)
	res := v.Validate()
	assert.Equal(t, "must have at least 2 items", res.Errors[0].Message)

	v, _ = validator.New(&s)
	validator.Map(v, &s.Labels).Length(2, 3)
	res = v.Validate()
	assert.Equal(t, "must have between 2 and 3 items", res.Errors[0].Message)

	s.Labels = map[string]string{"a": "1", "b": "2"}
	v, _ = validator.New(&s)
	validator.Map(v, &s.Labels).Max(1)
	res = v.Validate()
	assert.Equal(t, "must have at most 1 item", res.Errors[0].Message)

	s.Labels = map[string]string{}
	v, _ = validator.New(&s)
	validator.Map(v, &s.Labels).NotEmpty()
	res = v.Validate()
	assert.Equal(t, "must not be empty", res.Errors[0].Message)
}

func TestMapField_NilMap(t *testing.T) {
	s := deployTestStruct{Labels: nil}
	v, _ := validator.New(&s)
	validator.Map(v, &s.Labels).NotNil().NotEmpty()
	res := v.Validate()
	assert.Len(t, res.Errors, 2)
	assert.Equal(t, "labels", res.Errors[0].Field)
	assert.Equal(t, "must not be nil", res.Errors[0].Message)
	assert.Equal(t, "must not be empty", res.Errors[1].Message)
}

func TestMapField_NilPointerField(t *testing.T) {
	s := deployTestStruct{Extra: nil}

	v, _ := validator.New(&s)
	validator.Map(v, s.Extra).NotEmpty()
	assert.True(t, v.Validate().IsValid())

	v, _ = validator.New(&s)
	validator.Map(v, s.Extra).NotNil().WithName("extra")
	res := v.Validate()
	assert.False(t, res.IsValid())
	assert.Equal(t, "extra", res.Errors[0].Field)
	assert.Equal(t, "must not be nil", res.Errors[0].Message)
}

func TestMapField_SortedErrorOrder(t *testing.T) {
	s := deployTestStruct{Labels: map[string]string{"b": "", "a": "", "c": ""}}
	v, _ := validator.New(&s)
	validator.Map(v, &s.Labels).EachValue(rule.String().NotEmpty())
	res := v.Validate()
	assert.Len(t, res.Errors, 3)
	assert.Equal(t, "labels[a]", res.Errors[0].Field)
	assert.Equal(t, "labels[b]", res.Errors[1].Field)
	assert.Equal(t, "labels[c]", res.Errors[2].Field)
}

func TestMapField_IntKeysSortNumerically(t *testing.T) {
	s := deployTestStruct{Weights: map[int]float64{10: -1, 2: -1}}
	v, _ := validator.New(&s)
	validator.Map(v, &s.Weights).EachValue(rule.Float64().Positive())
	res := v.Validate()
	assert.Len(t, res.Errors, 2)
	assert.Equal(t, "weights[2]", res.Errors[0].Field)
	assert.Equal(t, "weights[10]", res.Errors[1].Field)
}

func TestMapField_Must(t *testing.T) {
	s := deployTestStruct{Labels: map[string]string{"team": "core"}}
	v, _ := validator.New(&s)
	validator.Map(v, &s.Labels).Must(func(m map[string]string) bool {
		return m["team"] == "platform"
	}).WithMessage("must belong to the platform team")
	res := v.Validate()
	assert.False(t, res.IsValid())
	assert.Equal(t, "labels", res.Errors[0].Field)
	assert.Equal(t, "must belong to the platform team", res.Errors[0].Message)
}

func TestMapField_WithName(t *testing.T) {
	s := deployTestStruct{Labels: map[string]string{}}
	v, _ := validator.New(&s)
	validator.Map(v, &s.Labels).NotEmpty().WithName("metadata.labels")
	res := v.Validate()
	assert.Equal(t, "metadata.labels", res.Errors[0].Field)
}

func TestMapOf_Standalone(t *testing.T) {
	payload := map[string]any{"labels": map[string]any{"env": ""}}
	v := validator.NewObject(payload)
	v.Field("labels", validator.MapOf(rule.String().NotEmpty()).NotEmpty())
	res := v.Validate()
	assert.False(t, res.IsValid())
	assert.Equal(t, "labels[env]", res.Errors[0].Field)
	assert.Equal(t, "must not be empty", res.Errors[0].Message)

	// A concretely typed map takes the reflection path.
	payload = map[string]any{"labels": map[string]string{"env": "prod"}}
	v = validator.NewObject(payload)
	v.Field("labels", validator.MapOf(rule.String().NotEmpty()))
	assert.True(t, v.Validate().IsValid())

	payload = map[string]any{"labels": 5}
	v = validator.NewObject(payload)
	v.Field("labels", validator.MapOf(rule.String().NotEmpty()))
	res = v.Validate()
	assert.Equal(t, "labels", res.Errors[0].Field)
	assert.Equal(t, "must be an object", res.Errors[0].Message)
}

func TestMapField_NamedMapType(t *testing.T) {
	type labels map[string]string
	type S struct {
		Labels labels `json:"labels"`
	}
	s := S{Labels: labels{"env": ""}}
	v, _ := validator.New(&s)
	validator.Map(v, &s.Labels).HasKey("env").EachValue(rule.String().NotEmpty())
	res := v.Validate()
	assert.Len(t, res.Errors, 1)
	assert.Equal(t, "labels[env]", res.Errors[0].Field)
}

func TestMapField_OtherKeyKinds(t *testing.T) {
	type S struct {
		ByPort map[uint16]string `json:"by_port"`
		Flags  map[bool]string   `json:"flags"`
	}
	s := S{
		ByPort: map[uint16]string{8080: "", 443: ""},
		Flags:  map[bool]string{true: "", false: ""},
	}
	v, _ := validator.New(&s)
	validator.Map(v, &s.ByPort).EachValue(rule.String().NotEmpty())
	validator.Map(v, &s.Flags).EachValue(rule.String().NotEmpty())
	res := v.Validate()
	assert.Len(t, res.Errors, 4)
	assert.Equal(t, "by_port[443]", res.Errors[0].Field)
	assert.Equal(t, "by_port[8080]", res.Errors[1].Field)
	assert.Equal(t, "flags[false]", res.Errors[2].Field)
	assert.Equal(t, "flags[true]", res.Errors[3].Field)
}

func TestMapOf_MissingKeySkipped(t *testing.T) {
	v := validator.NewObject(map[string]any{})
	v.Field("labels", validator.MapOf(rule.String().NotEmpty()).Min(1))
	assert.True(t, v.Validate().IsValid())
}

func TestMapOf_TypedNilMap(t *testing.T) {
	payload := map[string]any{"labels": map[string]string(nil)}
	v := validator.NewObject(payload)
	v.Field("labels", validator.MapOf(rule.String().NotEmpty()).NotNil())
	res := v.Validate()
	assert.False(t, res.IsValid())
	assert.Equal(t, "labels", res.Errors[0].Field)
	assert.Equal(t, "must not be nil", res.Errors[0].Message)
}

func TestMapField_Each_PointerValues(t *testing.T) {
	type image struct {
		Ref string `json:"ref"`
	}
	type deployment struct {
		Images map[string]*image `json:"images"`
	}

	d := deployment{Images: map[string]*image{"web": {Ref: ""}}}
	v, _ := validator.New(&d)
	validator.Map(v, &d.Images).Each(func(name string, img **image, sv *validator.Validator) {
		validator.String(sv, &(*img).Ref).NotEmpty()
	})
	result := v.Validate()
	assert.Len(t, result.Errors, 1)
	assert.Equal(t, "images[web].ref", result.Errors[0].Field)
}

func TestMapField_WithMessage_AfterEachVariants(t *testing.T) {
	type deployment struct {
		Labels map[string]string `json:"labels"`
	}

	// After EachKey, WithMessage replaces the whole "key ..." message.
	d := deployment{Labels: map[string]string{"": "x"}}
	v, _ := validator.New(&d)
	validator.Map(v, &d.Labels).EachKey(rule.String().NotEmpty()).WithMessage("bad key")
	result := v.Validate()
	assert.Len(t, result.Errors, 1)
	assert.Equal(t, "labels[]", result.Errors[0].Field)
	assert.Equal(t, "bad key", result.Errors[0].Message)

	d = deployment{Labels: map[string]string{"env": ""}}
	v, _ = validator.New(&d)
	validator.Map(v, &d.Labels).EachValue(rule.String().NotEmpty()).WithMessage("bad value")
	result = v.Validate()
	assert.Len(t, result.Errors, 1)
	assert.Equal(t, "bad value", result.Errors[0].Message)

	// A rule added after EachValue takes WithMessage back.
	v, _ = validator.New(&d)
	validator.Map(v, &d.Labels).EachValue(rule.String().NotEmpty()).Min(5).WithMessage("too few")
	result = v.Validate()
	assert.Len(t, result.Errors, 2)
	assert.Equal(t, "too few", result.Errors[0].Message)
	assert.Equal(t, "must not be empty", result.Errors[1].Message)
}

func TestMapField_PointerKeys_StableLabels(t *testing.T) {
	one, two := 1, 2
	type holder struct {
		M map[*int]string `json:"m"`
	}

	// Pointer keys are labelled and ordered by the value they point to, so
	// the output doesn't change with the allocator's addresses.
	h := holder{M: map[*int]string{&two: "", &one: ""}}
	v, _ := validator.New(&h)
	validator.Map(v, &h.M).EachValue(rule.String().NotEmpty())
	result := v.Validate()
	assert.Len(t, result.Errors, 2)
	assert.Equal(t, "m[1]", result.Errors[0].Field)
	assert.Equal(t, "m[2]", result.Errors[1].Field)

	// A nil pointer key sorts first and is labelled "nil".
	h = holder{M: map[*int]string{nil: "", &one: ""}}
	v, _ = validator.New(&h)
	validator.Map(v, &h.M).EachValue(rule.String().NotEmpty())
	result = v.Validate()
	assert.Len(t, result.Errors, 2)
	assert.Equal(t, "m[nil]", result.Errors[0].Field)
	assert.Equal(t, "m[1]", result.Errors[1].Field)
}

func TestMapField_WithMessage_AfterEach(t *testing.T) {
	type box struct {
		V string `json:"v"`
	}
	type holder struct {
		M map[string]box `json:"m"`
	}

	h := holder{M: map[string]box{"a": {V: ""}}}
	v, _ := validator.New(&h)
	validator.Map(v, &h.M).Each(func(k string, b *box, sv *validator.Validator) {
		validator.String(sv, &b.V).NotEmpty()
	}).WithMessage("bad entry")
	result := v.Validate()
	assert.Len(t, result.Errors, 1)
	assert.Equal(t, "m[a].v", result.Errors[0].Field)
	assert.Equal(t, "bad entry", result.Errors[0].Message)
}

func TestMapField_NaNKey_KeepsOtherKeysSorted(t *testing.T) {
	type holder struct {
		M map[float64]string `json:"m"`
	}

	// A NaN key must not break the sort's ordering contract for the rest;
	// repeat because map iteration hands sort a fresh permutation each time.
	for i := 0; i < 20; i++ {
		h := holder{M: map[float64]string{1: "", 2: "", math.NaN(): "", 3: ""}}
		v, _ := validator.New(&h)
		validator.Map(v, &h.M).EachValue(rule.String().NotEmpty())
		result := v.Validate()
		assert.Len(t, result.Errors, 4)
		assert.Equal(t, "m[1]", result.Errors[0].Field)
		assert.Equal(t, "m[2]", result.Errors[1].Field)
		assert.Equal(t, "m[3]", result.Errors[2].Field)
		assert.Equal(t, "m[NaN]", result.Errors[3].Field)
	}
}

func TestMapField_NilEachInputsPanic(t *testing.T) {
	type deployment struct {
		Labels map[string]string `json:"labels"`
	}
	d := deployment{}
	v, _ := validator.New(&d)

	assert.Panics(t, func() { validator.Map(v, &d.Labels).EachValue(nil) })
	assert.Panics(t, func() { validator.Map(v, &d.Labels).EachKey(nil) })
	assert.Panics(t, func() { validator.Map(v, &d.Labels).Each(nil) })
	assert.Panics(t, func() { validator.Map(v, &d.Labels).Must(nil) })
}
