package types

type Type struct {
	Name     string
	Code     TypeCode
	Size     int
	IsStruct bool
	Fields   []StructField
}

type StructField struct {
	Name string
	Type TypeCode
}

type TypeCode int
