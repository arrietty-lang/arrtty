package vm

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestVm_Execute_Math1(t *testing.T) {
	program := []*Fragment{
		// r1に15を代入
		NewOpcodeFragment(MOV),
		NewLiteralFragment(NewInt(15)),
		NewRegisterFragment(R1),
		// r2に15を代入
		NewOpcodeFragment(MOV),
		NewLiteralFragment(NewInt(15)),
		NewRegisterFragment(R2),
		// r1にr2を足す(r1 == 30)
		NewOpcodeFragment(ADD),
		NewRegisterFragment(R2),
		NewRegisterFragment(R1),

		// r2に20を代入(r2 == 20)
		NewOpcodeFragment(MOV),
		NewLiteralFragment(NewInt(20)),
		NewRegisterFragment(R2),
		// r1からr2を引く(30-20)
		NewOpcodeFragment(SUB),
		NewRegisterFragment(R2),
		NewRegisterFragment(R1),
	}
	vm := NewVm(program)
	err := vm.Execute()
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, NewLiteralFragment(NewInt(10)), vm.registers["R1"])
}
