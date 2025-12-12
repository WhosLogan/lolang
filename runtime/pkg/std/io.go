package std

import (
	"fmt"
	"shared/pkg/data"
	"shared/pkg/function"
	"shared/pkg/types"
)

var printFunc = stdFunc{
	Name: "print",
	Arguments: []function.Argument{
		{
			Index: 0,
			Name:  "str",
			Type:  types.LoString,
		},
	},
	Execute: func(values []data.Value) data.Value {
		fmt.Print(values[0].GetString())

		return data.Value{
			Type: types.LoVoid,
		}
	},
	ReturnType: types.LoVoid,
}

var printLnFunc = stdFunc{
	Name: "println",
	Arguments: []function.Argument{
		{
			Index: 0,
			Name:  "str",
			Type:  types.LoString,
		},
	},
	Execute: func(values []data.Value) data.Value {
		fmt.Println(values[0].GetString())

		return data.Value{
			Type: types.LoVoid,
		}
	},
	ReturnType: types.LoVoid,
}

var printIntFunc = stdFunc{
	Name: "printInt",
	Arguments: []function.Argument{
		{
			Index: 0,
			Name:  "n",
			Type:  types.LoInt,
		},
	},
	Execute: func(values []data.Value) data.Value {
		fmt.Println(values[0].GetInt())

		return data.Value{
			Type: types.LoVoid,
		}
	},
	ReturnType: types.LoVoid,
}