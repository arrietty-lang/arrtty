package vm

type Operation int

const (
	Illegal Operation = iota
	NOP
	SET
	ADD
	SUB
	CMP
	LT
	GT
	LTE
	GTE
	JMP
	JZ
	JNZ
	CALL
	RET
	CP
	PUSH
	POP

	EXIT
)
