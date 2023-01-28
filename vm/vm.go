package vm

import (
	"fmt"
)

type Vm struct {
	wayOut bool
	pc     int
	bp     int
	sp     int
	zf     int

	program   []*Fragment // 入力全体
	stack     []*Fragment
	registers map[string]*Fragment
	data      map[string]*Fragment // like section .data
}

func NewVm(program []*Fragment) *Vm {
	stack := make([]*Fragment, 10)
	return &Vm{
		wayOut:  false,
		pc:      0,
		bp:      0,
		sp:      len(stack) - 1,
		zf:      0,
		program: program,
		stack:   stack,
		registers: map[string]*Fragment{
			"R1": {},
			"R2": {},
			"R3": {},
		},
		data: map[string]*Fragment{},
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

func (v *Vm) addSPSafe(absoluteDiff int) error {
	if v.sp+absoluteDiff < 0 || len(v.stack)-1 < v.sp+absoluteDiff {
		return fmt.Errorf("access error")
	}
	v.sp += absoluteDiff
	return nil
}

func (v *Vm) subSPSafe(absoluteDiff int) error {
	if v.sp-absoluteDiff < 0 || len(v.stack)-1 < v.sp-absoluteDiff {
		return fmt.Errorf("access error")
	}
	v.sp -= absoluteDiff
	return nil
}

func (v *Vm) Execute() error {
	for !v.isProgramEof() {
		opcode := v.currentProgram()
		countOfOperand := opcode.CountOfOperand()
		operands := v.getArgs(countOfOperand)
		err := v.execute(opcode, operands)
		if err != nil {
			return err
		}
		if v.wayOut {
			break
		}
	}
	return nil
}

func (v *Vm) execute(opcode *Fragment, operands []*Fragment) error {
	switch *opcode.Opcode {
	case NOP:
		return v.nop(operands)
	case ADD:
		return v.add(operands)
	case SUB:
		return v.sub(operands)
	case MOV:
		return v.mov(operands)
	case CMP:
		return v.cmp(operands)
	case JMP:
		return v.jmp(operands)
	case JZ:
		return v.jz(operands)
	case PUSH:
		return v.push(operands)
	case POP:
		return v.pop(operands)
	case CALL:
		return v.call(operands)
	case RET:
		return v.ret(operands)
	case EXIT:
		return v.exit(operands)
	}
	v.wayOut = true
	return fmt.Errorf("unsupported opcode: %v", opcode.Opcode.String())
}

func (v *Vm) nop(operands []*Fragment) error {
	v.pc += 1 + len(operands)
	return nil
}

//func set

func (v *Vm) add(operands []*Fragment) error {
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
				return err
			}
			srcReg, err := v.getRegister(*src.Register)
			if err != nil {
				return err
			}
			if dstReg.Literal.Type == Int && srcReg.Literal.Type == Float {
				// src(float)をintにキャストする必要があるのでこれはエラー
				return fmt.Errorf("少数を整数として扱うことはできません")
			}
			if dstReg.Literal.Type == Float && srcReg.Literal.Type == Int {
				// intをfloatにキャストするのはok
				dstReg.Literal.F += float64(srcReg.Literal.I)
				return nil
			}
			if dstReg.Literal.Type == srcReg.Literal.Type && dstReg.Literal.Type == Int {
				dstReg.Literal.I += srcReg.Literal.I
				return nil
			}
			if dstReg.Literal.Type == srcReg.Literal.Type && dstReg.Literal.Type == Float {
				dstReg.Literal.F += srcReg.Literal.F
				return nil
			}
		}
		// todo : address, ...
	}
	return fmt.Errorf("unsupported")
}

func (v *Vm) sub(operands []*Fragment) error {
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
				return err
			}
			srcReg, err := v.getRegister(*src.Register)
			if err != nil {
				return err
			}
			if dstReg.Literal.Type == Int && srcReg.Literal.Type == Float {
				// src(float)をintにキャストする必要があるのでこれはエラー
				return fmt.Errorf("少数を整数として扱うことはできません")
			}
			if dstReg.Literal.Type == Float && srcReg.Literal.Type == Int {
				// intをfloatにキャストするのはok
				dstReg.Literal.F -= float64(srcReg.Literal.I)
				return nil
			}
			if dstReg.Literal.Type == srcReg.Literal.Type && dstReg.Literal.Type == Int {
				dstReg.Literal.I -= srcReg.Literal.I
				return nil
			}
			if dstReg.Literal.Type == srcReg.Literal.Type && dstReg.Literal.Type == Float {
				dstReg.Literal.F -= srcReg.Literal.F
				return nil
			}
		}
		// todo : address, ...
	}
	return fmt.Errorf("unsupported")
}

func (v *Vm) mov(operands []*Fragment) error {
	defer func() {
		v.pc += 1 + len(operands)
	}()
	src := operands[0]
	dst := operands[1]
	switch dst.Kind {
	case REGISTER:
		dstReg, err := v.getRegister(*dst.Register)
		if err != nil {
			return err
		}
		switch src.Kind {
		case REGISTER:
			srcReg, err := v.getRegister(*src.Register)
			if err != nil {
				return err
			}
			*dstReg = *srcReg
			return nil
		case LITERAL:
			dstReg.Kind = LITERAL
			dstReg.Literal = src.Literal
			return nil
		}
	}
	return fmt.Errorf("unsupported")
}

func (v *Vm) cmp(operands []*Fragment) error {
	defer func() {
		v.pc += 1 + len(operands)
	}()

	x1 := operands[0]
	x2 := operands[1]
	switch x1.Kind {
	case REGISTER:
		switch x2.Kind {
		case REGISTER:
			x1r, err := v.getRegister(*x1.Register)
			if err != nil {
				return err
			}
			x2r, err := v.getRegister(*x2.Register)
			if err != nil {
				return err
			}
			if isSameLiteral(x1r, x2r) {
				v.zf = 1
			} else {
				v.zf = 0
			}
			return nil
		case LITERAL:
			x1r, err := v.getRegister(*x1.Register)
			if err != nil {
				return err
			}
			if isSameLiteral(x1r, x2) {
				v.zf = 1
			} else {
				v.zf = 0
			}
			return nil
		}
	case LITERAL:
		switch x2.Kind {
		case REGISTER:
			x2r, err := v.getRegister(*x2.Register)
			if err != nil {
				return err
			}
			if isSameLiteral(x1, x2r) {
				v.zf = 1
			} else {
				v.zf = 0
			}
			return nil
		case LITERAL:
			if isSameLiteral(x1, x2) {
				v.zf = 1
			} else {
				v.zf = 0
			}
			return nil
		}
	}

	return fmt.Errorf("unsupported")
}

func (v *Vm) jmp(operands []*Fragment) error {
	x := operands[0]
	if x.Kind != LITERAL || x.Literal.Type != Int {
		return fmt.Errorf("unexpected")
	}
	v.pc = x.Literal.I
	return nil
}

func (v *Vm) jz(operands []*Fragment) error {
	// JUMP IF ZERO (ZFが1であればジャンプ)
	x := operands[0]
	if x.Kind != LITERAL || x.Literal.Type != Int {
		return fmt.Errorf("unexpected")
	}
	// ゼロなので飛ぶ
	if v.zf == 1 {
		v.pc = x.Literal.I
		return nil
	}
	// 飛ばなかった場合コマンドと引数を読み進める
	v.pc += 1 + len(operands)
	return nil
}

func (v *Vm) push(operands []*Fragment) error {
	defer func() {
		v.pc += 1 + len(operands)
	}()
	source := operands[0]
	switch source.Kind {
	case REGISTER:
		sourceValue, err := v.getRegister(*source.Register)
		if err != nil {
			return err
		}
		if err = v.subSPSafe(1); err != nil {
			return err
		}
		v.stack[v.sp] = sourceValue
		return nil
	case POINTER:
		switch *source.Pointer {
		case BP:
			sourceValue := NewLiteralFragment(NewInt(v.bp))
			if err := v.subSPSafe(1); err != nil {
				return err
			}
			v.stack[v.sp] = sourceValue
		case SP:
			sourceValue := NewLiteralFragment(NewInt(v.sp))
			if err := v.subSPSafe(1); err != nil {
				return err
			}
			v.stack[v.sp] = sourceValue
		}
	case ADDRESS:
		switch source.Address.Original {
		case BP:
			sourceValue := v.stack[v.bp+source.Address.Relative]
			if err := v.subSPSafe(1); err != nil {
				return err
			}
			v.stack[v.sp] = sourceValue
		case SP:
			sourceValue := v.stack[v.sp+source.Address.Relative]
			if err := v.subSPSafe(1); err != nil {
				return err
			}
			v.stack[v.sp] = sourceValue
		}
	}
	return fmt.Errorf("unsupported")
}

func (v *Vm) pop(operands []*Fragment) error {
	defer func() {
		v.pc += 1 + len(operands)
	}()
	dst := operands[0]
	switch dst.Kind {
	case REGISTER:
		dstReg, err := v.getRegister(*dst.Register)
		if err != nil {
			return err
		}
		*dstReg = *v.stack[v.sp]
		return v.addSPSafe(1)
	case POINTER:
		switch *dst.Pointer {
		case BP:
			value := v.stack[v.sp]
			if value.Kind != LITERAL || value.Literal.Type != Int {
				return fmt.Errorf("only integers can be assigned to bp")
			}
			v.bp = value.I
			return v.addSPSafe(1)
		case SP:
			value := v.stack[v.sp]
			if value.Kind != LITERAL || value.Literal.Type != Int {
				return fmt.Errorf("only integers can be assigned to sp")
			}
			v.sp = value.I
			return v.addSPSafe(1)
		}
	case ADDRESS:
		value := v.stack[v.sp]
		switch dst.Address.Original {
		case BP:
			v.stack[v.bp+dst.Address.Relative] = value
			return v.addSPSafe(1)
		case SP:
			v.stack[v.sp+dst.Address.Relative] = value
			return v.addSPSafe(1)
		}
	}
	return fmt.Errorf("unsupported")
}

func (v *Vm) call(operands []*Fragment) error {
	// 行き先がまともか？
	newLocation := operands[0]
	if newLocation.Kind != LITERAL || newLocation.Literal.Type != Int {
		return fmt.Errorf("unsupported")
	}
	// 戻ってくるところを保存
	if err := v.subSPSafe(1); err != nil {
		return err
	}
	v.stack[v.sp] = NewLiteralFragment(NewInt(v.pc + 2))
	v.pc = newLocation.Literal.I
	return nil
}

func (v *Vm) ret(operands []*Fragment) error {
	_ = operands
	rtnLocation := v.stack[v.sp]
	if err := v.addSPSafe(1); err != nil {
		return err
	}
	if rtnLocation.Kind != LITERAL || rtnLocation.Literal.Type != Int {
		return fmt.Errorf("unsupported")
	}
	v.pc = rtnLocation.Literal.I
	return nil
}

func (v *Vm) exit(operands []*Fragment) error {
	v.wayOut = true
	_ = operands
	return nil
}