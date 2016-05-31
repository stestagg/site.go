package log

import (
	"reflect"
	"fmt"
	"strings"
	"github.com/spf13/cast"
)

// Copyright 2011 The Go Authors. All rights reserved.
// indirect returns the value, after dereferencing as many times
// as necessary to reach the base type (or nil).
func indirect(a interface{}) interface{} {
	if a == nil {
		return nil
	}
	if t := reflect.TypeOf(a); t.Kind() != reflect.Ptr {
		// Avoid creating a reflect.Value if it's not a pointer.
		return a
	}
	v := reflect.ValueOf(a)
	for v.Kind() == reflect.Ptr && !v.IsNil() {
		v = v.Elem()
	}
	return v.Interface()
}


func arrayToString(a []interface{}) string {
	out := make([]string, len(a))
	for i, elm := range(a) {
		out[i] = ToString(elm)
	}
	return fmt.Sprintf("[%s]",  strings.Join(out, ", "))
}

func strArrayToString(a []string) string {
	out := make([]string, len(a))
	for i, elm := range(a) {
		out[i] = elm
	}
	return fmt.Sprintf("[%s]",  strings.Join(out, ", "))
}

func strMapToString(m map[string]interface{}) string {
	out := make([]string, len(m))
	i := 0
	for key, value := range(m) {
		out[i] = fmt.Sprintf("%s: %s", key, ToString(value))
		i += 1
	}
	return fmt.Sprintf("{%s}",  strings.Join(out, ", "))
}

func mapToString(m map[interface{}]interface{}) string {
	out := make([]string, len(m))
	i := 0
	for key, value := range(m) {
		out[i] = fmt.Sprintf("%s: %s", ToString(key), ToString(value))
		i += 1
	}
	return fmt.Sprintf("{%s}",  strings.Join(out, ", "))
}


func ToString(i interface{}) string {
	switch val := i.(type) {
	case error:
		return val.Error()
	}
	i = indirect(i)
	switch val := i.(type) {
	case []interface{}:
		return arrayToString(val)
	case []string:
		return strArrayToString(val)
	case map[string]interface{}:
		return strMapToString(val)
	case map[interface{}]interface{}:
		return mapToString(val)
	default:
		newval, err := cast.ToStringE(val)
		if err != nil {
			return fmt.Sprint(newval)
		}
		return newval

 	}
}