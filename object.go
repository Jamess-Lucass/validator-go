package validator

// ObjectValidator validates unstructured map[string]any data, declaring one rule per key.
type ObjectValidator struct {
	data  map[string]any
	rules []objectFieldRule
}

type objectFieldRule struct {
	key  string
	rule Rule
}

// NewObject returns a validator for unstructured data such as a decoded JSON
// payload. Declare rules with Field, then call Validate.
func NewObject(m map[string]any) *ObjectValidator {
	return &ObjectValidator{data: m}
}

// Field declares a rule for one key of the map. It panics if rule is nil.
func (mv *ObjectValidator) Field(key string, rule Rule) {
	if isNilRule(rule) {
		panic("validator: Field requires a non-nil rule")
	}
	mv.rules = append(mv.rules, objectFieldRule{key: key, rule: rule})
}

// Validate runs every field rule and returns the collected errors.
func (mv *ObjectValidator) Validate() *ValidationResult {
	return &ValidationResult{Errors: mv.execute()}
}

func (mv *ObjectValidator) execute() []ValidationError {
	var errs []ValidationError

	for _, fr := range mv.rules {
		for _, e := range fr.rule.Validate(mv.data[fr.key]) {
			errs = append(errs, ValidationError{
				Field:   joinFieldName(fr.key, e.Field),
				Message: e.Message,
			})
		}
	}

	return errs
}

// ObjectRule validates a nested unstructured object within NewObject data.
type ObjectRule struct {
	fn     func(*ObjectValidator)
	notNil bool
}

// Object declares rules for a nested object (map[string]any) inside Field:
//
//	v.Field("address", validator.Object(func(mv *validator.ObjectValidator) {
//		mv.Field("city", validator.String().NotEmpty())
//	}))
//
// For a typed map field on a struct (map[string]string and friends), use Map.
// It panics if fn is nil.
func Object(fn func(*ObjectValidator)) *ObjectRule {
	if fn == nil {
		panic("validator: Object requires a non-nil callback")
	}
	return &ObjectRule{fn: fn}
}

// NotNil fails when the key is missing or its value is nil.
func (r *ObjectRule) NotNil() *ObjectRule {
	r.notNil = true
	return r
}

// Validate implements Rule for standalone use; the value must be a
// string-keyed map.
func (r *ObjectRule) Validate(value any) []ValidationError {
	value = unwrapPointers(value)

	if value == nil {
		if r.notNil {
			return []ValidationError{{Message: msgNotNil}}
		}
		return nil
	}

	m, ok := coerceMap[string, any](value)
	if !ok {
		return []ValidationError{{Message: msgObject}}
	}

	mv := &ObjectValidator{data: m}
	r.fn(mv)

	return mv.execute()
}

// SliceOf runs each element of a slice value through the given rule, for use
// with NewObject fields.
func SliceOf(each Rule) *SliceRule[[]any, any] {
	return (&SliceRule[[]any, any]{}).EachValue(each)
}

// MapOf runs every value of a map through the given rule, for use with NewObject
// fields. It is to Map what SliceOf is to Slice.
func MapOf(each Rule) *MapRule[map[string]any, string, any] {
	return (&MapRule[map[string]any, string, any]{}).EachValue(each)
}
