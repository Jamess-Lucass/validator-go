package validator

import (
	"github.com/google/uuid"
)

type UUIDRule struct {
	fieldBase[uuid.UUID]
}

func (v *Validator) UUID(field *uuid.UUID) *UUIDRule {
	r := &UUIDRule{fieldBase[uuid.UUID]{v: v, fieldPtr: field}}
	v.addRule(r)
	return r
}

func UUID() *UUIDRule {
	return &UUIDRule{}
}

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

func (r *UUIDRule) NotNil() *UUIDRule {
	r.rules = append(r.rules, rule[uuid.UUID]{isNotNil: true, message: msgNotNil})
	return r
}

// NotEmpty fails on the nil UUID (all zeros), e.g. an unset identifier.
func (r *UUIDRule) NotEmpty() *UUIDRule {
	r.rules = append(r.rules, rule[uuid.UUID]{
		validate: func(val uuid.UUID) bool { return val != uuid.Nil },
		message:  msgNotEmpty,
	})
	return r
}

// Custom rules

func (r *UUIDRule) Must(fn func(uuid.UUID) bool) *UUIDRule {
	r.rules = append(r.rules, rule[uuid.UUID]{validate: fn, message: msgNotValid})
	return r
}

func (r *UUIDRule) WithMessage(msg string) *UUIDRule {
	if len(r.rules) > 0 {
		r.rules[len(r.rules)-1].message = msg
	}
	return r
}

func (r *UUIDRule) WithName(name string) *UUIDRule {
	r.name = name
	return r
}
