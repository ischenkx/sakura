package utils

import "reflect"

func CompareTypes(t1, t2 reflect.Type) bool {

	if t1.Kind() == reflect.Interface {
		if t2.Implements(t1) {
			return true
		}
	}

	if t2.Kind() == reflect.Interface {
		if t1.Implements(t2) {
			return true
		}
	}

	return t1.String() == t2.String()
}