package vm3

import (
	"fmt"
	"github.com/gookit/slog"
)

func (v *Vm) Push() error {
	defer func() {
		v.pc += 1 + PUSH.CountOfOperand()
	}()
	value := &v.program[v.pc+1]
	slog.Info("Push", "val", value.String())
	v._push(*value)
	return nil
}

func (v *Vm) Pop() error {
	defer func() {
		v.pc += 1 + POP.CountOfOperand()
	}()
	into := v.program[v.pc+1]
	value := v._pop()
	switch into.kind {
	case KOffset:
		loc, err := v.calculateOffset(into.offset)
		if err != nil {
			return err
		}
		v.stack[loc] = &value
		slog.Info("Pop", "kind", "offset", "into(sp)", loc, "into(addr)", value.offset.AddressString(), "val", value.String())
		return nil
	case KRegisterTag:
		v.registers[into.registerTag] = &value
		slog.Info("Pop", "kind", "register", "into(register)", into.registerTag, "val", value.String())
		return nil
	}
	return fmt.Errorf("popの値代入先が不正です: %s", into.kind)
}

func (v *Vm) add(from, to Literal) (Literal, error) {
	// [x] int += float
	// [o] float += int
	// [o] int += int
	// [o] float += float
	// other: error
	switch to.GetKind() {
	case KInt:
		switch from.GetKind() {
		case KInt:
			// [o] int += int
			return *NewLiteral(to.GetInt() + from.GetInt()), nil
		case KFloat:
			// [x] int += float
			return Literal{}, fmt.Errorf("許されていないペアの加算です: %s += %s", to.GetKind().String(), from.GetKind().String())
		default:
			return Literal{}, fmt.Errorf("加算元の型が不正です: %s += %s", to.GetKind().String(), from.GetKind().String())
		}
	case KFloat:
		switch from.GetKind() {
		case KInt:
			// [o] float += int
			return *NewLiteral(to.GetFloat() + float64(from.GetInt())), nil
		case KFloat:
			// [o] float += float
			return *NewLiteral(to.GetFloat() + from.GetFloat()), nil
		default:
			return Literal{}, fmt.Errorf("加算元の型が不正です: %s += %s", to.GetKind().String(), from.GetKind().String())
		}
	default:
		return Literal{}, fmt.Errorf("加算先の型が不正です: %s += %s", to.GetKind().String(), from.GetKind().String())
	}
}

func (v *Vm) sub(from, to Literal) (Literal, error) {
	// [x] int -= float
	// [o] float -= int
	// [o] int -= int
	// [o] float -= float
	// other: error
	switch to.GetKind() {
	case KInt:
		switch from.GetKind() {
		case KInt:
			// [o] int -= int
			return *NewLiteral(to.GetInt() - from.GetInt()), nil
		case KFloat:
			// [x] int -= float
			return Literal{}, fmt.Errorf("許されていないペアの減算です: %s -= %s", to.GetKind().String(), from.GetKind().String())
		default:
			return Literal{}, fmt.Errorf("減算元の型が不正です: %s -= %s", to.GetKind().String(), from.GetKind().String())
		}
	case KFloat:
		switch from.GetKind() {
		case KInt:
			// [o] float -= int
			return *NewLiteral(to.GetFloat() - float64(from.GetInt())), nil
		case KFloat:
			// [o] float -= float
			return *NewLiteral(to.GetFloat() - from.GetFloat()), nil
		default:
			return Literal{}, fmt.Errorf("減算元の型が不正です: %s -= %s", to.GetKind().String(), from.GetKind().String())
		}
	default:
		return Literal{}, fmt.Errorf("減算先の型が不正です: %s -= %s", to.GetKind().String(), from.GetKind().String())
	}
}

func (v *Vm) Add() error {
	defer func() {
		v.pc += 1 + ADD.CountOfOperand()
	}()
	from := v.program[v.pc+1]
	to := v.program[v.pc+2]
	switch to.kind {
	case KRegisterTag:
		switch from.kind {
		case KRegisterTag:
			// RTo += RFrom
			fromVal := *v.registers[from.registerTag]
			toVal := *v.registers[to.registerTag]
			slog.Info("ADD", "from", from.registerTag.String(), "to", to.registerTag.String())
			result, err := v.add(fromVal.literal, toVal.literal)
			if err != nil {
				return err
			}
			v.registers[to.registerTag] = NewLiteralData(result)
			return nil
		default:
			return fmt.Errorf("addはfrom: %sに対応していません", from.kind.String())
		}
	default:
		return fmt.Errorf("addはto: %sに対応していません", to.kind.String())
	}
	// return fmt.Errorf("予期しないエラー")
}

func (v *Vm) Sub() error {
	defer func() {
		v.pc += 1 + SUB.CountOfOperand()
	}()
	from := v.program[v.pc+1]
	to := v.program[v.pc+2]
	switch to.kind {
	case KRegisterTag:
		switch from.kind {
		case KRegisterTag:
			// RTo -= RFrom
			fromVal := *v.registers[from.registerTag]
			toVal := *v.registers[to.registerTag]
			slog.Info("SUB", "from", from.registerTag.String(), "to", to.registerTag.String())
			result, err := v.sub(fromVal.literal, toVal.literal)
			if err != nil {
				return err
			}
			v.registers[to.registerTag] = NewLiteralData(result)
			return nil
		default:
			return fmt.Errorf("fromはfrom: %sに対応していません", from.kind.String())
		}
	default:
		return fmt.Errorf("fromはto: %sに対応していません", to.kind.String())
	}
}

func (v *Vm) Mov() error {
	defer func() {
		v.pc += 1 + MOV.CountOfOperand()
	}()
	from := v.program[v.pc+1]
	to := v.program[v.pc+2]
	switch to.kind {
	case KRegisterTag:
		switch from.kind {
		case KRegisterTag:
			// RTo = RFrom
			fromVal := *v.registers[from.registerTag]
			v.registers[to.registerTag] = &fromVal
			return nil
		case KOffset:
			// RTo = FromOffset
			fromLoc, err := v.calculateOffset(from.offset)
			if err != nil {
				return err
			}
			fromVal := *v.stack[fromLoc]
			v.registers[to.registerTag] = &fromVal
			return nil
		default:
			return fmt.Errorf("代入元が不明です: %s", from.kind.String())
		}
	case KOffset:
		switch from.kind {
		case KRegisterTag:
			// ToOffset = RFrom
			fromVal := *v.registers[from.registerTag]
			toLoc, err := v.calculateOffset(to.offset)
			if err != nil {
				return err
			}
			v.stack[toLoc] = &fromVal
			return nil
		case KOffset:
			// ToOffset = FromOffset
			fromLoc, err := v.calculateOffset(from.offset)
			if err != nil {
				return err
			}
			fromVal := *v.stack[fromLoc]
			toLoc, err := v.calculateOffset(to.offset)
			if err != nil {
				return err
			}
			v.stack[toLoc] = &fromVal
			return nil
		default:
			return fmt.Errorf("代入元が不明です: %s", from.kind.String())
		}
	default:
		return fmt.Errorf("代入先が不明です: %s", to.kind.String())
	}
	//return fmt.Errorf("mov 不明なエラー")
}

func (v *Vm) Jmp() error {
	newLocLabel := v.program[v.pc+1]
	switch newLocLabel.kind {
	case KLabel:
		loc := v.labelLocation[newLocLabel.label.GetName()]
		v.pc = loc
		return nil
	default:
		return fmt.Errorf("ジャンプ先の型が不正です: %s", newLocLabel.kind.String())
	}
}

func (v *Vm) Jz() error {
	newLocLabel := v.program[v.pc+1]
	switch newLocLabel.kind {
	case KLabel:
		if v.zf != 1 {
			v.pc += 1 + JZ.CountOfOperand()
			return nil
		}
		loc := v.labelLocation[newLocLabel.label.GetName()]
		v.pc = loc
		return nil
	default:
		return fmt.Errorf("ジャンプ先の型が不正です: %s", newLocLabel.kind.String())
	}
}

func (v *Vm) lt(lhs, rhs Literal) (bool, error) {
	// [o] int < int
	// [o] int < float
	// [o] float < int
	switch lhs.GetKind() {
	case KInt:
		switch rhs.GetKind() {
		case KInt:
			// int < int
			return lhs.GetInt() < rhs.GetInt(), nil
		case KFloat:
			// int < float
			return float64(lhs.GetInt()) < rhs.GetFloat(), nil
		default:
			return false, fmt.Errorf("右辺の型が不正です: %s < %s", lhs.GetKind().String(), rhs.GetKind().String())
		}
	case KFloat:
		switch rhs.GetKind() {
		case KInt:
			// float < int
			return lhs.GetFloat() < float64(rhs.GetInt()), nil
		case KFloat:
			// float < float
			return lhs.GetFloat() < rhs.GetFloat(), nil
		default:
			return false, fmt.Errorf("右辺の型が不正です: %s < %s", lhs.GetKind().String(), rhs.GetKind().String())
		}
	default:
		return false, fmt.Errorf("左辺の型が不正です: %s < %s", lhs.GetKind().String(), rhs.GetKind().String())
	}
}

func (v *Vm) le(lhs, rhs Literal) (bool, error) {
	// [o] int <= int
	// [o] int <= float
	// [o] float <= int
	switch lhs.GetKind() {
	case KInt:
		switch rhs.GetKind() {
		case KInt:
			// int <= int
			return lhs.GetInt() <= rhs.GetInt(), nil
		case KFloat:
			// int <= float
			return float64(lhs.GetInt()) <= rhs.GetFloat(), nil
		default:
			return false, fmt.Errorf("右辺の型が不正です: %s <= %s", lhs.GetKind().String(), rhs.GetKind().String())
		}
	case KFloat:
		switch rhs.GetKind() {
		case KInt:
			// float <= int
			return lhs.GetFloat() <= float64(rhs.GetInt()), nil
		case KFloat:
			// float <= float
			return lhs.GetFloat() <= rhs.GetFloat(), nil
		default:
			return false, fmt.Errorf("右辺の型が不正です: %s <= %s", lhs.GetKind().String(), rhs.GetKind().String())
		}
	default:
		return false, fmt.Errorf("左辺の型が不正です: %s <= %s", lhs.GetKind().String(), rhs.GetKind().String())
	}
}

func (v *Vm) Lt() error {
	defer func() {
		v.pc += 1 + LT.CountOfOperand()
	}()

	// lh < rh
	lhs := v.program[v.pc+1]
	rhs := v.program[v.pc+2]

	switch lhs.kind {
	case KRegisterTag:
		switch rhs.kind {
		case KRegisterTag:
			// Rlh < Rrh
			lhsVal := *v.registers[lhs.registerTag]
			rhsVal := *v.registers[rhs.registerTag]

			lessThan, err := v.lt(lhsVal.literal, rhsVal.literal)
			if err != nil {
				return err
			}

			if lessThan {
				v.zf = 1
			} else {
				v.zf = 0
			}

			return nil
		default:
			return fmt.Errorf("右辺は不正な型です: %s < %s", lhs.kind.String(), rhs.kind.String())
		}
	default:
		return fmt.Errorf("左辺は不正な型です: %s < %s", lhs.kind.String(), rhs.kind.String())
	}
}

func (v *Vm) Le() error {
	defer func() {
		v.pc += 1 + LE.CountOfOperand()
	}()

	// lh < rh
	lhs := v.program[v.pc+1]
	rhs := v.program[v.pc+2]

	switch lhs.kind {
	case KRegisterTag:
		switch rhs.kind {
		case KRegisterTag:
			// Rlh <= Rrh
			lhsVal := *v.registers[lhs.registerTag]
			rhsVal := *v.registers[rhs.registerTag]

			lessThanOrEq, err := v.le(lhsVal.literal, rhsVal.literal)
			if err != nil {
				return err
			}

			if lessThanOrEq {
				v.zf = 1
			} else {
				v.zf = 0
			}

			return nil
		default:
			return fmt.Errorf("右辺は不正な型です: %s <= %s", lhs.kind.String(), rhs.kind.String())
		}
	default:
		return fmt.Errorf("左辺は不正な型です: %s <= %s", lhs.kind.String(), rhs.kind.String())
	}
}

func (v *Vm) Exit() error {
	v.exited = true
	return nil
}

func (v *Vm) cmp(lhs, rhs Literal) (bool, error) {
	switch lhs.GetKind() {
	case KInt:
		switch rhs.GetKind() {
		case KInt:
			// int == int
			return lhs.GetInt() == rhs.GetInt(), nil
		case KFloat:
			// int == float
			return float64(lhs.GetInt()) == rhs.GetFloat(), nil
		default:
			return false, fmt.Errorf("サポートしていません: %s == %s", lhs.GetKind().String(), rhs.GetKind().String())
		}
	case KFloat:
		switch rhs.GetKind() {
		case KInt:
			// float == int
			return lhs.GetFloat() == float64(rhs.GetInt()), nil
		case KFloat:
			// float == float
			return lhs.GetFloat() == rhs.GetFloat(), nil
		default:
			return false, fmt.Errorf("サポートしていません: %s == %s", lhs.GetKind().String(), rhs.GetKind().String())
		}
	default:
		return false, fmt.Errorf("サポートしていません: %s == %s", lhs.GetKind().String(), rhs.GetKind().String())
	}
}

func (v *Vm) Cmp() error {
	defer func() {
		v.pc += 1 + CMP.CountOfOperand()
	}()
	lhs := v.program[v.pc+1]
	rhs := v.program[v.pc+2]
	switch lhs.kind {
	case KRegisterTag:
		switch rhs.kind {
		case KRegisterTag:
			lhsVal := *v.registers[lhs.registerTag]
			rhsVal := *v.registers[rhs.registerTag]
			eq, err := v.cmp(lhsVal.literal, rhsVal.literal)
			if err != nil {
				return err
			}
			if eq {
				v.zf = 1
			} else {
				v.zf = 0
			}
			return nil
		default:
			return fmt.Errorf("サポートされていません: %s == %s", lhs.kind.String(), rhs.kind.String())
		}
	default:
		return fmt.Errorf("サポートされていません: %s == %s", lhs.kind.String(), rhs.kind.String())
	}
}

func (v *Vm) Call() error {
	newLoc := v.program[v.pc+1]
	switch newLoc.kind {
	case KLabel:
		loc, ok := v.labelLocation[newLoc.label.GetName()]
		if !ok {
			return fmt.Errorf("未定義ラベル: %s", newLoc.label.GetName())
		}
		v._push(*NewLiteralDataWithRaw(v.pc + 2))
		v.pc = loc
		return nil
	default:
		return fmt.Errorf("不正な宛先: %s", newLoc.label.GetName())
	}
}

func (v *Vm) Ret() error {
	newLoc := v._pop()
	if newLoc.kind != KLiteral || newLoc.literal.GetKind() != KInt {
		return fmt.Errorf("戻り先が不正です: pc=%d", newLoc.literal.GetInt())
	}
	v.pc = newLoc.literal.GetInt()
	return nil
}
