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
	registers map[string]*Fragment
}

func NewVm(program []*Fragment) *Vm {
	return &Vm{
		pc:      0,
		bp:      0,
		sp:      0,
		zf:      0,
		program: program,
		stack:   []*Fragment{},
		registers: map[string]*Fragment{
			"R1": {},
			"R2": {},
			"R3": {},
		},
	}
}

func (v *Vm) getRegister(r Register) (*Fragment, error) {
	reg, ok := v.registers[r.String()]
	if !ok {
		return nil, fmt.Errorf("unregistered register: %s", r.String())
	}
	return reg, nil
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

func (v *Vm) Execute() error {
	for !v.isProgramEof() {
		opcode := v.currentProgram()
		countOfOperand := opcode.CountOfOperand()
		operands := v.getArgs(countOfOperand)
		exit, err := v.execute(opcode, operands)
		if err != nil {
			return err
		}
		if exit {
			break
		}
		//v.pc += 1 + len(operands)
	}
	return nil
}

func (v *Vm) execute(opcode *Fragment, operands []*Fragment) (exit bool, err error) {
	switch *opcode.Opcode {
	case NOP:
		return v.nop(operands)
	case ADD:
		return v.add(operands)
	case SUB:
		return v.sub(operands)
	case MOV:
		return v.mov(operands)
	}
	return true, fmt.Errorf("unsupported opcode: %v", opcode.Opcode.String())
}

func (v *Vm) nop(operands []*Fragment) (bool, error) {
	v.pc += 1 + len(operands)
	return false, nil
}

//func set

func (v *Vm) add(operands []*Fragment) (bool, error) {
	defer func() {
		v.pc += 1 + len(operands)
	}()
	src := operands[0]
	dst := operands[1]
	switch dst.Kind {
	case REGISTER:
		switch src.Kind {
		case REGISTER:
			dstReg, err := v.getRegister(*dst.Register)
			if err != nil {
				return false, err
			}
			srcReg, err := v.getRegister(*src.Register)
			if err != nil {
				return false, err
			}
			if dstReg.Literal.Type == Int && srcReg.Literal.Type == Float {
				// src(float)をintにキャストする必要があるのでこれはエラー
				return false, fmt.Errorf("少数を整数として扱うことはできません")
			}
			if dstReg.Literal.Type == Float && srcReg.Literal.Type == Int {
				// intをfloatにキャストするのはok
				dstReg.Literal.F += float64(srcReg.Literal.I)
				return false, nil
			}
			if dstReg.Literal.Type == srcReg.Literal.Type && dstReg.Literal.Type == Int {
				dstReg.Literal.I += srcReg.Literal.I
				return false, nil
			}
			if dstReg.Literal.Type == srcReg.Literal.Type && dstReg.Literal.Type == Float {
				dstReg.Literal.F += srcReg.Literal.F
				return false, nil
			}
		}
		// todo : address, ...
	}
	return false, fmt.Errorf("unsupported")
}

func (v *Vm) sub(operands []*Fragment) (bool, error) {
	defer func() {
		v.pc += 1 + len(operands)
	}()
	src := operands[0]
	dst := operands[1]
	switch dst.Kind {
	case REGISTER:
		switch src.Kind {
		case REGISTER:
			dstReg, err := v.getRegister(*dst.Register)
			if err != nil {
				return false, err
			}
			srcReg, err := v.getRegister(*src.Register)
			if err != nil {
				return false, err
			}
			if dstReg.Literal.Type == Int && srcReg.Literal.Type == Float {
				// src(float)をintにキャストする必要があるのでこれはエラー
				return false, fmt.Errorf("少数を整数として扱うことはできません")
			}
			if dstReg.Literal.Type == Float && srcReg.Literal.Type == Int {
				// intをfloatにキャストするのはok
				dstReg.Literal.F -= float64(srcReg.Literal.I)
				return false, nil
			}
			if dstReg.Literal.Type == srcReg.Literal.Type && dstReg.Literal.Type == Int {
				dstReg.Literal.I -= srcReg.Literal.I
				return false, nil
			}
			if dstReg.Literal.Type == srcReg.Literal.Type && dstReg.Literal.Type == Float {
				dstReg.Literal.F -= srcReg.Literal.F
				return false, nil
			}
		}
		// todo : address, ...
	}
	return false, fmt.Errorf("unsupported")
}

func (v *Vm) mov(operands []*Fragment) (bool, error) {
	defer func() {
		v.pc += 1 + len(operands)
	}()
	src := operands[0]
	dst := operands[1]
	switch dst.Kind {
	case REGISTER:
		dstReg, err := v.getRegister(*dst.Register)
		if err != nil {
			return false, err
		}
		switch src.Kind {
		case REGISTER:
			srcReg, err := v.getRegister(*src.Register)
			if err != nil {
				return false, err
			}
			*dstReg = *srcReg
			return false, nil
		case LITERAL:
			dstReg.Kind = LITERAL
			dstReg.Literal = src.Literal
			return false, nil
		}
	}
	return false, fmt.Errorf("unsupported")
}
