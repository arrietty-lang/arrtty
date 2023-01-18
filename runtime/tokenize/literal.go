package tokenize

type LiteralKind int

const (
	LString LiteralKind = iota
	LInt
	LFloat
	LBool
	LNull
)

type Literal struct {
	Kind LiteralKind
	S    string
	I    int
	F    float64
	B    bool
}

func NewStringLiteral(s string) *Literal {
	return &Literal{
		Kind: LString,
		S:    s,
		I:    0,
		F:    0,
		B:    false,
	}
}

func NewIntLiteral(i int) *Literal {
	return &Literal{
		Kind: LInt,
		S:    "",
		I:    i,
		F:    0,
		B:    false,
	}
}

func NewFloatLiteral(f float64) *Literal {
	return &Literal{
		Kind: LFloat,
		S:    "",
		I:    0,
		F:    f,
		B:    false,
	}
}

func NewBoolLiteral(b bool) *Literal {
	return &Literal{
		Kind: LBool,
		S:    "",
		I:    0,
		F:    0,
		B:    b,
	}
}

func NewNullLiteral() *Literal {
	return &Literal{Kind: LNull}
}
