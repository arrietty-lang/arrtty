package parse

import (
	"fmt"
	"github.com/arrietty-lang/arrtty/preprocess/tokenize"
)

var token *tokenize.Token

func isEof() bool {
	return token.Kind == tokenize.Eof
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
		// "func" ident "(" funcParams ")" <funcReturns? "{">
		var returns *Node = nil
		if lcb := consumeKind(tokenize.Lcb); lcb == nil {
			// "func" ident "(" funcParams ")" <funcReturns>
			ret, err := funcReturns()
			if err != nil {
				return nil, err
			}
			returns = ret
			// "func" ident "(" funcParams ")" funcReturns <"{">
			_ = consumeKind(tokenize.Lcb)
		}
		// "func" ident "(" funcParams ")" funcReturns "{" <stmt>
		body, err := stmt()
		if err != nil {
			return nil, err
		}
		// "func" ident "(" funcParams ")" funcReturns "{" stmt <"}">
		_, err = expectKind(tokenize.Rcb)
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

func stmt() (*Node, error) {
	return nil, nil
}

func expr() (*Node, error) {
	return nil, nil
}

func assign() (*Node, error) {
	return nil, nil
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