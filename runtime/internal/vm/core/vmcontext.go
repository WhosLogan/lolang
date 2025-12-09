package core

import (
	"shared/pkg/data"
	"shared/pkg/function"
)

type VCtx struct {
	EntryPoint *function.Function
	Functions  map[int]function.Function
	CallStack  data.Stack[*FunctionCtx]
}

type FunctionCtx struct {
	Vm        *VCtx
	InstrPtr  int
	Function  *function.Function
	Stack     data.Stack[data.Value]
	Locals    map[int]data.Value
	Arguments []data.Value
	Running   bool
}

// Error halts the virtual machine and kills the running application
func (vm *VCtx) Error(err error) {
	panic(err)
}

func (vm *VCtx) Execute() {
	ctx := NewFunctionCtx(vm, vm.EntryPoint)
	ctx.Execute()
}

func NewFunctionCtx(vm *VCtx, fn *function.Function) *FunctionCtx {
	var locals map[int]data.Value
	for _, local := range fn.Locals {
		if !local.HasInitialValue {
			continue
		}

		locals[local.Index] = local.InitialValue
	}

	return &FunctionCtx{
		Vm:       vm,
		InstrPtr: 0,
		Function: fn,
		Stack:    data.Stack[data.Value]{},
		Locals:   locals,
	}
}

func (fn *FunctionCtx) Execute() {
	for fn.Running {
		ptr := fn.InstrPtr

		instr := fn.Function.Instructions[fn.InstrPtr]

		// Execute the instruction
		Handlers[instr.OpCode](fn)

		// If the ptr has been changed than a branch instruction was called and changed the flow
		if ptr != fn.InstrPtr {
			continue
		}

		fn.InstrPtr++
	}
}
