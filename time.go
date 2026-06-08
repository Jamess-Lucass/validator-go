package validator

import (
	"time"
)

type TimeRule struct {
	fieldBase[time.Time]
}

func (v *Validator) Time(field *time.Time) *TimeRule {
	r := &TimeRule{fieldBase[time.Time]{v: v, fieldPtr: field}}
	v.addRule(r)
	return r
}

func Time() *TimeRule {
	return &TimeRule{}
}

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

func (r *TimeRule) NotNil() *TimeRule {
	r.rules = append(r.rules, rule[time.Time]{isNotNil: true, message: msgNotNil})
	return r
}

func (r *TimeRule) NotEmpty() *TimeRule {
	r.rules = append(r.rules, rule[time.Time]{
		validate: func(val time.Time) bool { return !val.IsZero() },
		message:  msgNotEmpty,
	})
	return r
}

func (r *TimeRule) After(t time.Time) *TimeRule {
	r.rules = append(r.rules, rule[time.Time]{
		validate: func(val time.Time) bool { return val.After(t) },
		message:  "must be after " + t.Format(time.RFC3339),
	})
	return r
}

func (r *TimeRule) Before(t time.Time) *TimeRule {
	r.rules = append(r.rules, rule[time.Time]{
		validate: func(val time.Time) bool { return val.Before(t) },
		message:  "must be before " + t.Format(time.RFC3339),
	})
	return r
}

func (r *TimeRule) Must(fn func(time.Time) bool) *TimeRule {
	r.rules = append(r.rules, rule[time.Time]{validate: fn, message: msgNotValid})
	return r
}

func (r *TimeRule) WithMessage(msg string) *TimeRule {
	if len(r.rules) > 0 {
		r.rules[len(r.rules)-1].message = msg
	}
	return r
}

func (r *TimeRule) WithName(name string) *TimeRule {
	r.name = name
	return r
}
