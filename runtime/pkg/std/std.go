package std

import (
	"shared/pkg/data"
	"shared/pkg/function"
	"shared/pkg/types"
)

type stdFunc struct {
	Name       string
	Execute    func(values []data.Value) data.Value
	Arguments  []function.Argument
	ReturnType types.TypeCode
}

var StdLib = map[int]stdFunc{
	0: printFunc,
	1: printLnFunc,
	2: printIntFunc,
}
