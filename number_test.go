package validator_test

import (
	"math"
	"testing"

	validator "github.com/Jamess-Lucass/validator-go"
	"github.com/Jamess-Lucass/validator-go/rule"
	"github.com/stretchr/testify/assert"
)

func TestNumber_Int64Field(t *testing.T) {
	type Account struct {
		Balance int64 `json:"balance"`
	}
	a := Account{Balance: 5}
	v, _ := validator.New(&a)
	validator.Number(v, &a.Balance).Gte(10)
	res := v.Validate()
	assert.False(t, res.IsValid())
	assert.Equal(t, "balance", res.Errors[0].Field)
	assert.Equal(t, "must be greater than or equal to 10", res.Errors[0].Message)
}

func TestNumber_UintAndFloat32(t *testing.T) {
	type S struct {
		Count uint    `json:"count"`
		Ratio float32 `json:"ratio"`
	}
	s := S{Count: 0, Ratio: 2.5}
	v, _ := validator.New(&s)
	validator.Number(v, &s.Count).Positive()
	validator.Number(v, &s.Ratio).Lt(2.0)
	res := v.Validate()
	assert.Len(t, res.Errors, 2)
	assert.Equal(t, "count", res.Errors[0].Field)
	assert.Equal(t, "ratio", res.Errors[1].Field)
}

func TestNumber_CustomNamedType(t *testing.T) {
	type Age int
	type User struct {
		Age Age `json:"age"`
	}
	u := User{Age: 10}
	v, _ := validator.New(&u)
	validator.Number(v, &u.Age).Gte(18)
	res := v.Validate()
	assert.False(t, res.IsValid())
	assert.Equal(t, "age", res.Errors[0].Field)
}

func TestNumber_MultipleOf_Int64Precision(t *testing.T) {
	// 2^60 + 1 is odd, but float64 can't represent it exactly (> 2^53), so a
	// float-based modulo would wrongly call it even. Integer modulo is exact.
	type S struct {
		N int64 `json:"n"`
	}
	s := S{N: 1<<60 + 1}
	v, _ := validator.New(&s)
	validator.Number(v, &s.N).MultipleOf(2)
	assert.False(t, v.Validate().IsValid())

	s.N = 1 << 60
	v, _ = validator.New(&s)
	validator.Number(v, &s.N).MultipleOf(2)
	assert.True(t, v.Validate().IsValid())
}

func TestNumber_Standalone_CoercionAndParsing(t *testing.T) {
	assert.Empty(t, rule.Int().Gte(0).Validate(int8(5)))
	assert.Empty(t, rule.Int().Gte(0).Validate(uint64(5)))
	assert.Empty(t, rule.Int().Gte(0).Validate(float64(5)))

	// Numeric strings are parsed.
	assert.Empty(t, rule.Int().Gte(0).Validate("5"))
	assert.Empty(t, rule.Float64().Gte(0).Validate("2.5"))

	// Non-numeric strings and fractional floats for an int rule are rejected.
	errs := rule.Int().Gte(0).Validate("abc")
	assert.Len(t, errs, 1)
	assert.Equal(t, "must be an integer", errs[0].Message)

	errs = rule.Int().Gte(0).Validate(2.5)
	assert.Len(t, errs, 1)
	assert.Equal(t, "must be an integer", errs[0].Message)
}

func TestNumber_OverflowRejected(t *testing.T) {
	// A uint64 above MaxInt64 must be rejected, not wrapped to a negative int.
	errs := rule.Int().Negative().Validate(uint64(math.MaxUint64))
	assert.Len(t, errs, 1)
	assert.Equal(t, "must be an integer", errs[0].Message)

	// Same for huge floats: converting one to int64 is implementation-defined,
	// so it must never get that far.
	assert.Len(t, rule.Int().Gte(0).Validate(1e30), 1)
	assert.Empty(t, rule.Int().Gte(0).Validate(float64(1<<60)))

	// Narrow targets reject out-of-range input from every source kind.
	r8 := (rule.Number[int8]()).Gte(0)
	assert.Empty(t, r8.Validate(100))
	assert.Len(t, r8.Validate(300), 1)
	assert.Len(t, r8.Validate("300"), 1)

	u8 := (rule.Number[uint8]()).Gte(0)
	assert.Empty(t, u8.Validate("200"))
	assert.Len(t, u8.Validate(-1), 1)
	assert.Len(t, u8.Validate(256), 1)
	assert.Empty(t, u8.Validate(3.0))
	assert.Len(t, u8.Validate(-1.0), 1)
	assert.Len(t, u8.Validate(300.0), 1)

	// 1e300 fits a float64 but not a float32; it must not turn into +Inf.
	f32 := rule.Number[float32]()
	assert.Len(t, f32.Gte(0).Validate(1e300), 1)
	assert.Empty(t, (rule.Number[float32]()).Gte(0).Validate(2.5))
}

func TestNumber_MultipleOfZero(t *testing.T) {
	// Nothing is a multiple of zero; this must fail cleanly, not panic.
	type S struct {
		N int     `json:"n"`
		U uint    `json:"u"`
		F float64 `json:"f"`
	}
	s := S{N: 4, U: 4, F: 4}
	v, _ := validator.New(&s)
	validator.Number(v, &s.N).MultipleOf(0)
	validator.Number(v, &s.U).MultipleOf(0)
	validator.Number(v, &s.F).MultipleOf(0)
	res := v.Validate()
	assert.Len(t, res.Errors, 3)
	assert.Equal(t, "must be a multiple of 0", res.Errors[0].Message)

	v, _ = validator.New(&s)
	validator.Number(v, &s.U).MultipleOf(2)
	assert.True(t, v.Validate().IsValid())
}

func TestNumber_Uint64Range(t *testing.T) {
	// MaxUint64 doesn't fit an int64, but a uint64 rule should accept it,
	// including when parsed from a string.
	u64 := rule.Number[uint64]()
	assert.Empty(t, u64.Gte(0).Validate("18446744073709551615"))

	assert.Len(t, rule.Int().Gte(0).Validate("18446744073709551615"), 1)
}

func TestNumber_Finite(t *testing.T) {
	assert.Len(t, rule.Float64().Finite().Validate(math.NaN()), 1)
	assert.Len(t, rule.Float64().Finite().Validate(math.Inf(1)), 1)
	assert.Empty(t, rule.Float64().Finite().Validate(5.5))

	type Reading struct {
		Value float64 `json:"value"`
	}
	r := Reading{Value: math.NaN()}
	v, _ := validator.New(&r)
	validator.Number(v, &r.Value).Finite()
	res := v.Validate()
	assert.False(t, res.IsValid())
	assert.Equal(t, "value", res.Errors[0].Field)
	assert.Equal(t, "must be a finite number", res.Errors[0].Message)
}

func TestNumber_StringNaNInfRejected(t *testing.T) {
	assert.Len(t, rule.Float64().Validate("NaN"), 1)
	assert.Len(t, rule.Float64().Validate("+Inf"), 1)

	// A NaN/Inf float value is rejected too, so the value and string paths agree.
	assert.Len(t, rule.Float64().Validate(math.NaN()), 1)
	assert.Len(t, rule.Float64().Validate(math.Inf(1)), 1)
	assert.Len(t, rule.Float64().Validate(math.Inf(-1)), 1)
}

func TestStandalone_PointerValues(t *testing.T) {
	five := 5
	assert.Empty(t, rule.Int().Gte(3).Validate(&five))

	// A typed nil pointer behaves like nil: skipped unless NotNil.
	assert.Empty(t, rule.Int().Gte(3).Validate((*int)(nil)))

	errs := rule.Int().NotNil().Validate((*int)(nil))
	assert.Len(t, errs, 1)
	assert.Equal(t, "must not be nil", errs[0].Message)
}

type intTestStruct struct {
	Age      int  `json:"age"`
	Priority *int `json:"priority"`
}

func TestInt_NotEmpty(t *testing.T) {
	s := intTestStruct{Age: 0}
	v, _ := validator.New(&s)
	validator.Number(v, &s.Age).NotEmpty()
	result := v.Validate()
	assert.False(t, result.IsValid())

	s.Age = 5
	v, _ = validator.New(&s)
	validator.Number(v, &s.Age).NotEmpty()
	result = v.Validate()
	assert.True(t, result.IsValid())
}

func TestInt_Gte(t *testing.T) {
	s := intTestStruct{Age: 17}
	v, _ := validator.New(&s)
	validator.Number(v, &s.Age).Gte(18)
	result := v.Validate()
	assert.False(t, result.IsValid())

	s.Age = 18
	v, _ = validator.New(&s)
	validator.Number(v, &s.Age).Gte(18)
	result = v.Validate()
	assert.True(t, result.IsValid())
}

func TestInt_Lte(t *testing.T) {
	s := intTestStruct{Age: 121}
	v, _ := validator.New(&s)
	validator.Number(v, &s.Age).Lte(120)
	result := v.Validate()
	assert.False(t, result.IsValid())

	s.Age = 120
	v, _ = validator.New(&s)
	validator.Number(v, &s.Age).Lte(120)
	result = v.Validate()
	assert.True(t, result.IsValid())
}

func TestInt_Gt(t *testing.T) {
	s := intTestStruct{Age: 0}
	v, _ := validator.New(&s)
	validator.Number(v, &s.Age).Gt(0)
	result := v.Validate()
	assert.False(t, result.IsValid())

	s.Age = 1
	v, _ = validator.New(&s)
	validator.Number(v, &s.Age).Gt(0)
	result = v.Validate()
	assert.True(t, result.IsValid())
}

func TestInt_Lt(t *testing.T) {
	s := intTestStruct{Age: 100}
	v, _ := validator.New(&s)
	validator.Number(v, &s.Age).Lt(100)
	result := v.Validate()
	assert.False(t, result.IsValid())

	s.Age = 99
	v, _ = validator.New(&s)
	validator.Number(v, &s.Age).Lt(100)
	result = v.Validate()
	assert.True(t, result.IsValid())
}

func TestInt_Positive(t *testing.T) {
	s := intTestStruct{Age: -1}
	v, _ := validator.New(&s)
	validator.Number(v, &s.Age).Positive()
	result := v.Validate()
	assert.False(t, result.IsValid())

	s.Age = 1
	v, _ = validator.New(&s)
	validator.Number(v, &s.Age).Positive()
	result = v.Validate()
	assert.True(t, result.IsValid())
}

func TestInt_Negative(t *testing.T) {
	s := intTestStruct{Age: 1}
	v, _ := validator.New(&s)
	validator.Number(v, &s.Age).Negative()
	result := v.Validate()
	assert.False(t, result.IsValid())

	s.Age = -1
	v, _ = validator.New(&s)
	validator.Number(v, &s.Age).Negative()
	result = v.Validate()
	assert.True(t, result.IsValid())
}

func TestInt_MultipleOf(t *testing.T) {
	s := intTestStruct{Age: 7}
	v, _ := validator.New(&s)
	validator.Number(v, &s.Age).MultipleOf(3)
	result := v.Validate()
	assert.False(t, result.IsValid())

	s.Age = 9
	v, _ = validator.New(&s)
	validator.Number(v, &s.Age).MultipleOf(3)
	result = v.Validate()
	assert.True(t, result.IsValid())
}

func TestInt_Nonnegative(t *testing.T) {
	s := intTestStruct{Age: -1}
	v, _ := validator.New(&s)
	validator.Number(v, &s.Age).Nonnegative()
	assert.False(t, v.Validate().IsValid())

	s.Age = 0
	v, _ = validator.New(&s)
	validator.Number(v, &s.Age).Nonnegative()
	assert.True(t, v.Validate().IsValid())
}

func TestInt_Nonpositive(t *testing.T) {
	s := intTestStruct{Age: 1}
	v, _ := validator.New(&s)
	validator.Number(v, &s.Age).Nonpositive()
	assert.False(t, v.Validate().IsValid())

	s.Age = 0
	v, _ = validator.New(&s)
	validator.Number(v, &s.Age).Nonpositive()
	assert.True(t, v.Validate().IsValid())
}

func TestInt_WithName(t *testing.T) {
	s := intTestStruct{Priority: nil}
	v, _ := validator.New(&s)
	validator.Number(v, s.Priority).NotNil().WithName("priority")
	result := v.Validate()
	assert.False(t, result.IsValid())
	assert.Equal(t, "priority", result.Errors[0].Field)
}

func TestInt_NilSkip(t *testing.T) {
	s := intTestStruct{Priority: nil}
	v, _ := validator.New(&s)
	validator.Number(v, s.Priority).Gte(1)
	result := v.Validate()
	assert.True(t, result.IsValid())
}

func TestInt_NotNil(t *testing.T) {
	s := intTestStruct{Priority: nil}
	v, _ := validator.New(&s)
	validator.Number(v, s.Priority).NotNil()
	result := v.Validate()
	assert.False(t, result.IsValid())
}

func TestInt_Must(t *testing.T) {
	s := intTestStruct{Age: 13}
	v, _ := validator.New(&s)
	validator.Number(v, &s.Age).Must(func(val int) bool {
		return val%2 == 0
	}).WithMessage("must be even")
	result := v.Validate()
	assert.False(t, result.IsValid())
	assert.Equal(t, "must be even", result.Errors[0].Message)
}

func TestInt_Standalone(t *testing.T) {
	r := rule.Int().Gte(18)

	errs := r.Validate(20)
	assert.Empty(t, errs)

	errs = r.Validate(10)
	assert.Len(t, errs, 1)

	// JSON numbers decode as float64
	errs = r.Validate(float64(20))
	assert.Empty(t, errs)

	errs = r.Validate("not a number")
	assert.Len(t, errs, 1)
	assert.Equal(t, "must be an integer", errs[0].Message)
}

func TestInt_Coercion_AllNumericKinds(t *testing.T) {
	r := rule.Int().Gte(18)

	// Signed, unsigned, and integral floats of every width are accepted.
	assert.Empty(t, r.Validate(int8(20)))
	assert.Empty(t, r.Validate(int16(20)))
	assert.Empty(t, r.Validate(int32(20)))
	assert.Empty(t, r.Validate(int64(20)))
	assert.Empty(t, r.Validate(uint(20)))
	assert.Empty(t, r.Validate(uint8(20)))
	assert.Empty(t, r.Validate(uint64(20)))
	assert.Empty(t, r.Validate(float32(20)))
	assert.Empty(t, r.Validate(float64(20)))

	// A float with a fractional part is not an integer.
	errs := r.Validate(float64(20.5))
	assert.Len(t, errs, 1)
	assert.Equal(t, "must be an integer", errs[0].Message)
}

type float64TestStruct struct {
	Price float64  `json:"price"`
	Rate  *float64 `json:"rate"`
}

func TestFloat64_NotNil_WithName(t *testing.T) {
	s := float64TestStruct{Rate: nil}
	v, _ := validator.New(&s)
	validator.Number(v, s.Rate).NotNil().WithName("rate")
	result := v.Validate()
	assert.False(t, result.IsValid())
	assert.Equal(t, "rate", result.Errors[0].Field)
}

func TestFloat64_Gte(t *testing.T) {
	s := float64TestStruct{Price: -1.0}
	v, _ := validator.New(&s)
	validator.Number(v, &s.Price).Gte(0)
	result := v.Validate()
	assert.False(t, result.IsValid())

	s.Price = 0
	v, _ = validator.New(&s)
	validator.Number(v, &s.Price).Gte(0)
	result = v.Validate()
	assert.True(t, result.IsValid())
}

func TestFloat64_Positive(t *testing.T) {
	s := float64TestStruct{Price: 0}
	v, _ := validator.New(&s)
	validator.Number(v, &s.Price).Positive()
	result := v.Validate()
	assert.False(t, result.IsValid())

	s.Price = 0.01
	v, _ = validator.New(&s)
	validator.Number(v, &s.Price).Positive()
	result = v.Validate()
	assert.True(t, result.IsValid())
}

func TestFloat64_NotEmpty(t *testing.T) {
	s := float64TestStruct{Price: 0}
	v, _ := validator.New(&s)
	validator.Number(v, &s.Price).NotEmpty()
	result := v.Validate()
	assert.False(t, result.IsValid())

	s.Price = 5.5
	v, _ = validator.New(&s)
	validator.Number(v, &s.Price).NotEmpty()
	result = v.Validate()
	assert.True(t, result.IsValid())
}

func TestFloat64_Lt(t *testing.T) {
	s := float64TestStruct{Price: 5.0}
	v, _ := validator.New(&s)
	validator.Number(v, &s.Price).Lt(5.0)
	assert.False(t, v.Validate().IsValid())

	s.Price = 4.9
	v, _ = validator.New(&s)
	validator.Number(v, &s.Price).Lt(5.0)
	assert.True(t, v.Validate().IsValid())
}

func TestFloat64_Gt(t *testing.T) {
	s := float64TestStruct{Price: 5.0}
	v, _ := validator.New(&s)
	validator.Number(v, &s.Price).Gt(5.0)
	assert.False(t, v.Validate().IsValid())

	s.Price = 5.1
	v, _ = validator.New(&s)
	validator.Number(v, &s.Price).Gt(5.0)
	assert.True(t, v.Validate().IsValid())
}

func TestFloat64_Negative(t *testing.T) {
	s := float64TestStruct{Price: 1.0}
	v, _ := validator.New(&s)
	validator.Number(v, &s.Price).Negative()
	assert.False(t, v.Validate().IsValid())

	s.Price = -1.0
	v, _ = validator.New(&s)
	validator.Number(v, &s.Price).Negative()
	assert.True(t, v.Validate().IsValid())
}

func TestFloat64_Nonnegative(t *testing.T) {
	s := float64TestStruct{Price: -0.1}
	v, _ := validator.New(&s)
	validator.Number(v, &s.Price).Nonnegative()
	assert.False(t, v.Validate().IsValid())

	s.Price = 0
	v, _ = validator.New(&s)
	validator.Number(v, &s.Price).Nonnegative()
	assert.True(t, v.Validate().IsValid())
}

func TestFloat64_Nonpositive(t *testing.T) {
	s := float64TestStruct{Price: 0.1}
	v, _ := validator.New(&s)
	validator.Number(v, &s.Price).Nonpositive()
	assert.False(t, v.Validate().IsValid())

	s.Price = 0
	v, _ = validator.New(&s)
	validator.Number(v, &s.Price).Nonpositive()
	assert.True(t, v.Validate().IsValid())
}

func TestFloat64_MultipleOf(t *testing.T) {
	s := float64TestStruct{Price: 0.7}
	v, _ := validator.New(&s)
	validator.Number(v, &s.Price).MultipleOf(0.25)
	assert.False(t, v.Validate().IsValid())

	s.Price = 0.75
	v, _ = validator.New(&s)
	validator.Number(v, &s.Price).MultipleOf(0.25)
	assert.True(t, v.Validate().IsValid())

	// Large magnitude stays robust (a fixed epsilon would mis-handle this).
	s.Price = 1e12
	v, _ = validator.New(&s)
	validator.Number(v, &s.Price).MultipleOf(0.25)
	assert.True(t, v.Validate().IsValid())
}

func TestFloat64_Must(t *testing.T) {
	s := float64TestStruct{Price: 3.0}
	v, _ := validator.New(&s)
	validator.Number(v, &s.Price).Must(func(val float64) bool {
		return val >= 1.0 && val <= 2.0
	}).WithMessage("must be between 1 and 2")
	result := v.Validate()
	assert.False(t, result.IsValid())
	assert.Equal(t, "must be between 1 and 2", result.Errors[0].Message)
}

func TestFloat64_Standalone(t *testing.T) {
	r := rule.Float64().Gte(0).Lte(100)

	errs := r.Validate(50.5)
	assert.Empty(t, errs)

	errs = r.Validate(-1.0)
	assert.Len(t, errs, 1)

	errs = r.Validate("string")
	assert.Len(t, errs, 1)
	assert.Equal(t, "must be a number", errs[0].Message)
}

func TestNumber_OneOf(t *testing.T) {
	type order struct {
		Priority int `json:"priority"`
	}

	o := order{Priority: 4}
	v, _ := validator.New(&o)
	validator.Number(v, &o.Priority).OneOf(1, 2, 3)
	result := v.Validate()
	assert.False(t, result.IsValid())
	assert.Equal(t, "must be one of 1, 2, 3", result.Errors[0].Message)

	o.Priority = 2
	v, _ = validator.New(&o)
	validator.Number(v, &o.Priority).OneOf(1, 2, 3)
	assert.True(t, v.Validate().IsValid())
}

func TestFloat64_Coercion_AllNumericKinds(t *testing.T) {
	r := rule.Float64().Gte(0)

	assert.Empty(t, r.Validate(float32(1.5)))
	assert.Empty(t, r.Validate(float64(1.5)))
	assert.Empty(t, r.Validate(int(2)))
	assert.Empty(t, r.Validate(int64(2)))
	assert.Empty(t, r.Validate(uint(2)))

	errs := r.Validate(true)
	assert.Len(t, errs, 1)
	assert.Equal(t, "must be a number", errs[0].Message)
}

func TestNumber_FloatTargetRejectsInexactIntegers(t *testing.T) {
	// 2^53+1 is the first integer float64 cannot hold; silently rounding it
	// would run the checks against a value the caller never had.
	assert.Len(t, rule.Float64().Validate(int64(9007199254740993)), 1)
	assert.Len(t, rule.Float64().Validate(uint64(9007199254740993)), 1)
	assert.Len(t, rule.Float64().Validate("9007199254740993"), 1)
	assert.Empty(t, rule.Float64().Validate(int64(9007199254740992)))
	assert.Empty(t, rule.Float64().Validate("9007199254740992"))

	// float32 loses integer precision above 2^24.
	assert.Len(t, rule.Number[float32]().Validate(int64(16777217)), 1)
	assert.Empty(t, rule.Number[float32]().Validate(int64(16777216)))
}
