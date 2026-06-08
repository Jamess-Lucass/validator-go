package validator_test

import (
	"testing"

	validator "github.com/Jamess-Lucass/validator-go"
	"github.com/stretchr/testify/assert"
)

type stringTestStruct struct {
	Name  string  `json:"name"`
	Email string  `json:"email"`
	Phone *string `json:"phone"`
}

func TestString_NotEmpty(t *testing.T) {
	s := stringTestStruct{Name: ""}
	v, _ := validator.New(&s)
	v.String(&s.Name).NotEmpty()
	result := v.Validate()
	assert.False(t, result.IsValid())
	assert.Equal(t, "name", result.Errors[0].Field)
	assert.Equal(t, "must not be empty", result.Errors[0].Message)

	s.Name = "hello"
	v, _ = validator.New(&s)
	v.String(&s.Name).NotEmpty()
	result = v.Validate()
	assert.True(t, result.IsValid())
}

func TestString_Min(t *testing.T) {
	s := stringTestStruct{Name: "ab"}
	v, _ := validator.New(&s)
	v.String(&s.Name).Min(5)
	result := v.Validate()
	assert.False(t, result.IsValid())
	assert.Equal(t, "must be at least 5 characters", result.Errors[0].Message)

	s.Name = "abcde"
	v, _ = validator.New(&s)
	v.String(&s.Name).Min(5)
	result = v.Validate()
	assert.True(t, result.IsValid())
}

func TestString_Max(t *testing.T) {
	s := stringTestStruct{Name: "123456"}
	v, _ := validator.New(&s)
	v.String(&s.Name).Max(5)
	result := v.Validate()
	assert.False(t, result.IsValid())

	s.Name = "12345"
	v, _ = validator.New(&s)
	v.String(&s.Name).Max(5)
	result = v.Validate()
	assert.True(t, result.IsValid())
}

func TestString_Length(t *testing.T) {
	s := stringTestStruct{Name: "1234"}
	v, _ := validator.New(&s)
	v.String(&s.Name).Length(5, 10)
	result := v.Validate()
	assert.False(t, result.IsValid())

	s.Name = "12345"
	v, _ = validator.New(&s)
	v.String(&s.Name).Length(5, 10)
	result = v.Validate()
	assert.True(t, result.IsValid())

	s.Name = "12345678901"
	v, _ = validator.New(&s)
	v.String(&s.Name).Length(5, 10)
	result = v.Validate()
	assert.False(t, result.IsValid())
}

func TestString_Min_CountsRunesNotBytes(t *testing.T) {
	// "héllo" is 5 runes but 6 bytes (é is 2 bytes). Counting runes, Max(5)
	// passes; counting bytes it would fail.
	s := stringTestStruct{Name: "héllo"}
	v, _ := validator.New(&s)
	v.String(&s.Name).Min(5).Max(5)
	result := v.Validate()
	assert.True(t, result.IsValid())
}

func TestString_Min_SingularMessage(t *testing.T) {
	s := stringTestStruct{Name: ""}
	v, _ := validator.New(&s)
	v.String(&s.Name).Min(1)
	result := v.Validate()
	assert.False(t, result.IsValid())
	assert.Equal(t, "must be at least 1 character", result.Errors[0].Message)
}

func TestString_WithName(t *testing.T) {
	s := stringTestStruct{Phone: nil}
	v, _ := validator.New(&s)
	v.String(s.Phone).NotNil().WithName("phone")
	result := v.Validate()
	assert.False(t, result.IsValid())
	assert.Equal(t, "phone", result.Errors[0].Field)
}

func TestString_Email(t *testing.T) {
	s := stringTestStruct{Email: "test@example.com"}
	v, _ := validator.New(&s)
	v.String(&s.Email).Email()
	result := v.Validate()
	assert.True(t, result.IsValid())

	s.Email = "not-an-email"
	v, _ = validator.New(&s)
	v.String(&s.Email).Email()
	result = v.Validate()
	assert.False(t, result.IsValid())
	assert.Equal(t, "must be a valid email address", result.Errors[0].Message)
}

func TestString_Url(t *testing.T) {
	s := stringTestStruct{Name: "http://google.com"}
	v, _ := validator.New(&s)
	v.String(&s.Name).Url()
	result := v.Validate()
	assert.True(t, result.IsValid())

	s.Name = "not-a-url"
	v, _ = validator.New(&s)
	v.String(&s.Name).Url()
	result = v.Validate()
	assert.False(t, result.IsValid())
}

func TestString_Includes(t *testing.T) {
	s := stringTestStruct{Name: "X_test_X"}
	v, _ := validator.New(&s)
	v.String(&s.Name).Includes("test")
	result := v.Validate()
	assert.True(t, result.IsValid())

	s.Name = "X_Test_X"
	v, _ = validator.New(&s)
	v.String(&s.Name).Includes("test")
	result = v.Validate()
	assert.False(t, result.IsValid())
}

func TestString_StartsWith(t *testing.T) {
	s := stringTestStruct{Name: "test_X"}
	v, _ := validator.New(&s)
	v.String(&s.Name).StartsWith("test")
	result := v.Validate()
	assert.True(t, result.IsValid())

	s.Name = "Test_X"
	v, _ = validator.New(&s)
	v.String(&s.Name).StartsWith("test")
	result = v.Validate()
	assert.False(t, result.IsValid())
}

func TestString_EndsWith(t *testing.T) {
	s := stringTestStruct{Name: "X_test"}
	v, _ := validator.New(&s)
	v.String(&s.Name).EndsWith("test")
	result := v.Validate()
	assert.True(t, result.IsValid())

	s.Name = "X_Test"
	v, _ = validator.New(&s)
	v.String(&s.Name).EndsWith("test")
	result = v.Validate()
	assert.False(t, result.IsValid())
}

func TestString_NilSkip(t *testing.T) {
	s := stringTestStruct{Phone: nil}
	v, _ := validator.New(&s)
	v.String(s.Phone).Min(5)
	result := v.Validate()
	assert.True(t, result.IsValid())
}

func TestString_NotNil(t *testing.T) {
	s := stringTestStruct{Phone: nil}
	v, _ := validator.New(&s)
	v.String(s.Phone).NotNil()
	result := v.Validate()
	assert.False(t, result.IsValid())
	assert.Equal(t, "must not be nil", result.Errors[0].Message)

	phone := "1234567"
	s.Phone = &phone
	v, _ = validator.New(&s)
	v.String(s.Phone).NotNil().Min(5)
	result = v.Validate()
	assert.True(t, result.IsValid())
}

func TestString_Must(t *testing.T) {
	s := stringTestStruct{Email: "user@other.com"}
	v, _ := validator.New(&s)
	v.String(&s.Email).Must(func(val string) bool {
		return len(val) > 5 && val[len(val)-12:] == "@company.com"
	}).WithMessage("must be a company email")
	result := v.Validate()
	assert.False(t, result.IsValid())
	assert.Equal(t, "must be a company email", result.Errors[0].Message)

	s.Email = "user@company.com"
	v, _ = validator.New(&s)
	v.String(&s.Email).Must(func(val string) bool {
		return len(val) > 12 && val[len(val)-12:] == "@company.com"
	}).WithMessage("must be a company email")
	result = v.Validate()
	assert.True(t, result.IsValid())
}

func TestString_WithMessage(t *testing.T) {
	s := stringTestStruct{Name: ""}
	v, _ := validator.New(&s)
	v.String(&s.Name).NotEmpty().WithMessage("name is required")
	result := v.Validate()
	assert.False(t, result.IsValid())
	assert.Equal(t, "name is required", result.Errors[0].Message)
}

func TestString_MultipleRulesAllReported(t *testing.T) {
	s := stringTestStruct{Name: ""}
	v, _ := validator.New(&s)
	v.String(&s.Name).NotEmpty().Min(5)
	result := v.Validate()
	assert.Len(t, result.Errors, 2)
}

func TestString_JsonTagFieldName(t *testing.T) {
	s := stringTestStruct{Name: ""}
	v, _ := validator.New(&s)
	v.String(&s.Name).NotEmpty()
	result := v.Validate()
	assert.Equal(t, "name", result.Errors[0].Field)
}

func TestString_Standalone_Validate(t *testing.T) {
	rule := validator.String().NotEmpty().Email()

	errs := rule.Validate("test@example.com")
	assert.Empty(t, errs)

	errs = rule.Validate("")
	assert.Len(t, errs, 2)

	errs = rule.Validate(123)
	assert.Len(t, errs, 1)
	assert.Equal(t, "must be a string", errs[0].Message)

	errs = rule.Validate(nil)
	assert.Empty(t, errs)
}

func TestString_NotNil_WithMessage(t *testing.T) {
	s := stringTestStruct{Phone: nil}
	v, _ := validator.New(&s)
	v.String(s.Phone).NotNil().WithMessage("phone is required")
	result := v.Validate()
	assert.False(t, result.IsValid())
	assert.Equal(t, "phone is required", result.Errors[0].Message)
}

func TestString_WithMessage_DoesNotBreakNotNil(t *testing.T) {
	// Ensure that using WithMessage("must not be nil") on a non-NotNil rule
	// does NOT cause it to be treated as a NotNil rule
	s := stringTestStruct{Name: "hello"}
	v, _ := validator.New(&s)
	v.String(&s.Name).Min(10).WithMessage("must not be nil")
	result := v.Validate()
	assert.False(t, result.IsValid())
	assert.Equal(t, "must not be nil", result.Errors[0].Message)
}

func TestStruct_ValidatesCurrentState(t *testing.T) {
	s := stringTestStruct{Name: "ab"}

	v, _ := validator.New(&s)
	v.String(&s.Name).Min(5)

	// Mutate the struct after building rules but before validating
	s.Name = "abcde"

	result := v.Validate()
	assert.True(t, result.IsValid())
}

func TestStruct_ValidatesCurrentState_Pointer(t *testing.T) {
	short := "ab"
	s := stringTestStruct{Phone: &short}

	v, _ := validator.New(&s)
	v.String(s.Phone).Min(5)

	// Mutate the value through the pointer after building rules
	short = "abcdefg"

	result := v.Validate()
	assert.True(t, result.IsValid())
}
