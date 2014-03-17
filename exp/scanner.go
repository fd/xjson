package xjson

import (
	"fmt"
)

type Value interface {
	jsonValue()
}

type Null struct {
	buf []byte
	beg int
	end int
}
type Bool struct {
	buf []byte
	val bool
	beg int
	end int
}
type Number struct {
	buf []byte
	beg int
	end int
}
type String struct {
	buf []byte
	beg int
	end int
}
type Array struct {
	buf    []byte
	values []Value
	beg    int
	end    int
}
type Object struct {
	buf    []byte
	keys   []Value
	values []Value
	beg    int
	end    int
}

func (*Null) jsonValue()   {}
func (*Bool) jsonValue()   {}
func (*Number) jsonValue() {}
func (*String) jsonValue() {}
func (*Array) jsonValue()  {}
func (*Object) jsonValue() {}

type scanner struct {
	buf []byte
	pos int
	chr int
}

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
		return nil, s.err("unexpected byte %x", s.chr)
	}
}

func (s *scanner) scan_object() (*Object, error) {
	var (
		beg    int
		end    int
		keys   []Value
		values []Value
	)

	beg = s.pos
	if !s.scan_byte('{') {
		return nil, s.err("unexpected byte %x", s.chr)
	}

	s.skip_whitespace()

	for {
		key, err := s.scan_string()
		if err != nil {
			return nil, err
		}
		keys = append(keys, key)

		s.skip_whitespace()
		if !s.scan_byte(':') {
			return nil, s.err("unexpected byte %x", s.chr)
		}
		s.skip_whitespace()

		val, err := s.scan_value()
		if err != nil {
			return nil, err
		}
		values = append(values, val)

		s.skip_whitespace()
		if s.chr == ',' {
			s.next()
			s.skip_whitespace()
		} else if s.chr == '}' {
			s.next()
			break
		} else {
			return nil, s.err("unexpected byte %x", s.chr)
		}
	}

	end = s.pos
	return &Object{s.buf[beg:end], keys, values, beg, end}, nil
}

func (s *scanner) scan_array() (*Array, error) {
	var (
		beg    int
		end    int
		values []Value
	)

	beg = s.pos
	if !s.scan_byte('[') {
		return nil, s.err("unexpected byte %x", s.chr)
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
			return nil, s.err("unexpected byte %x", s.chr)
		}
	}

	end = s.pos
	return &Array{s.buf[beg:end], values, beg, end}, nil
}

func (s *scanner) scan_null() (*Null, error) {
	var (
		beg int
		end int
	)

	beg = s.pos
	if !s.scan_byte('n') {
		return nil, s.err("unexpected byte %x", s.chr)
	}
	if !s.scan_byte('u') {
		return nil, s.err("unexpected byte %x", s.chr)
	}
	if !s.scan_byte('l') {
		return nil, s.err("unexpected byte %x", s.chr)
	}
	if !s.scan_byte('l') {
		return nil, s.err("unexpected byte %x", s.chr)
	}
	end = s.pos

	return &Null{s.buf[beg:end], beg, end}, nil
}

func (s *scanner) scan_false() (*Bool, error) {
	var (
		beg int
		end int
	)

	beg = s.pos
	if !s.scan_byte('f') {
		return nil, s.err("unexpected byte %x", s.chr)
	}
	if !s.scan_byte('a') {
		return nil, s.err("unexpected byte %x", s.chr)
	}
	if !s.scan_byte('l') {
		return nil, s.err("unexpected byte %x", s.chr)
	}
	if !s.scan_byte('s') {
		return nil, s.err("unexpected byte %x", s.chr)
	}
	if !s.scan_byte('e') {
		return nil, s.err("unexpected byte %x", s.chr)
	}
	end = s.pos

	return &Bool{s.buf[beg:end], false, beg, end}, nil
}

func (s *scanner) scan_true() (*Bool, error) {
	var (
		beg int
		end int
	)

	beg = s.pos
	if !s.scan_byte('t') {
		return nil, s.err("unexpected byte %x", s.chr)
	}
	if !s.scan_byte('r') {
		return nil, s.err("unexpected byte %x", s.chr)
	}
	if !s.scan_byte('u') {
		return nil, s.err("unexpected byte %x", s.chr)
	}
	if !s.scan_byte('e') {
		return nil, s.err("unexpected byte %x", s.chr)
	}
	end = s.pos

	return &Bool{s.buf[beg:end], true, beg, end}, nil
}

func (s *scanner) scan_string() (*String, error) {
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
					return nil, s.err("unexpected byte %x", s.chr)
				}
			} else {
				return nil, s.err("unexpected byte %x", s.chr)
			}
		} else {
			s.next()
		}
	}

	end = s.pos
	return &String{s.buf[beg:end], beg, end}, nil
}

func (s *scanner) scan_number() (*Number, error) {
	var (
		beg int
		end int
	)

	beg = s.pos
	s.scan_byte('-')

	if s.scan_byte('0') {
		// ok
	} else if '1' <= s.chr && s.chr <= '9' {
		s.scan_dec_digits()
	} else {
		return nil, s.err("unexpected byte %x", s.chr)
	}

	if s.scan_byte('.') {
		if !s.scan_dec_digits() {
			return nil, s.err("unexpected byte %x", s.chr)
		}
	}

	if s.chr == 'e' || s.chr == 'E' {
		s.next()
		s.scan_sign()
		if !s.scan_dec_digits() {
			return nil, s.err("unexpected byte %x", s.chr)
		}
	}

	end = s.pos
	return &Number{s.buf[beg:end], beg, end}, nil
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
		if s.chr == ' ' || s.chr == '\t' || s.chr == '\r' || s.chr == '\n' {
			s.next()
		}
		break
	}
}

func (s *scanner) next() bool {
	if s.pos < len(s.buf) {
		s.chr = int(s.buf[s.pos])
		s.pos++
		return true
	}
	s.chr = -1
	return false
}
