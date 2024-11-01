package schema

import (
	"fmt"
	"reflect"
)

type ValidationResult struct {
	Errors []ValidationError
}

func (v *ValidationResult) IsValid() bool {
	return len(v.Errors) == 0
}

type ValidationError struct {
	Path    string `json:"path"`
	Message string `json:"message"`
}

func (m *ValidationError) Error() string {
	return m.Message
}

type ISchema interface {
	Parse(value any) *ValidationResult
}

type Schema[T any] struct {
	validators []Validator[T]
	isOptional bool
}

type Validator[T any] struct {
	MessageFunc  func(T) string
	ValidateFunc func(T) bool
}

func (s *Schema[T]) Refine(predicate func(T) bool) *Schema[T] {
	validator := Validator[T]{
		MessageFunc: func(value T) string {
			return "Invalid input"
		},
		ValidateFunc: predicate,
	}

	s.validators = append(s.validators, validator)

	return s
}

func (s *Schema[T]) Optional() *Schema[T] {
	s.isOptional = true
	return s
}

func (s *Schema[T]) Parse(value any) *ValidationResult {
	val, ok := value.(T)
	ptrVal, ptrOk := value.(*T)

	if !ptrOk && !ok {
		return &ValidationResult{Errors: []ValidationError{{Path: "", Message: fmt.Sprintf("Expected %s, received %T", reflect.TypeOf((*T)(nil)).Elem().String(), value)}}}
	}
	if ptrOk && ptrVal == nil && !s.isOptional {
		return &ValidationResult{Errors: []ValidationError{{Path: "", Message: "Value cannot be nil"}}}
	}

	res := &ValidationResult{Errors: []ValidationError{}}

	if s.isOptional && ptrVal == nil && ptrOk {
		return res
	}
	if !ok && ptrOk {
		val = *ptrVal
	}

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

func formatPath(key string, path string) string {
	if path != "" {
		return fmt.Sprintf("%s.%s", key, path)
	}

	return key
}
