package validator

type BoolRule struct {
	fieldBase[bool]
}

func (v *Validator) Bool(field *bool) *BoolRule {
	r := &BoolRule{fieldBase[bool]{v: v, fieldPtr: field}}
	v.addRule(r)
	return r
}

func Bool() *BoolRule {
	return &BoolRule{}
}

func (r *BoolRule) Validate(value any) []ValidationError {
	return validateStandalone(r.rules, value, "must be a boolean", func(v any) (bool, bool) {
		b, ok := v.(bool)
		return b, ok
	})
}

func (r *BoolRule) NotNil() *BoolRule {
	r.rules = append(r.rules, rule[bool]{isNotNil: true, message: msgNotNil})
	return r
}

func (r *BoolRule) NotEmpty() *BoolRule {
	r.rules = append(r.rules, rule[bool]{
		validate: func(val bool) bool { return val },
		message:  msgNotEmpty,
	})
	return r
}

func (r *BoolRule) Must(fn func(bool) bool) *BoolRule {
	r.rules = append(r.rules, rule[bool]{validate: fn, message: msgNotValid})
	return r
}

func (r *BoolRule) WithMessage(msg string) *BoolRule {
	if len(r.rules) > 0 {
		r.rules[len(r.rules)-1].message = msg
	}
	return r
}

func (r *BoolRule) WithName(name string) *BoolRule {
	r.name = name
	return r
}
