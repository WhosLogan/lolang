package vm

import (
	"shared/pkg/data"
	"shared/pkg/function"
)

type VCtx struct {
	EntryPoint *function.Function
	Functions  map[int]function.Function
	CallStack  data.Stack[*FunctionCtx]
}

type FunctionCtx struct {
	Vm        *VCtx
	InstrPtr  int
	Function  *function.Function
	Stack     data.Stack[data.Value]
	Locals    []data.Value
	Arguments []data.Value
}

// Error halts the virtual machine and kills the running application
func (vm *VCtx) Error(err error) {
	panic(err)
}
