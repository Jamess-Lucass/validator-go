package schema_test

import (
	"testing"

	schema "github.com/Jamess-Lucass/validator-go"
	"github.com/stretchr/testify/assert"
)

func TestString_Type(t *testing.T) {
	s := schema.String()

	assert.True(t, s.Parse("string").IsValid())

	assert.False(t, s.Parse(123).IsValid())
	assert.False(t, s.Parse(nil).IsValid())
	assert.False(t, s.Parse(map[string]int{
		"one": 1,
		"two": 2,
	}).IsValid())
	assert.False(t, s.Parse([]int{1, 2, 3}).IsValid())
	assert.False(t, s.Parse(0).IsValid())
}

func TestString_Path(t *testing.T) {
	s := schema.String().Min(4).Parse("123")

	assert.Len(t, s.Errors, 1)
	assert.Equal(t, "", s.Errors[0].Path)
}

func TestString_Max(t *testing.T) {
	s := schema.String().Max(5)

	assert.True(t, s.Parse("12345").IsValid())
	assert.True(t, s.Parse("1234").IsValid())

	assert.False(t, s.Parse("123456").IsValid())
}

func TestString_Min(t *testing.T) {
	s := schema.String().Min(5)

	assert.True(t, s.Parse("12345").IsValid())
	assert.True(t, s.Parse("123456").IsValid())

	assert.False(t, s.Parse("1234").IsValid())
}

func TestString_Length(t *testing.T) {
	s := schema.String().Length(5)

	assert.True(t, s.Parse("12345").IsValid())

	assert.False(t, s.Parse("123456").IsValid())
	assert.False(t, s.Parse("1234").IsValid())
}

func TestString_Url(t *testing.T) {
	s := schema.String().Url()

	assert.True(t, s.Parse("http://google.com").IsValid())
	assert.True(t, s.Parse("https://google.com/asdf?asdf=ljk3lk4&asdf=234#asdf").IsValid())

	assert.False(t, s.Parse("asdf").IsValid())
	assert.False(t, s.Parse("https:/").IsValid())
	assert.False(t, s.Parse("https").IsValid())
	assert.False(t, s.Parse("asdfj@lkjsdf.com").IsValid())
}

func TestString_Includes(t *testing.T) {
	s := schema.String().Includes("test")

	assert.True(t, s.Parse("X_test_X").IsValid())
	assert.True(t, s.Parse("test").IsValid())

	assert.False(t, s.Parse("Test").IsValid())
	assert.False(t, s.Parse("X_Test_X").IsValid())
	assert.False(t, s.Parse("TEST").IsValid())
	assert.False(t, s.Parse("3t3est").IsValid())
}

func TestString_StartsWith(t *testing.T) {
	s := schema.String().StartsWith("test")

	assert.True(t, s.Parse("test_X").IsValid())
	assert.True(t, s.Parse("test").IsValid())

	assert.False(t, s.Parse("Test").IsValid())
	assert.False(t, s.Parse("Test_X").IsValid())
	assert.False(t, s.Parse("TEST").IsValid())
	assert.False(t, s.Parse("teslt3").IsValid())
}

func TestString_EndsWith(t *testing.T) {
	s := schema.String().EndsWith("test")

	assert.True(t, s.Parse("X_test").IsValid())
	assert.True(t, s.Parse("test").IsValid())

	assert.False(t, s.Parse("Test").IsValid())
	assert.False(t, s.Parse("X_Test").IsValid())
	assert.False(t, s.Parse("TEST").IsValid())
	assert.False(t, s.Parse("3tes3t").IsValid())
}
