package schema

type IntSchema struct {
	checks []func(int) bool
}

func Int() *IntSchema {
	return &IntSchema{}
}

func (s *IntSchema) LessThan(min int) *IntSchema {
	s.checks = append(s.checks, func(value int) bool {
		return value < min
	})

	return s
}

func (s *IntSchema) Parse(value any) bool {
	val, ok := value.(int)
	if !ok {
		return false
	}

	for _, check := range s.checks {
		if !check(val) {
			return false
		}
	}

	return true
}
