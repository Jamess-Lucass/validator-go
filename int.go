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

func (s *IntSchema) Lte(min int) *IntSchema {
	validator := Validator[int]{
		MessageFunc: func(value int) string {
			return fmt.Sprintf("Int must be less than or equal to %d", min)
		},
		ValidateFunc: func(value int) bool {
			return value <= min
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

func (s *IntSchema) Gte(min int) *IntSchema {
	validator := Validator[int]{
		MessageFunc: func(value int) string {
			return fmt.Sprintf("Int must be greater than or equal to %d", min)
		},
		ValidateFunc: func(value int) bool {
			return value >= min
		},
	}

	s.validators = append(s.validators, validator)

	return s
}

func (s *IntSchema) Positive() *IntSchema {
	validator := Validator[int]{
		MessageFunc: func(value int) string {
			return "Int must be greater than 0"
		},
		ValidateFunc: func(value int) bool {
			return value > 0
		},
	}

	s.validators = append(s.validators, validator)

	return s
}

func (s *IntSchema) Nonnegative() *IntSchema {
	validator := Validator[int]{
		MessageFunc: func(value int) string {
			return "Int must be greater than or equal to 0"
		},
		ValidateFunc: func(value int) bool {
			return value >= 0
		},
	}

	s.validators = append(s.validators, validator)

	return s
}

func (s *IntSchema) Negative() *IntSchema {
	validator := Validator[int]{
		MessageFunc: func(value int) string {
			return "Int must be less than 0"
		},
		ValidateFunc: func(value int) bool {
			return value < 0
		},
	}

	s.validators = append(s.validators, validator)

	return s
}

func (s *IntSchema) Nonpositive() *IntSchema {
	validator := Validator[int]{
		MessageFunc: func(value int) string {
			return "Int must be less than or equal to 0"
		},
		ValidateFunc: func(value int) bool {
			return value <= 0
		},
	}

	s.validators = append(s.validators, validator)

	return s
}

func (s *IntSchema) MultipleOf(multiple int) *IntSchema {
	validator := Validator[int]{
		MessageFunc: func(value int) string {
			return fmt.Sprintf("Int must be a multiple of %d", multiple)
		},
		ValidateFunc: func(value int) bool {
			return value%multiple == 0
		},
	}

	s.validators = append(s.validators, validator)

	return s
}
