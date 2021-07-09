package utils

import "reflect"

var emptyInterface interface{}
var emptyInterfaceType = reflect.ValueOf(&emptyInterface).Type()

func CompareTypes(t1, t2 reflect.Type) bool {
	return t1.String() == t2.String()
}

func IsInterface(t reflect.Type) bool {
	return t.Kind() == reflect.Interface
}

func IsEmptyInterface(t reflect.Type) bool {
	if !IsInterface(t) {
		return false
	}
	return emptyInterfaceType.Implements(t)
}

func Implements(t1, t2 reflect.Type) bool {
	if t2.Kind() != reflect.Interface {
		return false
	}
	return t1.Implements(t2)
}
