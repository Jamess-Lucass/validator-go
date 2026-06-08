package validator

import (
	"fmt"
	"net/mail"
	"net/url"
	"strings"
	"unicode/utf8"
)

type StringRule struct {
	fieldBase[string]
}

func (v *Validator) String(field *string) *StringRule {
	r := &StringRule{fieldBase[string]{v: v, fieldPtr: field}}
	v.addRule(r)
	return r
}

func String() *StringRule {
	return &StringRule{}
}

func (r *StringRule) Validate(value any) []ValidationError {
	return validateStandalone(r.rules, value, "must be a string", func(v any) (string, bool) {
		s, ok := v.(string)
		return s, ok
	})
}

// Presence rules

func (r *StringRule) NotNil() *StringRule {
	r.rules = append(r.rules, rule[string]{isNotNil: true, message: msgNotNil})
	return r
}

func (r *StringRule) NotEmpty() *StringRule {
	r.rules = append(r.rules, rule[string]{
		validate: func(val string) bool { return val != "" },
		message:  msgNotEmpty,
	})
	return r
}

// Length rules

func (r *StringRule) Min(n int) *StringRule {
	r.rules = append(r.rules, rule[string]{
		validate: func(val string) bool { return utf8.RuneCountInString(val) >= n },
		message:  fmt.Sprintf("must be at least %d %s", n, pluralize(n, "character")),
	})
	return r
}

func (r *StringRule) Max(n int) *StringRule {
	r.rules = append(r.rules, rule[string]{
		validate: func(val string) bool { return utf8.RuneCountInString(val) <= n },
		message:  fmt.Sprintf("must be at most %d %s", n, pluralize(n, "character")),
	})
	return r
}

func (r *StringRule) Length(min, max int) *StringRule {
	r.rules = append(r.rules, rule[string]{
		validate: func(val string) bool {
			n := utf8.RuneCountInString(val)
			return n >= min && n <= max
		},
		message: fmt.Sprintf("must be between %d and %d characters", min, max),
	})
	return r
}

// Format rules

func (r *StringRule) Email() *StringRule {
	r.rules = append(r.rules, rule[string]{
		validate: func(val string) bool {
			addr, err := mail.ParseAddress(val)
			return err == nil && addr.Address == val
		},
		message: "must be a valid email address",
	})
	return r
}

func (r *StringRule) Url() *StringRule {
	r.rules = append(r.rules, rule[string]{
		validate: func(val string) bool {
			u, err := url.Parse(val)
			return err == nil && u.Scheme != "" && u.Host != ""
		},
		message: "must be a valid URL",
	})
	return r
}

// Content rules

func (r *StringRule) Includes(substr string) *StringRule {
	r.rules = append(r.rules, rule[string]{
		validate: func(val string) bool { return strings.Contains(val, substr) },
		message:  fmt.Sprintf("must contain \"%s\"", substr),
	})
	return r
}

func (r *StringRule) StartsWith(prefix string) *StringRule {
	r.rules = append(r.rules, rule[string]{
		validate: func(val string) bool { return strings.HasPrefix(val, prefix) },
		message:  fmt.Sprintf("must start with \"%s\"", prefix),
	})
	return r
}

func (r *StringRule) EndsWith(suffix string) *StringRule {
	r.rules = append(r.rules, rule[string]{
		validate: func(val string) bool { return strings.HasSuffix(val, suffix) },
		message:  fmt.Sprintf("must end with \"%s\"", suffix),
	})
	return r
}

// Custom rules

func (r *StringRule) Must(fn func(string) bool) *StringRule {
	r.rules = append(r.rules, rule[string]{validate: fn, message: msgNotValid})
	return r
}

func (r *StringRule) WithMessage(msg string) *StringRule {
	if len(r.rules) > 0 {
		r.rules[len(r.rules)-1].message = msg
	}
	return r
}

func (r *StringRule) WithName(name string) *StringRule {
	r.name = name
	return r
}
