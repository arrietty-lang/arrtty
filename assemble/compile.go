package assemble

import (
	"fmt"
	"github.com/arrietty-lang/arrtty/preprocess/analyze"
	"github.com/arrietty-lang/arrtty/preprocess/parse"
	"github.com/arrietty-lang/arrtty/preprocess/tokenize"
	"github.com/arrietty-lang/arrtty/vm"
)

var semOverall *analyze.Semantics
var currentNest int
var currentFnVariableBPs map[int]map[string]int

func searchBPDistFromVarName(variables map[int]map[string]int, nest int, varName string) int {
	for name, distance := range variables[nest] {
		if varName == name {
			return distance
		}
	}
	return -1
}

func defFunction(node *parse.Node) ([]vm.Fragment, error) {
	defFn := node.FuncDefField
	// bpをプッシュする前に戻り値の分だけぷっしゅしておく？
	// 関数として呼び出された場合に必要な命令を始めに入れとく

	var program []vm.Fragment
	program = append(program, vm.NewDefLabelFragment(defFn.Identifier.IdentField.Ident))

	if defFn.Identifier.IdentField.Ident != "main" {
		program = append(program, []vm.Fragment{
			vm.NewOpcodeFragment(vm.PUSH),
			vm.NewPointerFragment(vm.BP),

			vm.NewOpcodeFragment(vm.MOV),
			vm.NewPointerFragment(vm.SP),
			vm.NewPointerFragment(vm.BP),
		}...)
	}

	// 関数で使用されている変数(引数もむくむ)のBPからの距離
	var varDistFromBP = map[int]map[string]int{}
	var dist = 1
	for nest, variables := range semOverall.KnownValues[defFn.Identifier.IdentField.Ident] {
		varDistFromBP[nest] = map[string]int{}
		for varName := range variables {
			varDistFromBP[nest][varName] = dist
			dist++
		}
	}
	currentFnVariableBPs = varDistFromBP

	// 引数と変数を結びつける(代入によって)
	if defFn.Parameters != nil {
		program = append(program, []vm.Fragment{
			vm.NewOpcodeFragment(vm.SUB),
			vm.NewLiteralFragment(vm.NewInt(len(defFn.Parameters.PolynomialField.Values))),
			vm.NewPointerFragment(vm.SP),
		}...)
		for i, param := range defFn.Parameters.PolynomialField.Values {
			program = append(program, []vm.Fragment{
				vm.NewOpcodeFragment(vm.MOV),
				vm.NewAddressFragment(vm.NewAddress(vm.BP, i+1+1)), // iが0から始まるのでbp+0は意味ないので、1底上げ + BPと引数の間のスタックには戻る場所が保存されているのでその分の1
				vm.NewAddressFragment(vm.NewAddress(vm.BP, -searchBPDistFromVarName(varDistFromBP, 0, param.FuncParam.Identifier.IdentField.Ident))),
			}...)
		}
	}
	// こんな感じになってる
	// [ stack ]
	// | 引数2
	// | 引数1
	// | ret-pc
	// | bp

	for _, s := range defFn.Body.BlockField.Statements {
		f, err := stmt(s)
		if err != nil {
			return nil, err
		}
		program = append(program, f...)
	}

	//if defFn.Parameters != nil {
	//	program = append(program, []vm.Fragment{
	//		vm.NewOpcodeFragment(vm.ADD),
	//		vm.NewLiteralFragment(vm.NewInt(len(defFn.Parameters.PolynomialField.Values))),
	//		vm.NewPointerFragment(vm.SP),
	//	}...)
	//}

	// 現状復帰に必要な命令を最後に入れてあげる
	if defFn.Identifier.IdentField.Ident != "main" {
		program = append(program, []vm.Fragment{
			vm.NewOpcodeFragment(vm.MOV),
			vm.NewPointerFragment(vm.BP),
			vm.NewPointerFragment(vm.SP),

			vm.NewOpcodeFragment(vm.POP),
			vm.NewPointerFragment(vm.BP),
		}...)
		program = append(program, []vm.Fragment{
			vm.NewOpcodeFragment(vm.RET),
		}...)
	}

	return program, nil
}

func stmt(node *parse.Node) ([]vm.Fragment, error) {
	var program []vm.Fragment
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
			program = append(program, []vm.Fragment{
				vm.NewOpcodeFragment(vm.POP), vm.NewRegisterFragment(vm.R1),
			}...)
			program = append(program, []vm.Fragment{
				vm.NewOpcodeFragment(vm.MOV),
				vm.NewRegisterFragment(vm.R1),
				vm.NewRegisterFragment(returnRegs[i])}...)
		}
		return program, nil
	}
	return nil, fmt.Errorf("unsupport node kind")
}

func expr(node *parse.Node) ([]vm.Fragment, error) {
	return assign(node)
}

func assign(node *parse.Node) ([]vm.Fragment, error) {
	return equality(node)
}

func equality(node *parse.Node) ([]vm.Fragment, error) {
	return relation(node)
}

func relation(node *parse.Node) ([]vm.Fragment, error) {
	return add(node)
}

func add(node *parse.Node) ([]vm.Fragment, error) {
	var program []vm.Fragment
	switch node.Kind {
	case parse.NdAdd, parse.NdSub:
		var op vm.Opcode
		if node.Kind == parse.NdAdd {
			op = vm.ADD
		} else {
			op = vm.SUB
		}

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
		program = append(program, []vm.Fragment{
			vm.NewOpcodeFragment(vm.POP), vm.NewRegisterFragment(vm.R2), // 右辺
			vm.NewOpcodeFragment(vm.POP), vm.NewRegisterFragment(vm.R1), // 左辺
		}...)
		// 足し算/引き算
		program = append(program, []vm.Fragment{
			// r1 += r2
			vm.NewOpcodeFragment(op), vm.NewRegisterFragment(vm.R2), vm.NewRegisterFragment(vm.R1),
		}...)
		// 結果R1をスタックにプッシュ
		program = append(program, []vm.Fragment{
			vm.NewOpcodeFragment(vm.PUSH), vm.NewRegisterFragment(vm.R1),
		}...)
		return program, nil
	}
	return mul(node)
}

func mul(node *parse.Node) ([]vm.Fragment, error) {
	var program []vm.Fragment
	switch node.Kind {
	case parse.NdMul, parse.NdDiv:
		var op vm.Opcode
		if node.Kind == parse.NdMul {
			op = vm.MUL
		} else {
			op = vm.DIV
		}
		// 左辺を計算
		lhs, err := mul(node.BinaryField.Lhs)
		if err != nil {
			return nil, err
		}
		program = append(program, lhs...)
		// 右辺を計算
		rhs, err := mul(node.BinaryField.Rhs)
		if err != nil {
			return nil, err
		}
		program = append(program, rhs...)
		// 結果をスタックから取り出す
		program = append(program, []vm.Fragment{
			vm.NewOpcodeFragment(vm.POP), vm.NewRegisterFragment(vm.R2), // 右辺
			vm.NewOpcodeFragment(vm.POP), vm.NewRegisterFragment(vm.R1), // 左辺
		}...)
		// 掛け算割り算
		program = append(program, []vm.Fragment{
			// r1 += r2
			vm.NewOpcodeFragment(op), vm.NewRegisterFragment(vm.R2), vm.NewRegisterFragment(vm.R1),
		}...)
		// 結果R1をスタックにプッシュ
		program = append(program, []vm.Fragment{
			vm.NewOpcodeFragment(vm.PUSH), vm.NewRegisterFragment(vm.R1),
		}...)
		return program, nil

	}
	return unary(node)
}

func unary(node *parse.Node) ([]vm.Fragment, error) {
	return primary(node)
}

func primary(node *parse.Node) ([]vm.Fragment, error) {
	return access(node)
}

func access(node *parse.Node) ([]vm.Fragment, error) {
	return literal(node)
}

func literal(node *parse.Node) ([]vm.Fragment, error) {
	switch node.Kind {
	case parse.NdLiteral:
		switch node.LiteralField.Kind {
		case tokenize.LInt:
			return []vm.Fragment{
				vm.NewOpcodeFragment(vm.PUSH),
				vm.NewLiteralFragment(vm.NewInt(node.LiteralField.I)),
			}, nil
		}
	case parse.NdParenthesis:
		return expr(node.UnaryField.Value)
	case parse.NdCall:
		var program []vm.Fragment
		args := node.CallField.Args.PolynomialField.Values
		for i, j := 0, len(args)-1; i < j; i, j = i+1, j-1 {
			args[i], args[j] = args[j], args[i]
		}
		for _, arg := range args {
			// 計算結果はプッシュされるので、逆順に実行してあげるだけで良い..?
			p, err := expr(arg)
			if err != nil {
				return nil, err
			}
			program = append(program, p...)
		}
		program = append(program, []vm.Fragment{
			vm.NewOpcodeFragment(vm.CALL),
			vm.NewLabelFragment(node.CallField.Identifier.IdentField.Ident),
		}...)
		// 引数分spを加算
		program = append(program, []vm.Fragment{
			vm.NewOpcodeFragment(vm.ADD), vm.NewLiteralFragment(vm.NewInt(len(args))), vm.NewPointerFragment(vm.SP),
		}...)
		switch len(semOverall.KnownFunctions[node.CallField.Identifier.IdentField.Ident].Returns) {
		case 1:
			program = append(program, []vm.Fragment{
				vm.NewOpcodeFragment(vm.PUSH), vm.NewRegisterFragment(vm.R10),
			}...)
		case 2:
			program = append(program, []vm.Fragment{
				vm.NewOpcodeFragment(vm.PUSH), vm.NewRegisterFragment(vm.R11),
				vm.NewOpcodeFragment(vm.PUSH), vm.NewRegisterFragment(vm.R10),
			}...)
		default:
			return nil, fmt.Errorf("2つより多い戻り値には対応していません")
		}
		return program, nil
	case parse.NdIdent:
		loc := searchBPDistFromVarName(currentFnVariableBPs, currentNest, node.IdentField.Ident)
		if loc == -1 {
			return nil, fmt.Errorf("変数が定義されていません: %s", node.IdentField.Ident)
		}
		return []vm.Fragment{
			vm.NewOpcodeFragment(vm.PUSH),
			vm.NewAddressFragment(vm.NewAddress(vm.BP, -loc)),
		}, nil
	}
	return nil, fmt.Errorf("サポートされていないリテラルです")
}

func Compile(sem *analyze.Semantics) ([]vm.Fragment, error) {
	semOverall = sem
	currentNest = 0
	if len(sem.OutsideValues) != 0 || len(sem.OutsideFunctions) != 0 {
		return nil, fmt.Errorf("リンクが不完全です")
	}
	var program []vm.Fragment
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
	return program, nil
}
