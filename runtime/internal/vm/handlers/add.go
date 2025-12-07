package handlers

import (
	"errors"
	"runtime/internal/vm/function"
	"shared/pkg/data"
	"shared/pkg/types"
)

var add = Handler(func(ctx function.Ctx) {
	var first = ctx.Stack.Pop()
	var second = ctx.Stack.Pop()

	if first.Type == types.LoInt && second.Type == types.LoInt {
		v, err := data.NewValue(first.GetInt() + second.GetInt())
		if err != nil {
			ctx.Vm.Error(err)
		}
		ctx.Stack.Push(v)
	}

	ctx.Vm.Error(errors.New("unable to add specified type pattern"))
})
