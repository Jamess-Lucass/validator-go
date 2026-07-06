package validator_test

import (
	"testing"

	validator "github.com/Jamess-Lucass/validator-go"
	"github.com/Jamess-Lucass/validator-go/rule"
	"github.com/stretchr/testify/assert"
)

func TestComplex_Struct(t *testing.T) {
	type Signal struct {
		Amplitude complex128 `json:"amplitude"`
	}
	s := Signal{Amplitude: 0}
	v, _ := validator.New(&s)
	validator.Complex(v, &s.Amplitude).NotEmpty()
	res := v.Validate()
	assert.False(t, res.IsValid())
	assert.Equal(t, "amplitude", res.Errors[0].Field)
	assert.Equal(t, "must not be empty", res.Errors[0].Message)

	s.Amplitude = complex(1, 2)
	v, _ = validator.New(&s)
	validator.Complex(v, &s.Amplitude).NotEmpty()
	assert.True(t, v.Validate().IsValid())
}

func TestComplex_RuleSurface(t *testing.T) {
	type circuit struct {
		Z complex64 `json:"z"`
	}
	c := circuit{Z: 0}

	v, _ := validator.New(&c)
	validator.Complex(v, &c.Z).NotEmpty().WithMessage("impedance is required")
	res := v.Validate()
	assert.False(t, res.IsValid())
	assert.Equal(t, "z", res.Errors[0].Field)
	assert.Equal(t, "impedance is required", res.Errors[0].Message)

	v, _ = validator.New(&c)
	validator.Complex(v, &c.Z).NotEmpty().WithName("impedance")
	res = v.Validate()
	assert.Equal(t, "impedance", res.Errors[0].Field)

	r := rule.Complex64().Must(func(z complex64) bool {
		return real(z) > 0
	}).WithMessage("must have a positive real part")
	assert.Empty(t, r.Validate(complex64(complex(1, -1))))
	errs := r.Validate(complex64(complex(-1, 0)))
	assert.Len(t, errs, 1)
	assert.Equal(t, "must have a positive real part", errs[0].Message)

	assert.Len(t, rule.Complex128().NotNil().Validate(nil), 1)

	// Named complex types coerce by kind.
	type wave complex128
	assert.Empty(t, rule.Complex128().NotEmpty().Validate(wave(2+3i)))
}

func TestComplex_Standalone(t *testing.T) {
	r := rule.Complex128().NotEmpty()

	assert.Empty(t, r.Validate(complex(1, 1)))
	assert.Empty(t, r.Validate("1+2i")) // parsed
	assert.Empty(t, r.Validate(5.0))    // a real number is a valid complex

	assert.Len(t, r.Validate(complex128(0)), 1)

	errs := r.Validate("notcomplex")
	assert.Len(t, errs, 1)
	assert.Equal(t, "must be a complex number", errs[0].Message)
}

func TestComplex_RejectsOverflowAndNonFinite(t *testing.T) {
	r := rule.Complex64()

	// A magnitude that narrows to +Inf in complex64 is rejected, not wrapped.
	assert.Len(t, r.Validate(complex128(1e300)), 1)
	assert.Len(t, r.Validate(1e300), 1) // a real that overflows too
	assert.Len(t, r.Validate("1e300"), 1)

	// complex128 keeps the same value, so a finite magnitude stays valid.
	assert.Empty(t, rule.Complex128().Validate(complex128(1e300)))
}
