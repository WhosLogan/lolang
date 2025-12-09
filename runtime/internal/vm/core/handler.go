package core

import (
	"shared/pkg/opcodes"
)

type Handler func(ctx *FunctionCtx)

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

	// Branch operations
	opcodes.Br:  br,
	opcodes.Be:  be,
	opcodes.Bne: bne,
	opcodes.Bgt: bgt,
	opcodes.Blt: blt,

	// Equality comparisons
	opcodes.Cmp: cmp,
	opcodes.Clt: clt,
	opcodes.Cgt: cgt,

	// Stack control
	opcodes.Ldc8:  ldc8,
	opcodes.Ldarg: ldarg,
	opcodes.Pop:   pop,
	opcodes.Dup:   dup,

	// Misc
	opcodes.Nop: nop,
	opcodes.Ret: ret,
}
