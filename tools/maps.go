package tools

import (
	"fmt"
	"reflect"
)

//Преобразуем map Interface to String
func MapIterfaceToMapString(source map[interface{}]interface{}) map[string]string {
	data := make(map[string]string)
	for key, value := range source {
		strKey := fmt.Sprintf("%v", key)
		strValue := fmt.Sprintf("%v", value)
		data[strKey] = strValue
	}

	return data
}

//Сравниваем две мапки интерфейсов
func ComparedMapDataInterfaces(A interface{}, B interface{}) bool {
	a := reflect.ValueOf(A)
	b := reflect.ValueOf(B)
	if !reflect.DeepEqual(a.Interface(), b.Interface()) {
		return false
	}
	return true
}
