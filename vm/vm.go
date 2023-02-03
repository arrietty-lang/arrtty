package vm

import (
	"fmt"
	"log"
	"os"
)

type Vm struct {
	wayOut bool
	pc     int
	bp     int
	sp     int
	zf     int

	program   []Fragment // 入力全体
	stack     []Fragment
	registers map[string]Fragment
	data      map[string]Fragment // like section .data
	labels    map[string]int
}

func (v *Vm) Export() string {
	var result string
	pos := 0
	for {
		curt := v.program[pos]
		switch curt.Kind {
		case OPCODE:
			result += "\t" + curt.Opcode.String()
			for i := 1; i < curt.Opcode.CountOfOperand()+1; i++ {
				result += " " + v.program[pos+i].String()
				if i != curt.Opcode.CountOfOperand() {
					result += ","
				}
			}
			pos += 1 + curt.Opcode.CountOfOperand()
		case LABEL:
			result += curt.String()
			pos++
		default:
			log.Println(result)
			log.Fatalf("予期しない位置にデータを発見しました: %s", curt.String())
		}
		result += "\n"
		if len(v.program) <= pos {
			break
		}
	}
	return result
}

func NewVm(program []Fragment) *Vm {
	stack := make([]Fragment, 100)
	return &Vm{
		wayOut:  false,
		pc:      0,
		bp:      len(stack) - 1,
		sp:      len(stack) - 1,
		zf:      0,
		program: program,
		stack:   stack,
		registers: map[string]Fragment{
			"R1":  {},
			"R2":  {},
			"R3":  {},
			"R4":  {},
			"R5":  {},
			"R10": {},
			"R11": {},
			"EC":  {},
			"ED":  {},
			"EM":  {},
			"EP":  {},
			"EW":  {},
			"ER":  {},
		},
		data:   map[string]Fragment{},
		labels: map[string]int{},
	}
}

func (v *Vm) pushStack(f Fragment) error {
	v.sp--
	v.stack[v.sp] = f
	log.Printf("PUSH(INTO %v) %v", v.sp, f.String())
	return nil
}

func (v *Vm) popStack() Fragment {
	f := v.stack[v.sp]
	v.sp++
	log.Printf("POP(FROM %v) %v", v.sp-1, f.String())
	return f
}

func (v *Vm) getRegister(r Register) (Fragment, error) {
	reg, ok := v.registers[r.String()]
	if !ok {
		return Fragment{}, fmt.Errorf("unregistered register: %s", r.String())
	}
	return reg, nil
}

func (v *Vm) isProgramEof() bool {
	return v.pc >= len(v.program)
}

func (v *Vm) currentProgram() Fragment {
	return v.program[v.pc]
}

func (v *Vm) getArgs(n int) []Fragment {
	var args []Fragment
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

func (v *Vm) scan() error {
	for !v.isProgramEof() {
		f := v.currentProgram()
		if f.Kind == LABEL && f.Label.Define {
			v.labels[f.Id] = v.pc + 1
		}
		v.pc++
	}
	v.pc = 0
	return nil
}

func (v *Vm) ExitCode() int {
	f, err := v.getRegister(R10)
	if err != nil {
		log.Fatal(err)
	}
	if f.Kind == ILLEGAL {
		return 0
	}
	if f.Kind != LITERAL || f.Literal.Type != Int {
		log.Fatal("終了コードが不正な値です")
	}
	return f.Literal.I
}

func (v *Vm) Execute() error {
	_ = v.scan()

	entrypoint, ok := v.labels["main"]
	if !ok {
		return fmt.Errorf("label(main): not found")
	}
	v.pc = entrypoint

	for !v.isProgramEof() {
		opcode := v.currentProgram()
		if opcode.Kind == LABEL && opcode.Label.Define {
			v.pc++
			continue
		}
		countOfOperand := opcode.CountOfOperand()
		operands := v.getArgs(countOfOperand)
		opr := ""
		for _, o := range operands {
			opr += o.String() + ", "
		}
		log.Printf("[%v] %s\t%v", v.pc, opcode.String(), opr)
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

func (v *Vm) execute(opcode Fragment, operands []Fragment) error {
	switch *opcode.Opcode {
	case NOP:
		return v.nop(operands)
	case ADD:
		return v.add(operands)
	case SUB:
		return v.sub(operands)
	case MUL:
		return v.mul(operands)
	case DIV:
		return v.div(operands)
	case MOV:
		return v.mov(operands)
	case CMP:
		return v.cmp(operands)
	case LT:
		return v.lt(operands)
	case LE:
		return v.le(operands)
	case GT:
		return v.gt(operands)
	case GE:
		return v.ge(operands)
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
	case MSG:
		return v.msg(operands)
	case LEN:
		return v.len(operands)
	case SYSCALL:
		return v.syscall(operands)
	}
	v.wayOut = true
	return fmt.Errorf("unsupported opcode: %v", opcode.Opcode.String())
}

func (v *Vm) nop(operands []Fragment) error {
	v.pc += 1 + len(operands)
	return nil
}

//func set

func (v *Vm) add(operands []Fragment) error {
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
	case POINTER:
		switch src.Kind {
		case LITERAL:
			if src.Literal.Type != Int {
				return fmt.Errorf("ポインタに整数以外の即値を加えることはできません: %v", src.Literal.Type)
			}
			switch *dst.Pointer {
			case BP:
				v.bp += src.Literal.I
				return nil
			case SP:
				v.sp += src.Literal.I
				return nil
			}
		}
	}
	return fmt.Errorf("unsupported")
}

func (v *Vm) sub(operands []Fragment) error {
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
	case POINTER:
		switch src.Kind {
		case LITERAL:
			if src.Literal.Type != Int {
				return fmt.Errorf("ポインタから整数以外の即値を引くことはできません: %v", src.Literal.Type)
			}
			switch *dst.Pointer {
			case BP:
				v.bp -= src.Literal.I
				return nil
			case SP:
				v.sp -= src.Literal.I
				return nil
			}
		}
	}
	return fmt.Errorf("unsupported")
}

func (v *Vm) mul(operands []Fragment) error {
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
				dstReg.Literal.F *= float64(srcReg.Literal.I)
				return nil
			}
			if dstReg.Literal.Type == srcReg.Literal.Type && dstReg.Literal.Type == Int {
				dstReg.Literal.I *= srcReg.Literal.I
				return nil
			}
			if dstReg.Literal.Type == srcReg.Literal.Type && dstReg.Literal.Type == Float {
				dstReg.Literal.F *= srcReg.Literal.F
				return nil
			}
		}
		// todo : address, ...
	}
	return fmt.Errorf("unsupported")
}

func (v *Vm) div(operands []Fragment) error {
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
				dstReg.Literal.F /= float64(srcReg.Literal.I)
				return nil
			}
			if dstReg.Literal.Type == srcReg.Literal.Type && dstReg.Literal.Type == Int {
				dstReg.Literal.I /= srcReg.Literal.I
				return nil
			}
			if dstReg.Literal.Type == srcReg.Literal.Type && dstReg.Literal.Type == Float {
				dstReg.Literal.F /= srcReg.Literal.F
				return nil
			}
		}
		// todo : address, ...
	}
	return fmt.Errorf("unsupported")
}

func (v *Vm) mov(operands []Fragment) error {
	defer func() {
		v.pc += 1 + len(operands)
	}()
	src := operands[0]
	dst := operands[1]
	switch dst.Kind {
	case REGISTER:
		//dstReg, err := v.getRegister(*dst.Register)
		//if err != nil {
		//	return err
		//}
		switch src.Kind {
		case REGISTER:
			srcReg, err := v.getRegister(*src.Register)
			if err != nil {
				return err
			}
			v.registers[dst.Register.String()] = srcReg
			return nil
		case LITERAL:
			v.registers[dst.Register.String()] = src
			//dstReg.Literal = src.Literal
			return nil
		case VARIABLE:
			//dstReg.Kind = LITERAL
			value, ok := v.data[src.Variable.Name]
			if !ok {
				return fmt.Errorf("%s is not defined", src.Variable.Name)
			}
			v.registers[dst.Register.String()] = value
			return nil
		}
	case POINTER:
		switch src.Kind {
		case POINTER:
			switch *dst.Pointer {
			case BP:
				switch *src.Pointer {
				case SP:
					v.bp = v.sp
					return nil
				}
			case SP:
				switch *src.Pointer {
				case BP:
					v.sp = v.bp
					return nil
				}
			}
		}
	case ADDRESS:
		switch src.Kind {
		case ADDRESS:
			switch dst.Address.Original {
			case BP:
				switch src.Address.Original {
				case SP:
					v.stack[v.bp+dst.Address.Relative] = v.stack[v.sp+src.Address.Relative]
					return nil
				case BP:
					v.stack[v.bp+dst.Address.Relative] = v.stack[v.bp+src.Address.Relative]
					return nil
				}
			case SP:
				switch src.Address.Original {
				case BP:
					v.stack[v.sp+dst.Address.Relative] = v.stack[v.bp+src.Address.Relative]
					return nil
				case SP:
					v.stack[v.sp+dst.Address.Relative] = v.stack[v.sp+src.Address.Relative]
					return nil
				}
			}
		}
	}
	return fmt.Errorf("unsupported")
}

func (v *Vm) cmp(operands []Fragment) error {
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

func (v *Vm) lt(operands []Fragment) error {
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
			if x1r.Literal.Type != x2r.Literal.Type {
				return fmt.Errorf("must be same type")
			}
			var lessThan bool
			switch x1r.Literal.Type {
			case Int:
				lessThan = x1r.Literal.I < x2r.Literal.I
			case Float:
				lessThan = x1r.Literal.F < x2r.Literal.F
			default:
				return fmt.Errorf("lt is supporting only int and float")
			}
			if lessThan {
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
			if x1r.Literal.Type != x2.Literal.Type {
				return fmt.Errorf("must be same type")
			}
			var lessThan bool
			switch x1r.Literal.Type {
			case Int:
				lessThan = x1r.Literal.I < x2.Literal.I
			case Float:
				lessThan = x1r.Literal.F < x2.Literal.F
			default:
				return fmt.Errorf("lt is supporting only int and float")
			}
			if lessThan {
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
			if x1.Literal.Type != x2r.Literal.Type {
				return fmt.Errorf("must be same type")
			}
			var lessThan bool
			switch x1.Literal.Type {
			case Int:
				lessThan = x1.Literal.I < x2r.Literal.I
			case Float:
				lessThan = x1.Literal.F < x2r.Literal.F
			default:
				return fmt.Errorf("lt is supporting only int and float")
			}
			if lessThan {
				v.zf = 1
			} else {
				v.zf = 0
			}
			return nil
		case LITERAL:
			if x1.Literal.Type != x2.Literal.Type {
				return fmt.Errorf("must be same type")
			}
			var lessThan bool
			switch x1.Literal.Type {
			case Int:
				lessThan = x1.Literal.I < x2.Literal.I
			case Float:
				lessThan = x1.Literal.F < x2.Literal.F
			default:
				return fmt.Errorf("lt is supporting only int and float")
			}
			if lessThan {
				v.zf = 1
			} else {
				v.zf = 0
			}
			return nil
		}
	}

	return fmt.Errorf("unsupported")
}

func (v *Vm) le(operands []Fragment) error {
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
			if x1r.Literal.Type != x2r.Literal.Type {
				return fmt.Errorf("must be same type")
			}
			var lessThan bool
			switch x1r.Literal.Type {
			case Int:
				lessThan = x1r.Literal.I <= x2r.Literal.I
			case Float:
				lessThan = x1r.Literal.F <= x2r.Literal.F
			default:
				return fmt.Errorf("lt is supporting only int and float")
			}
			if lessThan {
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
			if x1r.Literal.Type != x2.Literal.Type {
				return fmt.Errorf("must be same type")
			}
			var lessThan bool
			switch x1r.Literal.Type {
			case Int:
				lessThan = x1r.Literal.I <= x2.Literal.I
			case Float:
				lessThan = x1r.Literal.F <= x2.Literal.F
			default:
				return fmt.Errorf("lt is supporting only int and float")
			}
			if lessThan {
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
			if x1.Literal.Type != x2r.Literal.Type {
				return fmt.Errorf("must be same type")
			}
			var lessThan bool
			switch x1.Literal.Type {
			case Int:
				lessThan = x1.Literal.I <= x2r.Literal.I
			case Float:
				lessThan = x1.Literal.F <= x2r.Literal.F
			default:
				return fmt.Errorf("lt is supporting only int and float")
			}
			if lessThan {
				v.zf = 1
			} else {
				v.zf = 0
			}
			return nil
		case LITERAL:
			if x1.Literal.Type != x2.Literal.Type {
				return fmt.Errorf("must be same type")
			}
			var lessThan bool
			switch x1.Literal.Type {
			case Int:
				lessThan = x1.Literal.I <= x2.Literal.I
			case Float:
				lessThan = x1.Literal.F <= x2.Literal.F
			default:
				return fmt.Errorf("lt is supporting only int and float")
			}
			if lessThan {
				v.zf = 1
			} else {
				v.zf = 0
			}
			return nil
		}
	}

	return fmt.Errorf("unsupported")
}

func (v *Vm) gt(operands []Fragment) error {
	tmp := operands[0]
	operands[0] = operands[1]
	operands[1] = tmp
	return v.lt(operands)
}

func (v *Vm) ge(operands []Fragment) error {
	tmp := operands[0]
	operands[0] = operands[1]
	operands[1] = tmp
	return v.le(operands)
}

func (v *Vm) jmp(operands []Fragment) error {
	x := operands[0]
	if x.Kind == LABEL {
		loc, ok := v.labels[x.Label.Id]
		if !ok {
			return fmt.Errorf("label(%s): not found", x.Label.Id)
		}
		v.pc = loc
		return nil
	}

	if x.Kind != LITERAL || x.Literal.Type != Int {
		return fmt.Errorf("unexpected")
	}
	v.pc = x.Literal.I
	return nil
}

func (v *Vm) jz(operands []Fragment) error {
	// JUMP IF ZERO (ZFが1であればジャンプ)
	x := operands[0]
	switch x.Kind {
	case LABEL:
		loc, ok := v.labels[x.Label.Id]
		if !ok {
			return fmt.Errorf("label(%s): not found", x.Label.Id)
		}
		if v.zf == 1 {
			v.pc = loc
			return nil
		}
	case LITERAL:
		if x.Literal.Type == Int {
			if v.zf == 1 {
				v.pc = x.Literal.I
				return nil
			}
		}
	}
	v.pc += 1 + len(operands)
	return nil
}

func (v *Vm) push(operands []Fragment) error {
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
		//if err = v.subSPSafe(1); err != nil {
		//	return err
		//}
		//v.stack[v.sp] = sourceValue
		_ = v.pushStack(sourceValue)
		return nil
	case POINTER:
		switch *source.Pointer {
		case BP:
			sourceValue := NewLiteralFragment(NewInt(v.bp))
			//if err := v.subSPSafe(1); err != nil {
			//	return err
			//}
			//v.stack[v.sp] = sourceValue
			_ = v.pushStack(sourceValue)
			return nil
		case SP:
			sourceValue := NewLiteralFragment(NewInt(v.sp))
			//if err := v.subSPSafe(1); err != nil {
			//	return err
			//}
			//v.stack[v.sp] = sourceValue
			_ = v.pushStack(sourceValue)
			return nil
		}
	case ADDRESS:
		switch source.Address.Original {
		case BP:
			sourceValue := v.stack[v.bp+source.Address.Relative]
			//if err := v.subSPSafe(1); err != nil {
			//	return err
			//}
			//v.stack[v.sp] = sourceValue
			_ = v.pushStack(sourceValue)
			return nil
		case SP:
			sourceValue := v.stack[v.sp+source.Address.Relative]
			//if err := v.subSPSafe(1); err != nil {
			//	return err
			//}
			//v.stack[v.sp] = sourceValue
			_ = v.pushStack(sourceValue)
			return nil
		}
	case LITERAL:
		//if err := v.subSPSafe(1); err != nil {
		//	return err
		//}
		//v.stack[v.sp] = source
		_ = v.pushStack(source)
		return nil
	}
	return fmt.Errorf("unsupported")
}

func (v *Vm) pop(operands []Fragment) error {
	defer func() {
		v.pc += 1 + len(operands)
	}()
	dst := operands[0]
	switch dst.Kind {
	case REGISTER:
		//dstReg, err := v.getRegister(*dst.Register)
		//if err != nil {
		//	return err
		//}
		//*dstReg = *v.stack[v.sp]
		r := v.popStack()
		v.registers[dst.Register.String()] = r
		return nil
	case POINTER:
		switch *dst.Pointer {
		case BP:
			//value := v.stack[v.sp]
			value := v.popStack()
			if value.Kind != LITERAL || value.Literal.Type != Int {
				return fmt.Errorf("only integers can be assigned to bp")
			}
			v.bp = value.I
			return nil
		case SP:
			//value := v.stack[v.sp]
			value := v.popStack()
			if value.Kind != LITERAL || value.Literal.Type != Int {
				return fmt.Errorf("only integers can be assigned to sp")
			}
			v.sp = value.I
			return nil
		}
	case ADDRESS:
		//value := v.stack[v.sp]
		value := v.popStack()
		switch dst.Address.Original {
		case BP:
			v.stack[v.bp+dst.Address.Relative] = value
			return nil
		case SP:
			v.stack[v.sp+dst.Address.Relative] = value
			return nil
		}
	}
	return fmt.Errorf("unsupported")
}

func (v *Vm) call(operands []Fragment) error {
	// 行き先がまともか？
	newLocation := operands[0]
	//if newLocation.Kind != LITERAL || newLocation.Literal.Type != Int {
	//	return fmt.Errorf("unsupported")
	//}
	var loc int
	switch newLocation.Kind {
	case LABEL:
		l, ok := v.labels[newLocation.Label.Id]
		if !ok {
			return fmt.Errorf("label(%s): not found", newLocation.Label.Id)
		}
		loc = l
	case LITERAL:
		if newLocation.Literal.Type == Int {
			loc = newLocation.Literal.I
		} else {
			return fmt.Errorf("unsupported")
		}
	default:
		return fmt.Errorf("unsupported")
	}
	// 戻ってくるところを保存
	//if err := v.subSPSafe(1); err != nil {
	//	return err
	//}
	//v.stack[v.sp] = NewLiteralFragment(NewInt(v.pc + 2))
	_ = v.pushStack(NewLiteralFragment(NewInt(v.pc + 2)))
	v.pc = loc
	return nil
}

func (v *Vm) ret(operands []Fragment) error {
	_ = operands
	//rtnLocation := v.stack[v.sp]
	//if err := v.addSPSafe(1); err != nil {
	//	return err
	//}
	rtnLocation := v.popStack()
	if rtnLocation.Kind != LITERAL || rtnLocation.Literal.Type != Int {
		return fmt.Errorf("unsupported")
	}
	v.pc = rtnLocation.Literal.I
	return nil
}

func (v *Vm) exit(operands []Fragment) error {
	v.wayOut = true
	_ = operands
	return nil
}

func (v *Vm) msg(operands []Fragment) error {
	defer func() {
		v.pc += 1 + len(operands)
	}()
	valLit := operands[0].Literal
	name := operands[1].Variable.Name
	if valLit.Type != String {
		return fmt.Errorf("msg value must be string")
	}
	v.data[name] = operands[0]
	return nil
}

func (v *Vm) len(operands []Fragment) error {
	defer func() {
		v.pc += 1 + len(operands)
	}()
	name := operands[0].Variable.Name
	dst, err := v.getRegister(*operands[1].Register)
	if err != nil {
		return err
	}
	variable, ok := v.data[name]
	if !ok {
		return fmt.Errorf("%s is not defined", name)
	}
	dst.Kind = LITERAL
	dst.Literal = NewInt(len([]byte(variable.Literal.S)))
	return nil
}

func (v *Vm) syscall(operands []Fragment) error {
	defer func() {
		v.pc += 1 + len(operands)
	}()
	op := SystemCall(operands[0].Literal.I)
	switch op {
	case WRITE:
		// 書き込むデータサイズ
		dataSize, err := v.getRegister(ED)
		if err != nil {
			return err
		}
		// 書き込むメッセージ全文
		message, err := v.getRegister(EW)
		if err != nil {
			return err
		}
		// 書き込むサイズ分だけ取り出す

		buf := []byte(message.Literal.S)
		if len(buf) > dataSize.Literal.I {
			buf = buf[:dataSize.Literal.I]
		}
		// 書き込み先の種類
		dst, err := v.getRegister(EM)
		if err != nil {
			return err
		}
		switch IODestination(dst.Literal.I) {
		case STDOUT:
			_, err = fmt.Fprint(os.Stdout, string(buf))
			if err != nil {
				return err
			}
			return nil
		case STDERR:
			_, err = fmt.Fprint(os.Stderr, string(buf))
			if err != nil {
				return err
			}
			return nil
		default:
			return fmt.Errorf("unsupported destination: %v", dst.Literal.I)
		}
	}
	return fmt.Errorf("unsupported system call")
}
