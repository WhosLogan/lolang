package handlers

import (
	"runtime/internal/vm/function"
)

var ldc8 = Handler(func(ctx function.Ctx) {
	value := ctx.Function.Instructions[ctx.InstrPtr].Operand
	ctx.Stack.Push(value)
})
