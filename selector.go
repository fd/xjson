package xjson

import (
	"fmt"
	"unicode"
)

type Selector interface {
	Value() Value
	String() string
}

type root_selector struct {
	value interface{}
}

func (i *root_selector) Value() Value {
	return Value{i.value, nil, i}
}

func (i *root_selector) String() string {
	return "$root"
}

type index_selector struct {
	value  interface{}
	idx    int
	parent Selector
}

func (i *index_selector) Value() Value {
	return Value{i.value, nil, i}
}

func (i *index_selector) String() string {
	return fmt.Sprintf("%s[%d]", i.parent, i.idx)
}

type key_selector struct {
	value  interface{}
	key    string
	parent Selector
}

func (i *key_selector) Value() Value {
	return Value{i.value, nil, i}
}

func (i *key_selector) String() string {
	if is_keyword(i.key) {
		return fmt.Sprintf("%s.%s", i.parent, i.key)
	} else {
		return fmt.Sprintf("%s[%q]", i.parent, i.key)
	}
}

func is_keyword(s string) bool {
	for i, r := range s {
		if i == 0 {
			if !unicode.IsLetter(r) && r != '_' {
				return false
			}
		} else {
			if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' {
				return false
			}
		}
	}
	return true
}
