package model

// Struct is the struct to be generated
type Struct struct {
	Type        string // Type is the type of the struct. e.g. Channel
	CustomTypes []*CustomType
	Fields      []*Field
}

// Field defines a single field of an object
type Field struct {
	Name string              // field name
	Type string              // type name as a string. e.g. int, float64 etc.
	Tags map[string][]string // tags associated with values
}

type CustomType struct {
	Name           string // type name e.g. ChannelType
	UnderlyingType string // like string, int etc.
}
