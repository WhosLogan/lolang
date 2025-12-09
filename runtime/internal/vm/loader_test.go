package vm

import (
	"shared/pkg/function"
	"shared/pkg/opcodes"
	"shared/pkg/types"
	"shared/pkg/vm"
	"testing"

	"github.com/vmihailenco/msgpack/v5"
)

func TestExecuteProgram(t *testing.T) {
	v := vm.Vm{
		EntryPoint: 0,
		Functions: map[int]function.Function{
			0: {
				Token: 0,
				Instructions: map[int]function.Instruction{
					0: {OpCode: opcodes.Nop},
					1: {OpCode: opcodes.Ret},
				},
				ReturnType: types.LoVoid,
			},
		},
	}

	data, err := msgpack.Marshal(v)
	if err != nil {
		t.Error(err)
	}

	err = ExecuteProgram(data)
	if err != nil {
		t.Error(err)
	}
}
