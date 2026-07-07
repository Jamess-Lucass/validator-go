package validator_test

import (
	"encoding/json"
	"testing"
	"time"

	validator "github.com/Jamess-Lucass/validator-go"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// Must is the same contract on every rule type: a nil function is a
// programmer error caught at declaration, not a nil deref at Validate time.
// String, Slice, and Map have the same check in their own test files.
func TestMust_NilPanicsOnEveryType(t *testing.T) {
	type target struct {
		N int        `json:"n"`
		B bool       `json:"b"`
		T time.Time  `json:"t"`
		U uuid.UUID  `json:"u"`
		C complex128 `json:"c"`
	}
	s := target{}
	v, _ := validator.New(&s)

	assert.Panics(t, func() { validator.Number(v, &s.N).Must(nil) })
	assert.Panics(t, func() { validator.Bool(v, &s.B).Must(nil) })
	assert.Panics(t, func() { validator.Time(v, &s.T).Must(nil) })
	assert.Panics(t, func() { validator.UUID(v, &s.U).Must(nil) })
	assert.Panics(t, func() { validator.Complex(v, &s.C).Must(nil) })
}

// Results usually go straight out of a JSON API, so the whole shape
// serialises lowercase.
func TestValidationResult_JSONShape(t *testing.T) {
	result := &validator.ValidationResult{
		Errors: []validator.ValidationError{{Field: "name", Message: "must not be empty"}},
	}
	b, err := json.Marshal(result)
	assert.NoError(t, err)
	assert.JSONEq(t, `{"errors":[{"field":"name","message":"must not be empty"}]}`, string(b))
}
