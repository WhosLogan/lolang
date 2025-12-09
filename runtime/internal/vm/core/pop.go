package core

var pop = Handler(func(ctx *FunctionCtx) {
	ctx.Stack.Pop()
})
