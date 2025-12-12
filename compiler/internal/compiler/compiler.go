package compiler

import (
	"fmt"
	"log"
	"shared/pkg/function"
	"shared/pkg/types"
	"shared/pkg/vm"

	"github.com/vmihailenco/msgpack/v5"
)

func Compile(src string) ([]byte, error) {
	prog, err := ParseString(src)
	if err != nil {
		return nil, fmt.Errorf("parse error %s", err.Error())
	}

	if err := CheckProgram(prog); err != nil {
		return nil, fmt.Errorf("type error %s", err.Error())
	}

	v := compileVm(prog)

	data, _ := msgpack.Marshal(v)

	return data, nil
}

func compileVm(prog *Program) *vm.Vm {
	m := new(vm.Vm)
	fnMap := initFunctions(prog, m)

	structMap := make(map[string]*StructDef)
	for _, st := range prog.Structs {
		structMap[st.Name] = st
	}

	m.StructTypes = []types.Type{}
	for _, t := range types.Types {
		if t.IsStruct {
			m.StructTypes = append(m.StructTypes, t)
		}
	}

	for _, fn := range prog.Functions {
		gen := &codegen{
			vm:      m,
			curFunc: fn,
			curVmFn: fnMap[fn.Name],
			structs: structMap,
		}

		err := gen.compileFunction()
		if err != nil {
			fmt.Println(err)
		}
	}

	return m
}

func initFunctions(prog *Program, m *vm.Vm) map[string]*function.Function {
	m.Functions = map[int]*function.Function{}
	funcMap := map[string]*function.Function{}

	for i, fn := range prog.Functions {
		ret := fn.RetType
		if ret == "" {
			ret = "void"
		}

		loType := types.GetTypeByName(ret)
		if loType == nil {
			log.Fatalf("type %s not found", ret)
		}

		newFn := &function.Function{
			Name:         fn.Name,
			Token:        function.CustomFunctionStart + i + 1,
			Instructions: nil,
			Locals:       nil,
			Arguments:    make([]function.Argument, len(fn.Params)),
			ReturnType:   loType.Code,
		}

		if fn.Name == "main" {
			newFn.Token = function.CustomFunctionStart
			m.EntryPoint = newFn.Token
		}

		for i, p := range fn.Params {
			paramType := types.GetTypeByName(p.Type)
			if paramType == nil {
				log.Fatalf("type %s not found", p.Type)
			}

			param := function.Argument{
				Index: i,
				Name:  p.Name,
				Type:  paramType.Code,
			}

			newFn.Arguments[i] = param
		}

		m.Functions[newFn.Token] = newFn

		funcMap[fn.Name] = newFn
	}

	return funcMap
}
