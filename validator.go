package validator

type ValidationResult struct {
	Errors []ValidationError
}

func (v *ValidationResult) IsValid() bool {
	return len(v.Errors) == 0
}

type ValidationError struct {
	Name         string
	errorMessage string
}

func (m *ValidationError) Error() string {
	return m.errorMessage
}

func (r *RuleConfig) WithMessage(message string) *RuleConfig {
	r.MessageFunc = func(r RuleConfig) string {
		return message
	}

	return r
}

type Rule struct {
	value any
	name  string
	rules []*RuleConfig
}

func RuleFor(value any, options ...*RuleConfig) *Rule {
	return &Rule{value: value, rules: options}
}

func (r *Rule) WithName(name string) *Rule {
	r.name = name

	return r
}

func (r *Rule) Validate() *ValidationResult {
	res := &ValidationResult{}

	for _, rule := range r.rules {
		rule.Value = r.value

		if !rule.ValidateFunc(*rule) {
			err := ValidationError{
				Name:         r.name,
				errorMessage: rule.MessageFunc(*rule),
			}

			res.Errors = append(res.Errors, err)
		}
	}

	return res
}

type MultipleRules []Rule

func Rules(options ...*Rule) *MultipleRules {
	multipleRules := make(MultipleRules, len(options))

	for i, opt := range options {
		multipleRules[i] = *opt
	}

	return &multipleRules
}

func (r *MultipleRules) Validate() *ValidationResult {
	res := &ValidationResult{}

	for _, rule := range *r {
		for _, err := range rule.Validate().Errors {
			res.Errors = append(res.Errors, err)
		}
	}

	return res
}
