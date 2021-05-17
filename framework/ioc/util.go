package ioc

import (
	"github.com/ischenkx/notify/internal/utils"
	"reflect"
)

func matchTypes(p reflect.Type, p2 reflect.Type) bool {

	if p == nil || p2 == nil {
		return false
	}

	if utils.CompareTypes(p, p2) {
		return true
	}

	interfaceRelation := false

	if p.Kind() == reflect.Interface {
		interfaceRelation = interfaceRelation||p2.Implements(p)
	}

	if p2.Kind() == reflect.Interface {
		interfaceRelation = interfaceRelation||p.Implements(p2)
	}

	return interfaceRelation
}
