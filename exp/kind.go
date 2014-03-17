package xjson

type Kind uint8

const (
	Null Kind = 1 + iota
	Bool
	Number
	String
	Array
	Object

	// special kind
	Error
)

func (k Kind) String() string {
	return kindStrings[k]
}

var kindStrings = map[Kind]string{
	Null:   "Null",
	Bool:   "Bool",
	Number: "Number",
	String: "String",
	Array:  "Array",
	Object: "Object",

	Error: "Error",
}
