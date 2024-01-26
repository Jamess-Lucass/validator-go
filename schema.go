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

func (s *Schema[T]) Parse(value any) *ValidationResult {
	val, ok := value.(T)
	if !ok {
		return &ValidationResult{Errors: []ValidationError{{Path: "", Message: fmt.Sprintf("Expected %s, received %T", reflect.TypeOf(val).String(), value)}}}
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
