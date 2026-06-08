package validator

// MapV validates unstructured map[string]any data using a builder pattern.
type MapV struct {
	data  map[string]any
	rules []mapFieldRule
}

type mapFieldRule struct {
	key  string
	rule Rule
}

func NewMap(m map[string]any) *MapV {
	return &MapV{data: m}
}

func (mv *MapV) Field(key string, rule Rule) {
	mv.rules = append(mv.rules, mapFieldRule{key: key, rule: rule})
}

func (mv *MapV) Validate() *ValidationResult {
	return &ValidationResult{Errors: mv.execute("")}
}

func (mv *MapV) execute(prefix string) []ValidationError {
	var errs []ValidationError

	for _, fr := range mv.rules {
		val, exists := mv.data[fr.key]
		if !exists {
			val = nil
		}

		fieldName := joinFieldName(prefix, fr.key)
		fieldErrs := fr.rule.Validate(val)

		for _, e := range fieldErrs {
			errs = append(errs, ValidationError{
				Field:   joinFieldName(fieldName, e.Field),
				Message: e.Message,
			})
		}
	}

	return errs
}

// MapRule implements Rule for nested map validation within Field().
type MapRule struct {
	fn     func(*MapV)
	notNil bool
}

func Map(fn func(*MapV)) *MapRule {
	return &MapRule{fn: fn}
}

func (r *MapRule) NotNil() *MapRule {
	r.notNil = true
	return r
}

func (r *MapRule) Validate(value any) []ValidationError {
	if value == nil {
		if r.notNil {
			return []ValidationError{{Message: msgNotNil}}
		}
		return nil
	}

	m, ok := value.(map[string]any)
	if !ok {
		return []ValidationError{{Message: msgObject}}
	}

	mv := &MapV{data: m}
	r.fn(mv)

	return mv.execute("")
}

// SliceOf validates a []any field within Field(), running each element through
// the given rule. It is a thin constructor over the canonical SliceRule so map
// and struct slice validation share one implementation.
func SliceOf(each Rule) *SliceRule[any] {
	return (&SliceRule[any]{}).EachValue(each)
}
