package handlers

import (
	"runtime/internal/vm"
)

var dup = Handler(func(ctx vm.FunctionCtx) {
	ctx.Stack.Push(ctx.Stack.Peek())
})
