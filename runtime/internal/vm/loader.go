package vm

import (
	"errors"
	"runtime/internal/vm/core"
	"shared/pkg/types"
	"shared/pkg/vm"

	"github.com/vmihailenco/msgpack/v5"
)

func ExecuteProgram(program []byte) error {
	var vmData vm.Vm

	err := msgpack.Unmarshal(program, &vmData)
	if err != nil {
		return errors.New("an invalid program has been provided")
	}

	for _, st := range vmData.StructTypes {
		if existing := types.GetTypeByCode(st.Code); existing == nil {
			types.Types = append(types.Types, st)
		}
	}

	core.RunVm(&vmData)

	return nil
}
