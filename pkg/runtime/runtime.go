package runtime

import (
	"fmt"
	"reflect"
)

// NullSafeValue wraps a value with null-safety checking
type NullSafeValue struct {
	Value interface{}
	IsNil bool
}

// NullCoalesce returns the value if not nil, otherwise returns the default
func (nsv NullSafeValue) Coalesce(defaultVal interface{}) interface{} {
	if nsv.IsNil {
		return defaultVal
	}
	return nsv.Value
}

// SafeCall attempts to call a method safely, returning nil if it fails
func SafeCall(obj interface{}, method string, args ... interface{}) interface{} {
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return nil
	}

	m := v.MethodByName(method)
	if !m. IsValid() {
		return nil
	}

	return m. Call(nil). Interface()
}

// SafeAccess safely accesses a field, returning nil if it doesn't exist
func SafeAccess(obj interface{}, field string) interface{} {
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return nil
	}

	f := v.FieldByName(field)
	if !f. IsValid() {
		return nil
	}

	return f. Interface()
}

// Assert checks that a value is not nil and of the expected type
func Assert(value interface{}, typename string) (interface{}, error) {
	if value == nil {
		return nil, fmt. Errorf("null safety violation: expected %s, got nil", typename)
	}

	actualType := reflect.TypeOf(value). String()
	if actualType != typename && ! isCompatible(actualType, typename) {
		return nil, fmt.Errorf("type error: expected %s, got %s", typename, actualType)
	}

	return value, nil
}

// isCompatible checks if two types are compatible
func isCompatible(actual, expected string) bool {
	if actual == expected {
		return true
	}
	if expected == "interface{}" {
		return true
	}
	return false
}

// CheckNil returns an error if value is nil
func CheckNil(value interface{}, name string) error {
	if value == nil {
		return fmt. Errorf("null safety violation: %s is nil", name)
	}
	return nil
}
