package vm

import "fmt"

type Pointer int

const (
	SP Pointer = iota
	BP
)

func (p Pointer) String() string {
	switch p {
	case SP:
		return "sp"
	case BP:
		return "bp"
	default:
		return "illegal"
	}
}

type Offset struct {
	pointer  Pointer
	relation int
}

func NewOffset(pointer Pointer, relation int) *Offset {
	return &Offset{
		pointer:  pointer,
		relation: relation,
	}
}

func (o *Offset) String() string {
	var s string
	if o.relation == 0 {
		s = "0"
	} else if o.relation < 0 {
		s = fmt.Sprintf("%d", o.relation)
	} else if o.relation > 0 {
		s = fmt.Sprintf("+%d", o.relation)
	}
	return fmt.Sprintf("Offset{ pointer: %s, relation: %s }", o.pointer.String(), s)
}
func (o *Offset) AddressString() string {
	var s string
	if o.relation == 0 {
		s = ""
	} else if o.relation < 0 {
		s = fmt.Sprintf("%d", o.relation)
	} else if o.relation > 0 {
		s = fmt.Sprintf("+%d", o.relation)
	}
	return fmt.Sprintf("[%s%s]", o.pointer.String(), s)
}

func (o *Offset) GetPointer() Pointer {
	return o.pointer
}
func (o *Offset) SetPointer(p Pointer) {
	o.pointer = p
}

func (o *Offset) GetRelation() int {
	return o.relation
}
func (o *Offset) SetRelation(r int) {
	o.relation = r
}
