package vm3

import "fmt"

type Label struct {
	define bool
	name   string
}

func NewLabel(define bool, name string) *Label {
	return &Label{
		define: define,
		name:   name,
	}
}

func (l *Label) String() string {
	return fmt.Sprintf("Label{ isDefine: %t, name: %s}", l.GetIsDefine(), l.GetName())
}

func (l *Label) GetName() string {
	return l.name
}
func (l *Label) SetName(name string) {
	l.name = name
}

func (l *Label) GetIsDefine() bool {
	return l.define
}
func (l *Label) SetIsDefine(b bool) {
	l.define = b
}
