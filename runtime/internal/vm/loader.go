package vm

import (
	"errors"
	"runtime/internal/vm/core"
	"shared/pkg/vm"

	"github.com/vmihailenco/msgpack/v5"
)

func ExecuteProgram(program []byte) error {
	var vmData vm.Vm

	err := msgpack.Unmarshal(program, &vmData)
	if err != nil {
		return errors.New("an invalid program has been provided")
	}

	core.NewVmCtx(&vmData).Execute()

	return nil
}
