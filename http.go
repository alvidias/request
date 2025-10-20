package request

import "reflect"

func IsStructMapOrSlice(data any) bool {
	if data == nil {
		return false
	}

	v := reflect.ValueOf(data)
	t := v.Type()

	if t.Kind() == reflect.Ptr {
		v = v.Elem()
		t = v.Type()
	}

	k := t.Kind()

	// check for []byte, not considered as json-serializable
	if k == reflect.Slice && t.Elem().Kind() == reflect.Uint8 {
		return false
	}

	return k == reflect.Map || k == reflect.Struct || k == reflect.Slice
}
