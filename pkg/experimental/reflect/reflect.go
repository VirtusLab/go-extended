package reflect

import (
	"fmt"
	"reflect"
)

// UnexpectedNilValue indicates a nil value is not expected in this context
type UnexpectedNilValue struct{}

func (*UnexpectedNilValue) Error() string {
	return "nil value"
}

// UnexpectedEmptyValue indicates an empty value is not expected in this context
type UnexpectedEmptyValue struct{}

func (*UnexpectedEmptyValue) Error() string {
	return "empty value"
}

// UnsupportedValueType indicates a value type is not supported in this context
type UnsupportedValueType struct {
	value interface{}
}

func (u *UnsupportedValueType) Error() string {
	return fmt.Sprintf("unsupported type of '%s' with value '%+v'", reflect.TypeOf(u.value), u.value)
}

// CheckEmpty checks if an interface{} is empty
// Returns error if the value is nil or empty
// Returns error if Kind is not Array, Chan, Map, Slice, or String
func CheckEmpty(value interface{}) error {
	length, err := Len(value)
	if err != nil {
		return err
	}

	if length == 0 {
		return &UnexpectedEmptyValue{}
	}

	return nil
}

// Len checks interface{} length
// Returns 0 if the value is nil
// Returns error if Kind is not Array, Chan, Map, Slice, or String
func Len(value interface{}) (int, error) {
	if value == nil {
		return 0, &UnexpectedNilValue{}
	}

	valueOfValue := reflect.ValueOf(value)
	k := valueOfValue.Kind()

	length := -1
	switch k {
	case reflect.Array:
		length = valueOfValue.Len()
	case reflect.Chan:
		length = valueOfValue.Len()
	case reflect.Map:
		length = valueOfValue.Len()
	case reflect.Slice:
		length = valueOfValue.Len()
	case reflect.String:
		length = valueOfValue.Len()
	default:
		return -1, &UnsupportedValueType{value}
	}

	return length, nil
}
