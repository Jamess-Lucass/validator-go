package schema

import (
	"fmt"
	"reflect"
	"strconv"
)

type ArraySchema[T any] struct {
	Schema[[]any]
	schema ISchema
}

var _ ISchema = (*ArraySchema[interface{}])(nil)

func Array[T any](s ISchema) *ArraySchema[[]T] {
	return &ArraySchema[[]T]{schema: s}
}

func (s *ArraySchema[T]) Refine(predicate func(T) bool) *ArraySchema[T] {
	validator := Validator[T]{
		MessageFunc: func(value T) string {
			return "Invalid input"
		},
		ValidateFunc: predicate,
	}

	s.validators = append(s.validators, validator)

	return s
}

func (s *ArraySchema[T]) Max(maxLength int) *ArraySchema[T] {
	validator := Validator[T]{
		MessageFunc: func(value T) string {
			return fmt.Sprintf("Array must contain at most %d element(s)", maxLength)
		},
		ValidateFunc: func(value T) bool {
			v := reflect.ValueOf(value)
			return v.Len() <= maxLength
		},
	}

	s.validators = append(s.validators, validator)

	return s
}

func (s *ArraySchema[T]) Min(minLength int) *ArraySchema[T] {
	validator := Validator[T]{
		MessageFunc: func(value T) string {
			return fmt.Sprintf("Array must contain at least %d element(s)", minLength)
		},
		ValidateFunc: func(value T) bool {
			v := reflect.ValueOf(value)
			return v.Len() >= minLength
		},
	}

	s.Schema.validators = append(s.validators, validator)

	return s
}

func (s *ArraySchema[T]) Parse(value any) *ValidationResult {
	v := reflect.ValueOf(value)

	if v.Kind() != reflect.Array && v.Kind() != reflect.Slice {
		return &ValidationResult{Errors: []ValidationError{{Path: "", Message: fmt.Sprintf("Expected array, got %T", value)}}}
	}

	// Parse array validations
	result := &ValidationResult{Errors: []ValidationError{}}

	for _, validator := range s.validators {
		if !validator.ValidateFunc(value) {
			err := ValidationError{
				Path:    "",
				Message: validator.MessageFunc(value),
			}

			result.Errors = append(result.Errors, err)
		}
	}

	for i := 0; i < v.Len(); i++ {
		item := v.Index(i).Interface()
		res := s.schema.Parse(item)

		if !res.IsValid() {
			for index, err := range res.Errors {
				res.Errors[index].Path = formatPath(strconv.Itoa(i), err.Path)
			}
			result.Errors = append(result.Errors, res.Errors...)
		}
	}

	return result
}
