package vm

type Variable struct {
	Name string
}

func (v *Variable) String() string {
	return v.Name
}

func NewVariable(name string) *Variable {
	return &Variable{Name: name}
}
