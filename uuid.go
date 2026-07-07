package validator

import (
	"github.com/google/uuid"
)

// UUIDRule validates a github.com/google/uuid value.
type UUIDRule struct {
	fieldBase[uuid.UUID]
}

// UUID validates a uuid.UUID field.
func UUID(v *Validator, field *uuid.UUID) *UUIDRule {
	r := &UUIDRule{fieldBase[uuid.UUID]{v: v, fieldPtr: field}}
	v.addRule(r)
	return r
}

// Validate implements Rule for standalone use; the value must be a uuid.UUID,
// a [16]byte, or a parseable string.
func (r *UUIDRule) Validate(value any) []ValidationError {
	return validateStandalone(r.rules, value, "must be a UUID", func(v any) (uuid.UUID, bool) {
		switch val := v.(type) {
		case uuid.UUID:
			return val, true
		case [16]byte:
			return val, true
		case string:
			id, err := uuid.Parse(val)
			if err != nil {
				return uuid.Nil, false
			}
			return id, true
		default:
			return uuid.Nil, false
		}
	})
}

// Presence rules

// NotNil fails when the field pointer or standalone value is nil.
func (r *UUIDRule) NotNil() *UUIDRule {
	r.rules = append(r.rules, rule[uuid.UUID]{isNotNil: true, message: msgNotNil})
	return r
}

// NotEmpty fails on the nil UUID (all zeros).
func (r *UUIDRule) NotEmpty() *UUIDRule {
	r.rules = append(r.rules, rule[uuid.UUID]{
		validate: func(val uuid.UUID) bool { return val != uuid.Nil },
		message:  msgNotEmpty,
	})
	return r
}

// Custom rules

// Must runs a custom check against the value. It panics if fn is nil.
func (r *UUIDRule) Must(fn func(uuid.UUID) bool) *UUIDRule {
	if fn == nil {
		panic("validator: Must requires a non-nil function")
	}
	r.rules = append(r.rules, rule[uuid.UUID]{validate: fn, message: msgNotValid})
	return r
}

// WithMessage replaces the previous rule's error message.
func (r *UUIDRule) WithMessage(msg string) *UUIDRule {
	if len(r.rules) > 0 {
		r.rules[len(r.rules)-1].message = msg
	}
	return r
}

// WithName overrides the field name used in errors.
func (r *UUIDRule) WithName(name string) *UUIDRule {
	r.name = name
	return r
}
