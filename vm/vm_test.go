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

func TestVm_Execute_IF(t *testing.T) {
	program := []*Fragment{
		// R1 = 1
		NewOpcodeFragment(MOV),        // 0
		NewLiteralFragment(NewInt(1)), // 1
		NewRegisterFragment(R1),       // 2
		// IF (R1 == 0) ZF = 1
		NewOpcodeFragment(CMP),        // 3
		NewRegisterFragment(R1),       // 4
		NewLiteralFragment(NewInt(0)), // 5
		// JUMP IF ZF = 1
		NewOpcodeFragment(JZ),          // 6
		NewLiteralFragment(NewInt(13)), // 7
		// ジャンプしなかった場合の処理
		NewOpcodeFragment(MOV),                 // 8
		NewLiteralFragment(NewString("hello")), // 9
		NewRegisterFragment(R2),                // 10
		// ジャンプした場合の処理が終わった場所にジャンプ
		NewOpcodeFragment(JMP),         // 11
		NewLiteralFragment(NewInt(16)), // 12
		// ジャンプした場合の処理
		NewOpcodeFragment(MOV),                 // 13
		NewLiteralFragment(NewString("world")), // 14
		NewRegisterFragment(R2),                // 15
		//
		NewOpcodeFragment(EXIT), // 16
	}
	// R1が１の場合
	program[1] = NewLiteralFragment(NewInt(1))
	vm := NewVm(program)
	err := vm.Execute()
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, NewLiteralFragment(NewString("hello")), vm.registers["R2"])

	// R1が0の場合
	program[1] = NewLiteralFragment(NewInt(0))
	vm = NewVm(program)
	err = vm.Execute()
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, NewLiteralFragment(NewString("world")), vm.registers["R2"])
}

func TestVm_Execute_Call(t *testing.T) {
	program := []*Fragment{
		NewOpcodeFragment(CALL),       // 0
		NewLiteralFragment(NewInt(6)), // 1

		NewOpcodeFragment(EXIT), // 2
		// 終了しているのでR2が999になることはない
		NewOpcodeFragment(MOV),          // 3
		NewLiteralFragment(NewInt(999)), // 4
		NewRegisterFragment(R2),         // 5

		NewOpcodeFragment(MOV),          // 6
		NewLiteralFragment(NewInt(100)), // 7
		NewRegisterFragment(R2),         // 8
		NewOpcodeFragment(RET),          // 9
	}

	vm := NewVm(program)
	err := vm.Execute()
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, NewLiteralFragment(NewInt(100)), vm.registers["R2"])
}
