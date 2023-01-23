package vm

import (
	"fmt"
)

type Vm struct {
	pc int
	bp int
	sp int
	zf int

	program   []*Fragment // 入力全体
	stack     []*Fragment
	registers []*Fragment
}

func NewVm(program []*Fragment) *Vm {
	return &Vm{
		pc:        0,
		bp:        0,
		sp:        0,
		zf:        0,
		program:   program,
		stack:     []*Fragment{},
		registers: []*Fragment{},
	}
}

func (v *Vm) isProgramEof() bool {
	return v.pc >= len(v.program)
}

func (v *Vm) currentProgram() *Fragment {
	return v.program[v.pc]
}

func (v *Vm) getArgs(n int) []*Fragment {
	var args []*Fragment
	for i := 0; i < n; i++ {
		// なぜ+1しているのかは不明
		// 参考: func operands()
		// https://github.com/x0y14/volume/blob/main/src/vvm/vvm.go
		args = append(args, v.program[v.pc+i+1])
	}
	return args
}

func (v *Vm) addSPSafe(positiveDiff int) error {
	if v.sp+positiveDiff < 0 || len(v.stack)-1 < v.sp+positiveDiff {
		return fmt.Errorf("access error")
	}
	v.sp += positiveDiff
	return nil
}

func (v *Vm) subSPSafe(positiveDiff int) error {
	if v.sp-positiveDiff < 0 || len(v.stack)-1 < v.sp-positiveDiff {
		return fmt.Errorf("access error")
	}
	v.sp -= positiveDiff
	return nil
}

func (v *Vm) Execute() {
	for !v.isProgramEof() {
		opcode := v.currentProgram()
		countOfOperand := opcode.CountOfOperand()
		operands := v.getArgs(countOfOperand)
		fmt.Println(opcode)
		fmt.Println(operands)
		v.pc += 1 + len(operands)
	}
}
