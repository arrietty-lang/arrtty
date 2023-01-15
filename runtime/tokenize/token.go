package tokenize

type Token struct {
	Kind TokenKind
	Pos  *Position

	S string
	F float64
	I int

	Next *Token
}
