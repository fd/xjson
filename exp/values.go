package xjson

import (
	"fmt"
	"sort"
	"strconv"
)

type Value interface {
	Kind() Kind

	String() string
	MaybeString() (string, bool)
	MustString() string

	Bool() bool
	MaybeBool() (bool, bool)
	MustBool() bool

	Float() float64
	MaybeFloat() (float64, bool)
	MustFloat() float64

	Int() int64
	MaybeInt() (int64, bool)
	MustInt() int64

	Uint() uint64
	MaybeUint() (uint64, bool)
	MustUint() uint64

	Interface() interface{}
	IsNil() bool
	Len() int

	Index(i int) Value
	MapIndex(key string) Value
	Path(parts ...interface{}) Value

	Selector() interface{}

	jsonValue()
}

type errorValue struct {
	err error
}
type zeroValue struct {
}
type nullValue struct {
	buf []byte
}
type boolValue struct {
	buf []byte
	val bool
}
type numberValue struct {
	buf   []byte
	flags numberFlags
}
type stringValue struct {
	buf []byte
	val string
}
type arrayValue struct {
	values []Value
}
type objectValue struct {
	members []objectMember
}
type objectMember struct {
	key   string
	value Value
}

var zero = &zeroValue{}

func (*errorValue) jsonValue()  {}
func (*zeroValue) jsonValue()   {}
func (*nullValue) jsonValue()   {}
func (*boolValue) jsonValue()   {}
func (*numberValue) jsonValue() {}
func (*stringValue) jsonValue() {}
func (*arrayValue) jsonValue()  {}
func (*objectValue) jsonValue() {}

func (*errorValue) Selector() interface{}  { return nil }
func (*zeroValue) Selector() interface{}   { return nil }
func (*nullValue) Selector() interface{}   { return nil }
func (*boolValue) Selector() interface{}   { return nil }
func (*numberValue) Selector() interface{} { return nil }
func (*stringValue) Selector() interface{} { return nil }
func (*arrayValue) Selector() interface{}  { return nil }
func (*objectValue) Selector() interface{} { return nil }

func (*errorValue) Kind() Kind  { return Error }
func (*zeroValue) Kind() Kind   { return Null }
func (*nullValue) Kind() Kind   { return Null }
func (*boolValue) Kind() Kind   { return Bool }
func (*numberValue) Kind() Kind { return Number }
func (*stringValue) Kind() Kind { return String }
func (*arrayValue) Kind() Kind  { return Array }
func (*objectValue) Kind() Kind { return Object }

func (x *errorValue) String() string  { panic(invalid_kind_error(String, x.Kind())) }
func (x *zeroValue) String() string   { return "" }
func (x *nullValue) String() string   { panic(invalid_kind_error(String, x.Kind())) }
func (x *boolValue) String() string   { panic(invalid_kind_error(String, x.Kind())) }
func (x *numberValue) String() string { panic(invalid_kind_error(String, x.Kind())) }
func (x *stringValue) String() string { return x.val }
func (x *arrayValue) String() string  { panic(invalid_kind_error(String, x.Kind())) }
func (x *objectValue) String() string { panic(invalid_kind_error(String, x.Kind())) }

func (x *errorValue) MaybeString() (string, bool)  { return "", false }
func (x *zeroValue) MaybeString() (string, bool)   { return "", false }
func (x *nullValue) MaybeString() (string, bool)   { return "", false }
func (x *boolValue) MaybeString() (string, bool)   { return "", false }
func (x *numberValue) MaybeString() (string, bool) { return "", false }
func (x *stringValue) MaybeString() (string, bool) { return x.val, true }
func (x *arrayValue) MaybeString() (string, bool)  { return "", false }
func (x *objectValue) MaybeString() (string, bool) { return "", false }

func (x *errorValue) MustString() string  { return "" }
func (x *zeroValue) MustString() string   { return "" }
func (x *nullValue) MustString() string   { return "" }
func (x *boolValue) MustString() string   { return "" }
func (x *numberValue) MustString() string { return "" }
func (x *stringValue) MustString() string { return x.val }
func (x *arrayValue) MustString() string  { return "" }
func (x *objectValue) MustString() string { return "" }

func (x *errorValue) Bool() bool  { panic(invalid_kind_error(Bool, x.Kind())) }
func (x *zeroValue) Bool() bool   { return false }
func (x *nullValue) Bool() bool   { panic(invalid_kind_error(Bool, x.Kind())) }
func (x *boolValue) Bool() bool   { return x.val }
func (x *numberValue) Bool() bool { panic(invalid_kind_error(Bool, x.Kind())) }
func (x *stringValue) Bool() bool { panic(invalid_kind_error(Bool, x.Kind())) }
func (x *arrayValue) Bool() bool  { panic(invalid_kind_error(Bool, x.Kind())) }
func (x *objectValue) Bool() bool { panic(invalid_kind_error(Bool, x.Kind())) }

func (x *errorValue) MaybeBool() (bool, bool)  { return false, false }
func (x *zeroValue) MaybeBool() (bool, bool)   { return false, false }
func (x *nullValue) MaybeBool() (bool, bool)   { return false, false }
func (x *boolValue) MaybeBool() (bool, bool)   { return x.val, true }
func (x *numberValue) MaybeBool() (bool, bool) { return false, false }
func (x *stringValue) MaybeBool() (bool, bool) { return false, false }
func (x *arrayValue) MaybeBool() (bool, bool)  { return false, false }
func (x *objectValue) MaybeBool() (bool, bool) { return false, false }

func (x *errorValue) MustBool() bool  { return false }
func (x *zeroValue) MustBool() bool   { return false }
func (x *nullValue) MustBool() bool   { return false }
func (x *boolValue) MustBool() bool   { return x.val }
func (x *numberValue) MustBool() bool { return false }
func (x *stringValue) MustBool() bool { return false }
func (x *arrayValue) MustBool() bool  { return false }
func (x *objectValue) MustBool() bool { return false }

func (x *errorValue) Float() float64  { panic(invalid_kind_error(Number, x.Kind())) }
func (x *zeroValue) Float() float64   { return 0 }
func (x *nullValue) Float() float64   { panic(invalid_kind_error(Number, x.Kind())) }
func (x *boolValue) Float() float64   { panic(invalid_kind_error(Number, x.Kind())) }
func (x *stringValue) Float() float64 { panic(invalid_kind_error(Number, x.Kind())) }
func (x *arrayValue) Float() float64  { panic(invalid_kind_error(Number, x.Kind())) }
func (x *objectValue) Float() float64 { panic(invalid_kind_error(Number, x.Kind())) }

func (x *numberValue) Float() float64 {
	f, _ := strconv.ParseFloat(string(x.buf), 64)
	return f
}

func (x *errorValue) MaybeFloat() (float64, bool)  { return 0, false }
func (x *zeroValue) MaybeFloat() (float64, bool)   { return 0, false }
func (x *nullValue) MaybeFloat() (float64, bool)   { return 0, false }
func (x *boolValue) MaybeFloat() (float64, bool)   { return 0, false }
func (x *numberValue) MaybeFloat() (float64, bool) { return x.Float(), true }
func (x *stringValue) MaybeFloat() (float64, bool) { return 0, false }
func (x *arrayValue) MaybeFloat() (float64, bool)  { return 0, false }
func (x *objectValue) MaybeFloat() (float64, bool) { return 0, false }

func (x *errorValue) MustFloat() float64  { return 0 }
func (x *zeroValue) MustFloat() float64   { return 0 }
func (x *nullValue) MustFloat() float64   { return 0 }
func (x *boolValue) MustFloat() float64   { return 0 }
func (x *numberValue) MustFloat() float64 { return x.Float() }
func (x *stringValue) MustFloat() float64 { return 0 }
func (x *arrayValue) MustFloat() float64  { return 0 }
func (x *objectValue) MustFloat() float64 { return 0 }

func (x *errorValue) Int() int64  { panic(invalid_kind_error(Number, x.Kind())) }
func (x *zeroValue) Int() int64   { return 0 }
func (x *nullValue) Int() int64   { panic(invalid_kind_error(Number, x.Kind())) }
func (x *boolValue) Int() int64   { panic(invalid_kind_error(Number, x.Kind())) }
func (x *stringValue) Int() int64 { panic(invalid_kind_error(Number, x.Kind())) }
func (x *arrayValue) Int() int64  { panic(invalid_kind_error(Number, x.Kind())) }
func (x *objectValue) Int() int64 { panic(invalid_kind_error(Number, x.Kind())) }

func (x *numberValue) Int() int64 {
	if x.flags&(numberHasExponent|numberHasFraction) > 0 {
		return int64(x.Float())
	}
	i, _ := strconv.ParseInt(string(x.buf), 10, 64)
	return i
}

func (x *errorValue) MaybeInt() (int64, bool)  { return 0, false }
func (x *zeroValue) MaybeInt() (int64, bool)   { return 0, false }
func (x *nullValue) MaybeInt() (int64, bool)   { return 0, false }
func (x *boolValue) MaybeInt() (int64, bool)   { return 0, false }
func (x *numberValue) MaybeInt() (int64, bool) { return x.Int(), true }
func (x *stringValue) MaybeInt() (int64, bool) { return 0, false }
func (x *arrayValue) MaybeInt() (int64, bool)  { return 0, false }
func (x *objectValue) MaybeInt() (int64, bool) { return 0, false }

func (x *errorValue) MustInt() int64  { return 0 }
func (x *zeroValue) MustInt() int64   { return 0 }
func (x *nullValue) MustInt() int64   { return 0 }
func (x *boolValue) MustInt() int64   { return 0 }
func (x *numberValue) MustInt() int64 { return x.Int() }
func (x *stringValue) MustInt() int64 { return 0 }
func (x *arrayValue) MustInt() int64  { return 0 }
func (x *objectValue) MustInt() int64 { return 0 }

func (x *errorValue) Uint() uint64  { panic(invalid_kind_error(Number, x.Kind())) }
func (x *zeroValue) Uint() uint64   { return 0 }
func (x *nullValue) Uint() uint64   { panic(invalid_kind_error(Number, x.Kind())) }
func (x *boolValue) Uint() uint64   { panic(invalid_kind_error(Number, x.Kind())) }
func (x *stringValue) Uint() uint64 { panic(invalid_kind_error(Number, x.Kind())) }
func (x *arrayValue) Uint() uint64  { panic(invalid_kind_error(Number, x.Kind())) }
func (x *objectValue) Uint() uint64 { panic(invalid_kind_error(Number, x.Kind())) }

func (x *numberValue) Uint() uint64 {
	if x.flags&(numberHasExponent|numberHasFraction) > 0 {
		return uint64(x.Float())
	}
	return uint64(x.Int())
}

func (x *errorValue) MaybeUint() (uint64, bool)  { return 0, false }
func (x *zeroValue) MaybeUint() (uint64, bool)   { return 0, false }
func (x *nullValue) MaybeUint() (uint64, bool)   { return 0, false }
func (x *boolValue) MaybeUint() (uint64, bool)   { return 0, false }
func (x *numberValue) MaybeUint() (uint64, bool) { return x.Uint(), true }
func (x *stringValue) MaybeUint() (uint64, bool) { return 0, false }
func (x *arrayValue) MaybeUint() (uint64, bool)  { return 0, false }
func (x *objectValue) MaybeUint() (uint64, bool) { return 0, false }

func (x *errorValue) MustUint() uint64  { return 0 }
func (x *zeroValue) MustUint() uint64   { return 0 }
func (x *nullValue) MustUint() uint64   { return 0 }
func (x *boolValue) MustUint() uint64   { return 0 }
func (x *numberValue) MustUint() uint64 { return x.Uint() }
func (x *stringValue) MustUint() uint64 { return 0 }
func (x *arrayValue) MustUint() uint64  { return 0 }
func (x *objectValue) MustUint() uint64 { return 0 }

func (x *errorValue) Interface() interface{} { return nil }
func (x *zeroValue) Interface() interface{}  { return nil }
func (x *nullValue) Interface() interface{}  { return nil }
func (x *boolValue) Interface() interface{}  { return x.Bool() }
func (x *numberValue) Interface() interface{} {
	if x.IsFloat() {
		return x.Float()
	}
	return x.Int()
}
func (x *stringValue) Interface() interface{} { return x.String() }
func (x *arrayValue) Interface() interface{}  { return nil }
func (x *objectValue) Interface() interface{} { return nil }

func (x *errorValue) IsNil() bool  { return false }
func (x *zeroValue) IsNil() bool   { return true }
func (x *nullValue) IsNil() bool   { return true }
func (x *boolValue) IsNil() bool   { return false }
func (x *numberValue) IsNil() bool { return false }
func (x *stringValue) IsNil() bool { return false }
func (x *arrayValue) IsNil() bool  { return len(x.values) == 0 }
func (x *objectValue) IsNil() bool { return len(x.members) == 0 }

func (x *errorValue) Len() int  { panic(fmt.Sprintf("%s has no len()", x.Kind())) }
func (x *zeroValue) Len() int   { return 0 }
func (x *nullValue) Len() int   { return 0 }
func (x *boolValue) Len() int   { panic(fmt.Sprintf("%s has no len()", x.Kind())) }
func (x *numberValue) Len() int { panic(fmt.Sprintf("%s has no len()", x.Kind())) }
func (x *stringValue) Len() int { return len(x.val) }
func (x *arrayValue) Len() int  { return len(x.values) }
func (x *objectValue) Len() int { return len(x.members) }

func (x *errorValue) Index(i int) Value  { return x }
func (x *zeroValue) Index(i int) Value   { return zero }
func (x *nullValue) Index(i int) Value   { return zero }
func (x *boolValue) Index(i int) Value   { panic(invalid_kind_error(Array, x.Kind())) }
func (x *numberValue) Index(i int) Value { panic(invalid_kind_error(Array, x.Kind())) }
func (x *stringValue) Index(i int) Value { panic(invalid_kind_error(Array, x.Kind())) }
func (x *objectValue) Index(i int) Value { panic(invalid_kind_error(Array, x.Kind())) }
func (x *arrayValue) Index(i int) Value {
	if i < len(x.values) {
		return x.values[i]
	}
	return zero
}

func (x *errorValue) MapIndex(key string) Value  { return x }
func (x *zeroValue) MapIndex(key string) Value   { return zero }
func (x *nullValue) MapIndex(key string) Value   { return zero }
func (x *boolValue) MapIndex(key string) Value   { panic(invalid_kind_error(Object, x.Kind())) }
func (x *numberValue) MapIndex(key string) Value { panic(invalid_kind_error(Object, x.Kind())) }
func (x *stringValue) MapIndex(key string) Value { panic(invalid_kind_error(Object, x.Kind())) }
func (x *arrayValue) MapIndex(key string) Value  { panic(invalid_kind_error(Object, x.Kind())) }
func (x *objectValue) MapIndex(key string) Value { return x.searchMember(key) }

func (x *errorValue) Path(parts ...interface{}) Value {
	return x
}
func (x *zeroValue) Path(parts ...interface{}) Value {
	if len(parts) == 0 {
		return x
	}
	return zero
}
func (x *nullValue) Path(parts ...interface{}) Value {
	if len(parts) == 0 {
		return x
	}
	return zero
}
func (x *boolValue) Path(parts ...interface{}) Value {
	if len(parts) == 0 {
		return x
	}
	return zero
}
func (x *numberValue) Path(parts ...interface{}) Value {
	if len(parts) == 0 {
		return x
	}
	return zero
}
func (x *stringValue) Path(parts ...interface{}) Value {
	if len(parts) == 0 {
		return x
	}
	return zero
}
func (x *arrayValue) Path(parts ...interface{}) Value {
	if len(parts) == 0 {
		return x
	}
	if idx, ok := parts[0].(int); ok {
		return x.Index(idx).Path(parts[1:]...)
	}
	return zero
}
func (x *objectValue) Path(parts ...interface{}) Value {
	if len(parts) == 0 {
		return x
	}
	if key, ok := parts[0].(string); ok {
		return x.MapIndex(key).Path(parts[1:]...)
	}
	return zero
}

func (x *numberValue) IsFloat() bool {
	return x.flags&(numberHasExponent|numberHasFraction) > 0
}

func (x *objectValue) searchMember(key string) Value {
	m := x.members
	i := sort.Search(len(m), func(i int) bool { return m[i].key >= key })
	if i < len(m) && m[i].key == key {
		return m[i].value
	} else {
		return zero
	}
}

type sortedObjectMembers []objectMember

func (l sortedObjectMembers) Len() int           { return len(l) }
func (l sortedObjectMembers) Less(i, j int) bool { return l[i].key < l[j].key }
func (l sortedObjectMembers) Swap(i, j int)      { l[i], l[j] = l[j], l[i] }

func invalid_kind_error(expected, actual Kind) error {
	return fmt.Errorf("%s is not a %s value", actual, expected)
}
