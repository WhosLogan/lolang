package vm

import (
	"shared/pkg/function"
	"shared/pkg/types"
)

type Vm struct {
	EntryPoint  int
	Functions   map[int]*function.Function
	StructTypes []types.Type
}
