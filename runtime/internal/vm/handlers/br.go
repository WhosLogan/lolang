package handlers

import "runtime/internal/vm"

var br = Handler(func(ctx vm.FunctionCtx) {
	instr := ctx.Function.Instructions[ctx.InstrPtr]
	token := instr.Operand.GetInt()

	ctx.InstrPtr = int(token)
})
