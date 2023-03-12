package analyze

import (
	"fmt"
	"github.com/arrietty-lang/arrtty/preprocess/parse"
	"github.com/arrietty-lang/arrtty/preprocess/tokenize"
	"strings"
)

var nest int
var knownFunction map[string]*FnDataType
var outsideFunction []*parse.Node

// var knownValues map[string]map[string][]*parse.DataType

var knownValues map[string]map[int]map[string][]*parse.DataType

var outsideValues []*parse.Node

func dataTypes(d *parse.DataType) []*parse.DataType {
	return []*parse.DataType{d}
}

func isSameType(x1 []*parse.DataType, x2 []*parse.DataType) bool {
	if len(x1) != len(x2) {
		return false
	}
	for i, x := range x1 {
		if x != x2[i] {
			return false
		}
	}
	return true
}

func isCalculable(x []*parse.DataType) bool {
	if len(x) != 1 {
		return false
	}
	switch x[0] {
	case parse.RuntimeInt, parse.RuntimeFloat:
		return true
	}
	return false
}

func isComparable(x []*parse.DataType) bool {
	if len(x) != 1 {
		return false
	}
	switch x[0] {
	case parse.RuntimeInt, parse.RuntimeFloat:
		return true
	}
	return false
}

func function(node *parse.Node) error {
	field := node.FuncDefField
	name := field.Identifier.IdentField.Ident
	//knownValues[name] = map[string][]*parse.DataType{}
	knownValues[name] = map[int]map[string][]*parse.DataType{}
	knownValues[name][nest] = map[string][]*parse.DataType{}

	var params []*parse.DataType
	// パラメータの型情報を取り出す
	if field.Parameters != nil {
		for _, paramNode := range field.Parameters.PolynomialField.Values {
			param := paramNode.FuncParam
			knownValues[name][nest][param.Identifier.IdentField.Ident] = dataTypes(param.DataType.DataTypeField.DataType)
			params = append(params, param.DataType.DataTypeField.DataType)
		}
	}

	// 戻り値の型を順番通りに準備
	var definedReturnTypes []*parse.DataType
	if field.Returns != nil {
		for _, returnTypeNode := range field.Returns.PolynomialField.Values {
			definedReturnTypes = append(definedReturnTypes, returnTypeNode.DataTypeField.DataType)
		}
	}
	knownFunction[name] = &FnDataType{
		Params:  params,
		Returns: definedReturnTypes,
	}
	// definedReturnTypesNode := NewFunctionNode(definedReturnTypes)
	// ブロックを解析して得られた実際の戻り値の型
	analyzedReturnTypes, err := stmt(field.Body, name)
	if err != nil {
		return err
	}
	if !isSameType(definedReturnTypes, analyzedReturnTypes) {
		return fmt.Errorf("戻り値の型が定義と一致しません")
	}

	if name == "main" && (!isSameType(analyzedReturnTypes, nil) && !isSameType(analyzedReturnTypes, dataTypes(parse.RuntimeInt))) {
		return fmt.Errorf("mainはInt, なし, 以外の戻り値の型をサポートしていません")
	}
	return nil
}

//
//func if_(node *parse.Node, functionName string) ([]*parse.DataType, error) {
//	var returnTypes []*parse.DataType
//	nest++
//	// if
//	knownValues[functionName][nest] = map[string][]*parse.DataType{}
//	for _, s := range node.IfElseField.IfBlock.BlockField.Statements {
//		rt, err := stmt(s, functionName)
//		if err != nil {
//			return nil, err
//		}
//		if s.Kind == parse.NdReturn {
//			if returnTypes == nil {
//				returnTypes = rt
//			} else {
//				if !isSameType(returnTypes, rt) {
//					return nil, fmt.Errorf("ブロック内で返却される戻り値が変化しています")
//				}
//			}
//		}
//	}
//	if !node.IfElseField.UseElse {
//		nest--
//		return returnTypes, nil
//	}
//
//	// else
//	if node.IfElseField.ElseBlock.Kind != parse.NdBlock {
//		rt, err := stmt(node.IfElseField.ElseBlock, functionName)
//		if err != nil {
//			return nil, err
//		}
//		if node.IfElseField.ElseBlock.Kind == parse.NdReturn {
//			if returnTypes == nil {
//				returnTypes = rt
//			} else {
//				if !isSameType(returnTypes, rt) {
//					return nil, fmt.Errorf("ブロック内で返却される戻り値が変化しています")
//				}
//			}
//		}
//	} else {
//		knownValues[functionName][nest] = map[string][]*parse.DataType{}
//		for _, s := range node.IfElseField.ElseBlock.BlockField.Statements {
//			rt, err := stmt(s, functionName)
//			if err != nil {
//				return nil, err
//			}
//			if s.Kind == parse.NdReturn {
//				if returnTypes == nil {
//					returnTypes = rt
//				} else {
//					if !isSameType(returnTypes, rt) {
//						return nil, fmt.Errorf("ブロック内で返却される戻り値が変化しています")
//					}
//				}
//			}
//		}
//	}
//	nest--
//	return returnTypes, nil
//}

func if_(node *parse.Node, functionName string) ([]*parse.DataType, error) {
	nest++
	for _, s := range node.BlockField.Statements {
		rt, err := stmt(s, functionName)
		if err != nil {
			return nil, err
		}
		if s.Kind == parse.NdReturn {
			if !isSameType(knownFunction[functionName].Returns, rt) {
				return nil, fmt.Errorf("期待される戻り値の型と一致しません")
			}
		}
	}
	nest--
	return nil, nil
}

func else_(node *parse.Node, functionName string) ([]*parse.DataType, error) {
	nest++
	if node.Kind != parse.NdBlock {
		_, err := stmt(node, functionName)
		if err != nil {
			return nil, err
		}
	} else {
		for _, s := range node.BlockField.Statements {
			rt, err := stmt(s, functionName)
			if err != nil {
				return nil, err
			}
			if s.Kind == parse.NdReturn {
				if !isSameType(knownFunction[functionName].Returns, rt) {
					return nil, fmt.Errorf("期待される戻り値の型と一致しません")
				}
			}
		}
	}
	nest--
	return nil, nil
}

func for_(node *parse.Node, functionName string) ([]*parse.DataType, error) {
	nest++
	for _, s := range node.ForField.Body.BlockField.Statements {
		rt, err := stmt(s, functionName)
		if err != nil {
			return nil, err
		}
		if s.Kind == parse.NdReturn {
			if !isSameType(knownFunction[functionName].Returns, rt) {
				return nil, fmt.Errorf("期待される戻り値の型と一致しません")
			}
		}
	}
	nest--
	return nil, nil
}

func stmt(node *parse.Node, functionName string) ([]*parse.DataType, error) {
	switch node.Kind {
	case parse.NdReturn:
		var returnTypes []*parse.DataType
		for _, v := range node.PolynomialField.Values {
			rt, err := expr(v, functionName)
			if err != nil {
				return nil, err
			}
			returnTypes = append(returnTypes, rt...)
		}
		//if node. != nil {
		//	return expr(node.ReturnField.Values, functionName)
		//}
		return returnTypes, nil
	case parse.NdIfElse:
		// IF
		_, err := if_(node.IfElseField.IfBlock, functionName)
		if err != nil {
			return nil, err
		}
		//if !isSameType(knownFunction[functionName], rt) {
		//	return nil, fmt.Errorf("戻り値の型が一致しません")
		//}
		if !node.IfElseField.UseElse {
			return nil, nil
		}
		// ELSE
		_, err = else_(node.IfElseField.ElseBlock, functionName)
		if err != nil {
			return nil, err
		}
		//if !isSameType(knownFunction[functionName], rt) {
		//	return nil, fmt.Errorf("戻り値の型が一致しません")
		//}
		return nil, nil
	case parse.NdFor:
		_, err := for_(node, functionName)
		if err != nil {
			return nil, err
		}
		return nil, nil
	case parse.NdBlock:
		var returnTypes []*parse.DataType
		for _, s := range node.BlockField.Statements {
			rt, err := stmt(s, functionName)
			if err != nil {
				return nil, err
			}
			if s.Kind == parse.NdReturn {
				if returnTypes == nil {
					returnTypes = rt
				} else {
					if !isSameType(returnTypes, rt) {
						return nil, fmt.Errorf("ブロック内で返却される戻り値が変化しています")
					}
				}
			}
		}
		return returnTypes, nil

	}
	return expr(node, functionName)
}

func expr(node *parse.Node, functionName string) ([]*parse.DataType, error) {
	return assign(node, functionName)
}

func assign(node *parse.Node, functionName string) ([]*parse.DataType, error) {
	switch node.Kind {
	case parse.NdVarDecl:
		name := node.VarDeclField.Identifier.IdentField.Ident
		typ := node.VarDeclField.Type.DataTypeField.DataType
		knownValues[functionName][nest][name] = dataTypes(typ)
		return dataTypes(typ), nil
	case parse.NdShortVarDecl:
		name := node.ShortVarDeclField.Identifier.IdentField.Ident
		typ, err := expr(node.ShortVarDeclField.Value, functionName)
		if err != nil {
			return nil, err
		}
		knownValues[functionName][nest][name] = typ
		return nil, nil
	case parse.NdAssign:
		// 型の変化なし
		defType, err := assign(node.AssignField.To, functionName)
		if err != nil {
			return nil, err
		}
		actualType, err := assign(node.AssignField.Value, functionName)
		if err != nil {
			return nil, err
		}
		if !isSameType(defType, actualType) {
			return nil, fmt.Errorf("代入された値と宣言の型が一致しません: %v <- %v", defType[0].Ident, actualType[0].Ident)
		}

		return nil, nil
	}
	return andor(node, functionName)
}

func andor(node *parse.Node, functionName string) ([]*parse.DataType, error) {
	switch node.Kind {
	case parse.NdAnd, parse.NdOr:
		lhs, err := andor(node, functionName)
		if err != nil {
			return nil, err
		}
		rhs, err := andor(node, functionName)
		if err != nil {
			return nil, err
		}
		if !isSameType(lhs, rhs) {
			return nil, fmt.Errorf("条件連結は同じ型のみで使用可能です: L:%s, R:%s", lhs[0].Ident, rhs[0].Ident)
		}
		if !isSameType(lhs, dataTypes(parse.RuntimeBool)) {
			return nil, fmt.Errorf("条件連結はBoolのみで使用可能です: L:%s, R:%s", lhs[0].Ident, rhs[0].Ident)
		}
		return dataTypes(parse.RuntimeBool), nil
	}
	return equality(node, functionName)
}

func equality(node *parse.Node, functionName string) ([]*parse.DataType, error) {
	switch node.Kind {
	case parse.NdEq, parse.NdNe:
		// todo : errとnilの関係
		lhs, err := equality(node, functionName)
		if err != nil {
			return nil, err
		}
		rhs, err := equality(node, functionName)
		if err != nil {
			return nil, err
		}
		if !isSameType(lhs, rhs) {
			return nil, fmt.Errorf("値比較は同じ型のみで使用可能です: L:%v, R:%v", lhs[0].Ident, rhs[0].Ident)
		}
		return dataTypes(parse.RuntimeBool), nil
	}
	return relational(node, functionName)
}

func relational(node *parse.Node, functionName string) ([]*parse.DataType, error) {
	switch node.Kind {
	case parse.NdLt, parse.NdLe, parse.NdGt, parse.NdGe:
		lhs, err := relational(node.BinaryField.Lhs, functionName)
		if err != nil {
			return nil, err
		}
		rhs, err := relational(node.BinaryField.Rhs, functionName)
		if err != nil {
			return nil, err
		}
		if !isSameType(lhs, rhs) {
			return nil, fmt.Errorf("大小比較は同じ型のみで使用できます: L:%v, R:%v", lhs[0].Ident, rhs[0].Ident)
		}
		if !isComparable(lhs) || !isComparable(rhs) {
			return nil, fmt.Errorf("大小比較は比較可能な型のみで使用できます(Int,Float) : L:%v, R:%v", lhs[0].Ident, rhs[0].Ident)
		}
		return dataTypes(parse.RuntimeBool), nil
	}
	return add(node, functionName)
}

func add(node *parse.Node, functionName string) ([]*parse.DataType, error) {
	switch node.Kind {
	case parse.NdAdd:
		lhs, err := add(node.BinaryField.Lhs, functionName)
		if err != nil {
			return nil, err
		}
		rhs, err := add(node.BinaryField.Rhs, functionName)
		if err != nil {
			return nil, err
		}
		// todo : intとfloatの計算について
		// Intはfloatに自動キャストできる.
		if !isSameType(lhs, rhs) {
			return nil, fmt.Errorf("計算は同じ型のみで行えます")
		}
		// + だけは 文字列を許可
		if (!isCalculable(lhs) || !isCalculable(rhs)) && !isSameType(lhs, dataTypes(parse.RuntimeString)) {
			return nil, fmt.Errorf("掛け算は計算可能な型(Int, Float)のみで使用できます: L:%v, R:%v", lhs[0].Ident, rhs[0].Ident)
		}
		return lhs, nil
	case parse.NdSub:
		lhs, err := add(node.BinaryField.Lhs, functionName)
		if err != nil {
			return nil, err
		}
		rhs, err := add(node.BinaryField.Rhs, functionName)
		if err != nil {
			return nil, err
		}
		// todo : intとfloatの計算について
		// Intはfloatに自動キャストできる.
		if !isSameType(lhs, rhs) {
			return nil, fmt.Errorf("計算は同じ型のみで行えます")
		}
		if !isCalculable(lhs) || !isCalculable(rhs) {
			return nil, fmt.Errorf("掛け算は計算可能な型(Int, Float)のみで使用できます: L:%v, R:%v", lhs[0].Ident, rhs[0].Ident)
		}
		return lhs, nil
	}
	return mul(node, functionName)
}

func mul(node *parse.Node, functionName string) ([]*parse.DataType, error) {
	switch node.Kind {
	case parse.NdMul, parse.NdDiv, parse.NdMod:
		lhs, err := mul(node.BinaryField.Lhs, functionName)
		if err != nil {
			return nil, err
		}
		rhs, err := mul(node.BinaryField.Rhs, functionName)
		if err != nil {
			return nil, err
		}
		// todo : intとfloatの計算について
		// Intはfloatに自動キャストできる.
		if !isSameType(lhs, rhs) {
			return nil, fmt.Errorf("計算は同じ型のみで行えます")
		}
		if !isCalculable(lhs) || !isCalculable(rhs) {
			return nil, fmt.Errorf("掛け算は計算可能な型(Int, Float)のみで使用できます: L:%v, R:%v", lhs[0].Ident, rhs[0].Ident)
		}
		return lhs, nil
	}
	return unary(node, functionName)
}

func unary(node *parse.Node, functionName string) ([]*parse.DataType, error) {
	switch node.Kind {
	case parse.NdNot:
		p, err := primary(node.UnaryField.Value, functionName)
		if err != nil {
			return nil, err
		}
		if !isSameType(p, []*parse.DataType{parse.RuntimeBool}) {
			return nil, fmt.Errorf("notはbool以外を値にできません: %v", p[0].Ident)
		}
		return p, nil
	}
	return primary(node, functionName)
}

func primary(node *parse.Node, functionName string) ([]*parse.DataType, error) {
	return access(node, functionName)
}

func access(node *parse.Node, functionName string) ([]*parse.DataType, error) {
	if node.Kind == parse.NdPrefix {
		return literal(node.PrefixField.Child, functionName)
	}
	return literal(node, functionName)
}

func literal(node *parse.Node, functionName string) ([]*parse.DataType, error) {
	switch node.Kind {
	case parse.NdParenthesis:
		return expr(node.UnaryField.Value, functionName)
	case parse.NdIdent:
		if node.IdentField.Ident == "true" || node.IdentField.Ident == "false" {
			return dataTypes(parse.RuntimeBool), nil
		}
		if node.IdentField.Ident == "nil" {
			return dataTypes(parse.RuntimeNil), nil
		}
		// todo : global変数
		// nestの深い値は、0になるまで遡って定義を調べることでif文などから関数内の値の方を参照する
		for i := nest; 0 <= i; i-- {
			typ, ok := knownValues[functionName][i][node.IdentField.Ident]
			if ok {
				return typ, nil
			}
		}
		if strings.Contains(node.IdentField.Ident, ".") {
			outsideValues = append(outsideValues, node)
			return dataTypes(parse.RuntimeUnknown), nil
		}
		typ, ok := knownValues["-global-"][0][node.IdentField.Ident]
		if !ok {
			return nil, fmt.Errorf("ana: %s is not defined", node.IdentField.Ident)
		}
		return typ, nil
	case parse.NdCall:
		// 期待する引数型
		typ, ok := knownFunction[node.CallField.Identifier.IdentField.Ident]
		if !ok {
			if strings.Contains(node.CallField.Identifier.IdentField.Ident, ".") {
				outsideFunction = append(outsideFunction, node)
				return dataTypes(parse.RuntimeUnknown), nil
			}
			return nil, fmt.Errorf("function %s is not defined", node.CallField.Identifier.IdentField.Ident) // ?
		}

		// 関数呼び出しで引数を渡さなかった場合、NILポインタが発生するのでチェックしてあげる
		if node.CallField.Args == nil {
			if !isSameType(typ.Params, nil) {
				return nil, fmt.Errorf("関数呼び出しの引数と与えられた型が異なります")
			}
			return typ.Returns, nil
		}

		//node.CallField.Args
		var args []*parse.DataType
		for _, arg := range node.CallField.Args.PolynomialField.Values {
			argT, err := expr(arg, functionName)
			if err != nil {
				return nil, err
			}
			args = append(args, argT...)
		}
		if !isSameType(typ.Params, args) {
			return nil, fmt.Errorf("関数呼び出しの引数と与えられた型が異なります")
		}
		return typ.Returns, nil
	default:
		switch node.LiteralField.Literal.Kind {
		case tokenize.LInt:
			return dataTypes(parse.RuntimeInt), nil
		case tokenize.LFloat:
			return dataTypes(parse.RuntimeFloat), nil
		case tokenize.LString:
			return dataTypes(parse.RuntimeString), nil
		case tokenize.LBool:
			return dataTypes(parse.RuntimeBool), nil
		case tokenize.LNil:
			return dataTypes(parse.RuntimeNil), nil
		}
	}
	return nil, fmt.Errorf("unknown data type")
}

func globalDecl(node *parse.Node) error {
	knownValues["-global-"][0][node.VarDeclField.Identifier.IdentField.Ident] = dataTypes(node.VarDeclField.Type.DataTypeField.DataType)
	return nil
}
func globalAssign(node *parse.Node) error {
	if err := globalDecl(node.AssignField.To); err != nil {
		return err
	}
	typ := knownValues["-global-"][0][node.AssignField.To.VarDeclField.Identifier.IdentField.Ident]

	if node.AssignField.Value.Kind != parse.NdLiteral {
		return fmt.Errorf("グローバル変数では即値以外を代入することはできません")
	}
	valType, err := expr(node.AssignField.Value, "-global-")
	if err != nil {
		return err
	}

	if !isSameType(typ, valType) {
		return fmt.Errorf("global 型の異なる値を代入することはできません")
	}

	return nil
}

func Analyze(nodes []*parse.Node) (*Semantics, error) {
	knownValues = map[string]map[int]map[string][]*parse.DataType{}
	knownValues["-global-"] = map[int]map[string][]*parse.DataType{}
	knownValues["-global-"][0] = map[string][]*parse.DataType{}
	knownFunction = map[string]*FnDataType{}
	outsideValues = []*parse.Node{}
	outsideFunction = []*parse.Node{}

	for _, node := range nodes {
		switch node.Kind {
		case parse.NdFuncDef:
			//ne
			if err := function(node); err != nil {
				return nil, err
			}
		case parse.NdVarDecl:
			if err := globalDecl(node); err != nil {
				return nil, err
			}
		case parse.NdAssign:
			if err := globalAssign(node); err != nil {
				return nil, err
			}
		}
	}
	return &Semantics{
		KnownValues:      knownValues,
		KnownFunctions:   knownFunction,
		OutsideValues:    outsideValues,
		OutsideFunctions: outsideFunction,
		Tree:             nodes,
	}, nil
}
