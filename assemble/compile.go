package assemble

import (
	"fmt"
	"github.com/arrietty-lang/arrtty/preprocess/analyze"
	"github.com/arrietty-lang/arrtty/preprocess/parse"
	"github.com/arrietty-lang/arrtty/vm"
)

var isMainReturned = false

func defFunction(node *parse.Node) ([]*vm.Fragment, error) {
	defFn := node.FuncDefField
	// bpをプッシュする前に戻り値の分だけぷっしゅしておく？
	// 関数として呼び出された場合に必要な命令を始めに入れとく
	program := []*vm.Fragment{
		vm.NewDefLabelFragment(defFn.Identifier.IdentField.Ident),

		vm.NewOpcodeFragment(vm.PUSH),
		vm.NewPointerFragment(vm.BP),

		vm.NewOpcodeFragment(vm.MOV),
		vm.NewPointerFragment(vm.SP),
		vm.NewPointerFragment(vm.BP),
	}

	// 引数分spを引く

	// 現状復帰に必要な命令を最後に入れてあげる
	defer func() {
		program = append(program, []*vm.Fragment{
			vm.NewOpcodeFragment(vm.MOV),
			vm.NewPointerFragment(vm.BP),
			vm.NewPointerFragment(vm.SP),

			vm.NewOpcodeFragment(vm.POP),
			vm.NewPointerFragment(vm.BP),

			vm.NewOpcodeFragment(vm.RET),
		}...)
	}()

	var isMain = false
	if defFn.Identifier.IdentField.Ident == "main" {
		isMain = true
	}
	for _, s := range defFn.Body.BlockField.Statements {
		f, err := stmt(s, isMain)
		if err != nil {
			return nil, err
		}
		program = append(program, f...)
	}

	return program, nil
}

func stmt(node *parse.Node, isMain bool) ([]*vm.Fragment, error) {
	var program []*vm.Fragment
	switch node.Kind {
	case parse.NdReturn:
		if isMain {
			program = append(program, []*vm.Fragment{
				vm.NewOpcodeFragment(vm.MOV),
				//
				vm.NewRegisterFragment(vm.R1),
			}...)
			isMainReturned = true
		}
	}
	return nil, fmt.Errorf("unsupport node kind")
}

func Compile(sem *analyze.Semantics) ([]*vm.Fragment, error) {
	if len(sem.OutsideValues) != 0 || len(sem.OutsideFunctions) != 0 {
		return nil, fmt.Errorf("リンクが不完全です")
	}

	program := []*vm.Fragment{
		vm.NewOpcodeFragment(vm.CALL),
		vm.NewLabelFragment("main"),
	}

	for _, n := range sem.Tree {
		switch n.Kind {
		case parse.NdFuncDef:
			frags, err := defFunction(n)
			if err != nil {
				return nil, err
			}
			program = append(program, frags...)
		}
	}

	if !isMainReturned {
		program = append(program, []*vm.Fragment{
			vm.NewOpcodeFragment(vm.MOV),
			vm.NewLiteralFragment(vm.NewInt(0)),
			vm.NewRegisterFragment(vm.R1),
		}...)
	}
	program = append(program, vm.NewOpcodeFragment(vm.EXIT))

	return program, nil
}
