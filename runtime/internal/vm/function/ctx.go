package function

import (
	"runtime/internal/vm"
	"shared/pkg/data"
	"shared/pkg/function"
)

type Ctx struct {
	Vm        *vm.VCtx
	Function  *function.Function
	Stack     data.Stack[data.Value]
	Locals    []data.Value
	Arguments []data.Value
}
