package handlers

import (
	"runtime/internal/vm/function"
	"shared/pkg/opcodes"
)

type Handler func(ctx function.Ctx)

var Handlers = map[opcodes.OpCode]Handler{
	// Arithmetic instructions
	opcodes.Add: add,
	opcodes.Sub: sub,
	opcodes.Div: div,
	opcodes.Mul: mul,
	opcodes.Rem: rem,
	opcodes.Xor: xor,
	opcodes.Not: not,
	opcodes.Or:  or,
	opcodes.Neg: neg,

	// Stack control
	opcodes.Ldc8:  ldc8,
	opcodes.Ldarg: ldarg,
	opcodes.Pop:   pop,
	opcodes.Dup:   dup,
}
