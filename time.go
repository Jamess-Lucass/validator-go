package validator

import (
	"time"
)

// TimeRule validates a time.Time field or value.
type TimeRule struct {
	fieldBase[time.Time]
}

// Time validates a time.Time field.
func Time(v *Validator, field *time.Time) *TimeRule {
	r := &TimeRule{fieldBase[time.Time]{v: v, fieldPtr: field}}
	v.addRule(r)
	return r
}

// Validate implements Rule for standalone use; the value must be a time.Time
// or an RFC3339 string.
func (r *TimeRule) Validate(value any) []ValidationError {
	return validateStandalone(r.rules, value, "must be a time", func(v any) (time.Time, bool) {
		switch val := v.(type) {
		case time.Time:
			return val, true
		case string:
			t, err := time.Parse(time.RFC3339, val)
			if err != nil {
				return time.Time{}, false
			}
			return t, true
		default:
			return time.Time{}, false
		}
	})
}

// NotNil fails when the field pointer or standalone value is nil.
func (r *TimeRule) NotNil() *TimeRule {
	r.rules = append(r.rules, rule[time.Time]{isNotNil: true, message: msgNotNil})
	return r
}

// NotEmpty fails on the zero time.
func (r *TimeRule) NotEmpty() *TimeRule {
	r.rules = append(r.rules, rule[time.Time]{
		validate: func(val time.Time) bool { return !val.IsZero() },
		message:  msgNotEmpty,
	})
	return r
}

// After requires the value to be after t.
func (r *TimeRule) After(t time.Time) *TimeRule {
	r.rules = append(r.rules, rule[time.Time]{
		validate: func(val time.Time) bool { return val.After(t) },
		message:  "must be after " + t.Format(time.RFC3339),
	})
	return r
}

// Before requires the value to be before t.
func (r *TimeRule) Before(t time.Time) *TimeRule {
	r.rules = append(r.rules, rule[time.Time]{
		validate: func(val time.Time) bool { return val.Before(t) },
		message:  "must be before " + t.Format(time.RFC3339),
	})
	return r
}

// Must runs a custom check against the value. It panics if fn is nil.
func (r *TimeRule) Must(fn func(time.Time) bool) *TimeRule {
	if fn == nil {
		panic("validator: Must requires a non-nil function")
	}
	r.rules = append(r.rules, rule[time.Time]{validate: fn, message: msgNotValid})
	return r
}

// WithMessage replaces the previous rule's error message.
func (r *TimeRule) WithMessage(msg string) *TimeRule {
	if len(r.rules) > 0 {
		r.rules[len(r.rules)-1].message = msg
	}
	return r
}

// WithName overrides the field name used in errors.
func (r *TimeRule) WithName(name string) *TimeRule {
	r.name = name
	return r
}
