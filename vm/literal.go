package vm

import (
	"github.com/arrietty-lang/arrtty/preprocess/parse"
	"github.com/arrietty-lang/arrtty/preprocess/tokenize"
	"strconv"
)

type LiteralType int

const (
	Int LiteralType = iota
	Float
	String
)

type Literal struct {
	Type LiteralType
	I    int
	F    float64
	S    string
}

func (l *Literal) String() string {
	switch l.Type {
	case Int:
		return strconv.Itoa(l.I)
	case Float:
		return strconv.FormatFloat(l.F, 'f', -1, 64)
	case String:
		return l.S
	}
	return ""
}

func FromLiteralField(field *parse.LiteralField) *Literal {
	switch field.Kind {
	case tokenize.LInt:
		return NewInt(field.I)
	case tokenize.LFloat:
		return NewFloat(field.F)
	case tokenize.LString:
		return NewString(field.S)
	case tokenize.LBool:
		//return new
	}
	return nil
}

func NewInt(i int) *Literal {
	return &Literal{
		Type: Int,
		I:    i,
		F:    0,
		S:    "",
	}
}

func NewFloat(f float64) *Literal {
	return &Literal{
		Type: Float,
		I:    0,
		F:    f,
		S:    "",
	}
}

func NewString(s string) *Literal {
	return &Literal{
		Type: String,
		I:    0,
		F:    0,
		S:    s,
	}
}

func isSameLiteral(x1, x2 Fragment) bool {
	if x1.Kind != LITERAL || x2.Kind != LITERAL {
		return false
	}
	if x1.Literal.Type != x2.Literal.Type {
		return false
	}
	switch x1.Type {
	case Int:
		return x1.Literal.I == x2.Literal.I
	case Float:
		return x1.Literal.F == x2.Literal.F
	case String:
		return x1.Literal.S == x2.Literal.S
	}
	return false
}
