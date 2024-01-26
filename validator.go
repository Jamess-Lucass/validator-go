package validator

type Rules struct {
	field   any
	IsValid bool
	Message string
	min     bool
}

func RuleFor(field any, options ...func(*Rules)) *Rules {
	rules := &Rules{field: field}
	for _, o := range options {
		o(rules)
	}
	return rules
}

func Min(min int) func(*Rules) {
	return func(s *Rules) {
		str, ok := s.field.(string)
		if !ok {
			s.IsValid = false
			return
		}

		s.IsValid = len(str) >= min
	}
}

func Max(max int) func(*Rules) {
	return func(s *Rules) {
		str, ok := s.field.(string)
		if !ok {
			s.IsValid = false
			return
		}

		s.IsValid = len(str) <= max
	}
}
