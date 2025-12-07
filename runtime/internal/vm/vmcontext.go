package vm

import "runtime/internal/vm/function"

type VCtx struct {
	Functions map[int]function.Function
}
