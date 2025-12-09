package core

var bgt = Handler(func(ctx *FunctionCtx) {
	v1 := ctx.Stack.Pop()
	v2 := ctx.Stack.Pop()

	if v1.GetInt() <= v2.GetInt() {
		return
	}

	instr := ctx.Function.Instructions[ctx.InstrPtr]
	token := int(instr.Operand.GetInt())

	ctx.InstrPtr = token
})
