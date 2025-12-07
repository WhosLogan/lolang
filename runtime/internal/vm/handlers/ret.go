package handlers

import (
	"runtime/internal/vm"
	"shared/pkg/types"
)

var ret = Handler(func(ctx vm.FunctionCtx) {
	if ctx.Vm.CallStack.Len() == 0 {
		return
	}

	prev := ctx.Vm.CallStack.Pop()
	if ctx.Function.ReturnType == types.LoVoid {
		return
	}

	prev.Stack.Push(ctx.Stack.Pop())
})
