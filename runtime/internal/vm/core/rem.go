package core

import (
	"errors"
	"shared/pkg/data"
	"shared/pkg/types"
)

var rem = Handler(func(ctx *FunctionCtx) {
	var first = ctx.Stack.Pop()
	var second = ctx.Stack.Pop()

	if first.Type == types.LoInt && second.Type == types.LoInt {
		v, err := data.NewValue(first.GetInt() % second.GetInt())
		if err != nil {
			ctx.Vm.Error(err)
		}
		ctx.Stack.Push(v)
		return
	}

	ctx.Vm.Error(errors.New("unable to get the remainder for the specified type pattern"))
})
