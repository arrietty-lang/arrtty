package vm

import "fmt"

type RegisterTag int

const (
	R1 RegisterTag = iota
	R2
	R3
	R10
	R11

	RSP
	RBP
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
	case R10:
		s = "R10"
	case R11:
		s = "R11"
	case RSP:
		s = "RSP"
	case RBP:
		s = "RBP"
	default:
		s = "illegal"
	}
	return fmt.Sprintf("RegisterTag{ %s }", s)
}
