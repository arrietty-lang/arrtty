package vm

type Opcode int

const (
	NOP Opcode = iota
	SET
	ADD
	SUB
	CMP
	LT
	GT
	LE
	GE
	JMP
	JZ
	JNZ
	JE
	JNE
	JL
	JLE
	JG
	JGE
	CALL
	RET
	MOV
	PUSH
	POP

	EXIT
)

func (o Opcode) CountOfOperand() int {
	switch o {
	case NOP:
		return 0
	case SET:
		return 2
	case ADD:
		return 2
	case SUB:
		return 2
	case CMP:
		return 2
	case LT:
		return 2
	case GT:
		return 2
	case LE:
		return 2
	case GE:
		return 2
	case JMP:
		return 1
	case JZ:
		return 1
	case JNZ:
		return 1
	case JE:
		return 2
	case JNE:
		return 2
	case JL:
		return 2
	case JLE:
		return 2
	case JG:
		return 2
	case JGE:
		return 2
	case CALL:
		return 1
	case RET:
		return 0
	case MOV:
		return 2
	case PUSH:
		return 1
	case POP:
		return 1
	case EXIT:
		return 1
	}
	return -1
}

var opcodes = [...]string{
	NOP:  "NOP",
	SET:  "SET",
	ADD:  "ADD",
	SUB:  "SUB",
	CMP:  "CMP",
	LT:   "LT",
	GT:   "GT",
	LE:   "LE",
	GE:   "GE",
	JMP:  "JMP",
	JZ:   "JZ",
	JNZ:  "JNZ",
	JE:   "JE",
	JNE:  "JNE",
	JL:   "JL",
	JLE:  "JLE",
	JG:   "JG",
	JGE:  "JGE",
	CALL: "CALL",
	RET:  "RET",
	MOV:  "MOV",
	PUSH: "PUSH",
	POP:  "POP",
	EXIT: "EXIT",
}

func (o Opcode) String() string {
	return opcodes[o]
}
