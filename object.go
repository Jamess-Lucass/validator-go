package schema

import (
	"log"
	"reflect"
)

type ObjectSchema struct {
	value map[string]ISchema
}

func Object(obj map[string]ISchema) *ObjectSchema {
	return &ObjectSchema{value: obj}
}

func (s *ObjectSchema) Parse(value any) bool {
	t := reflect.TypeOf(value)
	val := reflect.ValueOf(value)

	if t.Kind() != reflect.Struct {
		return false
	}

	for key, schema := range s.value {
		if _, ok := t.FieldByName(key); !ok {
			log.Printf("property by key '%s' was not found on struct '%v'\n", key, t.Name())
			return false
		}

		str := val.FieldByName(key).Interface()

		if !schema.Parse(str) {
			return false
		}
	}

	return true
}
