package schema

import (
	"fmt"
	"net/url"
	"strings"
)

type StringSchema struct {
	Schema[string]
}

var _ ISchema = (*StringSchema)(nil)

func String() *StringSchema {
	return &StringSchema{}
}

func (s *StringSchema) Max(maxLength int) *StringSchema {
	validator := Validator[string]{
		MessageFunc: func(value string) string {
			return fmt.Sprintf("String must contain at most %d character(s)", maxLength)
		},
		ValidateFunc: func(value string) bool {
			return len(value) <= maxLength
		},
	}

	s.validators = append(s.validators, validator)

	return s
}

func (s *StringSchema) Min(minLength int) *StringSchema {
	validator := Validator[string]{
		MessageFunc: func(value string) string {
			return fmt.Sprintf("String must contain at least %d character(s)", minLength)
		},
		ValidateFunc: func(value string) bool {
			return len(value) >= minLength
		},
	}

	s.validators = append(s.validators, validator)

	return s
}

func (s *StringSchema) Length(length int) *StringSchema {
	validator := Validator[string]{
		MessageFunc: func(value string) string {
			return fmt.Sprintf("String must contain exactly %d character(s)", length)
		},
		ValidateFunc: func(value string) bool {
			return len(value) == length
		},
	}

	s.validators = append(s.validators, validator)

	return s
}

func (s *StringSchema) Url() *StringSchema {
	validator := Validator[string]{
		MessageFunc: func(value string) string {
			return "Invalid url"
		},
		ValidateFunc: func(value string) bool {
			uri, err := url.ParseRequestURI(value)
			return err == nil && uri.Host != ""
		},
	}

	s.validators = append(s.validators, validator)

	return s
}

func (s *StringSchema) Includes(str string) *StringSchema {
	validator := Validator[string]{
		MessageFunc: func(value string) string {
			return fmt.Sprintf("Invalid input: must include \"%s\"", str)
		},
		ValidateFunc: func(value string) bool {
			return strings.Contains(value, str)
		},
	}

	s.validators = append(s.validators, validator)

	return s
}

func (s *StringSchema) StartsWith(str string) *StringSchema {
	validator := Validator[string]{
		MessageFunc: func(value string) string {
			return fmt.Sprintf("Invalid input: must start with \"%s\"", str)
		},
		ValidateFunc: func(value string) bool {
			return strings.HasPrefix(value, str)
		},
	}

	s.validators = append(s.validators, validator)

	return s
}

func (s *StringSchema) EndsWith(str string) *StringSchema {
	validator := Validator[string]{
		MessageFunc: func(value string) string {
			return fmt.Sprintf("Invalid input: must end with \"%s\"", str)
		},
		ValidateFunc: func(value string) bool {
			return strings.HasSuffix(value, str)
		},
	}

	s.validators = append(s.validators, validator)

	return s
}

func (s *StringSchema) Parse(value any) *ValidationResult {
	val, ok := value.(string)
	if !ok {
		return &ValidationResult{Errors: []ValidationError{{Path: "", Message: fmt.Sprintf("Expected string, received %T", value)}}}
	}

	res := &ValidationResult{}

	for _, validator := range s.validators {
		if !validator.ValidateFunc(val) {
			err := ValidationError{
				Path:    "",
				Message: validator.MessageFunc(val),
			}

			res.Errors = append(res.Errors, err)
		}
	}

	return res
}
