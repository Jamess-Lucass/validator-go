package schema

import (
	"fmt"
	"reflect"
)

type ObjectSchema struct {
	value map[string]ISchema
	Schema[map[string]interface{}]
}

var _ ISchema = (*ObjectSchema)(nil)

func Object(obj map[string]ISchema) *ObjectSchema {
	return &ObjectSchema{value: obj}
}

func (s *ObjectSchema) Refine(predicate func(map[string]interface{}) bool) *ObjectSchema {
	validator := Validator[map[string]interface{}]{
		MessageFunc: func(value map[string]interface{}) string {
			return "Invalid input"
		},
		ValidateFunc: predicate,
	}

	s.validators = append(s.validators, validator)

	return s
}

func (s *ObjectSchema) Parse(value any) *ValidationResult {
	t := reflect.TypeOf(value)
	val := reflect.ValueOf(value)

	if t.Kind() != reflect.Struct {
		return &ValidationResult{Errors: []ValidationError{{Path: "", Message: fmt.Sprintf("Expected struct, got %T", value)}}}
	}

	res := &ValidationResult{}

	for key, schema := range s.value {
		if _, ok := t.FieldByName(key); !ok {
			err := ValidationError{
				Path:    key,
				Message: "Required",
			}

			res.Errors = append(res.Errors, err)
			continue
		}

		str := val.FieldByName(key).Interface()

		result := schema.Parse(str)
		if !result.IsValid() {
			for _, err := range result.Errors {
				newError := ValidationError{
					Path:    key,
					Message: err.Message,
				}

				res.Errors = append(res.Errors, newError)
			}
		}
	}

	valueMap := StructToMap(value)

	for _, validator := range s.validators {
		if !validator.ValidateFunc(valueMap) {
			err := ValidationError{
				Path:    "",
				Message: validator.MessageFunc(valueMap),
			}

			res.Errors = append(res.Errors, err)
		}
	}

	return res
}

func StructToMap(item interface{}) map[string]interface{} {
	result := map[string]interface{}{}

	val := reflect.ValueOf(item)
	typ := reflect.TypeOf(item)

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		value := val.Field(i).Interface()

		result[field.Name] = value
	}

	return result
}
