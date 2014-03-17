package xjson

import (
	"fmt"
	"sort"
)

func Parse(b []byte) Value {
	s := new_scanner(b)
	v, err := s.scan_value()
	if err != nil {
		v = &errorValue{err}
	}
	return v
}

type scanner struct {
	buf []byte
	pos int
	chr int
}

type numberFlags uint8

const (
	numberHasFraction numberFlags = 1 << iota
	numberHasExponent
	numberIsNegative
)

func new_scanner(b []byte) *scanner {
	if len(b) == 0 {
		return &scanner{chr: -1}
	}
	return &scanner{buf: b, chr: int(b[0])}
}

func (s *scanner) err(format string, a ...interface{}) error {
	return fmt.Errorf("xjson: %s (pos=%d)", fmt.Sprintf(format, a...), s.pos)
}

func (s *scanner) scan_value() (Value, error) {
	s.skip_whitespace()

	switch s.chr {
	case 'n': // null
		return s.scan_null()
	case 't': // true
		return s.scan_true()
	case 'f': // false
		return s.scan_false()
	case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9': // number
		return s.scan_number()
	case '"': // string
		return s.scan_string()
	case '[': // array
		return s.scan_array()
	case '{': // object
		return s.scan_object()
	default:
		// parse error
		return nil, s.err("unexpected byte %q", s.chr)
	}
}

func (s *scanner) scan_object() (*objectValue, error) {
	var (
		beg     int
		end     int
		members []objectMember
	)

	beg = s.pos
	if !s.scan_byte('{') {
		return nil, s.err("unexpected byte %q", s.chr)
	}

	s.skip_whitespace()

	if s.chr == '}' {
		s.next()
		end = s.pos
		return &objectValue{nil}, nil
	}

	for {
		key_value, err := s.scan_string()
		if err != nil {
			return nil, err
		}
		key_bytes, ok := unquoteBytes(key_value.buf)
		if !ok {
			return nil, s.err("invalid string %q", key_value.buf)
		}
		key := string(key_bytes)

		s.skip_whitespace()
		if !s.scan_byte(':') {
			return nil, s.err("unexpected byte %q", s.chr)
		}
		s.skip_whitespace()

		val, err := s.scan_value()
		if err != nil {
			return nil, err
		}
		members = append(members, objectMember{key, val})

		s.skip_whitespace()
		if s.chr == ',' {
			s.next()
			s.skip_whitespace()
		} else if s.chr == '}' {
			s.next()
			break
		} else {
			return nil, s.err("unexpected byte %q", s.chr)
		}
	}

	sort.Sort(sortedObjectMembers(members))

	end = s.pos
	return &objectValue{members}, nil
}

func (s *scanner) scan_array() (*arrayValue, error) {
	var (
		beg    int
		end    int
		values []Value
	)

	beg = s.pos
	if !s.scan_byte('[') {
		return nil, s.err("unexpected byte %q", s.chr)
	}

	s.skip_whitespace()
	if s.chr == ']' {
		s.next()
		end = s.pos
		return &arrayValue{values}, nil
	}

	for {
		val, err := s.scan_value()
		if err != nil {
			return nil, err
		}
		values = append(values, val)

		s.skip_whitespace()
		if s.chr == ',' {
			s.next()
			s.skip_whitespace()
		} else if s.chr == ']' {
			s.next()
			break
		} else {
			return nil, s.err("unexpected byte %q", s.chr)
		}
	}

	end = s.pos
	return &arrayValue{values}, nil
}

func (s *scanner) scan_null() (*nullValue, error) {
	var (
		beg int
		end int
	)

	beg = s.pos
	if !s.scan_byte('n') {
		return nil, s.err("unexpected byte %q", s.chr)
	}
	if !s.scan_byte('u') {
		return nil, s.err("unexpected byte %q", s.chr)
	}
	if !s.scan_byte('l') {
		return nil, s.err("unexpected byte %q", s.chr)
	}
	if !s.scan_byte('l') {
		return nil, s.err("unexpected byte %q", s.chr)
	}
	end = s.pos

	return &nullValue{s.buf[beg:end]}, nil
}

func (s *scanner) scan_false() (*boolValue, error) {
	var (
		beg int
		end int
	)

	beg = s.pos
	if !s.scan_byte('f') {
		return nil, s.err("unexpected byte %q", s.chr)
	}
	if !s.scan_byte('a') {
		return nil, s.err("unexpected byte %q", s.chr)
	}
	if !s.scan_byte('l') {
		return nil, s.err("unexpected byte %q", s.chr)
	}
	if !s.scan_byte('s') {
		return nil, s.err("unexpected byte %q", s.chr)
	}
	if !s.scan_byte('e') {
		return nil, s.err("unexpected byte %q", s.chr)
	}
	end = s.pos

	return &boolValue{s.buf[beg:end], false}, nil
}

func (s *scanner) scan_true() (*boolValue, error) {
	var (
		beg int
		end int
	)

	beg = s.pos
	if !s.scan_byte('t') {
		return nil, s.err("unexpected byte %q", s.chr)
	}
	if !s.scan_byte('r') {
		return nil, s.err("unexpected byte %q", s.chr)
	}
	if !s.scan_byte('u') {
		return nil, s.err("unexpected byte %q", s.chr)
	}
	if !s.scan_byte('e') {
		return nil, s.err("unexpected byte %q", s.chr)
	}
	end = s.pos

	return &boolValue{s.buf[beg:end], true}, nil
}

func (s *scanner) scan_string() (*stringValue, error) {
	var (
		beg int
		end int
	)

	beg = s.pos

	if !s.scan_byte('"') {

	}

	for {
		if s.chr == '"' {
			s.next()
			break
		} else if s.chr == '\\' {
			s.next()
			if s.chr == '"' || s.chr == '\\' || s.chr == '/' || s.chr == 'b' || s.chr == 'f' || s.chr == 'n' || s.chr == 'r' || s.chr == 't' {
				s.next()
			} else if s.chr == 'u' {
				if !s.scan_hex_digits() {
					return nil, s.err("unexpected byte %q", s.chr)
				}
			} else {
				return nil, s.err("unexpected byte %q", s.chr)
			}
		} else {
			s.next()
		}
	}

	end = s.pos

	b, ok := unquoteBytes(s.buf[beg:end])
	if !ok {
		return nil, s.err("invalid string: %q", s.buf[beg:end])
	}

	return &stringValue{s.buf[beg:end], string(b)}, nil
}

func (s *scanner) scan_number() (*numberValue, error) {
	var (
		beg   int
		end   int
		flags numberFlags
	)

	beg = s.pos
	if s.scan_byte('-') {
		flags |= numberIsNegative
	}

	if s.scan_byte('0') {
		// ok
	} else if '1' <= s.chr && s.chr <= '9' {
		s.scan_dec_digits()
	} else {
		return nil, s.err("unexpected byte %q", s.chr)
	}

	if s.scan_byte('.') {
		flags |= numberHasFraction
		if !s.scan_dec_digits() {
			return nil, s.err("unexpected byte %q", s.chr)
		}
	}

	if s.chr == 'e' || s.chr == 'E' {
		s.next()
		flags |= numberHasExponent
		s.scan_sign()
		if !s.scan_dec_digits() {
			return nil, s.err("unexpected byte %q", s.chr)
		}
	}

	end = s.pos
	return &numberValue{s.buf[beg:end], flags}, nil
}

func (s *scanner) scan_byte(c byte) bool {
	if s.chr == int(c) {
		return s.next()
	}
	return false
}

func (s *scanner) scan_dec_digits() bool {
	ok := false
	for '0' <= s.chr && s.chr <= '9' {
		ok = true
		s.next()
	}
	return ok
}

func (s *scanner) scan_hex_digits() bool {
	ok := false
	for '0' <= s.chr && s.chr <= '9' || 'A' <= s.chr && s.chr <= 'F' || 'a' <= s.chr && s.chr <= 'f' {
		ok = true
		s.next()
	}
	return ok
}

func (s *scanner) scan_sign() bool {
	if s.chr == '+' || s.chr == '-' {
		return s.next()
	}
	return true
}

func (s *scanner) skip_whitespace() {
	for {
		if s.chr == ' ' || s.chr == '\t' || s.chr == '\f' || s.chr == '\r' || s.chr == '\n' {
			s.next()
		} else {
			break
		}
	}
}

func (s *scanner) next() bool {
	if s.pos == len(s.buf) {
		s.chr = -1
		return false
	}
	s.pos++
	if s.pos == len(s.buf) {
		s.chr = -1
		return false
	}
	s.chr = int(s.buf[s.pos])
	return true
}
