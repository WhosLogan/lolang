package compiler

import (
	"fmt"
	"log"
	"runtime/pkg/std"
	"shared/pkg/data"
	"shared/pkg/function"
	"shared/pkg/opcodes"
	"shared/pkg/types"
	"shared/pkg/vm"
)

type codegen struct {
	vm       *vm.Vm
	curFunc  *Function
	curVmFn  *function.Function
	nextSlot int
	structs  map[string]*StructDef
}

func (gen *codegen) compileFunction() error {
	gen.curVmFn.Instructions = map[int]function.Instruction{}

	paramSlots := make([]int, len(gen.curFunc.Params))
	for i, p := range gen.curFunc.Params {
		slot := gen.allocLocal(p.Name, p.Type)
		paramSlots[i] = slot
	}

	for i := 0; i < len(paramSlots); i++ {
		slot := paramSlots[i]
		gen.emit(opcodes.StLoc, data.MustNewValue(slot))
	}

	if err := gen.compileBlock(gen.curFunc.Body); err != nil {
		return err
	}

	if len(gen.curVmFn.Instructions) == 0 || gen.curVmFn.Instructions[len(gen.curVmFn.Instructions)-1].OpCode != opcodes.Ret {
		gen.emit(opcodes.Ret, data.Value{})
	}

	return nil
}

func (gen *codegen) allocLocal(name string, t string) int {
	for _, local := range gen.curVmFn.Locals {
		if local.Name == name {
			return local.Index
		}
	}

	localType := types.GetTypeByName(t)
	if localType == nil {
		log.Fatalf("unable to find local type: %s\n", t)
	}

	gen.curVmFn.Locals = append(gen.curVmFn.Locals, &function.Local{
		Index: gen.nextSlot,
		Name:  name,
		Type:  *localType,
	})

	gen.nextSlot++

	return gen.nextSlot - 1
}

func (gen *codegen) compileBlock(b *Block) error {
	for _, st := range b.Statements {
		if err := gen.compileStatement(st); err != nil {
			return err
		}
	}
	return nil
}

func (gen *codegen) emit(op opcodes.OpCode, operand data.Value) {
	idx := len(gen.curVmFn.Instructions)
	gen.curVmFn.Instructions[idx] = function.Instruction{
		OpCode:  op,
		Operand: operand,
	}
}

func (gen *codegen) compileStatement(s *Statement) error {
	switch {
	case s.Var != nil:
		v := s.Var
		slot := gen.allocLocal(v.Name, v.Type)
		if v.Init != nil {
			if err := gen.compileExpr(v.Init); err != nil {
				return err
			}

			gen.emit(opcodes.StLoc, data.MustNewValue(slot))
		}
		return nil

	case s.Return != nil:
		r := s.Return
		if r.Expr != nil {
			if err := gen.compileExpr(r.Expr); err != nil {
				return err
			}
		}
		gen.emit(opcodes.Ret, data.Value{})
		return nil

	case s.Expr != nil:
		if err := gen.compileExpr(s.Expr); err != nil {
			return err
		}
		return nil

	case s.For != nil:
		return gen.compileFor(s.For)
	}

	return nil
}

func (gen *codegen) compileFor(f *ForLoop) error {
	if f.Init != nil {
		if f.Init.Var != nil {
			v := f.Init.Var
			slot := gen.allocLocal(v.Name, v.Type)
			if v.Init != nil {
				if err := gen.compileExpr(v.Init); err != nil {
					return err
				}
				gen.emit(opcodes.StLoc, data.MustNewValue(slot))
			}
		} else if f.Init.Expr != nil {
			if err := gen.compileExpr(f.Init.Expr); err != nil {
				return err
			}
		}
	}

	condPC := len(gen.curVmFn.Instructions)

	if f.Cond != nil {
		if err := gen.compileExpr(f.Cond); err != nil {
			return err
		}
	} else {
	}

	jmpFalseIdx := len(gen.curVmFn.Instructions)
	gen.emit(opcodes.Bf, data.Value{})

	if err := gen.compileBlock(f.Body); err != nil {
		return err
	}

	if f.Post != nil {
		if err := gen.compileExpr(f.Post); err != nil {
			return err
		}
	}

	gen.emit(opcodes.Br, data.MustNewValue(condPC))

	// patch to here
	endPC := len(gen.curVmFn.Instructions)
	gen.curVmFn.Instructions[jmpFalseIdx] = function.Instruction{
		OpCode:  opcodes.Bf,
		Operand: data.MustNewValue(endPC),
	}

	return nil
}

func (gen *codegen) getLocalByName(name string) *function.Local {
	for _, local := range gen.curVmFn.Locals {
		if local.Name == name {
			return local
		}
	}

	return nil
}

func (gen *codegen) compileExpr(e *Expr) error {
	if e.Assign != nil {
		a := e.Assign

		local := gen.getLocalByName(a.Left)
		if local == nil {
			return fmt.Errorf("unknown local %q in assignment", a.Left)
		}

		if a.Fields != nil {
			gen.emit(opcodes.LdLoc, data.MustNewValue(local.Index))
			if err := gen.compileFieldAssignment(local.Type.Name, a.Fields, a.Right); err != nil {
				return err
			}
			gen.emit(opcodes.StLoc, data.MustNewValue(local.Index))
			return nil
		}

		if err := gen.compileSimpleExpr(a.Right); err != nil {
			return err
		}

		gen.emit(opcodes.StLoc, data.MustNewValue(local.Index))
		return nil
	}
	return gen.compileSimpleExpr(e.Simple)
}

func (gen *codegen) compileFieldAssignment(structTypeName string, access *FieldAccess, right *SimpleExpr) error {
	fieldIdx := gen.getFieldIndex(structTypeName, access.Field)
	if fieldIdx < 0 {
		return fmt.Errorf("unknown field %q in struct %s", access.Field, structTypeName)
	}

	if access.Next != nil {
		gen.emit(opcodes.LdFld, data.MustNewValue(fieldIdx))

		st := gen.structs[structTypeName]
		var nextTypeName string
		for _, f := range st.Fields {
			if f.Name == access.Field {
				nextTypeName = f.Type
				break
			}
		}
		return gen.compileFieldAssignment(nextTypeName, access.Next, right)
	}

	if err := gen.compileSimpleExpr(right); err != nil {
		return err
	}

	gen.emit(opcodes.StFld, data.MustNewValue(fieldIdx))
	return nil
}

func (gen *codegen) compileSimpleExpr(e *SimpleExpr) error {
	if err := gen.compileAddExpr(e.Left); err != nil {
		return err
	}
	if e.Op == nil {
		return nil
	}

	if err := gen.compileAddExpr(e.Right); err != nil {
		return err
	}

	switch *e.Op {
	case "==":
		gen.emit(opcodes.Cmp, data.Value{})
	case "!=":
		gen.emit(opcodes.Cmp, data.Value{})
		gen.emit(opcodes.Not, data.Value{})
	case "<":
		gen.emit(opcodes.Clt, data.Value{})
	case ">":
		gen.emit(opcodes.Cgt, data.Value{})
	case "<=":
		gen.emit(opcodes.Cle, data.Value{})
	case ">=":
		gen.emit(opcodes.Cge, data.Value{})
	default:
		return fmt.Errorf("unknown comparison op %q", *e.Op)
	}

	return nil
}

func (gen *codegen) compileAddExpr(e *AddExpr) error {
	if err := gen.compileMulExpr(e.Left); err != nil {
		return err
	}
	for _, r := range e.Rest {
		if err := gen.compileMulExpr(r.Expr); err != nil {
			return err
		}
		switch r.Op {
		case "+":
			gen.emit(opcodes.Add, data.Value{})
		case "-":
			gen.emit(opcodes.Sub, data.Value{})
		default:
			return fmt.Errorf("unknown add op %q", r.Op)
		}
	}
	return nil
}

func (gen *codegen) compileMulExpr(e *MulExpr) error {
	if err := gen.compilePrimary(e.Left); err != nil {
		return err
	}
	for _, r := range e.Rest {
		if err := gen.compilePrimary(r.Expr); err != nil {
			return err
		}
		switch r.Op {
		case "*":
			gen.emit(opcodes.Mul, data.Value{})
		case "/":
			gen.emit(opcodes.Div, data.Value{})
		default:
			return fmt.Errorf("unknown mul op %q", r.Op)
		}
	}
	return nil
}

func (gen *codegen) compilePrimary(p *Primary) error {
	switch {
	case p.Number != nil:
		gen.emit(opcodes.Ldc8, data.MustNewValue(*p.Number))
		return nil

	case p.String != nil:
		gen.emit(opcodes.LdStr, data.MustNewValue(*p.String))
		return nil

	case p.New != nil:
		return gen.compileNewExpr(p.New)

	case p.Ident != nil:
		return gen.compileIdentExpr(p.Ident)

	case p.Sub != nil:
		return gen.compileExpr(p.Sub)
	}

	return fmt.Errorf("invalid primary")
}

func (gen *codegen) compileNewExpr(n *NewExpr) error {
	structType := types.GetTypeByName(n.TypeName)
	if structType == nil {
		return fmt.Errorf("unknown struct type %q", n.TypeName)
	}

	gen.emit(opcodes.NewObj, data.MustNewValue(int(structType.Code)))

	for _, f := range n.Fields {
		if err := gen.compileExpr(f.Value); err != nil {
			return err
		}

		fieldIdx := gen.getFieldIndex(n.TypeName, f.Name)
		if fieldIdx < 0 {
			return fmt.Errorf("unknown field %q in struct %s", f.Name, n.TypeName)
		}

		gen.emit(opcodes.StFld, data.MustNewValue(fieldIdx))
	}

	return nil
}

func (gen *codegen) getFieldIndex(structName, fieldName string) int {
	st, ok := gen.structs[structName]
	if !ok {
		return -1
	}

	for i, f := range st.Fields {
		if f.Name == fieldName {
			return i
		}
	}

	return -1
}

func (gen *codegen) compileIdentExpr(id *IdentExpr) error {
	if id.Call != nil {
		return gen.compileCall(id)
	}

	local := gen.getLocalByName(id.Name)
	if local == nil {
		return fmt.Errorf("unknown variable %q", id.Name)
	}

	gen.emit(opcodes.LdLoc, data.MustNewValue(local.Index))

	if id.Access != nil {
		return gen.compileFieldAccess(local.Type.Name, id.Access)
	}

	return nil
}

func (gen *codegen) compileFieldAccess(structTypeName string, access *FieldAccess) error {
	fieldIdx := gen.getFieldIndex(structTypeName, access.Field)
	if fieldIdx < 0 {
		return fmt.Errorf("unknown field %q in struct %s", access.Field, structTypeName)
	}

	gen.emit(opcodes.LdFld, data.MustNewValue(fieldIdx))

	if access.Next != nil {
		st := gen.structs[structTypeName]
		if st == nil {
			return fmt.Errorf("unknown struct %q", structTypeName)
		}
		var nextTypeName string
		for _, f := range st.Fields {
			if f.Name == access.Field {
				nextTypeName = f.Type
				break
			}
		}
		return gen.compileFieldAccess(nextTypeName, access.Next)
	}

	return nil
}

func (gen *codegen) compileCall(id *IdentExpr) error {
	for _, arg := range id.Call.Args {
		if err := gen.compileExpr(arg); err != nil {
			return err
		}
	}

	token := 0
	ok := false
	for k, v := range std.StdLib {
		if v.Name == id.Name {
			token = k
			ok = true
			break
		}
	}

	for k, v := range gen.vm.Functions {
		if v.Name == id.Name {
			token = k
			ok = true
			break
		}
	}

	if !ok {
		return fmt.Errorf("unknown function %q in codegen", id.Name)
	}

	gen.emit(opcodes.Call, data.MustNewValue(token))
	return nil
}
