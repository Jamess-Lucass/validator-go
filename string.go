package schema

type StringSchema struct {
	checks []func(string) bool
}

func String() *StringSchema {
	return &StringSchema{}
}

func (s *StringSchema) Min(min int) *StringSchema {
	s.checks = append(s.checks, func(value string) bool {
		return len(value) >= min
	})

	return s
}

func (s *StringSchema) Parse(value any) bool {
	val, ok := value.(string)
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
