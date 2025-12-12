package core

import (
	"fmt"
	"shared/pkg/types"
)

func ldFld(ctx *functionCtx) {
	instr := ctx.Function.Instructions[ctx.InstrPtr]
	fieldIdx := int(instr.Operand.GetInt())
	structVal := ctx.Stack.Pop()

	structType := types.GetTypeByCode(structVal.Type)
	if structType == nil || !structType.IsStruct {
		ctx.Vm.error(fmt.Errorf("ldfld: expected struct type, got %d", structVal.Type))
		return
	}

	if fieldIdx >= len(structType.Fields) {
		ctx.Vm.error(fmt.Errorf("ldfld: field index %d out of range for struct %s", fieldIdx, structType.Name))
		return
	}

	fieldName := structType.Fields[fieldIdx].Name
	fieldVal, ok := structVal.GetField(fieldName)
	if !ok {
		ctx.Vm.error(fmt.Errorf("ldfld: field %s not found in struct", fieldName))
		return
	}

	ctx.Stack.Push(fieldVal)
}
