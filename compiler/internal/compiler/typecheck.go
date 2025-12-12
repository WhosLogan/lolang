package compiler

import (
	"fmt"
	"runtime/pkg/std"
	"shared/pkg/types"
)

func isValidType(t string) bool {
	return types.GetTypeByName(t) != nil
}

func isValidReturnType(t string) bool {
	if t == "void" {
		return true
	}
	return isValidType(t)
}

type funcType struct {
	Params []string
	Return string
}

type structType struct {
	Name   string
	Fields map[string]string
}

type checker struct {
	funcs         map[string]funcType
	stdlib        map[string]funcType
	structs       map[string]structType
	scopes        []map[string]string
	currentReturn string
}

func defaultStdLib() map[string]funcType {
	var lib = map[string]funcType{}

	for _, fn := range std.StdLib {
		var params []string

		for _, arg := range fn.Arguments {
			t := types.GetTypeByCode(arg.Type)
			if t == nil {
				panic(fmt.Sprintf("unknown std type %d", arg.Type))
			}

			params = append(params, t.Name)
		}

		retType := types.GetTypeByCode(fn.ReturnType)
		if retType == nil {
			panic(fmt.Sprintf("unknown std return type %d", fn.ReturnType))
		}

		lib[fn.Name] = funcType{Params: params, Return: retType.Name}
	}

	return lib
}

func CheckProgram(prog *Program) error {
	c := &checker{
		funcs:   make(map[string]funcType),
		stdlib:  defaultStdLib(),
		structs: make(map[string]structType),
	}

	for _, st := range prog.Structs {
		if _, exists := c.structs[st.Name]; exists {
			return fmt.Errorf("duplicate struct %q", st.Name)
		}

		fields := make(map[string]string)
		for _, f := range st.Fields {
			if _, exists := fields[f.Name]; exists {
				return fmt.Errorf("struct %s: duplicate field %q", st.Name, f.Name)
			}
			fields[f.Name] = f.Type
		}

		c.structs[st.Name] = structType{
			Name:   st.Name,
			Fields: fields,
		}
	}

	for _, st := range prog.Structs {
		if types.GetTypeByName(st.Name) == nil {
			types.RegisterStructType(st.Name, nil)
		}
	}

	for _, st := range prog.Structs {
		var typeFields []types.StructField
		for _, f := range st.Fields {
			fieldType := types.GetTypeByName(f.Type)
			if fieldType == nil {
				return fmt.Errorf("struct %s: unknown field type %q", st.Name, f.Type)
			}
			typeFields = append(typeFields, types.StructField{
				Name: f.Name,
				Type: fieldType.Code,
			})
		}
		registeredType := types.GetTypeByName(st.Name)
		if registeredType != nil {
			registeredType.Fields = typeFields
		}
	}

	for _, st := range prog.Structs {
		for _, f := range st.Fields {
			if !isValidType(f.Type) {
				return fmt.Errorf("struct %s: unknown field type %q", st.Name, f.Type)
			}
		}
	}

	for _, fn := range prog.Functions {
		if _, exists := c.funcs[fn.Name]; exists {
			return fmt.Errorf("duplicate function %q", fn.Name)
		}

		ret := fn.RetType
		if ret == "" {
			ret = "void"
		}
		if !isValidReturnType(ret) {
			return fmt.Errorf("function %s: invalid return type %q", fn.Name, ret)
		}

		paramTypes := make([]string, len(fn.Params))
		for i, p := range fn.Params {
			if !isValidType(p.Type) {
				return fmt.Errorf("function %s: unknown param type %q", fn.Name, p.Type)
			}
			paramTypes[i] = p.Type
		}

		c.funcs[fn.Name] = funcType{
			Params: paramTypes,
			Return: ret,
		}
	}

	for _, fn := range prog.Functions {
		ft := c.funcs[fn.Name]
		if err := c.checkFunction(fn, ft.Return); err != nil {
			return fmt.Errorf("in function %s: %w", fn.Name, err)
		}
	}

	return nil
}

func (c *checker) pushScope() {
	c.scopes = append(c.scopes, make(map[string]string))
}

func (c *checker) popScope() {
	c.scopes = c.scopes[:len(c.scopes)-1]
}

func (c *checker) declare(name, typ string) error {
	scope := c.scopes[len(c.scopes)-1]
	if _, exists := scope[name]; exists {
		return fmt.Errorf("variable %q already declared", name)
	}
	scope[name] = typ
	return nil
}

func (c *checker) lookup(name string) (string, bool) {
	for i := len(c.scopes) - 1; i >= 0; i-- {
		if t, ok := c.scopes[i][name]; ok {
			return t, true
		}
	}
	return "", false
}

func (c *checker) checkFunction(fn *Function, retType string) error {
	c.pushScope()
	defer c.popScope()

	prev := c.currentReturn
	c.currentReturn = retType
	defer func() { c.currentReturn = prev }()

	for _, p := range fn.Params {
		if err := c.declare(p.Name, p.Type); err != nil {
			return err
		}
	}

	return c.checkBlock(fn.Body)
}

func (c *checker) checkBlock(b *Block) error {
	for _, st := range b.Statements {
		if err := c.checkStatement(st); err != nil {
			return err
		}
	}
	return nil
}

func (c *checker) checkStatement(s *Statement) error {
	switch {
	case s.For != nil:
		return c.checkFor(s.For)

	case s.Var != nil:
		v := s.Var
		if !isValidType(v.Type) {
			return fmt.Errorf("unknown type %q", v.Type)
		}
		if v.Init != nil {
			t, err := c.exprType(v.Init)
			if err != nil {
				return err
			}
			if t != v.Type {
				return fmt.Errorf("cannot assign %s to variable %s of type %s", t, v.Name, v.Type)
			}
		}
		return c.declare(v.Name, v.Type)

	case s.Return != nil:
		r := s.Return
		if c.currentReturn == "" {
			c.currentReturn = "void"
		}
		if r.Expr == nil {
			if c.currentReturn != "void" {
				return fmt.Errorf("missing return value, function expects %s", c.currentReturn)
			}
			return nil
		}

		t, err := c.exprType(r.Expr)
		if err != nil {
			return err
		}
		if c.currentReturn == "void" {
			return fmt.Errorf("cannot return a value from void function")
		}
		if t != c.currentReturn {
			return fmt.Errorf("cannot return %s from function returning %s", t, c.currentReturn)
		}
		return nil

	case s.Expr != nil:
		_, err := c.exprType(s.Expr)
		return err

	default:
		return nil
	}
}

func (c *checker) checkFor(f *ForLoop) error {
	c.pushScope()
	defer c.popScope()

	if f.Init != nil {
		if f.Init.Var != nil {
			v := f.Init.Var
			if !isValidType(v.Type) {
				return fmt.Errorf("unknown type %q", v.Type)
			}
			if v.Init != nil {
				t, err := c.exprType(v.Init)
				if err != nil {
					return err
				}
				if t != v.Type {
					return fmt.Errorf("cannot assign %s to %s of type %s",
						t, v.Name, v.Type)
				}
			}
			if err := c.declare(v.Name, v.Type); err != nil {
				return err
			}
		} else if f.Init.Expr != nil {
			if _, err := c.exprType(f.Init.Expr); err != nil {
				return err
			}
		}
	}

	if f.Cond != nil {
		t, err := c.exprType(f.Cond)
		if err != nil {
			return err
		}
		if t != "bool" {
			return fmt.Errorf("for condition must be bool, got %s", t)
		}
	}

	if f.Post != nil {
		if _, err := c.exprType(f.Post); err != nil {
			return err
		}
	}

	return c.checkBlock(f.Body)
}

func (c *checker) exprType(e *Expr) (string, error) {
	if e.Assign != nil {
		a := e.Assign
		rightType, err := c.simpleExprType(a.Right)
		if err != nil {
			return "", err
		}
		leftType, ok := c.lookup(a.Left)
		if !ok {
			return "", fmt.Errorf("assignment to undeclared variable %q", a.Left)
		}

		if a.Fields != nil {
			fieldType, err := c.fieldAccessType(leftType, a.Fields)
			if err != nil {
				return "", err
			}
			if fieldType != rightType {
				return "", fmt.Errorf("cannot assign %s to field of type %s", rightType, fieldType)
			}
			return fieldType, nil
		}

		if leftType != rightType {
			return "", fmt.Errorf("cannot assign %s to %s of type %s", rightType, a.Left, leftType)
		}
		return leftType, nil
	}
	return c.simpleExprType(e.Simple)
}

func (c *checker) simpleExprType(e *SimpleExpr) (string, error) {
	leftType, err := c.addExprType(e.Left)
	if err != nil {
		return "", err
	}
	if e.Op == nil {
		return leftType, nil
	}
	rightType, err := c.addExprType(e.Right)
	if err != nil {
		return "", err
	}
	if leftType != rightType {
		return "", fmt.Errorf("type mismatch in comparison: %s vs %s", leftType, rightType)
	}
	return "bool", nil
}

func (c *checker) addExprType(e *AddExpr) (string, error) {
	t, err := c.mulExprType(e.Left)
	if err != nil {
		return "", err
	}
	for _, r := range e.Rest {
		rt, err := c.mulExprType(r.Expr)
		if err != nil {
			return "", err
		}
		switch r.Op {
		case "+":
			if t != rt || (t != "int" && t != "string") {
				return "", fmt.Errorf("invalid + between %s and %s", t, rt)
			}
		case "-", "*", "/":
			if t != "int" || rt != "int" {
				return "", fmt.Errorf("operator %s only defined on int, got %s and %s", r.Op, t, rt)
			}
		}
	}
	return t, nil
}

func (c *checker) mulExprType(e *MulExpr) (string, error) {
	t, err := c.primaryType(e.Left)
	if err != nil {
		return "", err
	}
	for _, r := range e.Rest {
		rt, err := c.primaryType(r.Expr)
		if err != nil {
			return "", err
		}
		if t != "int" || rt != "int" {
			return "", fmt.Errorf("operator %s only defined on int, got %s and %s", r.Op, t, rt)
		}
	}
	return t, nil
}

func (c *checker) primaryType(p *Primary) (string, error) {
	switch {
	case p.Number != nil:
		return "int", nil
	case p.String != nil:
		return "string", nil
	case p.New != nil:
		return c.newExprType(p.New)
	case p.Ident != nil:
		return c.identType(p.Ident)
	case p.Sub != nil:
		return c.exprType(p.Sub)
	default:
		return "", fmt.Errorf("invalid expression")
	}
}

func (c *checker) newExprType(n *NewExpr) (string, error) {
	st, ok := c.structs[n.TypeName]
	if !ok {
		return "", fmt.Errorf("unknown struct type %q", n.TypeName)
	}

	providedFields := make(map[string]bool)
	for _, f := range n.Fields {
		if _, exists := st.Fields[f.Name]; !exists {
			return "", fmt.Errorf("struct %s has no field %q", n.TypeName, f.Name)
		}
		if providedFields[f.Name] {
			return "", fmt.Errorf("duplicate field %q in struct initializer", f.Name)
		}
		providedFields[f.Name] = true

		valType, err := c.exprType(f.Value)
		if err != nil {
			return "", err
		}
		if valType != st.Fields[f.Name] {
			return "", fmt.Errorf("field %s: expected %s, got %s", f.Name, st.Fields[f.Name], valType)
		}
	}

	return n.TypeName, nil
}

func (c *checker) identType(id *IdentExpr) (string, error) {
	if id.Call != nil {
		if fnType, ok := c.funcs[id.Name]; ok {
			return c.checkCall(id.Name, fnType, id.Call.Args)
		}
		if fnType, ok := c.stdlib[id.Name]; ok {
			return c.checkCall(id.Name, fnType, id.Call.Args)
		}
		return "", fmt.Errorf("call to unknown function %q", id.Name)
	}

	t, ok := c.lookup(id.Name)
	if !ok {
		return "", fmt.Errorf("use of undeclared variable %q", id.Name)
	}

	if id.Access != nil {
		return c.fieldAccessType(t, id.Access)
	}

	return t, nil
}

func (c *checker) fieldAccessType(baseType string, access *FieldAccess) (string, error) {
	st, ok := c.structs[baseType]
	if !ok {
		return "", fmt.Errorf("cannot access field on non-struct type %q", baseType)
	}

	fieldType, ok := st.Fields[access.Field]
	if !ok {
		return "", fmt.Errorf("struct %s has no field %q", baseType, access.Field)
	}

	if access.Next != nil {
		return c.fieldAccessType(fieldType, access.Next)
	}

	return fieldType, nil
}

func (c *checker) checkCall(name string, fnType funcType, args []*Expr) (string, error) {
	if len(fnType.Params) != len(args) {
		return "", fmt.Errorf("function %s: expected %d args, got %d",
			name, len(fnType.Params), len(args))
	}
	for i, arg := range args {
		at, err := c.exprType(arg)
		if err != nil {
			return "", err
		}
		if at != fnType.Params[i] {
			return "", fmt.Errorf("function %s arg %d: expected %s, got %s",
				name, i+1, fnType.Params[i], at)
		}
	}

	return fnType.Return, nil
}
