package validator

import (
	"fmt"
	"math"
	"reflect"
	"sort"
)

// MapRule validates a map field. Named map types (type Labels
// map[string]string) work too. Entries are visited in sorted key order so
// error output is stable from run to run.
type MapRule[M ~map[K]V, K comparable, V any] struct {
	fieldBase[M]
	eachFn      func(key K, value *V, sv *Validator)
	eachValRule Rule
	eachKeyRule Rule

	// msgTarget and eachMark route WithMessage to whichever check was added last.
	eachFnMsg  string
	eachValMsg string
	eachKeyMsg string
	msgTarget  *string
	eachMark   int
}

// Map validates a map field. Like Slice, it is a free function because Go
// methods cannot have type parameters. Pass the field address:
//
//	validator.Map(v, &d.Labels).NotEmpty().EachValue(validator.String().NotEmpty())
func Map[M ~map[K]V, K comparable, V any](v *Validator, field *M) *MapRule[M, K, V] {
	r := &MapRule[M, K, V]{fieldBase: fieldBase[M]{v: v, fieldPtr: field}}
	v.addRule(r)
	return r
}

func (r *MapRule[M, K, V]) execute(prefix string) []ValidationError {
	fieldName := r.resolveName(prefix)

	if r.fieldPtr == nil {
		return executeRules(r.rules, nil, fieldName, true)
	}

	val := *r.fieldPtr

	errs := containerErrors(r.rules, val, val == nil, fieldName)
	errs = append(errs, r.entryErrors(val, fieldName)...)

	return errs
}

// Validate implements Rule for standalone use. Each callbacks run here too,
// since entry names come from the key rather than struct reflection.
func (r *MapRule[M, K, V]) Validate(value any) []ValidationError {
	value = unwrapPointers(value)

	if value == nil {
		return executeRules(r.rules, nil, "", true)
	}

	m, ok := coerceMap[K, V](value)
	if !ok {
		return []ValidationError{{Message: msgObject}}
	}

	rv := reflect.ValueOf(value)
	errs := containerErrors(r.rules, M(m), rv.Kind() == reflect.Map && rv.IsNil(), "")
	errs = append(errs, r.entryErrors(M(m), "")...)

	return errs
}

func (r *MapRule[M, K, V]) entryErrors(m M, fieldName string) []ValidationError {
	if r.eachFn == nil && r.eachValRule == nil && r.eachKeyRule == nil {
		return nil
	}

	var errs []ValidationError
	for _, k := range sortedKeys(map[K]V(m)) {
		entryName := fieldName + "[" + keyLabel(k) + "]"

		if r.eachKeyRule != nil {
			for _, e := range r.eachKeyRule.Validate(k) {
				errs = append(errs, ValidationError{
					Field:   entryName,
					Message: orMsg(r.eachKeyMsg, "key "+e.Message),
				})
			}
		}

		if r.eachValRule != nil {
			for _, e := range r.eachValRule.Validate(m[k]) {
				errs = append(errs, ValidationError{
					Field:   joinFieldName(entryName, e.Field),
					Message: orMsg(r.eachValMsg, e.Message),
				})
			}
		}

		if r.eachFn != nil {
			// Map values aren't addressable, so the callback gets a pointer
			// to a copy.
			val := m[k]
			sv := &Validator{structPtr: &val, prefix: entryName}
			r.eachFn(k, &val, sv)
			for _, e := range sv.Validate().Errors {
				e.Message = orMsg(r.eachFnMsg, e.Message)
				errs = append(errs, e)
			}
		}
	}

	return errs
}

// Presence rules

// NotNil fails on a nil map, which is distinct from an empty one.
func (r *MapRule[M, K, V]) NotNil() *MapRule[M, K, V] {
	r.rules = append(r.rules, rule[M]{isNotNil: true, message: msgNotNil})
	return r
}

// NotEmpty fails when the map has no entries.
func (r *MapRule[M, K, V]) NotEmpty() *MapRule[M, K, V] {
	r.rules = append(r.rules, rule[M]{
		validate: func(val M) bool { return len(val) > 0 },
		message:  msgNotEmpty,
	})
	return r
}

// Count rules

// Min requires at least n entries.
func (r *MapRule[M, K, V]) Min(n int) *MapRule[M, K, V] {
	r.rules = append(r.rules, rule[M]{
		validate: func(val M) bool { return len(val) >= n },
		message:  msgMinItems(n),
	})
	return r
}

// Max allows at most n entries.
func (r *MapRule[M, K, V]) Max(n int) *MapRule[M, K, V] {
	r.rules = append(r.rules, rule[M]{
		validate: func(val M) bool { return len(val) <= n },
		message:  msgMaxItems(n),
	})
	return r
}

// Length requires between min and max entries.
func (r *MapRule[M, K, V]) Length(min, max int) *MapRule[M, K, V] {
	r.rules = append(r.rules, rule[M]{
		validate: func(val M) bool { return len(val) >= min && len(val) <= max },
		message:  msgItemsBetween(min, max),
	})
	return r
}

// Key rules

// HasKey requires the key to be present; chain it to require several.
func (r *MapRule[M, K, V]) HasKey(key K) *MapRule[M, K, V] {
	r.rules = append(r.rules, rule[M]{
		validate: func(val M) bool {
			_, ok := val[key]
			return ok
		},
		message: fmt.Sprintf("must contain key \"%v\"", key),
	})
	return r
}

// Iteration rules

// Each runs the callback for every entry in sorted key order. Errors raised
// on sv are prefixed with the entry, like "containers[web].image". The value
// is a copy, since map values aren't addressable. It panics if fn is nil.
func (r *MapRule[M, K, V]) Each(fn func(key K, value *V, sv *Validator)) *MapRule[M, K, V] {
	if fn == nil {
		panic("validator: Each requires a non-nil callback")
	}
	r.eachFn = fn
	r.eachFnMsg = ""
	r.msgTarget = &r.eachFnMsg
	r.eachMark = len(r.rules)
	return r
}

// EachValue validates every value against the given rule, naming errors by
// entry: "labels[env]: must not be empty". It panics if rule is nil.
func (r *MapRule[M, K, V]) EachValue(rule Rule) *MapRule[M, K, V] {
	if isNilRule(rule) {
		panic("validator: EachValue requires a non-nil rule")
	}
	r.eachValRule = rule
	r.eachValMsg = ""
	r.msgTarget = &r.eachValMsg
	r.eachMark = len(r.rules)
	return r
}

// EachKey validates every key against the given rule. Error messages are
// prefixed with "key" to tell them apart from value errors on the same entry:
// "labels[e]: key must be at least 2 characters". It panics if rule is nil.
func (r *MapRule[M, K, V]) EachKey(rule Rule) *MapRule[M, K, V] {
	if isNilRule(rule) {
		panic("validator: EachKey requires a non-nil rule")
	}
	r.eachKeyRule = rule
	r.eachKeyMsg = ""
	r.msgTarget = &r.eachKeyMsg
	r.eachMark = len(r.rules)
	return r
}

// Custom rules

// Must runs a custom check against the whole map. It panics if fn is nil.
func (r *MapRule[M, K, V]) Must(fn func(M) bool) *MapRule[M, K, V] {
	if fn == nil {
		panic("validator: Must requires a non-nil function")
	}
	r.rules = append(r.rules, rule[M]{validate: fn, message: msgNotValid})
	return r
}

// WithMessage replaces the previous rule's error message. Directly after
// Each, EachValue, or EachKey it instead replaces the message of every error
// they produce.
func (r *MapRule[M, K, V]) WithMessage(msg string) *MapRule[M, K, V] {
	if r.msgTarget != nil && len(r.rules) == r.eachMark {
		*r.msgTarget = msg
		return r
	}
	if len(r.rules) > 0 {
		r.rules[len(r.rules)-1].message = msg
	}
	return r
}

// WithName overrides the field name used in errors.
func (r *MapRule[M, K, V]) WithName(name string) *MapRule[M, K, V] {
	r.name = name
	return r
}

// coerceMap converts a value into map[K]V, directly when the types line up
// and entry by entry via reflection when they don't, like a map[string]string
// handed to a rule built for map[string]any.
func coerceMap[K comparable, V any](value any) (map[K]V, bool) {
	if m, ok := value.(map[K]V); ok {
		return m, true
	}

	rv := reflect.ValueOf(value)
	if rv.Kind() != reflect.Map {
		return nil, false
	}

	out := make(map[K]V, rv.Len())
	iter := rv.MapRange()
	for iter.Next() {
		k, ok := iter.Key().Interface().(K)
		if !ok {
			return nil, false
		}
		v, ok := iter.Value().Interface().(V)
		if !ok {
			return nil, false
		}
		out[k] = v
	}

	return out, true
}

func sortedKeys[K comparable, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool { return lessKey(keys[i], keys[j]) })
	return keys
}

func lessKey(a, b any) bool {
	av, bv := reflect.ValueOf(a), reflect.ValueOf(b)
	if !av.IsValid() || !bv.IsValid() || av.Kind() != bv.Kind() {
		// Mixed concrete types (or a nil) behind an interface key type.
		return fmt.Sprint(a) < fmt.Sprint(b)
	}

	switch av.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return av.Int() < bv.Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return av.Uint() < bv.Uint()
	case reflect.Float32, reflect.Float64:
		// NaN sorts last; a plain < would break sort's ordering contract.
		af, bf := av.Float(), bv.Float()
		if math.IsNaN(af) {
			return false
		}
		if math.IsNaN(bf) {
			return true
		}
		return af < bf
	case reflect.String:
		return av.String() < bv.String()
	case reflect.Bool:
		return !av.Bool() && bv.Bool()
	case reflect.Pointer:
		// Addresses change between runs, so order by label, nil first.
		if av.IsNil() || bv.IsNil() {
			return av.IsNil() && !bv.IsNil()
		}
		return keyLabel(a) < keyLabel(b)
	default:
		return fmt.Sprint(a) < fmt.Sprint(b)
	}
}

// keyLabel renders a key for error field names; pointer keys show the value
// they point to.
func keyLabel(k any) string {
	rv := reflect.ValueOf(k)
	for rv.Kind() == reflect.Pointer {
		if rv.IsNil() {
			return "nil"
		}
		rv = rv.Elem()
	}

	if !rv.IsValid() {
		return fmt.Sprint(k)
	}

	return fmt.Sprint(rv.Interface())
}
