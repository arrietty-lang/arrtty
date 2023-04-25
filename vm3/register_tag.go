package vm3

import "fmt"

type RegisterTag int

const (
	R1 RegisterTag = iota
	R2
	R3
)

func (r RegisterTag) String() string {
	var s string
	switch r {
	case R1:
		s = "R1"
	case R2:
		s = "R2"
	case R3:
		s = "R3"
	default:
		s = "illegal"
	}
	return fmt.Sprintf("RegisterTag{ %s }", s)
}
