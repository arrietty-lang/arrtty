package vm2

import "github.com/arrietty-lang/arrtty/vm"

func (v *Vm) Nop() error {
	return nil
}

func (v *Vm) Add(target, value vm.Fragment) error {
	return nil
}

func (v *Vm) Sub(target, value vm.Fragment) error {
	return nil
}

func (v *Vm) Mul(target, value vm.Fragment) error {
	return nil
}

func (v *Vm) Div(target, value vm.Fragment) error {
	return nil
}

func (v *Vm) Mov(from, to vm.Fragment) error {
	return nil
}

func (v *Vm) Cmp(x, y vm.Fragment) error {
	return nil
}

func (v *Vm) Lt(small, big vm.Fragment) error {
	return nil
}

func (v *Vm) Le(small, big vm.Fragment) error {
	return nil
}

func (v *Vm) Gt(big, small vm.Fragment) error {
	return nil
}

func (v *Vm) Ge(big, small vm.Fragment) error {
	return nil
}

func (v *Vm) Jmp(to vm.Fragment) error {
	return nil
}

func (v *Vm) Jz(to vm.Fragment) error {
	// jump if zero(flag == 1)
	return nil
}

func (v *Vm) Push(value vm.Fragment) error {
	return nil
}

func (v *Vm) Pop(to vm.Fragment) error {
	return nil
}

func (v *Vm) Call(fn vm.Fragment) error {
	return nil
}

func (v *Vm) Ret() error {
	return nil
}

func (v *Vm) Exit() error {
	return nil
}
