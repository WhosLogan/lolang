package internal

import (
	"errors"
	"fmt"
	"shared/pkg/function"
	"shared/pkg/types"
	"shared/pkg/vm"
	"strconv"

	"github.com/vmihailenco/msgpack/v5"
)

func Disassemble(program []byte) error {
	var vmData vm.Vm

	err := msgpack.Unmarshal(program, &vmData)
	if err != nil {
		return errors.New("an invalid program has been provided")
	}

	for i, fn := range vmData.Functions {
		fmt.Printf("Fn #%d: %s\n", i, fn.Name)

		for j := range len(fn.Instructions) {
			fmt.Printf("Instr #%d: Opcode: %s Operand: %s\n", j, fn.Instructions[j].OpCode.String(), tryGetOperand(fn.Instructions[j]))
		}
	}

	return nil
}

func tryGetOperand(instr function.Instruction) string {
	if instr.Operand.Data == nil {
		return "nil"
	}

	switch instr.Operand.Type {
	case types.LoBool:
		if instr.Operand.GetBool() {
			return "true"
		} else {
			return "false"
		}
	case types.LoInt:
		return strconv.Itoa(int(instr.Operand.GetInt()))
	case types.LoString:
		return instr.Operand.GetString()
	default:
		return "N/A"
	}
}
