package tokenize

import "log"

type Token struct {
	Kind TokenKind
	Pos  *Position

	Literal *Literal
	Ident   string

	Next *Token
}

func NewToken(kind TokenKind, pos *Position, literal *Literal) *Token {
	return &Token{
		Kind:    kind,
		Pos:     pos,
		Literal: literal,
		Next:    nil,
	}
}

func NewIdentToken(pos *Position, s string) *Token {
	return &Token{
		Kind:    Ident,
		Pos:     pos,
		Literal: nil,
		Ident:   s,
		Next:    nil,
	}
}

func NewCommentToken(pos *Position, s string) *Token {
	return &Token{
		Kind:    Comment,
		Pos:     pos,
		Literal: nil,
		Ident:   s,
		Next:    nil,
	}
}

func NewWhiteToken(pos *Position, s string) *Token {
	return &Token{
		Kind:    White,
		Pos:     pos,
		Literal: nil,
		Ident:   s,
		Next:    nil,
	}
}

func NewNewlineToken(pos *Position, s string) *Token {
	return &Token{
		Kind:    Newline,
		Pos:     pos,
		Literal: nil,
		Ident:   s,
		Next:    nil,
	}
}

func NewOpSymbolChain(cur *Token, pos *Position, str string) *Token {
	var tok *Token
	switch str {
	case "(":
		tok = NewToken(Lrb, pos, nil)
	case ")":
		tok = NewToken(Rrb, pos, nil)
	case "[":
		tok = NewToken(Lsb, pos, nil)
	case "]":
		tok = NewToken(Rsb, pos, nil)
	case "{":
		tok = NewToken(Lcb, pos, nil)
	case "}":
		tok = NewToken(Rcb, pos, nil)
	case ".":
		tok = NewToken(Dot, pos, nil)
	case ",":
		tok = NewToken(Comma, pos, nil)
	case ":":
		tok = NewToken(Colon, pos, nil)
	case ";":
		tok = NewToken(Semi, pos, nil)

	case "+":
		tok = NewToken(Add, pos, nil)
	case "-":
		tok = NewToken(Sub, pos, nil)
	case "*":
		tok = NewToken(Mul, pos, nil)
	case "/":
		tok = NewToken(Div, pos, nil)
	case "%":
		tok = NewToken(Mod, pos, nil)

	case "==":
		tok = NewToken(Eq, pos, nil)
	case "!=":
		tok = NewToken(Ne, pos, nil)
	case ">":
		tok = NewToken(Gt, pos, nil)
	case "<":
		tok = NewToken(Lt, pos, nil)
	case ">=":
		tok = NewToken(Ge, pos, nil)
	case "<=":
		tok = NewToken(Le, pos, nil)

	case "=":
		tok = NewToken(Assign, pos, nil)
	case "+=":
		tok = NewToken(AddAssign, pos, nil)
	case "-=":
		tok = NewToken(SubAssign, pos, nil)
	case "*=":
		tok = NewToken(MulAssign, pos, nil)
	case "/=":
		tok = NewToken(DivAssign, pos, nil)
	case "%=":
		tok = NewToken(ModAssign, pos, nil)
	case ":=":
		tok = NewToken(ColonAssign, pos, nil)

	case "&&":
		tok = NewToken(And, pos, nil)
	case "||":
		tok = NewToken(Or, pos, nil)
	case "!":
		tok = NewToken(Not, pos, nil)
	default:
		log.Fatalf("unsupported operator/symbol: %s", str)
	}
	cur.Next = tok
	return tok
}

// IsLiteralIdent アイデントがリテラルに使用されているキーワードか
func IsLiteralIdent(s string) bool {
	switch s {
	case "true":
		return true
	case "false":
		return true
	case "null":
		return true
	}
	return false
}

func NewIdentChain(cur *Token, position *Position, s string) *Token {
	tok := NewIdentToken(position, s)
	cur.Next = tok
	return tok
}

func NewEofChain(cur *Token, pos *Position) *Token {
	tok := NewToken(Eof, pos, nil)
	cur.Next = tok
	return tok
}

func NewLiteralChain(cur *Token, pos *Position, literal *Literal) *Token {
	var tokenKind TokenKind
	switch literal.Kind {
	case LString:
		tokenKind = String
	case LInt:
		tokenKind = Int
	case LFloat:
		tokenKind = Float
	case LBool:
		tokenKind = Bool
	case LNull:
		tokenKind = Null
	}
	tok := NewToken(tokenKind, pos, literal)
	cur.Next = tok
	return tok
}

func NewCommentChain(cur *Token, pos *Position, s string) *Token {
	tok := NewCommentToken(pos, s)
	cur.Next = tok
	return tok
}

func NewWhiteChain(cur *Token, pos *Position, s string) *Token {
	tok := NewWhiteToken(pos, s)
	cur.Next = tok
	return tok
}

func NewNewlineChain(cur *Token, pos *Position, s string) *Token {
	tok := NewNewlineToken(pos, s)
	cur.Next = tok
	return tok
}
