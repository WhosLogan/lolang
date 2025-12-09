package handlers

import (
	"runtime/internal/vm"
	"shared/pkg/data"
)

var clt = Handler(func(ctx vm.FunctionCtx) {
	v1 := ctx.Stack.Pop()
	v2 := ctx.Stack.Pop()

	lt := false

	if v1.GetInt() < v2.GetInt() {
		lt = true
	}

	v, _ := data.NewValue(lt)
	ctx.Stack.Push(v)
})
