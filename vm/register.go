package vm

type Register int

const (
	R1 Register = iota
	R2
	R3
	R4
	R5
)

var regs = [...]string{
	R1: "R1",
	R2: "R2",
	R3: "R3",
	R4: "R4",
	R5: "R5",
}

func (r Register) String() string {
	return regs[r]
}
