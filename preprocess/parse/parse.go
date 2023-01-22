package parse

import (
	"fmt"
	"github.com/arrietty-lang/arrtty/preprocess/tokenize"
)

var token *tokenize.Token

func isEof() bool {
	return token.Kind == tokenize.Eof
}

func peekKind(kind tokenize.TokenKind) *tokenize.Token {
	if token.Kind == kind {
		return token
	}
	return nil
}

func consumeKind(kind tokenize.TokenKind) *tokenize.Token {
	if token.Kind == kind {
		tok := token
		token = token.Next
		return tok
	}
	return nil
}

func consumeIdent(s string) *tokenize.Token {
	if token.Kind == tokenize.Ident && s == token.Literal.S {
		tok := token
		token = token.Next
		return tok
	}
	return nil
}

func expectKind(kind tokenize.TokenKind) (*tokenize.Token, error) {
	if token.Kind == kind {
		tok := token
		token = token.Next
		return tok, nil
	}
	return nil, fmt.Errorf("unexpected token: %v", token.Kind.String())
}

func program() ([]*Node, error) {
	var nodes []*Node
	if !isEof() {
		n, err := toplevel()
		if err != nil {
			return nil, err
		}
		nodes = append(nodes, n)
	}
	return nodes, nil
}

func toplevel() (*Node, error) {
	// コメント
	if t := consumeKind(tokenize.Comment); t != nil {
		return NewCommentNode(t.Pos, t.Literal.S), nil
	}

	// 関数定義
	if t := consumeIdent("func"); t != nil {
		// "func" <ident>
		id, err := expectKind(tokenize.Ident)
		if err != nil {
			return nil, err
		}
		// "func" ident <"(">
		_, err = expectKind(tokenize.Lrb)
		if err != nil {
			return nil, err
		}
		// "func" ident "(" <funcParams?>
		var params *Node
		if consumeKind(tokenize.Rrb) == nil {
			params, err = funcParams()
			if err != nil {
				return nil, err
			}
			// "func" ident "(" funcParams? <")">
			_, err = expectKind(tokenize.Rrb)
			if err != nil {
				return nil, err
			}
		}
		// "func" ident "(" funcParams ")" <funcReturns? block
		var returns *Node = nil
		if peekKind(tokenize.Lcb) == nil {
			returns, err = funcReturns()
			if err != nil {
				return nil, err
			}
		}

		// "{" stmt "}"のときのデータ
		//if lcb := consumeKind(tokenize.Lcb); lcb == nil {
		//	// "func" ident "(" funcParams ")" <funcReturns>
		//	ret, err := funcReturns()
		//	if err != nil {
		//		return nil, err
		//	}
		//	returns = ret
		//	// "func" ident "(" funcParams ")" funcReturns <"{">
		//	_ = consumeKind(tokenize.Lcb)
		//}
		//// "func" ident "(" funcParams ")" funcReturns "{" <stmt>
		//body, err := stmt()
		//if err != nil {
		//	return nil, err
		//}
		//// "func" ident "(" funcParams ")" funcReturns "{" stmt <"}">
		//_, err = expectKind(tokenize.Rcb)
		//if err != nil {
		//	return nil, err
		//}

		body, err := stmt()
		if err != nil {
			return nil, err
		}

		return NewFuncDefNode(t.Pos,
			NewIdentNode(id.Pos, id.Literal.S),
			params,
			returns,
			body), nil
	}

	// import
	if t := consumeIdent("import"); t != nil {
		// "import" <target>
		target, err := expectKind(tokenize.String)
		if err != nil {
			return nil, err
		}
		return NewImportNode(t.Pos, target.Literal.S), nil
	}

	// 変数定義
	if c := consumeIdent("var"); c != nil {
		// "var" <ident>
		id, err := expectKind(tokenize.Ident)
		if err != nil {
			return nil, err
		}
		// "var" ident <types>
		typ, err := types()
		if err != nil {
			return nil, err
		}
		// "var" ident types "="?
		if a := consumeKind(tokenize.Assign); a == nil {
			// var decl
			return NewVarDeclNode(c.Pos, NewIdentNode(id.Pos, id.Literal.S), typ), nil
		}
		// var assign
		// "var" ident types "=" <andor>
		value, err := andor()
		if err != nil {
			return nil, err
		}
		return NewAssignNode(c.Pos, NewVarDeclNode(c.Pos, NewIdentNode(id.Pos, id.Literal.S), typ), value), nil
	}

	return nil, nil
}

//func block() (*Node, error) {
//	lcb, err := expectKind(tokenize.Lcb)
//	if err != nil {
//		return nil, err
//	}
//
//	var statements []*Node
//
//	for consumeKind(tokenize.Rcb) == nil {
//		statement, err := stmt()
//		if err != nil {
//			return nil, err
//		}
//		statements = append(statements, statement)
//	}
//
//	return NewBlockNode(lcb.Pos, statements), nil
//}

func stmt() (*Node, error) {
	// コメント
	if comment := consumeKind(tokenize.Comment); comment != nil {
		return NewCommentNode(comment.Pos, comment.Literal.S), nil
	}

	// block
	if lcb := consumeKind(tokenize.Lcb); lcb != nil {
		var statements []*Node
		for consumeKind(tokenize.Rcb) == nil {
			statement, err := stmt()
			if err != nil {
				return nil, err
			}
			statements = append(statements, statement)
		}
		return NewBlockNode(lcb.Pos, statements), nil
	}

	// return
	// 行終端の";"を消しちゃったからexpr?が判別できないかも。
	// とりあえず"}"が存在するかで判断をする
	if return_ := consumeIdent("return"); return_ != nil {
		if peekKind(tokenize.Rcb) != nil {
			return NewReturnNode(return_.Pos, nil), nil
		}
		value, err := expr()
		if err != nil {
			return nil, err
		}
		return NewReturnNode(return_.Pos, value), nil
	}

	// if else
	if if_ := consumeIdent("if"); if_ != nil {
		cond, err := expr()
		if err != nil {
			return nil, err
		}
		ifBlock, err := stmt()
		if err != nil {
			return nil, err
		}
		// 続いてelseがなかったら
		if consumeIdent("else") == nil { // elseあった場合はここで消費される
			return NewIfElseNode(if_.Pos, false, cond, ifBlock, nil), nil
		}
		elseBlock, err := stmt()
		if err != nil {
			return nil, err
		}
		return NewIfElseNode(if_.Pos, true, cond, ifBlock, elseBlock), nil
	}

	// for
	if for_ := consumeIdent("for"); for_ != nil {
		// for {}
		if peekKind(tokenize.Lcb) != nil {
			body, err := stmt()
			if err != nil {
				return nil, err
			}
			return NewForNode(for_.Pos, nil, nil, nil, body), nil
		}

		var init *Node
		var cond *Node
		var loop *Node
		// init
		if peekKind(tokenize.Lcb) == nil {
			i, err := expr()
			if err != nil {
				return nil, err
			}
			init = i
		}
		// cond
		if peekKind(tokenize.Lcb) == nil {
			c, err := expr()
			if err != nil {
				return nil, err
			}
			cond = c
		}
		// loop
		if peekKind(tokenize.Lcb) == nil {
			l, err := expr()
			if err != nil {
				return nil, err
			}
			loop = l
		}

		// 3このときは正常
		// 2このときはinit, condのじ順番で埋められるので正常
		// 1このときはinitだけ埋められるけど、これはcondであるべきなので修正する
		if init != nil && cond == nil && loop == nil {
			cond = init
			init = nil
		}

		// loop block
		body, err := stmt()
		if err != nil {
			return nil, err
		}
		return NewForNode(for_.Pos, init, cond, loop, body), nil
	}

	return expr()
}

func expr() (*Node, error) {
	return assign()
}

func assign() (*Node, error) {
	// var
	if var_ := consumeIdent("var"); var_ != nil {
		id, err := expectKind(tokenize.Ident)
		if err != nil {
			return nil, err
		}
		typ, err := types()
		if err != nil {
			return nil, err
		}
		idNode := NewIdentNode(id.Pos, id.Literal.S)
		// イコール、代入がなかった場合
		if eq := consumeKind(tokenize.Eq); eq == nil {
			return NewVarDeclNode(var_.Pos, idNode, typ), nil
		}
		// 代入あった場合
		value, err := andor()
		if err != nil {
			return nil, err
		}
		return NewAssignNode(var_.Pos, NewVarDeclNode(var_.Pos, idNode, typ), value), nil
	}

	// assign系
	andor_, err := andor()
	if err != nil {
		return nil, err
	}
	// 代入
	if consumeKind(tokenize.Assign) != nil {
		value, err := andor()
		if err != nil {
			return nil, err
		}
		return NewAssignNode(andor_.Pos, andor_, value), nil
	}
	// 簡略代入
	if consumeKind(tokenize.ColonAssign) != nil {
		value, err := andor()
		if err != nil {
			return nil, err
		}
		return NewShortVarDeclNode(andor_.Pos, andor_, value), nil
	}

	return andor()
}

func andor() (*Node, error) {
	n, err := equality()
	if err != nil {
		return nil, err
	}
	for {
		if and := consumeKind(tokenize.And); and != nil {
			rhs, err := equality()
			if err != nil {
				return nil, err
			}
			n = NewBinaryNode(NdAnd, and.Pos, n, rhs)
		} else if or := consumeKind(tokenize.Or); or != nil {
			rhs, err := equality()
			if err != nil {
				return nil, err
			}
			n = NewBinaryNode(NdOr, or.Pos, n, rhs)
		} else {
			break
		}
	}
	return n, nil
}

func equality() (*Node, error) {
	n, err := relational()
	if err != nil {
		return nil, err
	}
	for {
		if eq := consumeKind(tokenize.Eq); eq != nil {
			rhs, err := relational()
			if err != nil {
				return nil, err
			}
			n = NewBinaryNode(NdEq, eq.Pos, n, rhs)
		} else if ne := consumeKind(tokenize.Ne); ne != nil {
			rhs, err := relational()
			if err != nil {
				return nil, err
			}
			n = NewBinaryNode(NdNe, ne.Pos, n, rhs)
		} else {
			break
		}
	}
	return n, nil
}

func relational() (*Node, error) {
	n, err := add()
	if err != nil {
		return nil, err
	}
	for {
		if lt := consumeKind(tokenize.Lt); lt != nil {
			rhs, err := add()
			if err != nil {
				return nil, err
			}
			n = NewBinaryNode(NdLt, lt.Pos, n, rhs)
		} else if le := consumeKind(tokenize.Le); le != nil {
			rhs, err := add()
			if err != nil {
				return nil, err
			}
			n = NewBinaryNode(NdLe, le.Pos, n, rhs)
		} else if gt := consumeKind(tokenize.Gt); gt != nil {
			rhs, err := add()
			if err != nil {
				return nil, err
			}
			n = NewBinaryNode(NdGt, gt.Pos, n, rhs)
		} else if ge := consumeKind(tokenize.Ge); ge != nil {
			rhs, err := add()
			if err != nil {
				return nil, err
			}
			n = NewBinaryNode(NdGe, ge.Pos, n, rhs)
		} else {
			break
		}
	}
	return n, nil
}

func add() (*Node, error) {
	n, err := mul()
	if err != nil {
		return nil, err
	}
	for {
		if plus := consumeKind(tokenize.Add); plus != nil {
			rhs, err := mul()
			if err != nil {
				return nil, err
			}
			n = NewBinaryNode(NdAdd, plus.Pos, n, rhs)
		} else if minus := consumeKind(tokenize.Sub); minus != nil {
			rhs, err := mul()
			if err != nil {
				return nil, err
			}
			n = NewBinaryNode(NdSub, minus.Pos, n, rhs)
		} else {
			break
		}
	}
	return n, nil
}

func mul() (*Node, error) {
	n, err := unary()
	if err != nil {
		return nil, err
	}
	for {
		if star := consumeKind(tokenize.Mul); star != nil {
			rhs, err := unary()
			if err != nil {
				return nil, err
			}
			n = NewBinaryNode(NdMul, star.Pos, n, rhs)
		} else if div := consumeKind(tokenize.Div); div != nil {
			rhs, err := unary()
			if err != nil {
				return nil, err
			}
			n = NewBinaryNode(NdDiv, div.Pos, n, rhs)
		} else if mod := consumeKind(tokenize.Mod); mod != nil {
			rhs, err := unary()
			if err != nil {
				return nil, err
			}
			n = NewBinaryNode(NdMod, mod.Pos, n, rhs)
		} else {
			break
		}
	}
	return n, nil
}

func unary() (*Node, error) {
	if plus := consumeKind(tokenize.Add); plus != nil {
		return primary()
	} else if minus := consumeKind(tokenize.Sub); minus != nil {
		v, err := primary()
		if err != nil {
			return nil, err
		}
		return NewBinaryNode(
			NdSub,
			minus.Pos,
			NewLiteralNode(minus.Pos, tokenize.NewIntLiteral(0)),
			v), nil
	} else if not := consumeKind(tokenize.Not); not != nil {
		v, err := primary()
		if err != nil {
			return nil, err
		}
		return NewUnaryNode(NdNot, not.Pos, v), nil
	}
	return primary()
}

func primary() (*Node, error) {
	return access()
}

func access() (*Node, error) {
	return literal()
}

func literal() (*Node, error) {
	// "(" expr ")"
	if lrb := consumeKind(tokenize.Lrb); lrb != nil {
		expression, err := expr()
		if err != nil {
			return nil, err
		}
		_, err = expectKind(tokenize.Rrb)
		if err != nil {
			return nil, err
		}
		return NewUnaryNode(NdParenthesis, expression.Pos, expression), nil
	}

	if id := consumeKind(tokenize.Ident); id != nil {
		// call
		if lrb := consumeKind(tokenize.Lrb); lrb != nil {
			if consumeKind(tokenize.Rrb) != nil {
				return NewCallNode(id.Pos, NewIdentNode(id.Pos, id.Literal.S), nil), nil
			}
			args, err := callArgs()
			if err != nil {
				return nil, err
			}
			_, err = expectKind(tokenize.Rrb)
			if err != nil {
				return nil, err
			}
			return NewCallNode(id.Pos, NewIdentNode(id.Pos, id.Literal.S), args), nil
		}
		// ident
		return NewIdentNode(id.Pos, id.Literal.S), nil
	}

	if i := consumeKind(tokenize.Int); i != nil {
		return NewLiteralNode(i.Pos, i.Literal), nil
	}
	if f := consumeKind(tokenize.Float); f != nil {
		return NewLiteralNode(f.Pos, f.Literal), nil
	}
	if s := consumeKind(tokenize.String); s != nil {
		return NewLiteralNode(s.Pos, s.Literal), nil
	}
	if b := consumeKind(tokenize.Bool); b != nil {
		return NewLiteralNode(b.Pos, b.Literal), nil
	}
	if n := consumeKind(tokenize.Nil); n != nil {
		return NewLiteralNode(n.Pos, n.Literal), nil
	}

	return nil, fmt.Errorf("[%d:%d] unexpected token: %v",
		token.Pos.LineNo, token.Pos.Lat, token.Kind.String())
}

func types() (*Node, error) {
	id, err := expectKind(tokenize.Ident)
	if err != nil {
		return nil, err
	}
	return NewDataTypeNode(id.Pos, GetDataTypeByIdent(id.Literal.S)), nil
}

func callArgs() (*Node, error) {
	var args []*Node
	first, err := expr()
	if err != nil {
		return nil, err
	}
	args = append(args, first)
	for consumeKind(tokenize.Comma) != nil {
		arg, err := expr()
		if err != nil {
			return nil, err
		}
		args = append(args, arg)
	}
	return NewPolynomialNode(NdArgs, first.Pos, args), nil
}

func funcParams() (*Node, error) {
	var params []*Node

	firstIdT, err := expectKind(tokenize.Ident)
	if err != nil {
		return nil, err
	}
	firstId := NewIdentNode(firstIdT.Pos, firstIdT.Literal.S)
	firstType, err := types()
	if err != nil {
		return nil, err
	}
	params = append(params, NewFuncParamNode(firstIdT.Pos, firstId, firstType))

	for consumeKind(tokenize.Comma) != nil {
		idt, err := expectKind(tokenize.Ident)
		if err != nil {
			return nil, err
		}
		id := NewIdentNode(idt.Pos, idt.Literal.S)
		typ, err := types()
		if err != nil {
			return nil, err
		}
		params = append(params, NewFuncParamNode(idt.Pos, id, typ))
	}

	return NewPolynomialNode(NdParams, firstIdT.Pos, params), nil
}

func funcReturns() (*Node, error) {
	var returnTypes []*Node

	lrb := consumeKind(tokenize.Lrb)
	if lrb == nil {
		typ, err := types()
		if err != nil {
			return nil, err
		}
		return NewPolynomialNode(NdReturn, typ.Pos, []*Node{typ}), nil
	}

	for consumeKind(tokenize.Rrb) == nil {
		typ, err := types()
		if err != nil {
			return nil, err
		}
		returnTypes = append(returnTypes, typ)
		if consumeKind(tokenize.Comma) == nil {
			_, err = expectKind(tokenize.Rrb)
			if err != nil {
				return nil, err
			}
		}
	}

	return NewPolynomialNode(NdReturn, lrb.Pos, returnTypes), nil
}
