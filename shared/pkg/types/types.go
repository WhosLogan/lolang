package types

const (
	LoInt = TypeCode(iota)
	LoVoid
	LoString
	LoBool
	LoStructStart = TypeCode(10000)
)

var Types = []Type{
	{
		Name: "int",
		Code: LoInt,
		Size: 8,
	},
	{
		Name: "bool",
		Code: LoBool,
		Size: 1,
	},
	{
		Name: "void",
		Code: LoVoid,
		Size: 0,
	},
	{
		Name: "string",
		Code: LoString,
	},
}

var nextStructCode = LoStructStart

func RegisterStructType(name string, fields []StructField) *Type {
	t := Type{
		Name:     name,
		Code:     nextStructCode,
		IsStruct: true,
		Fields:   fields,
	}
	nextStructCode++
	Types = append(Types, t)
	return &Types[len(Types)-1]
}

func GetTypeByName(name string) *Type {
	for i := range Types {
		if Types[i].Name == name {
			return &Types[i]
		}
	}

	return nil
}

func GetTypeByCode(code TypeCode) *Type {
	for i := range Types {
		if Types[i].Code == code {
			return &Types[i]
		}
	}

	return nil
}

func ResetStructTypes() {
	newTypes := make([]Type, 0)
	for _, t := range Types {
		if !t.IsStruct {
			newTypes = append(newTypes, t)
		}
	}
	Types = newTypes
	nextStructCode = LoStructStart
}
