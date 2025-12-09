package vm

import (
	"shared/pkg/data"
	"shared/pkg/function"
	"shared/pkg/opcodes"
	"shared/pkg/types"
	"shared/pkg/vm"
	"testing"

	"github.com/vmihailenco/msgpack/v5"
)

func TestExecuteProgram(t *testing.T) {
	num1, _ := data.NewValue(int64(5))
	num2, _ := data.NewValue(int64(15))
	num3, _ := data.NewValue(int64(1000001))

	v := vm.Vm{
		EntryPoint: 1000000,
		Functions: map[int]function.Function{
			1000000: {
				Token: 1000000,
				Instructions: map[int]function.Instruction{
					0: {OpCode: opcodes.Call, Operand: num3},
					1: {OpCode: opcodes.Ret},
				},
				ReturnType: types.LoVoid,
			},
			1000001: {
				Token: 1000001,
				Instructions: map[int]function.Instruction{
					0: {OpCode: opcodes.Ldc8, Operand: num1},
					1: {OpCode: opcodes.Ldc8, Operand: num2},
					2: {OpCode: opcodes.Mul},
					3: {OpCode: opcodes.Pop},
					4: {OpCode: opcodes.Ret},
				},
				ReturnType: types.LoVoid,
			},
		},
	}

	vmd, err := msgpack.Marshal(v)
	if err != nil {
		t.Error(err)
	}

	err = ExecuteProgram(vmd)
	if err != nil {
		t.Error(err)
	}
}
