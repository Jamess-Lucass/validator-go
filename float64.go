package schema

import (
	"fmt"
	"math"
)

type Float64Schema struct {
	Schema[float64]
}

var _ ISchema = (*Float64Schema)(nil)

func Float64() *Float64Schema {
	return &Float64Schema{}
}

func (s *Float64Schema) Lt(min float64) *Float64Schema {
	validator := Validator[float64]{
		MessageFunc: func(value float64) string {
			return fmt.Sprintf("float64 must be less than %f", min)
		},
		ValidateFunc: func(value float64) bool {
			return value < min
		},
	}

	s.validators = append(s.validators, validator)

	return s
}

func (s *Float64Schema) Lte(min float64) *Float64Schema {
	validator := Validator[float64]{
		MessageFunc: func(value float64) string {
			return fmt.Sprintf("float64 must be less than or equal to %f", min)
		},
		ValidateFunc: func(value float64) bool {
			return value <= min
		},
	}

	s.validators = append(s.validators, validator)

	return s
}

func (s *Float64Schema) Gt(min float64) *Float64Schema {
	validator := Validator[float64]{
		MessageFunc: func(value float64) string {
			return fmt.Sprintf("float64 must be greater than %f", min)
		},
		ValidateFunc: func(value float64) bool {
			return value > min
		},
	}

	s.validators = append(s.validators, validator)

	return s
}

func (s *Float64Schema) Gte(min float64) *Float64Schema {
	validator := Validator[float64]{
		MessageFunc: func(value float64) string {
			return fmt.Sprintf("float64 must be greater than or equal to %f", min)
		},
		ValidateFunc: func(value float64) bool {
			return value >= min
		},
	}

	s.validators = append(s.validators, validator)

	return s
}

func (s *Float64Schema) Positive() *Float64Schema {
	validator := Validator[float64]{
		MessageFunc: func(value float64) string {
			return "float64 must be greater than 0"
		},
		ValidateFunc: func(value float64) bool {
			return value > 0
		},
	}

	s.validators = append(s.validators, validator)

	return s
}

func (s *Float64Schema) Nonnegative() *Float64Schema {
	validator := Validator[float64]{
		MessageFunc: func(value float64) string {
			return "float64 must be greater than or equal to 0"
		},
		ValidateFunc: func(value float64) bool {
			return value >= 0
		},
	}

	s.validators = append(s.validators, validator)

	return s
}

func (s *Float64Schema) Negative() *Float64Schema {
	validator := Validator[float64]{
		MessageFunc: func(value float64) string {
			return "float64 must be less than 0"
		},
		ValidateFunc: func(value float64) bool {
			return value < 0
		},
	}

	s.validators = append(s.validators, validator)

	return s
}

func (s *Float64Schema) Nonpositive() *Float64Schema {
	validator := Validator[float64]{
		MessageFunc: func(value float64) string {
			return "float64 must be less than or equal to 0"
		},
		ValidateFunc: func(value float64) bool {
			return value <= 0
		},
	}

	s.validators = append(s.validators, validator)

	return s
}

func (s *Float64Schema) MultipleOf(multiple float64) *Float64Schema {
	validator := Validator[float64]{
		MessageFunc: func(value float64) string {
			return fmt.Sprintf("float64 must be a multiple of %f", multiple)
		},
		ValidateFunc: func(value float64) bool {
			return math.Mod(value, multiple) == 0.0
		},
	}

	s.validators = append(s.validators, validator)

	return s
}
