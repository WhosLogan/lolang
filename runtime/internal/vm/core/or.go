package core

import (
	"errors"
	"shared/pkg/data"
	"shared/pkg/types"
)

var or = Handler(func(ctx *FunctionCtx) {
	var first = ctx.Stack.Pop()
	var second = ctx.Stack.Pop()

	if first.Type == types.LoInt && second.Type == types.LoInt {
		v, err := data.NewValue(first.GetInt() | second.GetInt())
		if err != nil {
			ctx.Vm.Error(err)
		}
		ctx.Stack.Push(v)
	}

	ctx.Vm.Error(errors.New("unable to or the specified type pattern"))
})
