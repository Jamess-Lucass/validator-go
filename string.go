package schema

import "fmt"

type StringSchema struct {
	Schema[string]
}

var _ ISchema = (*StringSchema)(nil)

func String() *StringSchema {
	return &StringSchema{}
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
