package vm

import "testing"

func TestVm_Execute(t *testing.T) {
	program := []*Fragment{
		NewOpcodeFragment(CMP),
		NewLiteralFragment(NewInt(1)),
		NewLiteralFragment(NewInt(2)),
		NewOpcodeFragment(NOP),
		NewOpcodeFragment(JMP),
		NewAddressFragment(NewAddress(50, 100)),
	}
	vm := NewVm(program)
	vm.Execute()
}
