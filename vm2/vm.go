package vm2

import "github.com/arrietty-lang/arrtty/vm"

type Vm struct {
	stack []vm.Fragment
	lines [][]vm.Fragment
}

func NewVm(lines [][]vm.Fragment) *Vm {
	v := Vm{
		stack: []vm.Fragment{},
		lines: lines,
	}
	return &v
}

func (v *Vm) labelScan() error {
	return nil
}

func (v *Vm) Execute() error {
	return nil
}
