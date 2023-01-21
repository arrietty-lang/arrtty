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
		// "func" ident "(" <funcParams>
		params, err := funcParams()
		if err != nil {
			return nil, err
		}
		// "func" ident "(" funcParams <")">
		_, err = expectKind(tokenize.Rrb)
		if err != nil {
			return nil, err
		}
		// "func" ident "(" funcParams ")" <funcReturns? block
		var returns *Node = nil

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
	return andor()
}

func andor() (*Node, error) {
	return nil, nil
}

func equality() (*Node, error) {
	return nil, nil
}

func relational() (*Node, error) {
	return nil, nil
}

func add() (*Node, error) {
	return nil, nil
}

func mul() (*Node, error) {
	return nil, nil
}

func unary() (*Node, error) {
	return nil, nil
}

func primary() (*Node, error) {
	return nil, nil
}

func access() (*Node, error) {
	return nil, nil
}

func literal() (*Node, error) {
	return nil, nil
}

func types() (*Node, error) {
	return nil, nil
}

func callArgs() (*Node, error) {
	return nil, nil
}

func funcParams() (*Node, error) {
	return nil, nil
}

func funcReturns() (*Node, error) {
	return nil, nil
}
