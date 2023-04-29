package assemble2

import (
	"fmt"
	"github.com/arrietty-lang/arrtty/preprocess/analyze"
	"github.com/arrietty-lang/arrtty/preprocess/parse"
	"github.com/arrietty-lang/arrtty/preprocess/tokenize"
	"github.com/arrietty-lang/arrtty/vm3"
	"math/rand"
)

var currentFunctionName string

var semOverall *analyze.Semantics
var currentNest int
var currentFnVariableBPs map[int]map[string]int

var globals []string
var dataSection []vm3.Data

func init() {
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func searchBPDistFromVarName(variables map[int]map[string]int, nest int, varName string) int {
	for name, distance := range variables[nest] {
		if varName == name {
			return distance
		}
	}
	return -1
}

func defFunction(node *parse.Node) ([]vm3.Data, error) {
	defFn := node.FuncDefField
	currentFunctionName = defFn.Identifier.IdentField.Ident
	// bpをプッシュする前に戻り値の分だけぷっしゅしておく？
	// 関数として呼び出された場合に必要な命令を始めに入れとく

	var program []vm3.Data
	program = append(program,
		*vm3.NewLabelData(*vm3.NewLabel(true, defFn.Identifier.IdentField.Ident)))

	if defFn.Identifier.IdentField.Ident != "main" {
		program = append(program, []vm3.Data{
			*vm3.NewOpcodeData(vm3.PUSH), *vm3.NewRegisterTagData(vm3.RBP),
			*vm3.NewOpcodeData(vm3.MOV), *vm3.NewRegisterTagData(vm3.RSP), *vm3.NewRegisterTagData(vm3.RBP),
		}...)
	}

	// 関数で使用されている変数(引数もむくむ)のBPからの距離
	var totalVariables = 0
	var varDistFromBP = map[int]map[string]int{}
	var dist = 1
	for nest, variables := range semOverall.KnownValues[defFn.Identifier.IdentField.Ident] {
		varDistFromBP[nest] = map[string]int{}
		for varName := range variables {
			varDistFromBP[nest][varName] = dist
			dist++
			totalVariables++
		}
	}
	currentFnVariableBPs = varDistFromBP

	// 関数内で使用される変数の数だけSPを下げる(変数用の領域確保)
	program = append(program, []vm3.Data{
		*vm3.NewOpcodeData(vm3.PUSH), *vm3.NewLiteralDataWithRaw(totalVariables),
		*vm3.NewOpcodeData(vm3.POP), *vm3.NewRegisterTagData(vm3.R1),
		*vm3.NewOpcodeData(vm3.SUB), *vm3.NewRegisterTagData(vm3.R1), *vm3.NewRegisterTagData(vm3.RSP),
	}...)

	// 引数と変数を結びつける(代入によって)
	if defFn.Parameters != nil {
		for i, param := range defFn.Parameters.PolynomialField.Values {
			relation_ := -searchBPDistFromVarName(varDistFromBP, 0, param.FuncParam.Identifier.IdentField.Ident)
			program = append(program, []vm3.Data{
				*vm3.NewOpcodeData(vm3.MOV), *vm3.NewOffsetData(*vm3.NewOffset(vm3.BP, i+2)), *vm3.NewOffsetData(*vm3.NewOffset(vm3.BP, relation_)),
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

	return program, nil
}

func stmt(node *parse.Node) ([]vm3.Data, error) {
	var program []vm3.Data
	switch node.Kind {
	case parse.NdReturn:
		if len(node.PolynomialField.Values) > 2 {
			return nil, fmt.Errorf("２つ以上の戻り値は現在サポートされていません")
		}
		// 戻り値の準備
		returnRegs := []vm3.RegisterTag{vm3.R10, vm3.R11}
		for i, rn := range node.PolynomialField.Values {
			rv, err := expr(rn)
			if err != nil {
				return nil, err
			}
			program = append(program, rv...)
			// 計算結果をR1に移動
			program = append(program, []vm3.Data{
				*vm3.NewOpcodeData(vm3.POP), *vm3.NewRegisterTagData(vm3.R1),
				*vm3.NewOpcodeData(vm3.MOV), *vm3.NewRegisterTagData(vm3.R1), *vm3.NewRegisterTagData(returnRegs[i]),
			}...)
		}
		// リターン本文
		// メインだけはリターンで何も返さない..?
		if currentFunctionName != "main" {
			program = append(program, []vm3.Data{
				*vm3.NewOpcodeData(vm3.MOV), *vm3.NewRegisterTagData(vm3.RBP), *vm3.NewRegisterTagData(vm3.RSP),
				*vm3.NewOpcodeData(vm3.POP), *vm3.NewRegisterTagData(vm3.RBP),
				*vm3.NewOpcodeData(vm3.RET),
			}...)
		}
		return program, nil
	case parse.NdAssign:
		val, err := expr(node.AssignField.Value)
		if err != nil {
			return nil, err
		}
		var loc int
		var ident string
		switch node.AssignField.To.Kind {
		case parse.NdVarDecl:
			ident = node.AssignField.To.VarDeclField.Identifier.IdentField.Ident
			loc = searchBPDistFromVarName(currentFnVariableBPs, currentNest, ident)
		case parse.NdIdent:
			ident = node.AssignField.To.IdentField.Ident
			loc = searchBPDistFromVarName(currentFnVariableBPs, currentNest, ident)
		}
		// 変数の中身をスタックにプッシュ
		program = append(program, val...)
		// 関数内の変数
		if loc != -1 {
			program = append(program, []vm3.Data{
				// valの結果を取り出す
				*vm3.NewOpcodeData(vm3.POP), *vm3.NewRegisterTagData(vm3.R1),
				// 変数の場所に格納
				*vm3.NewOpcodeData(vm3.MOV), *vm3.NewRegisterTagData(vm3.R1), *vm3.NewOffsetData(*vm3.NewOffset(vm3.BP, -loc)),
			}...)
			return program, nil
		}

		// グローバル変数ならば
		for _, g := range globals {
			if g == ident {
				program = append(program, []vm3.Data{
					// valの結果を取り出す
					*vm3.NewOpcodeData(vm3.POP), *vm3.NewRegisterTagData(vm3.R1),
					// 変数の場所に格納
					*vm3.NewOpcodeData(vm3.MOV), *vm3.NewRegisterTagData(vm3.R1), *vm3.NewLabelData(*vm3.NewLabel(false, ident)),
				}...)
				return program, nil
			}
		}

		return nil, fmt.Errorf("変数の位置を特定できませんでした: %v", ident)
	case parse.NdVarDecl:
		// アナライズされて関数のはじめにspまとめて引かれているので特にすることはないはず
		return nil, nil
	case parse.NdIfElse:
		// 条件を計算
		cond, err := expr(node.IfElseField.Cond)
		if err != nil {
			return nil, err
		}
		program = append(program, cond...)
		// 結果を取り出す
		program = append(program, []vm3.Data{
			*vm3.NewOpcodeData(vm3.POP), *vm3.NewRegisterTagData(vm3.R1),
		}...)

		// 各ジャンプ先のラベルを用意
		ifBlockLabel := "if_if_block_" + RandStringRunes(20)
		endLabel := "if_end_" + RandStringRunes(20)

		// 条件に合致したらIFブロックへ飛ぶ
		program = append(program, []vm3.Data{
			*vm3.NewOpcodeData(vm3.JZ), *vm3.NewLabelData(*vm3.NewLabel(false, ifBlockLabel)),
		}...)

		// エルスを使用していればエルスのブロックを展開
		if node.IfElseField.UseElse {
			elseBlock, err := stmt(node.IfElseField.ElseBlock)
			if err != nil {
				return nil, err
			}
			program = append(program, elseBlock...)
		}
		// エンドラベルに飛ぶ命令を挿入(条件に合致しなかった場合、IFを読み飛ばすため)
		program = append(program, []vm3.Data{
			*vm3.NewOpcodeData(vm3.JMP), *vm3.NewLabelData(*vm3.NewLabel(false, endLabel)),
		}...)

		// if-block
		program = append(program, *vm3.NewLabelData(*vm3.NewLabel(true, ifBlockLabel)))
		ifBlock, err := stmt(node.IfElseField.IfBlock)
		if err != nil {
			return nil, err
		}
		program = append(program, ifBlock...)

		// 最終的なジャンプ先
		program = append(program, *vm3.NewLabelData(*vm3.NewLabel(true, endLabel)))
		return program, nil
	case parse.NdBlock:
		for _, n := range node.BlockField.Statements {
			f, err := stmt(n)
			if err != nil {
				return nil, err
			}
			program = append(program, f...)
		}
		return program, nil
	}
	return expr(node)
}

func expr(node *parse.Node) ([]vm3.Data, error) {
	return assign(node)
}

func assign(node *parse.Node) ([]vm3.Data, error) {
	return equality(node)
}

func equality(node *parse.Node) ([]vm3.Data, error) {
	return relation(node)
}

func relation(node *parse.Node) ([]vm3.Data, error) {
	var program []vm3.Data
	switch node.Kind {
	case parse.NdLt:
		lhs, err := relation(node.BinaryField.Lhs)
		if err != nil {
			return nil, err
		}
		rhs, err := relation(node.BinaryField.Rhs)
		if err != nil {
			return nil, err
		}
		program = append(program, lhs...)
		program = append(program, rhs...)

		program = append(program, []vm3.Data{
			*vm3.NewOpcodeData(vm3.POP), *vm3.NewRegisterTagData(vm3.R2),
			*vm3.NewOpcodeData(vm3.POP), *vm3.NewRegisterTagData(vm3.R1),
			*vm3.NewOpcodeData(vm3.LT), *vm3.NewRegisterTagData(vm3.R1), *vm3.NewRegisterTagData(vm3.R2),
		}...)

		return program, nil
	case parse.NdLe:
		lhs, err := relation(node.BinaryField.Lhs)
		if err != nil {
			return nil, err
		}
		rhs, err := relation(node.BinaryField.Rhs)
		if err != nil {
			return nil, err
		}
		program = append(program, lhs...)
		program = append(program, rhs...)

		program = append(program, []vm3.Data{
			*vm3.NewOpcodeData(vm3.POP), *vm3.NewRegisterTagData(vm3.R2),
			*vm3.NewOpcodeData(vm3.POP), *vm3.NewRegisterTagData(vm3.R1),
			*vm3.NewOpcodeData(vm3.LE), *vm3.NewRegisterTagData(vm3.R1), *vm3.NewRegisterTagData(vm3.R2),
		}...)

		return program, nil
	case parse.NdGt:
		lhs, err := relation(node.BinaryField.Lhs)
		if err != nil {
			return nil, err
		}
		rhs, err := relation(node.BinaryField.Rhs)
		if err != nil {
			return nil, err
		}
		program = append(program, lhs...)
		program = append(program, rhs...)

		program = append(program, []vm3.Data{
			*vm3.NewOpcodeData(vm3.POP), *vm3.NewRegisterTagData(vm3.R2),
			*vm3.NewOpcodeData(vm3.POP), *vm3.NewRegisterTagData(vm3.R1),
			*vm3.NewOpcodeData(vm3.LT), *vm3.NewRegisterTagData(vm3.R2), *vm3.NewRegisterTagData(vm3.R1),
		}...)
		return program, nil
	case parse.NdGe:
		lhs, err := relation(node.BinaryField.Lhs)
		if err != nil {
			return nil, err
		}
		rhs, err := relation(node.BinaryField.Rhs)
		if err != nil {
			return nil, err
		}
		program = append(program, lhs...)
		program = append(program, rhs...)

		program = append(program, []vm3.Data{
			*vm3.NewOpcodeData(vm3.POP), *vm3.NewRegisterTagData(vm3.R2),
			*vm3.NewOpcodeData(vm3.POP), *vm3.NewRegisterTagData(vm3.R1),
			*vm3.NewOpcodeData(vm3.LE), *vm3.NewRegisterTagData(vm3.R2), *vm3.NewRegisterTagData(vm3.R1),
		}...)

		return program, nil
	}
	return add(node)
}

func add(node *parse.Node) ([]vm3.Data, error) {
	var program []vm3.Data
	switch node.Kind {
	case parse.NdAdd, parse.NdSub:
		var op vm3.Opcode
		if node.Kind == parse.NdAdd {
			op = vm3.ADD
		} else {
			op = vm3.SUB
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
		program = append(program, []vm3.Data{
			*vm3.NewOpcodeData(vm3.POP), *vm3.NewRegisterTagData(vm3.R2),
			*vm3.NewOpcodeData(vm3.POP), *vm3.NewRegisterTagData(vm3.R1),
		}...)
		// 足し算/引き算
		program = append(program, []vm3.Data{
			// r1 += r2
			*vm3.NewOpcodeData(op), *vm3.NewRegisterTagData(vm3.R2), *vm3.NewRegisterTagData(vm3.R1),
		}...)
		// 結果R1をスタックにプッシュ
		program = append(program, []vm3.Data{
			*vm3.NewOpcodeData(vm3.PUSH), *vm3.NewRegisterTagData(vm3.R1),
		}...)
		return program, nil
	}
	return mul(node)
}

func mul(node *parse.Node) ([]vm3.Data, error) {
	var program []vm3.Data
	switch node.Kind {
	case parse.NdMul, parse.NdDiv:
		var op vm3.Opcode
		if node.Kind == parse.NdMul {
			op = vm3.MUL
		} else {
			op = vm3.DIV
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
		program = append(program, []vm3.Data{
			*vm3.NewOpcodeData(vm3.POP), *vm3.NewRegisterTagData(vm3.R2),
			*vm3.NewOpcodeData(vm3.POP), *vm3.NewRegisterTagData(vm3.R1),
		}...)
		// 掛け算割り算
		program = append(program, []vm3.Data{
			// r1 += r2
			*vm3.NewOpcodeData(op), *vm3.NewRegisterTagData(vm3.R2), *vm3.NewRegisterTagData(vm3.R1),
		}...)
		// 結果R1をスタックにプッシュ
		program = append(program, []vm3.Data{
			*vm3.NewOpcodeData(vm3.PUSH), *vm3.NewRegisterTagData(vm3.R1),
		}...)
		return program, nil

	}
	return unary(node)
}

func unary(node *parse.Node) ([]vm3.Data, error) {
	return primary(node)
}

func primary(node *parse.Node) ([]vm3.Data, error) {
	return access(node)
}

func access(node *parse.Node) ([]vm3.Data, error) {
	return literal(node)
}

func literalFromField(l *parse.LiteralField) *vm3.Literal {
	switch l.Kind {
	case tokenize.LInt:
		return vm3.NewLiteral(l.I)
	case tokenize.LFloat:
		return vm3.NewLiteral(l.F)
	case tokenize.LString:
		return vm3.NewLiteral(l.S)
	case tokenize.LBool:
		if l.B {
			return vm3.NewLiteral(1)
		}
		return vm3.NewLiteral(0)
	default:
		return nil
	}
}

func literal(node *parse.Node) ([]vm3.Data, error) {
	switch node.Kind {
	case parse.NdLiteral:
		return []vm3.Data{
			*vm3.NewOpcodeData(vm3.PUSH), *vm3.NewLiteralData(*literalFromField(node.LiteralField)),
		}, nil
	case parse.NdParenthesis:
		return expr(node.UnaryField.Value)
	case parse.NdCall:
		var program []vm3.Data
		// 引数があるかチェック
		if node.CallField.Args != nil {
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
			program = append(program, []vm3.Data{
				*vm3.NewOpcodeData(vm3.CALL), *vm3.NewLabelData(*vm3.NewLabel(false, node.CallField.Identifier.IdentField.Ident)),
			}...)
			// 引数分spを加算
			program = append(program, []vm3.Data{
				*vm3.NewOpcodeData(vm3.PUSH), *vm3.NewLiteralData(*vm3.NewLiteral(len(args))),
				*vm3.NewOpcodeData(vm3.POP), *vm3.NewRegisterTagData(vm3.R1),
				*vm3.NewOpcodeData(vm3.ADD), *vm3.NewRegisterTagData(vm3.R1), *vm3.NewRegisterTagData(vm3.RSP),
			}...)
		}

		// 戻り地に関する記述..?
		switch len(semOverall.KnownFunctions[node.CallField.Identifier.IdentField.Ident].Returns) {
		case 1:
			program = append(program, []vm3.Data{
				*vm3.NewOpcodeData(vm3.PUSH), *vm3.NewRegisterTagData(vm3.R10),
			}...)
		case 2:
			program = append(program, []vm3.Data{
				*vm3.NewOpcodeData(vm3.PUSH), *vm3.NewRegisterTagData(vm3.R11),
				*vm3.NewOpcodeData(vm3.PUSH), *vm3.NewRegisterTagData(vm3.R10),
			}...)
		default:
			return nil, fmt.Errorf("2つより多い戻り値には対応していません")
		}
		return program, nil
	case parse.NdIdent:
		loc := searchBPDistFromVarName(currentFnVariableBPs, currentNest, node.IdentField.Ident)
		if loc == -1 {
			for _, g := range globals {
				if node.IdentField.Ident == g {
					// global変数として存在する
					return []vm3.Data{
						*vm3.NewOpcodeData(vm3.MOV), *vm3.NewLabelData(*vm3.NewLabel(false, node.IdentField.Ident)), *vm3.NewRegisterTagData(vm3.R1),
						*vm3.NewOpcodeData(vm3.PUSH), *vm3.NewRegisterTagData(vm3.R1),
					}, nil
				}
			}
			return nil, fmt.Errorf("変数が定義されていません: %s", node.IdentField.Ident)
		}
		return []vm3.Data{
			*vm3.NewOpcodeData(vm3.PUSH), *vm3.NewOffsetData(*vm3.NewOffset(vm3.BP, -loc)),
		}, nil
	}
	return nil, fmt.Errorf("サポートされていないリテラルです")
}

func Compile(sem *analyze.Semantics) ([]vm3.Data, error) {
	semOverall = sem
	currentNest = 0
	if len(sem.OutsideValues) != 0 || len(sem.OutsideFunctions) != 0 {
		return nil, fmt.Errorf("リンクが不完全です")
	}
	var program []vm3.Data
	for _, n := range sem.Tree {
		switch n.Kind {
		case parse.NdFuncDef:
			frags, err := defFunction(n)
			if err != nil {
				return nil, err
			}
			program = append(program, frags...)
		case parse.NdAssign:
			globals = append(globals, n.AssignField.To.VarDeclField.Identifier.IdentField.Ident)
			dataSection = append(dataSection, []vm3.Data{
				*vm3.NewOpcodeData(vm3.MOV), *vm3.NewLiteralData(*literalFromField(n.AssignField.Value.LiteralField)), *vm3.NewLabelData(*vm3.NewLabel(false, n.AssignField.To.VarDeclField.Identifier.IdentField.Ident)),
			}...)
		case parse.NdVarDecl:
			globals = append(globals, n.VarDeclField.Identifier.IdentField.Ident)
		}
	}

	program = append(program, *vm3.NewLabelData(*vm3.NewLabel(true, ".data")))
	program = append(program, dataSection...)
	return program, nil
}
