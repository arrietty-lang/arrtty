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
