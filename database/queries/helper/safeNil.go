package helper

import "reflect"

func SafeNil(v interface{}) interface{} {
	if v == nil {
		return nil
	}

	val := reflect.ValueOf(v)

	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return nil
		}
		return val.Elem().Interface()
	}

	return v
}
