package vm

type Fragment struct {
	*Opcode
	*Literal
	*Address
}

func NewOpcodeFragment(opcode Opcode) *Fragment {
	return &Fragment{
		Opcode:  &opcode,
		Literal: nil,
	}
}

func NewLiteralFragment(literal *Literal) *Fragment {
	return &Fragment{
		Opcode:  nil,
		Literal: literal,
	}
}

func NewAddressFragment(address *Address) *Fragment {
	return &Fragment{
		Opcode:  nil,
		Literal: nil,
		Address: address,
	}
}
