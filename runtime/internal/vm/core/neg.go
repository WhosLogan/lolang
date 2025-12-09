package core

import (
	"errors"
	"shared/pkg/data"
	"shared/pkg/types"
)

var neg = Handler(func(ctx *FunctionCtx) {
	var num = ctx.Stack.Pop()

	if num.Type == types.LoInt {
		v, err := data.NewValue(-num.GetInt())
		if err != nil {
			ctx.Vm.Error(err)
		}
		ctx.Stack.Push(v)
	}

	ctx.Vm.Error(errors.New("unable to negate the specified type"))
})
