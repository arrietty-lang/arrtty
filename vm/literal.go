package vm

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

func isSameLiteral(x1, x2 *Fragment) bool {
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