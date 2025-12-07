package handlers

import (
	"runtime/internal/vm"
)

var nop = Handler(func(ctx vm.FunctionCtx) {
	// Do nothing
})
