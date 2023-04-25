package vm3

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestVm_Push(t *testing.T) {
	stackSize := 10
	push := []Data{
		*NewOpcodeData(PUSH), *NewLiteralDataWithRaw(1),
		*NewOpcodeData(PUSH), *NewLiteralDataWithRaw(1.2),
		*NewOpcodeData(PUSH), *NewLiteralDataWithRaw("1.23"),
	}
	virtualMachine := NewVm(push, stackSize)
	err := virtualMachine.Execute()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, NewLiteralDataWithRaw(1), virtualMachine.stack[stackSize-2])
	assert.Equal(t, NewLiteralDataWithRaw(1.2), virtualMachine.stack[stackSize-3])
	assert.Equal(t, NewLiteralDataWithRaw("1.23"), virtualMachine.stack[stackSize-4])
}

func TestVm_Pop(t *testing.T) {
	stackSize := 10
	push := []Data{
		*NewOpcodeData(PUSH), *NewLiteralDataWithRaw(1),
		*NewOpcodeData(PUSH), *NewLiteralDataWithRaw(1.2),
		*NewOpcodeData(POP), *NewRegisterTagData(R1),
		*NewOpcodeData(POP), *NewRegisterTagData(R2),
	}
	virtualMachine := NewVm(push, stackSize)
	err := virtualMachine.Execute()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, NewLiteralDataWithRaw(1.2), virtualMachine.registers[R1])
	assert.Equal(t, NewLiteralDataWithRaw(1), virtualMachine.registers[R2])

	nonNilStacks := 0
	for i := 0; i < stackSize; i++ {
		if virtualMachine.stack[i] != nil {
			nonNilStacks++
		}
	}
	assert.Equal(t, 0, nonNilStacks)
}

func TestVm_Add(t *testing.T) {
	stackSize := 10
	add := []Data{
		*NewOpcodeData(PUSH), *NewLiteralDataWithRaw(10), // R1
		*NewOpcodeData(PUSH), *NewLiteralDataWithRaw(33), // R2
		*NewOpcodeData(POP), *NewRegisterTagData(R2), // stack[8] -> R2
		*NewOpcodeData(POP), *NewRegisterTagData(R1), // stack[7] -> R1
		*NewOpcodeData(ADD), *NewRegisterTagData(R2), *NewRegisterTagData(R1), // R1 += R2
	}

	virtualMachine := NewVm(add, stackSize)
	err := virtualMachine.Execute()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, NewLiteralDataWithRaw(43), virtualMachine.registers[R1])
	assert.Equal(t, NewLiteralDataWithRaw(33), virtualMachine.registers[R2])
	nonNilStacks := 0
	for i := 0; i < stackSize; i++ {
		if virtualMachine.stack[i] != nil {
			nonNilStacks++
		}
	}
	assert.Equal(t, 0, nonNilStacks)
}

func TestVm_Add2(t *testing.T) {
	stackSize := 10
	add := []Data{
		*NewOpcodeData(PUSH), *NewLiteralDataWithRaw(33), // R1
		*NewOpcodeData(POP), *NewRegisterTagData(R1), // stack[8] -> R1
		*NewOpcodeData(ADD), *NewRegisterTagData(R1), *NewRegisterTagData(R1), // R1 += R1
	}

	virtualMachine := NewVm(add, stackSize)
	err := virtualMachine.Execute()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, NewLiteralDataWithRaw(66), virtualMachine.registers[R1])
	nonNilStacks := 0
	for i := 0; i < stackSize; i++ {
		if virtualMachine.stack[i] != nil {
			nonNilStacks++
		}
	}
	assert.Equal(t, 0, nonNilStacks)
}

func TestVm_Sub(t *testing.T) {
	stackSize := 10
	add := []Data{
		*NewOpcodeData(PUSH), *NewLiteralDataWithRaw(33), // R1
		*NewOpcodeData(PUSH), *NewLiteralDataWithRaw(10), // R2
		*NewOpcodeData(POP), *NewRegisterTagData(R2), // stack[8] -> R2
		*NewOpcodeData(POP), *NewRegisterTagData(R1), // stack[7] -> R1
		*NewOpcodeData(SUB), *NewRegisterTagData(R2), *NewRegisterTagData(R1), // R1 -= R2
	}

	virtualMachine := NewVm(add, stackSize)
	err := virtualMachine.Execute()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, NewLiteralDataWithRaw(23), virtualMachine.registers[R1])
	assert.Equal(t, NewLiteralDataWithRaw(10), virtualMachine.registers[R2])
	nonNilStacks := 0
	for i := 0; i < stackSize; i++ {
		if virtualMachine.stack[i] != nil {
			nonNilStacks++
		}
	}
	assert.Equal(t, 0, nonNilStacks)
}

func TestVm_Sub2(t *testing.T) {
	stackSize := 10
	add := []Data{
		*NewOpcodeData(PUSH), *NewLiteralDataWithRaw(33), // R1
		*NewOpcodeData(POP), *NewRegisterTagData(R1), // stack[8] -> R1
		*NewOpcodeData(SUB), *NewRegisterTagData(R1), *NewRegisterTagData(R1), // R1 -= R1
	}

	virtualMachine := NewVm(add, stackSize)
	err := virtualMachine.Execute()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, NewLiteralDataWithRaw(0), virtualMachine.registers[R1])
	nonNilStacks := 0
	for i := 0; i < stackSize; i++ {
		if virtualMachine.stack[i] != nil {
			nonNilStacks++
		}
	}
	assert.Equal(t, 0, nonNilStacks)
}

func TestVm_Mov(t *testing.T) {
	stackSize := 10
	mov := []Data{
		*NewOpcodeData(PUSH), *NewLiteralDataWithRaw(9), // R1
		*NewOpcodeData(POP), *NewRegisterTagData(R1), // stack[8] -> R1
		*NewOpcodeData(MOV), *NewRegisterTagData(R1), *NewRegisterTagData(R2), // R2 = R1

		*NewOpcodeData(PUSH), *NewLiteralDataWithRaw(2), // R1
		*NewOpcodeData(POP), *NewRegisterTagData(R1), // stack[8] -> R1
		*NewOpcodeData(ADD), *NewRegisterTagData(R1), *NewRegisterTagData(R2), // R2 += R1
	}

	virtualMachine := NewVm(mov, stackSize)
	err := virtualMachine.Execute()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, NewLiteralDataWithRaw(2), virtualMachine.registers[R1])
	assert.Equal(t, NewLiteralDataWithRaw(11), virtualMachine.registers[R2])

	nonNilStacks := 0
	for i := 0; i < stackSize; i++ {
		if virtualMachine.stack[i] != nil {
			nonNilStacks++
		}
	}
	assert.Equal(t, 0, nonNilStacks)
}

func TestVm_Mov2(t *testing.T) {
	stackSize := 10
	mov := []Data{
		*NewOpcodeData(PUSH), *NewLiteralDataWithRaw(2), // R3
		*NewOpcodeData(POP), *NewRegisterTagData(R3), // stack[8] -> R3

		*NewOpcodeData(PUSH), *NewLiteralDataWithRaw(3), // R1
		*NewOpcodeData(POP), *NewRegisterTagData(R1), // stack[8] -> R1

		*NewOpcodeData(MOV), *NewRegisterTagData(R1), *NewRegisterTagData(R2), // R2 = R1
		*NewOpcodeData(ADD), *NewRegisterTagData(R3), *NewRegisterTagData(R1), // R1 += R3
		*NewOpcodeData(ADD), *NewRegisterTagData(R1), *NewRegisterTagData(R3), // // R3 += R1
	}

	virtualMachine := NewVm(mov, stackSize)
	err := virtualMachine.Execute()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, NewLiteralDataWithRaw(5), virtualMachine.registers[R1])
	assert.Equal(t, NewLiteralDataWithRaw(3), virtualMachine.registers[R2])
	assert.Equal(t, NewLiteralDataWithRaw(7), virtualMachine.registers[R3])

	nonNilStacks := 0
	for i := 0; i < stackSize; i++ {
		if virtualMachine.stack[i] != nil {
			nonNilStacks++
		}
	}
	assert.Equal(t, 0, nonNilStacks)
}

func TestVm_Lt(t *testing.T) {
	stackSize := 10
	lt := []Data{
		*NewOpcodeData(PUSH), *NewLiteralDataWithRaw(2), // R1
		*NewOpcodeData(POP), *NewRegisterTagData(R1), // stack[8] -> R1
		*NewOpcodeData(PUSH), *NewLiteralDataWithRaw(3), // R2
		*NewOpcodeData(POP), *NewRegisterTagData(R2), // stack[8] -> R2
		*NewOpcodeData(LT), *NewRegisterTagData(R1), *NewRegisterTagData(R2), // R1 < R2
	}

	virtualMachine := NewVm(lt, stackSize)
	err := virtualMachine.Execute()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 1, virtualMachine.zf)

	nonNilStacks := 0
	for i := 0; i < stackSize; i++ {
		if virtualMachine.stack[i] != nil {
			nonNilStacks++
		}
	}
	assert.Equal(t, 0, nonNilStacks)
}

func TestVm_Lt2(t *testing.T) {
	stackSize := 10
	lt := []Data{
		*NewOpcodeData(PUSH), *NewLiteralDataWithRaw(3), // R1
		*NewOpcodeData(POP), *NewRegisterTagData(R1), // stack[8] -> R1
		*NewOpcodeData(PUSH), *NewLiteralDataWithRaw(3), // R2
		*NewOpcodeData(POP), *NewRegisterTagData(R2), // stack[8] -> R2
		*NewOpcodeData(LT), *NewRegisterTagData(R1), *NewRegisterTagData(R2), // R1 < R2
	}

	virtualMachine := NewVm(lt, stackSize)
	err := virtualMachine.Execute()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 0, virtualMachine.zf)

	nonNilStacks := 0
	for i := 0; i < stackSize; i++ {
		if virtualMachine.stack[i] != nil {
			nonNilStacks++
		}
	}
	assert.Equal(t, 0, nonNilStacks)
}

func TestVm_Le(t *testing.T) {
	stackSize := 10
	lt := []Data{
		*NewOpcodeData(PUSH), *NewLiteralDataWithRaw(2), // R1
		*NewOpcodeData(POP), *NewRegisterTagData(R1), // stack[8] -> R1
		*NewOpcodeData(PUSH), *NewLiteralDataWithRaw(3), // R2
		*NewOpcodeData(POP), *NewRegisterTagData(R2), // stack[8] -> R2
		*NewOpcodeData(LE), *NewRegisterTagData(R1), *NewRegisterTagData(R2), // R1 <= R2
	}

	virtualMachine := NewVm(lt, stackSize)
	err := virtualMachine.Execute()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 1, virtualMachine.zf)

	nonNilStacks := 0
	for i := 0; i < stackSize; i++ {
		if virtualMachine.stack[i] != nil {
			nonNilStacks++
		}
	}
	assert.Equal(t, 0, nonNilStacks)
}

func TestVm_Le2(t *testing.T) {
	stackSize := 10
	lt := []Data{
		*NewOpcodeData(PUSH), *NewLiteralDataWithRaw(3), // R1
		*NewOpcodeData(POP), *NewRegisterTagData(R1), // stack[8] -> R1
		*NewOpcodeData(PUSH), *NewLiteralDataWithRaw(3), // R2
		*NewOpcodeData(POP), *NewRegisterTagData(R2), // stack[8] -> R2
		*NewOpcodeData(LE), *NewRegisterTagData(R1), *NewRegisterTagData(R2), // R1 <= R2
	}

	virtualMachine := NewVm(lt, stackSize)
	err := virtualMachine.Execute()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 1, virtualMachine.zf)

	nonNilStacks := 0
	for i := 0; i < stackSize; i++ {
		if virtualMachine.stack[i] != nil {
			nonNilStacks++
		}
	}
	assert.Equal(t, 0, nonNilStacks)
}

func TestVm_Jmp(t *testing.T) {
	stackSize := 10
	jmp := []Data{
		*NewOpcodeData(PUSH), *NewLiteralDataWithRaw(3), // R1
		*NewOpcodeData(POP), *NewRegisterTagData(R1), // stack[8] -> R1
		*NewOpcodeData(JMP), *NewLabelData(*NewLabel(false, "afterExit")),
		*NewOpcodeData(EXIT),
		*NewLabelData(*NewLabel(true, "afterExit")),
		*NewOpcodeData(PUSH), *NewLiteralDataWithRaw(10), // R1
		*NewOpcodeData(POP), *NewRegisterTagData(R1), // stack[8] -> R1
	}

	virtualMachine := NewVm(jmp, stackSize)
	err := virtualMachine.Execute()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, NewLiteralDataWithRaw(10), virtualMachine.registers[R1])

	nonNilStacks := 0
	for i := 0; i < stackSize; i++ {
		if virtualMachine.stack[i] != nil {
			nonNilStacks++
		}
	}
	assert.Equal(t, 0, nonNilStacks)
}

func TestVm_Jmp2(t *testing.T) {
	stackSize := 10
	jmp := []Data{
		*NewOpcodeData(PUSH), *NewLiteralDataWithRaw(1), // R1
		*NewOpcodeData(POP), *NewRegisterTagData(R1), // stack[8] -> R1
		*NewOpcodeData(PUSH), *NewLiteralDataWithRaw(3), // R2
		*NewOpcodeData(POP), *NewRegisterTagData(R2), // stack[8] -> R2
		*NewOpcodeData(LT), *NewRegisterTagData(R1), *NewRegisterTagData(R2), // R1 < R2

		*NewOpcodeData(JZ), *NewLabelData(*NewLabel(false, "afterExit")),
		*NewOpcodeData(EXIT),
		*NewLabelData(*NewLabel(true, "afterExit")),

		*NewOpcodeData(PUSH), *NewLiteralDataWithRaw(10), // R1
		*NewOpcodeData(POP), *NewRegisterTagData(R1), // stack[8] -> R1
	}

	virtualMachine := NewVm(jmp, stackSize)
	err := virtualMachine.Execute()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, NewLiteralDataWithRaw(10), virtualMachine.registers[R1])

	nonNilStacks := 0
	for i := 0; i < stackSize; i++ {
		if virtualMachine.stack[i] != nil {
			nonNilStacks++
		}
	}
	assert.Equal(t, 0, nonNilStacks)
}

func TestVm_Jmp3(t *testing.T) {
	stackSize := 10
	jmp := []Data{
		*NewOpcodeData(PUSH), *NewLiteralDataWithRaw(3), // R1
		*NewOpcodeData(POP), *NewRegisterTagData(R1), // stack[8] -> R1
		*NewOpcodeData(PUSH), *NewLiteralDataWithRaw(1), // R2
		*NewOpcodeData(POP), *NewRegisterTagData(R2), // stack[8] -> R2
		*NewOpcodeData(LT), *NewRegisterTagData(R1), *NewRegisterTagData(R2), // R1 < R2

		*NewOpcodeData(JZ), *NewLabelData(*NewLabel(false, "afterExit")),
		*NewOpcodeData(EXIT),
		*NewLabelData(*NewLabel(true, "afterExit")),

		*NewOpcodeData(PUSH), *NewLiteralDataWithRaw(10), // R1
		*NewOpcodeData(POP), *NewRegisterTagData(R1), // stack[8] -> R1
	}

	virtualMachine := NewVm(jmp, stackSize)
	err := virtualMachine.Execute()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, NewLiteralDataWithRaw(3), virtualMachine.registers[R1])

	nonNilStacks := 0
	for i := 0; i < stackSize; i++ {
		if virtualMachine.stack[i] != nil {
			nonNilStacks++
		}
	}
	assert.Equal(t, 0, nonNilStacks)
}

func TestVm_Cmp(t *testing.T) {
	stackSize := 10
	jmp := []Data{
		*NewOpcodeData(PUSH), *NewLiteralDataWithRaw(3), // R1
		*NewOpcodeData(POP), *NewRegisterTagData(R1), // stack[8] -> R1
		*NewOpcodeData(PUSH), *NewLiteralDataWithRaw(1), // R2
		*NewOpcodeData(POP), *NewRegisterTagData(R2), // stack[8] -> R2

		*NewOpcodeData(CMP), *NewRegisterTagData(R1), *NewRegisterTagData(R2),
	}

	virtualMachine := NewVm(jmp, stackSize)
	err := virtualMachine.Execute()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 0, virtualMachine.zf)

	nonNilStacks := 0
	for i := 0; i < stackSize; i++ {
		if virtualMachine.stack[i] != nil {
			nonNilStacks++
		}
	}
	assert.Equal(t, 0, nonNilStacks)
}

func TestVm_Cmp2(t *testing.T) {
	stackSize := 10
	jmp := []Data{
		*NewOpcodeData(PUSH), *NewLiteralDataWithRaw(1), // R1
		*NewOpcodeData(POP), *NewRegisterTagData(R1), // stack[8] -> R1
		*NewOpcodeData(PUSH), *NewLiteralDataWithRaw(1), // R2
		*NewOpcodeData(POP), *NewRegisterTagData(R2), // stack[8] -> R2

		*NewOpcodeData(CMP), *NewRegisterTagData(R1), *NewRegisterTagData(R2),
	}

	virtualMachine := NewVm(jmp, stackSize)
	err := virtualMachine.Execute()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 1, virtualMachine.zf)

	nonNilStacks := 0
	for i := 0; i < stackSize; i++ {
		if virtualMachine.stack[i] != nil {
			nonNilStacks++
		}
	}
	assert.Equal(t, 0, nonNilStacks)
}
