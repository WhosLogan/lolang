package core

var ldarg = Handler(func(ctx *FunctionCtx) {
	idxOp := ctx.Function.Instructions[ctx.InstrPtr].Operand
	idx := idxOp.GetInt()

	ctx.Stack.Push(ctx.Arguments[idx])
})
