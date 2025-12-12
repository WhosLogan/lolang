package core

import (
	"shared/pkg/data"
	"shared/pkg/types"
)

func newObj(ctx *functionCtx) {
	instr := ctx.Function.Instructions[ctx.InstrPtr]
	typeCode := types.TypeCode(instr.Operand.GetInt())
	structVal := data.NewStructValue(typeCode)
	ctx.Stack.Push(structVal)
}
