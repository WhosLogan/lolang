package handlers

import (
	"runtime/internal/vm"
)

var ldc8 = Handler(func(ctx vm.FunctionCtx) {
	value := ctx.Function.Instructions[ctx.InstrPtr].Operand
	ctx.Stack.Push(value)
})
