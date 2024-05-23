package utils

import "reflect"

func IsNil(i interface{}) bool {
	if i == nil {
		return true
	}
	switch reflect.TypeOf(i).Kind() {
	case reflect.Ptr, reflect.Map, reflect.Array, reflect.Chan, reflect.Slice:
		return reflect.ValueOf(i).IsNil()
	}
	return false
}

func NS(i interface{}) string {
	if i == nil {
		return ""
	}
	switch reflect.TypeOf(i).Kind() {
	case reflect.String:
		return reflect.ValueOf(i).String()
	}
	return ""
}
