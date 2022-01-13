package helpers

import "reflect"

func ReflectValueToInterface(value reflect.Value) interface{} {
	return value.Interface()
}

func SliceReflectValuesToInterfaces(values []reflect.Value) []interface{} {
	data := make([]interface{}, len(values))
	for i, obj := range values {
		data[i] = ReflectValueToInterface(obj)
	}
	return data
}
