package util

import (
	"fmt"
	"reflect"
)

// StructToMap ...
func StructToMap(value interface{}) (params map[string]string) {
	params = make(map[string]string)
	if value == nil {
		return
	}
	t := reflect.TypeOf(value)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	v := reflect.Indirect(reflect.ValueOf(value))

	if t.Kind() == reflect.Struct {
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			fvalue := v.FieldByName(f.Name)
			if IsNil(fvalue) {
				continue
			}
			fvalue = reflect.Indirect(fvalue)
			xvalue := fvalue.Interface()
			if fvalue.Kind() == reflect.Struct {
				for k, v := range StructToMap(xvalue) {
					params[k] = interfaceToString(v)
				}
			} else {
				xname := f.Tag.Get("json")
				params[xname] = interfaceToString(xvalue)
			}
		}
	}
	return
}

// IsNil is nil
func IsNil(v reflect.Value) bool {
	if v.Kind() == reflect.Ptr {
		return v.IsNil()
	}
	return false
}

func interfaceToString(src interface{}) string {
	if src == nil {
		return ""
	}
	switch src.(type) {
	case string:
		return src.(string)
	case int, int8, int32, int64:
	case uint8, uint16, uint32, uint64:
	case float32, float64:
		return fmt.Sprint(src)
	}
	return ""
}
