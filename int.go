package schema

import "fmt"

type IntSchema struct {
	Schema[int]
}

var _ ISchema = (*IntSchema)(nil)

func Int() *IntSchema {
	return &IntSchema{}
}

func (s *IntSchema) Lt(min int) *IntSchema {
	validator := Validator[int]{
		MessageFunc: func(value int) string {
			return fmt.Sprintf("Int must be less than %d", min)
		},
		ValidateFunc: func(value int) bool {
			return value < min
		},
	}

	s.validators = append(s.validators, validator)

	return s
}

func (s *IntSchema) Gt(min int) *IntSchema {
	validator := Validator[int]{
		MessageFunc: func(value int) string {
			return fmt.Sprintf("Int must be greater than %d", min)
		},
		ValidateFunc: func(value int) bool {
			return value > min
		},
	}

	s.validators = append(s.validators, validator)

	return s
}

func (s *IntSchema) Parse(value any) *ValidationResult {
	val, ok := value.(int)
	if !ok {
		return &ValidationResult{Errors: []ValidationError{{Path: "", Message: fmt.Sprintf("Expected int, received %T", value)}}}
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
