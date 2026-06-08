package validator_test

import (
	"testing"

	validator "github.com/Jamess-Lucass/validator-go"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type uuidTestStruct struct {
	ID      uuid.UUID  `json:"id"`
	OwnerID *uuid.UUID `json:"owner_id"`
}

func TestUUID_NotEmpty(t *testing.T) {
	s := uuidTestStruct{ID: uuid.Nil}
	v, _ := validator.New(&s)
	v.UUID(&s.ID).NotEmpty()
	result := v.Validate()
	assert.False(t, result.IsValid())
	assert.Equal(t, "id", result.Errors[0].Field)
	assert.Equal(t, "must not be empty", result.Errors[0].Message)

	s.ID = uuid.New()
	v, _ = validator.New(&s)
	v.UUID(&s.ID).NotEmpty()
	result = v.Validate()
	assert.True(t, result.IsValid())
}

func TestUUID_NilSkip(t *testing.T) {
	s := uuidTestStruct{OwnerID: nil}
	v, _ := validator.New(&s)
	v.UUID(s.OwnerID).NotEmpty()
	result := v.Validate()
	assert.True(t, result.IsValid())
}

func TestUUID_NotNil(t *testing.T) {
	s := uuidTestStruct{OwnerID: nil}
	v, _ := validator.New(&s)
	v.UUID(s.OwnerID).NotNil().WithName("owner_id")
	result := v.Validate()
	assert.False(t, result.IsValid())
	assert.Equal(t, "owner_id", result.Errors[0].Field)
	assert.Equal(t, "must not be nil", result.Errors[0].Message)
}

func TestUUID_Must(t *testing.T) {
	id := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	s := uuidTestStruct{ID: id}
	v, _ := validator.New(&s)
	v.UUID(&s.ID).Must(func(val uuid.UUID) bool {
		return val.Version() == 4
	}).WithMessage("must be a v4 UUID")
	result := v.Validate()
	assert.False(t, result.IsValid())
	assert.Equal(t, "must be a v4 UUID", result.Errors[0].Message)
}

func TestUUID_Standalone(t *testing.T) {
	rule := validator.UUID().NotEmpty()

	errs := rule.Validate(uuid.New())
	assert.Empty(t, errs)

	// Valid UUID string is parsed and accepted.
	errs = rule.Validate("550e8400-e29b-41d4-a716-446655440000")
	assert.Empty(t, errs)

	// Raw 16-byte array.
	errs = rule.Validate([16]byte{1})
	assert.Empty(t, errs)

	// All-zero array is the nil UUID, so NotEmpty fails.
	errs = rule.Validate([16]byte{})
	assert.Len(t, errs, 1)
	assert.Equal(t, "must not be empty", errs[0].Message)

	// Invalid string.
	errs = rule.Validate("not-a-uuid")
	assert.Len(t, errs, 1)
	assert.Equal(t, "must be a UUID", errs[0].Message)

	// Wrong type.
	errs = rule.Validate(123)
	assert.Len(t, errs, 1)
	assert.Equal(t, "must be a UUID", errs[0].Message)

	// nil with no NotNil rule is skipped.
	errs = rule.Validate(nil)
	assert.Empty(t, errs)
}

func TestUUID_Map(t *testing.T) {
	m := map[string]any{"id": "not-a-uuid"}
	v := validator.NewMap(m)
	v.Field("id", validator.UUID())
	result := v.Validate()
	assert.False(t, result.IsValid())
	assert.Equal(t, "id", result.Errors[0].Field)
	assert.Equal(t, "must be a UUID", result.Errors[0].Message)
}
