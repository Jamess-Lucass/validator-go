package schema

import (
	"fmt"
)

type LiteralSchema[T any] struct {
	Schema[T]
}

var _ ISchema = (*LiteralSchema[interface{}])(nil)

func Literal[T any](value T) *LiteralSchema[T] {
	validator := Validator[T]{
		MessageFunc: func(val T) string {
			return fmt.Sprintf("Invalid literal value, expected \"%v\"", value)
		},
		ValidateFunc: func(val T) bool {
			return any(value) == any(val)
		},
	}

	return &LiteralSchema[T]{
		Schema: Schema[T]{
			validators: []Validator[T]{
				validator,
			},
		},
	}
}
