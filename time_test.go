package validator_test

import (
	"testing"
	"time"

	validator "github.com/Jamess-Lucass/validator-go"
	"github.com/Jamess-Lucass/validator-go/rule"
	"github.com/stretchr/testify/assert"
)

type timeTestStruct struct {
	ExpiresAt   time.Time  `json:"expires_at"`
	ScheduledAt *time.Time `json:"scheduled_at"`
}

func TestTime_After(t *testing.T) {
	past := time.Now().Add(-1 * time.Hour)
	s := timeTestStruct{ExpiresAt: past}
	now := time.Now()
	v, _ := validator.New(&s)
	validator.Time(v, &s.ExpiresAt).After(now)
	result := v.Validate()
	assert.False(t, result.IsValid())

	s.ExpiresAt = time.Now().Add(1 * time.Hour)
	v, _ = validator.New(&s)
	validator.Time(v, &s.ExpiresAt).After(now)
	result = v.Validate()
	assert.True(t, result.IsValid())
}

func TestTime_Before(t *testing.T) {
	future := time.Now().Add(1 * time.Hour)
	s := timeTestStruct{ExpiresAt: future}
	now := time.Now()
	v, _ := validator.New(&s)
	validator.Time(v, &s.ExpiresAt).Before(now)
	result := v.Validate()
	assert.False(t, result.IsValid())

	s.ExpiresAt = time.Now().Add(-1 * time.Hour)
	v, _ = validator.New(&s)
	validator.Time(v, &s.ExpiresAt).Before(now)
	result = v.Validate()
	assert.True(t, result.IsValid())
}

func TestTime_NotEmpty(t *testing.T) {
	s := timeTestStruct{}
	v, _ := validator.New(&s)
	validator.Time(v, &s.ExpiresAt).NotEmpty()
	result := v.Validate()
	assert.False(t, result.IsValid())

	s.ExpiresAt = time.Now()
	v, _ = validator.New(&s)
	validator.Time(v, &s.ExpiresAt).NotEmpty()
	result = v.Validate()
	assert.True(t, result.IsValid())
}

func TestTime_NilSkip(t *testing.T) {
	s := timeTestStruct{ScheduledAt: nil}
	v, _ := validator.New(&s)
	validator.Time(v, s.ScheduledAt).After(time.Now())
	result := v.Validate()
	assert.True(t, result.IsValid())
}

func TestTime_NotNil(t *testing.T) {
	s := timeTestStruct{ScheduledAt: nil}
	v, _ := validator.New(&s)
	validator.Time(v, s.ScheduledAt).NotNil().WithName("scheduled_at")
	result := v.Validate()
	assert.False(t, result.IsValid())
	assert.Equal(t, "scheduled_at", result.Errors[0].Field)
}

func TestTime_Must(t *testing.T) {
	cutoff := time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)
	s := timeTestStruct{ExpiresAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)}
	v, _ := validator.New(&s)
	validator.Time(v, &s.ExpiresAt).Must(func(val time.Time) bool {
		return val.After(cutoff)
	}).WithMessage("must be after 2022")
	result := v.Validate()
	assert.False(t, result.IsValid())
	assert.Equal(t, "must be after 2022", result.Errors[0].Message)
}

func TestTime_Standalone(t *testing.T) {
	cutoff := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	r := rule.Time().NotEmpty().After(cutoff)

	errs := r.Validate(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC))
	assert.Empty(t, errs)

	// RFC3339 strings are parsed.
	errs = r.Validate("2024-06-01T00:00:00Z")
	assert.Empty(t, errs)

	errs = r.Validate("2020-01-01T00:00:00Z")
	assert.Len(t, errs, 1)

	// Non-RFC3339 string.
	errs = r.Validate("01/01/2024")
	assert.Len(t, errs, 1)
	assert.Equal(t, "must be a time", errs[0].Message)

	// Wrong type.
	errs = r.Validate(123)
	assert.Len(t, errs, 1)
	assert.Equal(t, "must be a time", errs[0].Message)
}

func TestTime_Map(t *testing.T) {
	m := map[string]any{"created": "2024-01-01T00:00:00Z"}
	v := validator.NewObject(m)
	v.Field("created", rule.Time().NotEmpty().Before(time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)))
	result := v.Validate()
	assert.False(t, result.IsValid())
	assert.Equal(t, "created", result.Errors[0].Field)
}
