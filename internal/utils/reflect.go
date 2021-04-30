package utils

import "reflect"

func CompareTypes(t1, t2 reflect.Type) bool {
	return t1.String() == t2.String()
}