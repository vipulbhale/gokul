package controller

import (
	"reflect"
)

var (
	mapOfControllerNameToControllerObj = make(map[string]reflect.Type)
)

func RegisterControllers() {
}

func New(name string) (interface{}, bool) {
	t, ok := mapOfControllerNameToControllerObj[name]
	if !ok {
		return nil, false
	}
	v := reflect.New(t)
	return v.Interface(), true
}
