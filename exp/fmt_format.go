package xjson

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
)

func format_str(f fmt.State, c rune) string {
	format := "%"

	if f.Flag('+') {
		format += "+"
	}
	if f.Flag('-') {
		format += "-"
	}
	if f.Flag('#') {
		format += "#"
	}
	if f.Flag(' ') {
		format += " "
	}
	if f.Flag('0') {
		format += "0"
	}

	if w, ok := f.Width(); ok {
		format += strconv.Itoa(w)
	}

	if w, ok := f.Precision(); ok {
		format += "." + strconv.Itoa(w)
	}

	format += string([]rune{c})

	return format
}

func (x *nullValue) Format(f fmt.State, c rune) {
	if c == 'j' {
		if f.Flag('#') {
			f.Write(x.buf)
			return
		} else {
			var (
				buf bytes.Buffer
				err = json.Compact(&buf, x.buf)
			)
			if err != nil {
				panic(err)
			}
			buf.WriteTo(f)
			return
		}
	}
	fmt.Fprintf(f, format_str(f, c), nil)
}

func (x *boolValue) Format(f fmt.State, c rune) {
	if c == 'j' {
		if f.Flag('#') {
			f.Write(x.buf)
			return
		} else {
			var (
				buf bytes.Buffer
				err = json.Compact(&buf, x.buf)
			)
			if err != nil {
				panic(err)
			}
			buf.WriteTo(f)
			return
		}
	}
	fmt.Fprintf(f, format_str(f, c), x.val)
}

func (x *numberValue) Format(f fmt.State, c rune) {
	if c == 'j' {
		if f.Flag('#') {
			f.Write(x.buf)
			return
		} else {
			var (
				buf bytes.Buffer
				err = json.Compact(&buf, x.buf)
			)
			if err != nil {
				panic(err)
			}
			buf.WriteTo(f)
			return
		}
	}
	if c == 'b' || c == 'e' || c == 'E' || c == 'f' || c == 'g' || c == 'G' {
		fmt.Fprintf(f, format_str(f, c), x.Float())
	} else {
		fmt.Fprintf(f, format_str(f, c), x.Int())
	}
}

func (x *stringValue) Format(f fmt.State, c rune) {
	if c == 'j' {
		if f.Flag('#') {
			f.Write(x.buf)
			return
		} else {
			var (
				buf bytes.Buffer
				err = json.Compact(&buf, x.buf)
			)
			if err != nil {
				panic(err)
			}
			buf.WriteTo(f)
			return
		}
	}
	fmt.Fprintf(f, format_str(f, c), x.val)
}

func (x *arrayValue) Format(f fmt.State, c rune) {
	if c == 'j' {
		if f.Flag('#') {
			f.Write(x.buf)
			return
		} else {
			var (
				buf bytes.Buffer
				err = json.Compact(&buf, x.buf)
			)
			if err != nil {
				panic(err)
			}
			buf.WriteTo(f)
			return
		}
	}
	fmt.Fprintf(f, format_str(f, c), x.values)
}

func (x *objectValue) Format(f fmt.State, c rune) {
	if c == 'j' {
		if f.Flag('#') {
			f.Write(x.buf)
			return
		} else {
			var (
				buf bytes.Buffer
				err = json.Compact(&buf, x.buf)
			)
			if err != nil {
				panic(err)
			}
			buf.WriteTo(f)
			return
		}
	}
	fmt.Fprintf(f, format_str(f, c), x.members)
}
