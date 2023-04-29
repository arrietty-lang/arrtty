package vm

import (
	"fmt"
	"github.com/gookit/slog"
)

type Vm struct {
	program       []Data
	pc            int
	sp            int
	bp            int
	zf            int
	stack         []*Data
	registers     map[RegisterTag]*Data
	labelLocation map[string]int
	data          map[string]*Data
	exited        bool
}

func NewVm(program []Data, stackSize int) *Vm {
	var stack = make([]*Data, stackSize)
	var registers = map[RegisterTag]*Data{}
	var labelLocation = map[string]int{}
	var data = map[string]*Data{}

	for i, d := range program {
		slog.Debug(i, d.String())
	}

	return &Vm{
		program:       program,
		pc:            0,
		sp:            stackSize - 1,
		bp:            0,
		zf:            0,
		stack:         stack,
		registers:     registers,
		labelLocation: labelLocation,
		data:          data,
		exited:        false,
	}
}

func (v *Vm) labelScan() error {
	for i, p := range v.program {
		if p.kind == KLabel && p.label.GetIsDefine() {
			v.labelLocation[p.label.GetName()] = i
			slog.Debug("labelScan", p.label.GetName(), i)
		}
	}
	return nil
}

func (v *Vm) calculateOffset(offset Offset) (int, error) {
	switch offset.GetPointer() {
	case SP:
		return v.sp + offset.GetRelation(), nil
	case BP:
		return v.bp + offset.GetRelation(), nil
	}
	return 0, fmt.Errorf("不正なポインタです")
}

func (v *Vm) _push(d Data) {
	v.sp--
	v.stack[v.sp] = &d
	//slog.Debug("_push", "into(sp)", v.sp)
}
func (v *Vm) _pop() Data {
	//slog.Debug("_pop", "from(sp)", v.sp)
	d := *v.stack[v.sp]
	v.stack[v.sp] = nil
	v.sp++
	return d
}

func (v *Vm) getRSP() int {
	return v.sp
}
func (v *Vm) setRSP(i int) {
	v.sp = i
}

func (v *Vm) getRBP() int {
	return v.bp
}
func (v *Vm) setRBP(i int) {
	v.bp = i
}

func (v *Vm) GetRegisterByTag(tag RegisterTag) (*Data, bool) {
	switch tag {
	case RSP:
		return NewLiteralDataWithRaw(v.sp), true
	case RBP:
		return NewLiteralDataWithRaw(v.bp), true
	default:
		d, ok := v.registers[tag]
		if d.literal.GetKind() == KString {
			slog.Debug("string returned")
		}
		return d, ok
	}
}

func (v *Vm) SetRegisterByTag(tag RegisterTag, data *Data) error {
	if data.literal.GetKind() == KString {
		slog.Debug("string set.")
	}
	switch tag {
	case RSP:
		if data.kind != KLiteral || data.literal.GetKind() != KInt {
			return fmt.Errorf("SPに%s.%sを代入することはできません", data.kind.String(), data.literal.GetKind())
		}
		v.sp = data.literal.GetInt()
		return nil
	case RBP:
		if data.kind != KLiteral || data.literal.GetKind() != KInt {
			return fmt.Errorf("BPに%s.%sを代入することはできません", data.kind.String(), data.literal.GetKind())
		}
		v.bp = data.literal.GetInt()
		return nil
	default:
		v.registers[tag] = data
		return nil
	}
}

func (v *Vm) GetDataByLabel(label string) (Data, bool) {
	d, ok := v.data[label]
	return *d, ok
}

func (v *Vm) ExitCode() (int, error) {
	d := v.registers[R10]
	if d == nil {
		return 0, nil
	}
	if d.kind != KLiteral || d.literal.GetKind() != KInt {
		return 0, fmt.Errorf("終了コードが不正な値です: %s", d.String())
	}
	return d.literal.GetInt(), nil
}

func (v *Vm) Execute() error {
	err := v.labelScan()
	if err != nil {
		return err
	}

	entryPoint, ok := v.labelLocation["main"]
	if !ok {
		return fmt.Errorf("main label not found")
	}
	v.pc = entryPoint

	for v.pc < len(v.program) && !v.exited {
		slog.Debug("Execute", "pc", v.pc, "data", v.program[v.pc].String())
		if v.program[v.pc].kind == KLabel {
			v.pc++
			continue
		} else if v.program[v.pc].kind != KOpcode {
			return fmt.Errorf("pcはopcodeを予想しましたが、%sが発見されました", v.program[v.pc].kind.String())
		}
		switch v.program[v.pc].opcode {
		case PUSH:
			err := v.Push()
			if err != nil {
				return err
			}
		case POP:
			err := v.Pop()
			if err != nil {
				return err
			}
		case ADD:
			err := v.Add()
			if err != nil {
				return err
			}
		case SUB:
			err := v.Sub()
			if err != nil {
				return err
			}
		case MOV:
			err := v.Mov()
			if err != nil {
				return err
			}
		case LT:
			err := v.Lt()
			if err != nil {
				return err
			}
		case JMP:
			err := v.Jmp()
			if err != nil {
				return err
			}
		case JZ:
			err := v.Jz()
			if err != nil {
				return err
			}
		case EXIT:
			err := v.Exit()
			if err != nil {
				return err
			}
		case LE:
			err := v.Le()
			if err != nil {
				return err
			}
		case CMP:
			err := v.Cmp()
			if err != nil {
				return err
			}
		case CALL:
			err := v.Call()
			if err != nil {
				return err
			}
		case RET:
			err := v.Ret()
			if err != nil {
				return err
			}
		default:
			return fmt.Errorf("サポートされていない操作です: %s", v.program[v.pc].opcode.String())
		}
	}
	return nil
}
