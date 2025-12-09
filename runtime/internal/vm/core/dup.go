package core

var dup = Handler(func(ctx *FunctionCtx) {
	ctx.Stack.Push(ctx.Stack.Peek())
})
