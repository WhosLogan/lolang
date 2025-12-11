package opcodes

//go:generate stringer -type=OpCode
type OpCode int

// Arithmetic operations
const (
	Add OpCode = iota + 100
	Sub
	Mul
	Div
	Rem
	Xor
	Not
	Or
	Neg
)

// Branch Operations
const (
	Br OpCode = iota + 200
	Be
	Bne
	Blt
	Bgt
	Bt
	Bf
)

// Equality
const (
	Cmp OpCode = iota + 300
	Clt
	Cgt
	Cle
	Cge
)

// Stack Control
const (
	Ldc8 OpCode = iota + 400
	LdArg
	Pop
	Dup
	LdLoc
	StLoc
	LdStr
)

// Misc
const (
	Nop OpCode = iota + 500
	Ret
	Call
)
