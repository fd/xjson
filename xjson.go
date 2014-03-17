package xjson

import (
	"encoding/json"
	"fmt"
)

type Value struct {
	inner    interface{}
	err      error
	selector Selector
}

type Kind int

const (
	Null Kind = iota
	Bool
	Number
	String
	Array
	Object

	// special kind
	Error
)

func Parse(data []byte) Value {
	var (
		v interface{}
	)
	if err := json.Unmarshal(data, &v); err != nil {
		v = err
	}
	return ValueOf(v)
}

func ValueOf(x interface{}) Value {
	var (
		v Value
	)

	switch i := x.(type) {
	case nil:
		v.inner = x
	case error:
		v.inner = x
	case bool:
		v.inner = x

	case int:
		v.inner = int64(i)
	case int8:
		v.inner = int64(i)
	case int16:
		v.inner = int64(i)
	case int32:
		v.inner = int64(i)
	case int64:
		v.inner = x
	case uint8:
		v.inner = int64(i)
	case uint16:
		v.inner = int64(i)
	case uint32:
		v.inner = int64(i)
	case uint64:
		v.inner = int64(i)

	case float32:
		v.inner = float64(i)
	case float64:
		v.inner = x

	case string:
		v.inner = x
	case []interface{}:
		v.inner = x
	case map[string]interface{}:
		v.inner = x
	default:
		return ValueOf(type_conflict_error(x, "json type", &root_selector{}))
	}

	v.selector = &root_selector{v.inner}

	return v
}

func (x Value) Selector() Selector {
	if x.selector == nil {
		return &root_selector{nil}
	}
	return x.selector
}

func (x Value) Kind() Kind {
	i, err := x.Interface()
	if err != nil {
		return Error
	}

	switch i.(type) {
	case nil:
		return Null
	case bool:
		return Bool
	case int64, float64:
		return Number
	case string:
		return String
	case []interface{}:
		return Array
	case map[string]interface{}:
		return Object
	default:
		panic("should not happen!")
	}
}

func (x Value) Interface() (interface{}, error) {
	if x.err != nil {
		return nil, x.err
	}
	return x.inner, nil
}

func (x Value) Bool() (bool, error) {
	i, err := x.Interface()
	if err != nil {
		return false, err
	}
	if v, ok := i.(bool); ok {
		return v, nil
	}
	return false, type_conflict_error(i, "json bool", x.selector)
}

func (x Value) Int64() (int64, error) {
	i, err := x.Interface()
	if err != nil {
		return 0, err
	}
	if v, ok := i.(int64); ok {
		return v, nil
	}
	if v, ok := i.(float64); ok {
		return int64(v), nil
	}
	return 0, type_conflict_error(x.inner, "json number", x.selector)
}

func (x Value) Uint64() (uint64, error) {
	i, err := x.Interface()
	if err != nil {
		return 0, err
	}
	if v, ok := i.(int64); ok {
		return uint64(v), nil
	}
	if v, ok := i.(float64); ok {
		return uint64(v), nil
	}
	return 0, type_conflict_error(x.inner, "json number", x.selector)
}

func (x Value) Float64() (float64, error) {
	i, err := x.Interface()
	if err != nil {
		return 0, err
	}
	if v, ok := i.(float64); ok {
		return v, nil
	}
	if v, ok := i.(int64); ok {
		return float64(v), nil
	}
	return 0, type_conflict_error(x.inner, "json number", x.selector)
}

func (x Value) String() (string, error) {
	i, err := x.Interface()
	if err != nil {
		return "", err
	}
	if v, ok := i.(string); ok {
		return v, nil
	}
	return "", type_conflict_error(x.inner, "json string", x.selector)
}

func (x Value) Array() ([]interface{}, error) {
	i, err := x.Interface()
	if err != nil {
		return nil, err
	}
	if v, ok := i.([]interface{}); ok {
		return v, nil
	}
	return nil, type_conflict_error(x.inner, "json array", x.selector)
}

func (x Value) Object() (map[string]interface{}, error) {
	i, err := x.Interface()
	if err != nil {
		return nil, err
	}
	if v, ok := i.(map[string]interface{}); ok {
		return v, nil
	}
	return nil, type_conflict_error(x.inner, "json object", x.selector)
}

func (x Value) MustBool() bool {
	v, _ := x.Bool()
	return v
}

func (x Value) MustInt64() int64 {
	v, _ := x.Int64()
	return v
}

func (x Value) MustUint64() uint64 {
	v, _ := x.Uint64()
	return v
}

func (x Value) MustFloat64() float64 {
	v, _ := x.Float64()
	return v
}

func (x Value) MustString() string {
	v, _ := x.String()
	return v
}

func (x Value) MustArray() []interface{} {
	v, _ := x.Array()
	return v
}

func (x Value) MustObject() map[string]interface{} {
	v, _ := x.Object()
	return v
}

func (x Value) GetIndex(idx int) Value {
	a, err := x.Array()
	if err != nil {
		return Value{nil, err, &index_selector{err, idx, x.selector}}
	}
	if idx >= len(a) {
		err = fmt.Errorf("xjson: index out of range")
		sel := &index_selector{err, idx, x.selector}
		err = &selector_error{err, sel}
		return Value{nil, err, sel}
	}
	v := a[idx]
	return Value{v, nil, &index_selector{v, idx, x.selector}}
}

func (x Value) Get(key string) Value {
	o, err := x.Object()
	if err != nil {
		return Value{nil, err, &key_selector{err, key, x.selector}}
	}
	v, found := o[key]
	if !found {
		err = fmt.Errorf("xjson: key not found")
		sel := &key_selector{err, key, x.selector}
		err = &selector_error{err, sel}
		return Value{nil, err, sel}
	}
	return Value{v, nil, &key_selector{v, key, x.selector}}
}

func (x Value) GetPath(parts ...interface{}) Value {
	for _, part := range parts {
		switch y := part.(type) {
		case int:
			x = x.GetIndex(y)
		case string:
			x = x.Get(y)
		default:
			panic("GetPath() expects string and int arguments")
		}
	}
	return x
}

func (x *Value) MarshalJSON() ([]byte, error) {
	i, err := x.Interface()
	if err != nil {
		return nil, err
	}

	return json.Marshal(i)
}

func (x *Value) UnmarshalJSON(b []byte) error {
	v := Parse(b)
	*x = v
	return nil
}

func type_conflict_error(x interface{}, expected_type string, sel Selector) error {
	return &selector_error{fmt.Errorf("xjson: %T is not a %s", x, expected_type), sel}
}

type selector_error struct {
	err      error
	selector Selector
}

func (s *selector_error) Error() string {
	return fmt.Sprintf("%s (at: %s)", s.err, s.selector)
}

func (x Value) Len() int {
	switch x.Kind() {
	case Array:
		return len(x.MustArray())
	case Object:
		return len(x.MustObject())
	default:
		return 0
	}
}
