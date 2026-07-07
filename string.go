package validator

import (
	"fmt"
	"net/mail"
	"net/url"
	"reflect"
	"regexp"
	"strings"
	"unicode/utf8"
)

// StringRule validates a string field or value. Named string types
// (type ID string) work too.
type StringRule[T ~string] struct {
	fieldBase[T]
}

// String validates a string field, plain or named:
//
//	validator.String(v, &s.Name).NotEmpty().Min(2)
func String[T ~string](v *Validator, field *T) *StringRule[T] {
	r := &StringRule[T]{fieldBase[T]{v: v, fieldPtr: field}}
	v.addRule(r)
	return r
}

// Validate implements Rule for standalone use; the value must be a string or
// a named string type.
func (r *StringRule[T]) Validate(value any) []ValidationError {
	return validateStandalone(r.rules, value, "must be a string", func(v any) (T, bool) {
		if s, ok := v.(T); ok {
			return s, true
		}
		if rv := reflect.ValueOf(v); rv.Kind() == reflect.String {
			return T(rv.String()), true
		}
		return "", false
	})
}

// Presence rules

// NotNil fails when the field pointer or standalone value is nil.
func (r *StringRule[T]) NotNil() *StringRule[T] {
	r.rules = append(r.rules, rule[T]{isNotNil: true, message: msgNotNil})
	return r
}

// NotEmpty fails on "".
func (r *StringRule[T]) NotEmpty() *StringRule[T] {
	r.rules = append(r.rules, rule[T]{
		validate: func(val T) bool { return val != "" },
		message:  msgNotEmpty,
	})
	return r
}

// Length rules

// Min requires at least n characters, counting runes rather than bytes.
func (r *StringRule[T]) Min(n int) *StringRule[T] {
	r.rules = append(r.rules, rule[T]{
		validate: func(val T) bool { return utf8.RuneCountInString(string(val)) >= n },
		message:  fmt.Sprintf("must be at least %d %s", n, pluralise(n, "character")),
	})
	return r
}

// Max allows at most n characters, counting runes rather than bytes.
func (r *StringRule[T]) Max(n int) *StringRule[T] {
	r.rules = append(r.rules, rule[T]{
		validate: func(val T) bool { return utf8.RuneCountInString(string(val)) <= n },
		message:  fmt.Sprintf("must be at most %d %s", n, pluralise(n, "character")),
	})
	return r
}

// Length requires between min and max characters, counting runes.
func (r *StringRule[T]) Length(min, max int) *StringRule[T] {
	r.rules = append(r.rules, rule[T]{
		validate: func(val T) bool {
			n := utf8.RuneCountInString(string(val))
			return n >= min && n <= max
		},
		message: fmt.Sprintf("must be between %d and %d characters", min, max),
	})
	return r
}

// Format rules

// Email requires a valid email address.
func (r *StringRule[T]) Email() *StringRule[T] {
	r.rules = append(r.rules, rule[T]{
		validate: func(val T) bool {
			addr, err := mail.ParseAddress(string(val))
			return err == nil && addr.Address == string(val)
		},
		message: "must be a valid email address",
	})
	return r
}

// URL requires an absolute URL with a scheme and host.
func (r *StringRule[T]) URL() *StringRule[T] {
	r.rules = append(r.rules, rule[T]{
		validate: func(val T) bool {
			u, err := url.Parse(string(val))
			return err == nil && u.Scheme != "" && u.Host != ""
		},
		message: "must be a valid URL",
	})
	return r
}

// Matches requires the value to match the regular expression. The pattern is
// compiled once when the rule is declared; like regexp.MustCompile, an
// invalid pattern panics.
func (r *StringRule[T]) Matches(pattern string) *StringRule[T] {
	re := regexp.MustCompile(pattern)
	r.rules = append(r.rules, rule[T]{
		validate: func(val T) bool { return re.MatchString(string(val)) },
		message:  fmt.Sprintf("must match \"%s\"", pattern),
	})
	return r
}

// Content rules

// Includes requires the substring to be present.
func (r *StringRule[T]) Includes(substr string) *StringRule[T] {
	r.rules = append(r.rules, rule[T]{
		validate: func(val T) bool { return strings.Contains(string(val), substr) },
		message:  fmt.Sprintf("must contain \"%s\"", substr),
	})
	return r
}

// StartsWith requires the given prefix.
func (r *StringRule[T]) StartsWith(prefix string) *StringRule[T] {
	r.rules = append(r.rules, rule[T]{
		validate: func(val T) bool { return strings.HasPrefix(string(val), prefix) },
		message:  fmt.Sprintf("must start with \"%s\"", prefix),
	})
	return r
}

// EndsWith requires the given suffix.
func (r *StringRule[T]) EndsWith(suffix string) *StringRule[T] {
	r.rules = append(r.rules, rule[T]{
		validate: func(val T) bool { return strings.HasSuffix(string(val), suffix) },
		message:  fmt.Sprintf("must end with \"%s\"", suffix),
	})
	return r
}

// OneOf requires the value to be one of the given values, e.g. the members of
// a string enum.
func (r *StringRule[T]) OneOf(values ...T) *StringRule[T] {
	r.rules = append(r.rules, rule[T]{
		validate: func(val T) bool {
			for _, v := range values {
				if val == v {
					return true
				}
			}
			return false
		},
		message: msgOneOf(values, true),
	})
	return r
}

// Custom rules

// Must runs a custom check against the value. It panics if fn is nil.
func (r *StringRule[T]) Must(fn func(T) bool) *StringRule[T] {
	if fn == nil {
		panic("validator: Must requires a non-nil function")
	}
	r.rules = append(r.rules, rule[T]{validate: fn, message: msgNotValid})
	return r
}

// WithMessage replaces the previous rule's error message.
func (r *StringRule[T]) WithMessage(msg string) *StringRule[T] {
	if len(r.rules) > 0 {
		r.rules[len(r.rules)-1].message = msg
	}
	return r
}

// WithName overrides the field name used in errors.
func (r *StringRule[T]) WithName(name string) *StringRule[T] {
	r.name = name
	return r
}
