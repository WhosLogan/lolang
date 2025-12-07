package handlers

import "runtime/internal/vm/function"

var ldarg = Handler(func(ctx function.Ctx) {
	idxOp := ctx.Function.Instructions[ctx.InstrPtr].Operand
	idx := idxOp.GetInt()

	ctx.Stack.Push(ctx.Arguments[idx])
})
