package assemble2

import (
	"fmt"
	"github.com/arrietty-lang/arrtty/preprocess/analyze"
	"github.com/arrietty-lang/arrtty/preprocess/parse"
	"github.com/arrietty-lang/arrtty/preprocess/tokenize"
	"github.com/arrietty-lang/arrtty/vm"
	"math/rand"
)

var semOverall *analyze.Semantics

var (
	currentFuncName        string
	currentFuncVariableBPs map[int]map[string]int // 現在の関数の、変数のBP
	currentNest            int
)

var (
	globalVariables map[string]*vm.Data
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

// ベースポインタを探す
func searchBPDistFromVarName(variables map[int]map[string]int, nest int, varName string) (int, bool) {
	for name, distance := range variables[nest] {
		if varName == name {
			return distance, true
		}
	}
	return 0, false
}

func toplevel(node *parse.Node) ([]vm.Data, error) {
	switch node.Kind {
	case parse.NdFuncDef:
		return tlFuncDef(node)
	case parse.NdImport:
		return tlImport(node)
	case parse.NdVarDecl:
		return tlGlobalVarDecl(node)
	case parse.NdAssign:
		return tlGlobalAssign(node)
	default:
		return nil, fmt.Errorf("tlはこれをサポートしていません: %v", node.String())
	}
}

func tlFuncDef(node *parse.Node) ([]vm.Data, error) {
	field := node.FuncDefField
	currentFuncName = field.Identifier.IdentField.Ident

	var program []vm.Data
	// 関数の定義、関数ラベルの作成
	program = append(program, []vm.Data{
		*vm.NewLabelData(*vm.NewLabel(true, currentFuncName)),
	}...)
	// メイン関数でなければ、関数終了時の戻り場所を記録
	if currentFuncName != "main" {
		program = append(program, []vm.Data{
			*vm.NewOpcodeData(vm.PUSH), *vm.NewRegisterTagData(vm.RBP),
			*vm.NewOpcodeData(vm.MOV), *vm.NewRegisterTagData(vm.RSP), *vm.NewRegisterTagData(vm.RBP),
		}...)
	}
	// 関数で使用されている変数の相対位置を事前に計算
	totalVars := 0
	distance := 1
	varsFromBPWithNest := map[int]map[string]int{} // 浅いものが深いものにアクセスできないようにするために分ける
	for nest, variables := range semOverall.KnownValues[currentFuncName] {
		varsFromBPWithNest[nest] = map[string]int{}
		for varName := range variables {
			varsFromBPWithNest[nest][varName] = distance
			distance++
			totalVars++
		}
	}
	currentFuncVariableBPs = varsFromBPWithNest

	// 関数内で使用される変数の数だけSPを下げる(変数領域の確保)
	program = append(program, []vm.Data{
		*vm.NewOpcodeData(vm.PUSH), *vm.NewLiteralDataWithRaw(totalVars),
		*vm.NewOpcodeData(vm.POP), *vm.NewRegisterTagData(vm.R1),
		*vm.NewOpcodeData(vm.SUB), *vm.NewRegisterTagData(vm.R1), *vm.NewRegisterTagData(vm.RSP),
	}...)

	// 引数と変数を結びつける(代入によって)
	if field.Parameters != nil {
		for i, param := range field.Parameters.PolynomialField.Values {
			r, ok := searchBPDistFromVarName(varsFromBPWithNest, 0, param.FuncParam.Identifier.IdentField.Ident)
			if !ok {
				return nil, fmt.Errorf("未定義変数: %s", param.FuncParam.Identifier.IdentField.Ident)
			}
			program = append(program, []vm.Data{
				*vm.NewOpcodeData(vm.MOV), *vm.NewOffsetData(*vm.NewOffset(vm.BP, i+2)), *vm.NewOffsetData(*vm.NewOffset(vm.BP, -r)),
			}...)
		}
	}

	// 関数の中身を展開
	for _, s := range field.Body.BlockField.Statements {
		f, err := stmt(s)
		if err != nil {
			return nil, err
		}
		program = append(program, f...)
	}

	return program, nil
}

func tlImport(node *parse.Node) ([]vm.Data, error) {
	return nil, fmt.Errorf("unimplemented")
}

func tlGlobalVarDecl(node *parse.Node) ([]vm.Data, error) {
	// アナライズされて関数のはじめにspまとめて引かれているので特にすることはないはず...?
	return nil, nil
}

func tlGlobalAssign(node *parse.Node) ([]vm.Data, error) {
	// todo
	field := node.AssignField
	// 代入したい値を取り出す
	value, err := expr(field.Value)
	if err != nil {
		return nil, err
	}
	// 一個であることを確認
	if len(value) != 1 {
		return nil, fmt.Errorf("グローバル変数は複数の値の代入に対応していません: %d", len(value))
	}
	v := value[0]
	// 即値であることを確認
	if v.GetKind() != vm.KLiteral {
		return nil, fmt.Errorf("グローバル変数は即値ではないものを代入できません: %v", v.String())
	}
	// 代入先の名前を取得
	var ident string
	switch node.AssignField.To.Kind {
	case parse.NdVarDecl:
		ident = node.AssignField.To.VarDeclField.Identifier.IdentField.Ident
	case parse.NdIdent:
		ident = node.AssignField.To.IdentField.Ident
	default:
		return nil, fmt.Errorf("代入先は対応していません: %s", node.String())
	}
	// 代入先が定義されているか確認
	_, ok := globalVariables[ident]
	if !ok {
		return nil, fmt.Errorf("未定義グローバル変数: %s", ident)
	}

	// グローバル変数として値を保存
	globalVariables[ident] = &v
	return []vm.Data{}, nil
}

func stmt(node *parse.Node) ([]vm.Data, error) {
	var program []vm.Data
	switch node.Kind {
	case parse.NdReturn:
		if len(node.PolynomialField.Values) > 2 {
			return nil, fmt.Errorf("２つ以上の戻り値は現在サポートされていません")
		}
		// 戻り値の準備
		returnRegs := []vm.RegisterTag{vm.R10, vm.R11}
		for i, rn := range node.PolynomialField.Values {
			rv, err := expr(rn)
			if err != nil {
				return nil, err
			}
			program = append(program, rv...)
			// 計算結果をR1に移動
			program = append(program, []vm.Data{
				*vm.NewOpcodeData(vm.POP), *vm.NewRegisterTagData(vm.R1),
				*vm.NewOpcodeData(vm.MOV), *vm.NewRegisterTagData(vm.R1), *vm.NewRegisterTagData(returnRegs[i]),
			}...)
		}
		// リターン本文
		// メインだけはリターンで何も返さない..?
		if currentFuncName != "main" {
			program = append(program, []vm.Data{
				*vm.NewOpcodeData(vm.MOV), *vm.NewRegisterTagData(vm.RBP), *vm.NewRegisterTagData(vm.RSP),
				*vm.NewOpcodeData(vm.POP), *vm.NewRegisterTagData(vm.RBP),
				*vm.NewOpcodeData(vm.RET),
			}...)
		}
		return program, nil
	case parse.NdIfElse:
		// 条件を計算
		cond, err := expr(node.IfElseField.Cond)
		if err != nil {
			return nil, err
		}
		program = append(program, cond...)

		// 各ジャンプ先のラベルを用意
		ifBlockLabel := "if_if_block_" + RandStringRunes(20)
		endLabel := "if_end_" + RandStringRunes(20)

		// 条件に合致したらIFブロックへ飛ぶ
		program = append(program, []vm.Data{
			*vm.NewOpcodeData(vm.JZ), *vm.NewLabelData(*vm.NewLabel(false, ifBlockLabel)),
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
		program = append(program, []vm.Data{
			*vm.NewOpcodeData(vm.JMP), *vm.NewLabelData(*vm.NewLabel(false, endLabel)),
		}...)

		// if-block
		program = append(program, *vm.NewLabelData(*vm.NewLabel(true, ifBlockLabel)))
		ifBlock, err := stmt(node.IfElseField.IfBlock)
		if err != nil {
			return nil, err
		}
		program = append(program, ifBlock...)

		// 最終的なジャンプ先
		program = append(program, *vm.NewLabelData(*vm.NewLabel(true, endLabel)))
		return program, nil
	case parse.NdFor:
		return nil, fmt.Errorf("unimplemented")
	case parse.NdBlock:
		for _, n := range node.BlockField.Statements {
			f, err := stmt(n)
			if err != nil {
				return nil, err
			}
			program = append(program, f...)
		}
		return program, nil
	default:
		return expr(node)
	}
}

func expr(node *parse.Node) ([]vm.Data, error) {
	return assign(node)
}

func assign(node *parse.Node) ([]vm.Data, error) {
	var program []vm.Data
	switch node.Kind {
	case parse.NdVarDecl:
		// 多分何もしなくていい
		return []vm.Data{}, nil
	case parse.NdShortVarDecl:
		return nil, fmt.Errorf("unimplemented")
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
			l, ok := searchBPDistFromVarName(currentFuncVariableBPs, currentNest, ident)
			if !ok {
				return nil, fmt.Errorf("未定義変数: %s", ident)
			}
			loc = l
		case parse.NdIdent:
			ident = node.AssignField.To.IdentField.Ident
			l, ok := searchBPDistFromVarName(currentFuncVariableBPs, currentNest, ident)
			if !ok {
				return nil, fmt.Errorf("未定義変数: %s", ident)
			}
			loc = l
		}
		// 変数の中身をスタックにプッシュ
		program = append(program, val...)
		// 関数内の変数
		if loc != -1 {
			program = append(program, []vm.Data{
				// valの結果を取り出す
				*vm.NewOpcodeData(vm.POP), *vm.NewRegisterTagData(vm.R1),
				// 変数の場所に格納
				*vm.NewOpcodeData(vm.MOV), *vm.NewRegisterTagData(vm.R1), *vm.NewOffsetData(*vm.NewOffset(vm.BP, -loc)),
			}...)
			return program, nil
		}

		// グローバル変数ならば
		for name := range globalVariables {
			if name == ident {
				program = append(program, []vm.Data{
					// valの結果を取り出す
					*vm.NewOpcodeData(vm.POP), *vm.NewRegisterTagData(vm.R1),
					// 変数の場所に格納
					*vm.NewOpcodeData(vm.MOV), *vm.NewRegisterTagData(vm.R1), *vm.NewLabelData(*vm.NewLabel(false, ident)),
				}...)
				return program, nil
			}
		}

		return nil, fmt.Errorf("変数の位置を特定できませんでした: %v", ident)
	default:
		return andor(node)
	}
}

func andor(node *parse.Node) ([]vm.Data, error) {
	switch node.Kind {
	case parse.NdAnd:
		return nil, fmt.Errorf("unimplemented")
	case parse.NdOr:
		return nil, fmt.Errorf("unimplemented")
	default:
		return equality(node)
	}
}

func equality(node *parse.Node) ([]vm.Data, error) {
	switch node.Kind {
	case parse.NdEq:
		return nil, fmt.Errorf("unimplemented")
	case parse.NdNe:
		return nil, fmt.Errorf("unimplemented")
	default:
		return relational(node)
	}
}

func relational(node *parse.Node) ([]vm.Data, error) {
	var program []vm.Data
	switch node.Kind {
	case parse.NdLt:
		lhs, err := relational(node.BinaryField.Lhs)
		if err != nil {
			return nil, err
		}
		rhs, err := relational(node.BinaryField.Rhs)
		if err != nil {
			return nil, err
		}
		program = append(program, lhs...)
		program = append(program, rhs...)

		program = append(program, []vm.Data{
			*vm.NewOpcodeData(vm.POP), *vm.NewRegisterTagData(vm.R2),
			*vm.NewOpcodeData(vm.POP), *vm.NewRegisterTagData(vm.R1),
			*vm.NewOpcodeData(vm.LT), *vm.NewRegisterTagData(vm.R1), *vm.NewRegisterTagData(vm.R2),
		}...)

		return program, nil
	case parse.NdLe:
		lhs, err := relational(node.BinaryField.Lhs)
		if err != nil {
			return nil, err
		}
		rhs, err := relational(node.BinaryField.Rhs)
		if err != nil {
			return nil, err
		}
		program = append(program, lhs...)
		program = append(program, rhs...)

		program = append(program, []vm.Data{
			*vm.NewOpcodeData(vm.POP), *vm.NewRegisterTagData(vm.R2),
			*vm.NewOpcodeData(vm.POP), *vm.NewRegisterTagData(vm.R1),
			*vm.NewOpcodeData(vm.LE), *vm.NewRegisterTagData(vm.R1), *vm.NewRegisterTagData(vm.R2),
		}...)

		return program, nil
	case parse.NdGt:
		lhs, err := relational(node.BinaryField.Lhs)
		if err != nil {
			return nil, err
		}
		rhs, err := relational(node.BinaryField.Rhs)
		if err != nil {
			return nil, err
		}
		program = append(program, lhs...)
		program = append(program, rhs...)

		program = append(program, []vm.Data{
			*vm.NewOpcodeData(vm.POP), *vm.NewRegisterTagData(vm.R2),
			*vm.NewOpcodeData(vm.POP), *vm.NewRegisterTagData(vm.R1),
			*vm.NewOpcodeData(vm.LT), *vm.NewRegisterTagData(vm.R2), *vm.NewRegisterTagData(vm.R1),
		}...)
		return program, nil
	case parse.NdGe:
		lhs, err := relational(node.BinaryField.Lhs)
		if err != nil {
			return nil, err
		}
		rhs, err := relational(node.BinaryField.Rhs)
		if err != nil {
			return nil, err
		}
		program = append(program, lhs...)
		program = append(program, rhs...)

		program = append(program, []vm.Data{
			*vm.NewOpcodeData(vm.POP), *vm.NewRegisterTagData(vm.R2),
			*vm.NewOpcodeData(vm.POP), *vm.NewRegisterTagData(vm.R1),
			*vm.NewOpcodeData(vm.LE), *vm.NewRegisterTagData(vm.R2), *vm.NewRegisterTagData(vm.R1),
		}...)

		return program, nil
	default:
		return add(node)
	}
}

func add(node *parse.Node) ([]vm.Data, error) {
	var program []vm.Data
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
		program = append(program, []vm.Data{
			*vm.NewOpcodeData(vm.POP), *vm.NewRegisterTagData(vm.R2),
			*vm.NewOpcodeData(vm.POP), *vm.NewRegisterTagData(vm.R1),
		}...)
		// 足し算/引き算
		program = append(program, []vm.Data{
			// r1 += r2
			*vm.NewOpcodeData(op), *vm.NewRegisterTagData(vm.R2), *vm.NewRegisterTagData(vm.R1),
		}...)
		// 結果R1をスタックにプッシュ
		program = append(program, []vm.Data{
			*vm.NewOpcodeData(vm.PUSH), *vm.NewRegisterTagData(vm.R1),
		}...)
		return program, nil
	}
	return mul(node)
}

func mul(node *parse.Node) ([]vm.Data, error) {
	var program []vm.Data
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
		program = append(program, []vm.Data{
			*vm.NewOpcodeData(vm.POP), *vm.NewRegisterTagData(vm.R2),
			*vm.NewOpcodeData(vm.POP), *vm.NewRegisterTagData(vm.R1),
		}...)
		// 掛け算割り算
		program = append(program, []vm.Data{
			// r1 += r2
			*vm.NewOpcodeData(op), *vm.NewRegisterTagData(vm.R2), *vm.NewRegisterTagData(vm.R1),
		}...)
		// 結果R1をスタックにプッシュ
		program = append(program, []vm.Data{
			*vm.NewOpcodeData(vm.PUSH), *vm.NewRegisterTagData(vm.R1),
		}...)
		return program, nil

	}
	return unary(node)
}

func unary(node *parse.Node) ([]vm.Data, error) {
	switch node.Kind {
	case parse.NdPlus:
		return nil, fmt.Errorf("unimplemented")
	case parse.NdMinus:
		return nil, fmt.Errorf("unimplemented")
	case parse.NdNot:
		return nil, fmt.Errorf("unimplemented")
	default:
		return primary(node)
	}
}

func primary(node *parse.Node) ([]vm.Data, error) {
	return access(node)
}

func access(node *parse.Node) ([]vm.Data, error) {
	switch node.Kind {
	case parse.NdAccess:
		return nil, fmt.Errorf("unimplemented")
	default:
		return literal(node)
	}
}

func literalFromField(l *parse.LiteralField) *vm.Literal {
	switch l.Kind {
	case tokenize.LInt:
		return vm.NewLiteral(l.I)
	case tokenize.LFloat:
		return vm.NewLiteral(l.F)
	case tokenize.LString:
		return vm.NewLiteral(l.S)
	case tokenize.LBool:
		if l.B {
			return vm.NewLiteral(1)
		}
		return vm.NewLiteral(0)
	default:
		return nil
	}
}

func literal(node *parse.Node) ([]vm.Data, error) {
	switch node.Kind {
	case parse.NdLiteral:
		return []vm.Data{
			*vm.NewOpcodeData(vm.PUSH), *vm.NewLiteralData(*literalFromField(node.LiteralField)),
		}, nil
	case parse.NdParenthesis:
		return expr(node.UnaryField.Value)
	case parse.NdCall:
		var program []vm.Data
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
			program = append(program, []vm.Data{
				*vm.NewOpcodeData(vm.CALL), *vm.NewLabelData(*vm.NewLabel(false, node.CallField.Identifier.IdentField.Ident)),
			}...)
			// 引数分spを加算
			program = append(program, []vm.Data{
				*vm.NewOpcodeData(vm.PUSH), *vm.NewLiteralData(*vm.NewLiteral(len(args))),
				*vm.NewOpcodeData(vm.POP), *vm.NewRegisterTagData(vm.R1),
				*vm.NewOpcodeData(vm.ADD), *vm.NewRegisterTagData(vm.R1), *vm.NewRegisterTagData(vm.RSP),
			}...)
		}

		// 戻り地に関する記述..?
		switch len(semOverall.KnownFunctions[node.CallField.Identifier.IdentField.Ident].Returns) {
		case 1:
			program = append(program, []vm.Data{
				*vm.NewOpcodeData(vm.PUSH), *vm.NewRegisterTagData(vm.R10),
			}...)
		case 2:
			program = append(program, []vm.Data{
				*vm.NewOpcodeData(vm.PUSH), *vm.NewRegisterTagData(vm.R11),
				*vm.NewOpcodeData(vm.PUSH), *vm.NewRegisterTagData(vm.R10),
			}...)
		default:
			return nil, fmt.Errorf("2つより多い戻り値には対応していません")
		}
		return program, nil
	case parse.NdIdent:
		loc, ok := searchBPDistFromVarName(currentFuncVariableBPs, currentNest, node.IdentField.Ident)
		if !ok {
			return nil, fmt.Errorf("未定義変数: %s", node.IdentField.Ident)
		}
		if loc == -1 {
			for name := range globalVariables {
				if node.IdentField.Ident == name {
					// global変数として存在する
					return []vm.Data{
						*vm.NewOpcodeData(vm.MOV), *vm.NewLabelData(*vm.NewLabel(false, node.IdentField.Ident)), *vm.NewRegisterTagData(vm.R1),
						*vm.NewOpcodeData(vm.PUSH), *vm.NewRegisterTagData(vm.R1),
					}, nil
				}
			}
			return nil, fmt.Errorf("変数が定義されていません: %s", node.IdentField.Ident)
		}
		return []vm.Data{
			*vm.NewOpcodeData(vm.PUSH), *vm.NewOffsetData(*vm.NewOffset(vm.BP, -loc)),
		}, nil
	default:
		return nil, fmt.Errorf("想定されていないリテラル: %s", node.String())
	}
}

func Compile(sem *analyze.Semantics) ([]vm.Data, error) {
	semOverall = sem
	// グローバル変数一覧を用意
	globalVariables = map[string]*vm.Data{}
	for ident := range semOverall.KnownValues["-global-"][0] {
		globalVariables[ident] = nil
	}

	currentNest = 0
	if len(semOverall.OutsideValues) != 0 || len(semOverall.OutsideFunctions) != 0 {
		return nil, fmt.Errorf("リンクが不完全です")
	}

	var program []vm.Data
	for _, node := range semOverall.Tree {
		p, err := toplevel(node)
		if err != nil {
			return nil, err
		}
		program = append(program, p...)
	}

	// 代入済みのグローバル変数を.dataに
	// todo : 代入されていないグローバル変数の扱い(vmがdataSectionに容量を確保できれば良いので、戻り値として名前だけあげれば良い...?)

	//program = append(program, *vm.NewLabelData(*vm.NewLabel(true, ".data")))
	//program = append(program, g...) // todo : global
	//program = append(program, []vm.Data{
	//	*vm.NewOpcodeData(vm.MOV), *vm.NewLiteralDataWithRaw(1), *vm.NewLabelData(*vm.NewLabel(false, "placeholder")),
	//}...)
	return program, nil
}
