package tools

import (
	"fmt"
	"reflect"
)

//Заполняем структуру из мапы интерфейсов
func StructFillFromMapInterface(m map[string]interface{}, s interface{}) error {
	structValue := reflect.ValueOf(s).Elem()

	for fieldName, fieldValue := range m {
		structFieldValue := structValue.FieldByName(fieldName)
		if !structFieldValue.IsValid() {
			return fmt.Errorf("No such field: %v in object", fieldName)
		}
		if !structFieldValue.CanSet() {
			return fmt.Errorf("Can not set field: %v, value: %v", fieldName, fieldValue)
		}
		value := reflect.ValueOf(fieldValue)
		if structFieldValue.Type() != value.Type() {
			return fmt.Errorf("Provided value: %v, value type: %v does not match object field type: %v", fieldValue, value.Type(), structFieldValue.Type())
		}

		structFieldValue.Set(value)

	}
	return nil
}

//Конвертируем структуру в мапку интерфейсов
func StructToMap(obj interface{}) map[string]interface{} {

	res := map[string]interface{}{}
	if obj == nil {
		return res
	}
	v := reflect.TypeOf(obj)
	reflectValue := reflect.ValueOf(obj)
	reflectValue = reflect.Indirect(reflectValue)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	for i := 0; i < v.NumField(); i++ {
		tag := v.Field(i).Tag.Get("struct-to-map")
		field := reflectValue.Field(i).Interface()
		if tag != "" && tag != "-" {
			if v.Field(i).Type.Kind() == reflect.Struct {
				res[tag] = StructToMap(field)
			} else {
				res[tag] = field
			}
		}
	}
	return res
}
