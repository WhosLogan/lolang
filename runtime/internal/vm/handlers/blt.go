package handlers

import "runtime/internal/vm"

var blt = Handler(func(ctx vm.FunctionCtx) {
	v1 := ctx.Stack.Pop()
	v2 := ctx.Stack.Pop()

	if v1.GetInt() > v2.GetInt() {
		return
	}

	instr := ctx.Function.Instructions[ctx.InstrPtr]
	token := int(instr.Operand.GetInt())

	ctx.InstrPtr = token
})
