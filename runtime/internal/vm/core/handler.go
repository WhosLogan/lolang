package core

import (
	"shared/pkg/opcodes"
)

type handler func(ctx *functionCtx)

var handlers = map[opcodes.OpCode]handler{
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
	opcodes.Bt:  bt,
	opcodes.Bf:  bf,

	// Equality comparisons
	opcodes.Cmp: cmp,
	opcodes.Clt: clt,
	opcodes.Cgt: cgt,
	opcodes.Cle: cle,
	opcodes.Cge: cge,

	// Stack control
	opcodes.Ldc8:  ldc8,
	opcodes.Pop:   pop,
	opcodes.Dup:   dup,
	opcodes.LdLoc: ldLoc,
	opcodes.StLoc: stLoc,
	opcodes.LdStr: ldStr,

	// Struct operations
	opcodes.NewObj: newObj,
	opcodes.LdFld:  ldFld,
	opcodes.StFld:  stFld,

	// Misc
	opcodes.Nop: nop,
	opcodes.Ret: ret,
}

// Avoid cyclic dependency
func init() {
	handlers[opcodes.Call] = call
}
