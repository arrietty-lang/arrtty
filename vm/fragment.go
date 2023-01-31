package vm

type FragmentKind int

const (
	_ FragmentKind = iota
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

func NewOpcodeFragment(opcode Opcode) *Fragment {
	return &Fragment{
		Kind:    OPCODE,
		Opcode:  &opcode,
		Literal: nil,
	}
}

func NewLiteralFragment(literal *Literal) *Fragment {
	return &Fragment{
		Kind:    LITERAL,
		Opcode:  nil,
		Literal: literal,
	}
}

func NewAddressFragment(address *Address) *Fragment {
	return &Fragment{
		Kind:    ADDRESS,
		Opcode:  nil,
		Literal: nil,
		Address: address,
	}
}

func NewRegisterFragment(reg Register) *Fragment {
	return &Fragment{
		Kind:     REGISTER,
		Opcode:   nil,
		Literal:  nil,
		Address:  nil,
		Register: &reg,
	}
}

func NewVariableFragment(v *Variable) *Fragment {
	return &Fragment{
		Kind:     VARIABLE,
		Variable: v,
	}
}

func NewDefLabelFragment(id string) *Fragment {
	return &Fragment{Kind: LABEL,
		Label: &Label{
			Id:     id,
			Define: true,
		}}
}
func NewLabelFragment(id string) *Fragment {
	return &Fragment{
		Kind:  LABEL,
		Label: &Label{Id: id, Define: false},
	}
}
