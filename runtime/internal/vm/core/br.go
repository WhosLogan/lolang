package core

var br = Handler(func(ctx *FunctionCtx) {
	instr := ctx.Function.Instructions[ctx.InstrPtr]
	token := instr.Operand.GetInt()

	ctx.InstrPtr = int(token)
})
