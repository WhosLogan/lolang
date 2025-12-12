package core

import (
	"fmt"
	"shared/pkg/types"
)

func stFld(ctx *functionCtx) {
	instr := ctx.Function.Instructions[ctx.InstrPtr]
	fieldIdx := int(instr.Operand.GetInt())
	val := ctx.Stack.Pop()
	structVal := ctx.Stack.Pop()

	structType := types.GetTypeByCode(structVal.Type)
	if structType == nil || !structType.IsStruct {
		ctx.Vm.error(fmt.Errorf("stfld: expected struct type, got %d", structVal.Type))
		return
	}

	if fieldIdx >= len(structType.Fields) {
		ctx.Vm.error(fmt.Errorf("stfld: field index %d out of range for struct %s", fieldIdx, structType.Name))
		return
	}

	fieldName := structType.Fields[fieldIdx].Name
	structVal.SetField(fieldName, val)
	ctx.Stack.Push(structVal)
}
