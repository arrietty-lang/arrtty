package tokenize

type TokenKind int

const (
	_ TokenKind = iota
	Eof

	Comment
	White
	Newline

	Ident
	Float
	Int
	String
	RawString
	True
	False
	Null

	Lrb
	Rrb
	Lsb
	Rsb
	Lcb
	Rcb
	Dot
	Comma
	Colon
	Semi

	Add
	Sub
	Mul
	Div
	Mod

	Eq
	Ne
	Gt
	Lt
	Ge
	Le

	Assign
	AddAssign
	SubAssign
	MulAssign
	DivAssign
	ModAssign
	ColonAssign

	And
	Or
	Not
)

var tokenKinds = [...]string{
	Eof:         "Eof",
	Comment:     "Comment",
	White:       "White",
	Newline:     "Newline",
	Ident:       "Ident",
	Float:       "Float",
	Int:         "Int",
	String:      "String",
	RawString:   "RawString",
	True:        "True",
	False:       "False",
	Null:        "Null",
	Lrb:         "Lrb",
	Rrb:         "Rrb",
	Lsb:         "Lsb",
	Rsb:         "Rsb",
	Lcb:         "Lcb",
	Rcb:         "Rcb",
	Dot:         "Dot",
	Comma:       "Comma",
	Colon:       "Colon",
	Semi:        "Semi",
	Add:         "Add",
	Sub:         "Sub",
	Mul:         "Mul",
	Div:         "Div",
	Mod:         "Mod",
	Eq:          "Eq",
	Ne:          "Ne",
	Gt:          "Gt",
	Lt:          "Lt",
	Ge:          "Ge",
	Le:          "Le",
	Assign:      "Assign",
	AddAssign:   "AddAssign",
	SubAssign:   "SubAssign",
	MulAssign:   "MulAssign",
	DivAssign:   "DivAssign",
	ModAssign:   "ModAssign",
	ColonAssign: "ColonAssign",
	And:         "And",
	Or:          "Or",
	Not:         "Not",
}

func (t TokenKind) String() string {
	return tokenKinds[t]
}
