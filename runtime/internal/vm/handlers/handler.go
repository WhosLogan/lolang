package handlers

import (
	"runtime/internal/vm/function"
	"shared/pkg/opcodes"
)

type Handler func(ctx function.Ctx)

var Handlers = map[opcodes.OpCode]Handler{
	opcodes.Pop: pop,
	opcodes.Add: add,
}
