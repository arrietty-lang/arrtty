package assemble

import (
	"fmt"
	"github.com/arrietty-lang/arrtty/preprocess/analyze"
	"github.com/arrietty-lang/arrtty/preprocess/parse"
	"github.com/arrietty-lang/arrtty/preprocess/tokenize"
	"github.com/arrietty-lang/arrtty/vm"
)

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

	for _, s := range defFn.Body.BlockField.Statements {
		f, err := stmt(s)
		if err != nil {
			return nil, err
		}
		program = append(program, f...)
	}

	// 現状復帰に必要な命令を最後に入れてあげる
	program = append(program, []*vm.Fragment{
		vm.NewOpcodeFragment(vm.MOV),
		vm.NewPointerFragment(vm.BP),
		vm.NewPointerFragment(vm.SP),

		vm.NewOpcodeFragment(vm.POP),
		vm.NewPointerFragment(vm.BP),

		//vm.NewOpcodeFragment(vm.RET),
	}...)

	return program, nil
}

func stmt(node *parse.Node) ([]*vm.Fragment, error) {
	var program []*vm.Fragment
	switch node.Kind {
	case parse.NdReturn:
		if len(node.PolynomialField.Values) > 2 {
			return nil, fmt.Errorf("２つ以上の戻り値は現在サポートされていません")
		}
		returnRegs := []vm.Register{vm.R10, vm.R11}
		for i, rn := range node.PolynomialField.Values {
			rv, err := expr(rn)
			if err != nil {
				return nil, err
			}
			program = append(program, rv...)
			// 計算結果をR1に移動
			program = append(program, []*vm.Fragment{
				vm.NewOpcodeFragment(vm.POP), vm.NewRegisterFragment(vm.R1),
			}...)
			program = append(program, []*vm.Fragment{
				vm.NewOpcodeFragment(vm.MOV),
				vm.NewRegisterFragment(vm.R1),
				vm.NewRegisterFragment(returnRegs[i])}...)
		}
		return program, nil
	}
	return nil, fmt.Errorf("unsupport node kind")
}

func expr(node *parse.Node) ([]*vm.Fragment, error) {
	return assign(node)
}

func assign(node *parse.Node) ([]*vm.Fragment, error) {
	return equality(node)
}

func equality(node *parse.Node) ([]*vm.Fragment, error) {
	return relation(node)
}

func relation(node *parse.Node) ([]*vm.Fragment, error) {
	return add(node)
}

func add(node *parse.Node) ([]*vm.Fragment, error) {
	var program []*vm.Fragment
	switch node.Kind {
	case parse.NdAdd:
		// 左辺を計算
		lhs, err := add(node.BinaryField.Lhs)
		if err != nil {
			return nil, err
		}
		program = append(program, lhs...)
		// 右辺を計算
		rhs, err := add(node.BinaryField.Rhs)
		if err != nil {
			return nil, err
		}
		program = append(program, rhs...)
		// 結果をスタックから取り出す
		program = append(program, []*vm.Fragment{
			vm.NewOpcodeFragment(vm.POP), vm.NewRegisterFragment(vm.R2), // 右辺
			vm.NewOpcodeFragment(vm.POP), vm.NewRegisterFragment(vm.R1), // 左辺
		}...)
		// 足し算
		program = append(program, []*vm.Fragment{
			// r1 += r2
			vm.NewOpcodeFragment(vm.ADD), vm.NewRegisterFragment(vm.R2), vm.NewRegisterFragment(vm.R1),
		}...)
		// 結果R1をスタックにプッシュ
		program = append(program, []*vm.Fragment{
			vm.NewOpcodeFragment(vm.PUSH), vm.NewRegisterFragment(vm.R1),
		}...)
		return program, nil
	}
	return mul(node)
}

func mul(node *parse.Node) ([]*vm.Fragment, error) {
	return unary(node)
}

func unary(node *parse.Node) ([]*vm.Fragment, error) {
	return primary(node)
}

func primary(node *parse.Node) ([]*vm.Fragment, error) {
	return access(node)
}

func access(node *parse.Node) ([]*vm.Fragment, error) {
	return literal(node)
}

func literal(node *parse.Node) ([]*vm.Fragment, error) {
	switch node.Kind {
	case parse.NdLiteral:
		switch node.LiteralField.Kind {
		case tokenize.LInt:
			return []*vm.Fragment{
				vm.NewOpcodeFragment(vm.PUSH),
				vm.NewLiteralFragment(vm.NewInt(node.LiteralField.I)),
			}, nil
		}
	}
	return nil, fmt.Errorf("サポートされていないリテラルです")
}

func Compile(sem *analyze.Semantics) ([]*vm.Fragment, error) {
	if len(sem.OutsideValues) != 0 || len(sem.OutsideFunctions) != 0 {
		return nil, fmt.Errorf("リンクが不完全です")
	}

	program := []*vm.Fragment{
		//vm.NewOpcodeFragment(vm.CALL),
		//vm.NewLabelFragment("main"),
		//vm.NewOpcodeFragment(vm.EXIT),
	}
	//var program []*vm.Fragment

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

	//if !isMainReturned {
	//	program = append(program, []*vm.Fragment{
	//		vm.NewOpcodeFragment(vm.MOV),
	//		vm.NewLiteralFragment(vm.NewInt(0)),
	//		vm.NewRegisterFragment(vm.R1),
	//	}...)
	//}

	return program, nil
}
