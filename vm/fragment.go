package vm

type FragmentKind int

const (
	ILLEGAL FragmentKind = iota
	OPCODE
	LITERAL
	ADDRESS
	REGISTER
	POINTER
	VARIABLE
	LABEL
)

type Fragment struct {
	Kind FragmentKind
	*Opcode
	*Literal
	*Address
	*Register
	*Pointer
	*Variable
	*Label
}

func (f Fragment) String() string {
	switch f.Kind {
	case OPCODE:
		return f.Opcode.String()
	case LITERAL:
		return f.Literal.String()
	case ADDRESS:
		return f.Address.String()
	case REGISTER:
		return f.Register.String()
	case POINTER:
		return f.Pointer.String()
	case VARIABLE:
		return f.Variable.String()
	case LABEL:
		return f.Label.String()
	}
	return ""
}

func NewOpcodeFragment(opcode Opcode) Fragment {
	return Fragment{
		Kind:    OPCODE,
		Opcode:  &opcode,
		Literal: nil,
	}
}

func NewLiteralFragment(literal *Literal) Fragment {
	return Fragment{
		Kind:    LITERAL,
		Opcode:  nil,
		Literal: literal,
	}
}

func NewAddressFragment(address *Address) Fragment {
	return Fragment{
		Kind:    ADDRESS,
		Opcode:  nil,
		Literal: nil,
		Address: address,
	}
}

func NewRegisterFragment(reg Register) Fragment {
	return Fragment{
		Kind:     REGISTER,
		Opcode:   nil,
		Literal:  nil,
		Address:  nil,
		Register: &reg,
	}
}

func NewVariableFragment(v *Variable) Fragment {
	return Fragment{
		Kind:     VARIABLE,
		Variable: v,
	}
}

func NewDefLabelFragment(id string) Fragment {
	return Fragment{Kind: LABEL,
		Label: &Label{
			Id:     id,
			Define: true,
		}}
}

func NewLabelFragment(id string) Fragment {
	return Fragment{
		Kind:  LABEL,
		Label: &Label{Id: id, Define: false},
	}
}

func NewPointerFragment(pointer Pointer) Fragment {
	return Fragment{
		Kind:    POINTER,
		Pointer: &pointer,
	}
}
