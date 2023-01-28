package vm

type Variable struct {
	Name string
}

func NewVariable(name string) *Variable {
	return &Variable{Name: name}
}
