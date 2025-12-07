package handlers

import (
	"runtime/internal/vm"
)

var ldarg = Handler(func(ctx vm.FunctionCtx) {
	idxOp := ctx.Function.Instructions[ctx.InstrPtr].Operand
	idx := idxOp.GetInt()

	ctx.Stack.Push(ctx.Arguments[idx])
})
