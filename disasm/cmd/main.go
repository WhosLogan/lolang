package main

import (
	"disasm/internal"
	"os"
)

func main() {
	data, _ := os.ReadFile(os.Args[1])
	_ = internal.Disassemble(data)
}
