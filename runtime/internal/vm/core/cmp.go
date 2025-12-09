package core

import (
	"shared/pkg/data"
)

var cmp = Handler(func(ctx *FunctionCtx) {
	v1 := ctx.Stack.Pop()
	v2 := ctx.Stack.Pop()

	eq := false

	if v1.Equal(&v2) {
		eq = true
	}

	v, _ := data.NewValue(eq)
	ctx.Stack.Push(v)
})
