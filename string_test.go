package validator_test

import (
	"testing"

	validator "github.com/Jamess-Lucass/validator-go"
	"github.com/Jamess-Lucass/validator-go/rule"
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
	validator.String(v, &s.Name).NotEmpty()
	result := v.Validate()
	assert.False(t, result.IsValid())
	assert.Equal(t, "name", result.Errors[0].Field)
	assert.Equal(t, "must not be empty", result.Errors[0].Message)

	s.Name = "hello"
	v, _ = validator.New(&s)
	validator.String(v, &s.Name).NotEmpty()
	result = v.Validate()
	assert.True(t, result.IsValid())
}

func TestString_Min(t *testing.T) {
	s := stringTestStruct{Name: "ab"}
	v, _ := validator.New(&s)
	validator.String(v, &s.Name).Min(5)
	result := v.Validate()
	assert.False(t, result.IsValid())
	assert.Equal(t, "must be at least 5 characters", result.Errors[0].Message)

	s.Name = "abcde"
	v, _ = validator.New(&s)
	validator.String(v, &s.Name).Min(5)
	result = v.Validate()
	assert.True(t, result.IsValid())
}

func TestString_Max(t *testing.T) {
	s := stringTestStruct{Name: "123456"}
	v, _ := validator.New(&s)
	validator.String(v, &s.Name).Max(5)
	result := v.Validate()
	assert.False(t, result.IsValid())

	s.Name = "12345"
	v, _ = validator.New(&s)
	validator.String(v, &s.Name).Max(5)
	result = v.Validate()
	assert.True(t, result.IsValid())
}

func TestString_Length(t *testing.T) {
	s := stringTestStruct{Name: "1234"}
	v, _ := validator.New(&s)
	validator.String(v, &s.Name).Length(5, 10)
	result := v.Validate()
	assert.False(t, result.IsValid())

	s.Name = "12345"
	v, _ = validator.New(&s)
	validator.String(v, &s.Name).Length(5, 10)
	result = v.Validate()
	assert.True(t, result.IsValid())

	s.Name = "12345678901"
	v, _ = validator.New(&s)
	validator.String(v, &s.Name).Length(5, 10)
	result = v.Validate()
	assert.False(t, result.IsValid())
}

func TestString_Min_CountsRunesNotBytes(t *testing.T) {
	// "héllo" is 5 runes but 6 bytes (é is 2 bytes). Counting runes, Max(5)
	// passes; counting bytes it would fail.
	s := stringTestStruct{Name: "héllo"}
	v, _ := validator.New(&s)
	validator.String(v, &s.Name).Min(5).Max(5)
	result := v.Validate()
	assert.True(t, result.IsValid())
}

func TestString_Min_SingularMessage(t *testing.T) {
	s := stringTestStruct{Name: ""}
	v, _ := validator.New(&s)
	validator.String(v, &s.Name).Min(1)
	result := v.Validate()
	assert.False(t, result.IsValid())
	assert.Equal(t, "must be at least 1 character", result.Errors[0].Message)
}

func TestString_WithName(t *testing.T) {
	s := stringTestStruct{Phone: nil}
	v, _ := validator.New(&s)
	validator.String(v, s.Phone).NotNil().WithName("phone")
	result := v.Validate()
	assert.False(t, result.IsValid())
	assert.Equal(t, "phone", result.Errors[0].Field)
}

func TestString_Email(t *testing.T) {
	s := stringTestStruct{Email: "test@example.com"}
	v, _ := validator.New(&s)
	validator.String(v, &s.Email).Email()
	result := v.Validate()
	assert.True(t, result.IsValid())

	s.Email = "not-an-email"
	v, _ = validator.New(&s)
	validator.String(v, &s.Email).Email()
	result = v.Validate()
	assert.False(t, result.IsValid())
	assert.Equal(t, "must be a valid email address", result.Errors[0].Message)
}

func TestString_URL(t *testing.T) {
	s := stringTestStruct{Name: "http://google.com"}
	v, _ := validator.New(&s)
	validator.String(v, &s.Name).URL()
	result := v.Validate()
	assert.True(t, result.IsValid())

	s.Name = "not-a-url"
	v, _ = validator.New(&s)
	validator.String(v, &s.Name).URL()
	result = v.Validate()
	assert.False(t, result.IsValid())
}

func TestString_Includes(t *testing.T) {
	s := stringTestStruct{Name: "X_test_X"}
	v, _ := validator.New(&s)
	validator.String(v, &s.Name).Includes("test")
	result := v.Validate()
	assert.True(t, result.IsValid())

	s.Name = "X_Test_X"
	v, _ = validator.New(&s)
	validator.String(v, &s.Name).Includes("test")
	result = v.Validate()
	assert.False(t, result.IsValid())
}

func TestString_StartsWith(t *testing.T) {
	s := stringTestStruct{Name: "test_X"}
	v, _ := validator.New(&s)
	validator.String(v, &s.Name).StartsWith("test")
	result := v.Validate()
	assert.True(t, result.IsValid())

	s.Name = "Test_X"
	v, _ = validator.New(&s)
	validator.String(v, &s.Name).StartsWith("test")
	result = v.Validate()
	assert.False(t, result.IsValid())
}

func TestString_EndsWith(t *testing.T) {
	s := stringTestStruct{Name: "X_test"}
	v, _ := validator.New(&s)
	validator.String(v, &s.Name).EndsWith("test")
	result := v.Validate()
	assert.True(t, result.IsValid())

	s.Name = "X_Test"
	v, _ = validator.New(&s)
	validator.String(v, &s.Name).EndsWith("test")
	result = v.Validate()
	assert.False(t, result.IsValid())
}

func TestString_NilSkip(t *testing.T) {
	s := stringTestStruct{Phone: nil}
	v, _ := validator.New(&s)
	validator.String(v, s.Phone).Min(5)
	result := v.Validate()
	assert.True(t, result.IsValid())
}

func TestString_NotNil(t *testing.T) {
	s := stringTestStruct{Phone: nil}
	v, _ := validator.New(&s)
	validator.String(v, s.Phone).NotNil()
	result := v.Validate()
	assert.False(t, result.IsValid())
	assert.Equal(t, "must not be nil", result.Errors[0].Message)

	phone := "1234567"
	s.Phone = &phone
	v, _ = validator.New(&s)
	validator.String(v, s.Phone).NotNil().Min(5)
	result = v.Validate()
	assert.True(t, result.IsValid())
}

func TestString_Must(t *testing.T) {
	s := stringTestStruct{Email: "user@other.com"}
	v, _ := validator.New(&s)
	validator.String(v, &s.Email).Must(func(val string) bool {
		return len(val) > 5 && val[len(val)-12:] == "@company.com"
	}).WithMessage("must be a company email")
	result := v.Validate()
	assert.False(t, result.IsValid())
	assert.Equal(t, "must be a company email", result.Errors[0].Message)

	s.Email = "user@company.com"
	v, _ = validator.New(&s)
	validator.String(v, &s.Email).Must(func(val string) bool {
		return len(val) > 12 && val[len(val)-12:] == "@company.com"
	}).WithMessage("must be a company email")
	result = v.Validate()
	assert.True(t, result.IsValid())
}

func TestString_WithMessage(t *testing.T) {
	s := stringTestStruct{Name: ""}
	v, _ := validator.New(&s)
	validator.String(v, &s.Name).NotEmpty().WithMessage("name is required")
	result := v.Validate()
	assert.False(t, result.IsValid())
	assert.Equal(t, "name is required", result.Errors[0].Message)
}

func TestString_MultipleRulesAllReported(t *testing.T) {
	s := stringTestStruct{Name: ""}
	v, _ := validator.New(&s)
	validator.String(v, &s.Name).NotEmpty().Min(5)
	result := v.Validate()
	assert.Len(t, result.Errors, 2)
}

func TestString_JsonTagFieldName(t *testing.T) {
	s := stringTestStruct{Name: ""}
	v, _ := validator.New(&s)
	validator.String(v, &s.Name).NotEmpty()
	result := v.Validate()
	assert.Equal(t, "name", result.Errors[0].Field)
}

func TestString_Standalone_Validate(t *testing.T) {
	r := rule.String().NotEmpty().Email()

	errs := r.Validate("test@example.com")
	assert.Empty(t, errs)

	errs = r.Validate("")
	assert.Len(t, errs, 2)

	errs = r.Validate(123)
	assert.Len(t, errs, 1)
	assert.Equal(t, "must be a string", errs[0].Message)

	errs = r.Validate(nil)
	assert.Empty(t, errs)
}

func TestString_NamedStringType(t *testing.T) {
	type code string
	assert.Empty(t, rule.String().Min(2).Validate(code("ab")))

	errs := rule.String().Min(2).Validate(code("a"))
	assert.Len(t, errs, 1)
	assert.Equal(t, "must be at least 2 characters", errs[0].Message)
}

func TestString_NotNil_WithMessage(t *testing.T) {
	s := stringTestStruct{Phone: nil}
	v, _ := validator.New(&s)
	validator.String(v, s.Phone).NotNil().WithMessage("phone is required")
	result := v.Validate()
	assert.False(t, result.IsValid())
	assert.Equal(t, "phone is required", result.Errors[0].Message)
}

func TestString_WithMessage_DoesNotBreakNotNil(t *testing.T) {
	// Ensure that using WithMessage("must not be nil") on a non-NotNil rule
	// does NOT cause it to be treated as a NotNil rule
	s := stringTestStruct{Name: "hello"}
	v, _ := validator.New(&s)
	validator.String(v, &s.Name).Min(10).WithMessage("must not be nil")
	result := v.Validate()
	assert.False(t, result.IsValid())
	assert.Equal(t, "must not be nil", result.Errors[0].Message)
}

func TestStruct_ValidatesCurrentState(t *testing.T) {
	s := stringTestStruct{Name: "ab"}

	v, _ := validator.New(&s)
	validator.String(v, &s.Name).Min(5)

	// Mutate the struct after building rules but before validating
	s.Name = "abcde"

	result := v.Validate()
	assert.True(t, result.IsValid())
}

func TestString_OneOf(t *testing.T) {
	s := stringTestStruct{Name: "green"}
	v, _ := validator.New(&s)
	validator.String(v, &s.Name).OneOf("red", "blue")
	result := v.Validate()
	assert.False(t, result.IsValid())
	assert.Equal(t, "must be one of \"red\", \"blue\"", result.Errors[0].Message)

	s.Name = "red"
	v, _ = validator.New(&s)
	validator.String(v, &s.Name).OneOf("red", "blue")
	assert.True(t, v.Validate().IsValid())
}

func TestString_Matches(t *testing.T) {
	s := stringTestStruct{Name: "PRJ-123"}
	v, _ := validator.New(&s)
	validator.String(v, &s.Name).Matches(`^PRJ-\d+$`)
	assert.True(t, v.Validate().IsValid())

	s.Name = "prj-123"
	v, _ = validator.New(&s)
	validator.String(v, &s.Name).Matches(`^PRJ-\d+$`)
	result := v.Validate()
	assert.False(t, result.IsValid())
	assert.Equal(t, `must match "^PRJ-\d+$"`, result.Errors[0].Message)
}

func TestString_Matches_InvalidPatternPanics(t *testing.T) {
	s := stringTestStruct{}
	v, _ := validator.New(&s)
	assert.Panics(t, func() { validator.String(v, &s.Name).Matches("(") })
}

func TestString_NamedStringType_StructMode(t *testing.T) {
	type status string
	type ticket struct {
		Status status `json:"status"`
	}

	// A named string field attaches without a cast, and Must receives the
	// named type.
	tk := ticket{Status: "archived"}
	v, _ := validator.New(&tk)
	validator.String(v, &tk.Status).OneOf("open", "closed").Must(func(s status) bool {
		return s != "archived"
	})
	result := v.Validate()
	assert.Len(t, result.Errors, 2)
	assert.Equal(t, "status", result.Errors[0].Field)
	assert.Equal(t, "must be one of \"open\", \"closed\"", result.Errors[0].Message)

	tk.Status = "open"
	v, _ = validator.New(&tk)
	validator.String(v, &tk.Status).OneOf("open", "closed")
	assert.True(t, v.Validate().IsValid())
}

func TestString_OneOf_MessageEdgeCases(t *testing.T) {
	// An empty allowlist (a spread of an empty slice) fails every value with
	// the generic message rather than a dangling "must be one of ".
	errs := rule.String().OneOf().Validate("x")
	assert.Len(t, errs, 1)
	assert.Equal(t, "is not valid", errs[0].Message)

	errs = rule.String().OneOf(`a"b`).Validate("x")
	assert.Len(t, errs, 1)
	assert.Equal(t, `must be one of "a\"b"`, errs[0].Message)
}

func TestString_Must_NilPanics(t *testing.T) {
	s := stringTestStruct{}
	v, _ := validator.New(&s)
	assert.Panics(t, func() { validator.String(v, &s.Name).Must(nil) })
}

func TestStruct_ValidatesCurrentState_Pointer(t *testing.T) {
	short := "ab"
	s := stringTestStruct{Phone: &short}

	v, _ := validator.New(&s)
	validator.String(v, s.Phone).Min(5)

	// Mutate the value through the pointer after building rules
	short = "abcdefg"

	result := v.Validate()
	assert.True(t, result.IsValid())
}
