package schema

import (
	"fmt"
	"reflect"
	"strconv"
)

type ArraySchema struct {
	Schema[[]interface{}]
	schema ISchema
}

var _ ISchema = (*ArraySchema)(nil)

func Array(s ISchema) *ArraySchema {
	return &ArraySchema{schema: s}
}

func (s *ArraySchema) Max(maxLength int) *ArraySchema {
	validator := Validator[[]interface{}]{
		MessageFunc: func(value []interface{}) string {
			return fmt.Sprintf("Array must contain at most %d element(s)", maxLength)
		},
		ValidateFunc: func(value []interface{}) bool {
			return len(value) <= maxLength
		},
	}

	s.validators = append(s.validators, validator)

	return s
}

func (s *ArraySchema) Min(minLength int) *ArraySchema {
	validator := Validator[[]interface{}]{
		MessageFunc: func(value []interface{}) string {
			return fmt.Sprintf("Array must contain at least %d element(s)", minLength)
		},
		ValidateFunc: func(value []interface{}) bool {
			return len(value) >= minLength
		},
	}

	s.validators = append(s.validators, validator)

	return s
}

func (s *ArraySchema) Parse(value any) *ValidationResult {
	t := reflect.TypeOf(value)

	if t == nil || (t.Kind() != reflect.Array && t.Kind() != reflect.Slice) {
		return &ValidationResult{Errors: []ValidationError{{Path: "", Message: fmt.Sprintf("Expected array, got %T", value)}}}
	}

	v := reflect.ValueOf(value)

	val := make([]interface{}, v.Len())
	for i := 0; i < v.Len(); i++ {
		val[i] = v.Index(i).Interface()
	}

	// Parse array validations
	result := &ValidationResult{Errors: []ValidationError{}}

	for _, validator := range s.validators {
		if !validator.ValidateFunc(val) {
			err := ValidationError{
				Path:    "",
				Message: validator.MessageFunc(val),
			}

			result.Errors = append(result.Errors, err)
		}
	}

	// Parse schema validations within array for each item
	for i := 0; i < len(val); i++ {
		res := s.schema.Parse(val[i])

		if !res.IsValid() {
			for index, err := range res.Errors {
				res.Errors[index].Path = formatPath(strconv.Itoa(i), err.Path)
			}

			result.Errors = append(result.Errors, res.Errors...)
		}
	}

	return result
}
