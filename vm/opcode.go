package vm

import "fmt"

type Opcode int

const (
	NOP Opcode = iota
	// ADD `add x1 x2`でx2 += x1
	ADD
	// SUB `sub x1 x2`でx2 -= x1
	SUB
	MUL
	DIV
	// CMP `cmp x1 x2`で一致したらZF=1, そうでなければZF=0
	CMP
	LT
	GT
	LE
	GE
	// JMP `jmp x`でpc=x, xはpcの絶対位置
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
	// MSG `msg r '...'`でrに'...'を代入
	MSG
	// LEN 文字の長さBYTEを取得する
	LEN
	// SYSCALL kernel call
	SYSCALL

	EXIT
)

func (o Opcode) CountOfOperand() int {
	switch o {
	case NOP:
		return 0
	case ADD:
		return 2
	case SUB:
		return 2
	case MUL:
		return 2
	case DIV:
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
		return 0
	case MSG:
		return 2
	case LEN:
		return 2
	case SYSCALL:
		return 1
	}
	return -1
}

var opcodes = [...]string{
	NOP:     "NOP",
	ADD:     "ADD",
	SUB:     "SUB",
	MUL:     "MUL",
	DIV:     "DIV",
	CMP:     "CMP",
	LT:      "LT",
	GT:      "GT",
	LE:      "Le",
	GE:      "GE",
	JMP:     "JMP",
	JZ:      "JZ",
	JNZ:     "JNZ",
	JE:      "JE",
	JNE:     "JNE",
	JL:      "JL",
	JLE:     "JLE",
	JG:      "JG",
	JGE:     "JGE",
	CALL:    "CALL",
	RET:     "RET",
	MOV:     "MOV",
	PUSH:    "PUSH",
	POP:     "POP",
	EXIT:    "EXIT",
	MSG:     "MSG",
	LEN:     "LEN",
	SYSCALL: "SYSCALL",
}

func (o Opcode) String() string {
	return fmt.Sprintf("Opcode{ %s }", opcodes[o])
}
