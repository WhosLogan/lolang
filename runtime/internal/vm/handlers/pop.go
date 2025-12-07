package handlers

import (
	"runtime/internal/vm"
)

var pop = Handler(func(ctx vm.FunctionCtx) {
	ctx.Stack.Pop()
})
