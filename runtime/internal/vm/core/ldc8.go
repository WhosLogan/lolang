package core

var ldc8 = Handler(func(ctx *FunctionCtx) {
	value := ctx.Function.Instructions[ctx.InstrPtr].Operand
	ctx.Stack.Push(value)
})
