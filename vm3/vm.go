package vm3

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
	exited        bool
}

func NewVm(program []Data, stackSize int) *Vm {
	var stack = make([]*Data, stackSize)
	var registers = map[RegisterTag]*Data{}
	var labelLocation = map[string]int{}

	return &Vm{
		program:       program,
		pc:            0,
		sp:            stackSize - 1,
		bp:            0,
		zf:            0,
		stack:         stack,
		registers:     registers,
		labelLocation: labelLocation,
		exited:        false,
	}
}

func (v *Vm) labelScan() error {
	for i, p := range v.program {
		if p.kind == KLabel {
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

func (v *Vm) Execute() error {
	err := v.labelScan()
	if err != nil {
		return err
	}
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
		default:
			return fmt.Errorf("サポートされていない操作です: %s", v.program[v.pc].opcode.String())
		}
	}
	return nil
}
