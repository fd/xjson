package xjson

import (
	"reflect"
)

func (x Value) Unwrap(i interface{}) error {
	return x.UnwrapValue(reflect.ValueOf(i))
}

func (x Value) UnwrapValue(v reflect.Value) error {
	if x.Kind() == Error {
		return x.err
	}
	if x.Kind() == Null {
		return nil
	}

	// drill down on the pointers
	for v.Kind() == reflect.Ptr {
		if v.IsNil() {
			v.Set(reflect.New(v.Type().Elem()))
		}
		v = v.Elem()
	}

	switch x.Kind() {
	case Bool:
		v.SetBool(x.MustBool())
	case Number:
		switch v.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			v.SetInt(x.MustInt64())
		case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			v.SetUint(x.MustUint64())
		case reflect.Float32, reflect.Float64:
			v.SetFloat(x.MustFloat64())
		default:
			v.Set(reflect.ValueOf(x.inner))
		}
	case String:
		v.SetString(x.MustString())
	case Array:
		l := x.Len()
		s := reflect.MakeSlice(v.Type(), l, l)

		for i := 0; i < l; i++ {
			x.GetIndex(i).UnwrapValue(s.Index(i))
		}

		v.Set(s)
	case Object:
		m := reflect.MakeMap(v.Type())

		for key := range x.MustObject() {
			e := reflect.New(m.Type().Elem())
			x.Get(key).UnwrapValue(e)
			m.SetMapIndex(reflect.ValueOf(key), e.Elem())
		}

		v.Set(m)
	}

	return nil
}
