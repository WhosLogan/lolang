package function

import (
	"shared/pkg/data"
	"shared/pkg/opcodes"
	"shared/pkg/types"
)

type Function struct {
	Token        int // Identifier of the function similar to mdtoken in MSIL
	Instructions map[int]Instruction
	Locals       []Local
	Arguments    []Argument
	ReturnType   types.Type
}

type Instruction struct {
	OpCode  opcodes.OpCode
	Operand data.Value
}

type Argument struct {
	Index int
	Name  string
	Type  types.Type
}

type Local struct {
	Index int
	Name  string
	Type  types.Type
}
