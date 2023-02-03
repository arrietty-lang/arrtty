package vm

type Pointer int

const (
	SP Pointer = iota
	BP
)

func (p *Pointer) String() string {
	switch *p {
	case SP:
		return "sp"
	case BP:
		return "bp"
	}
	return ""
}
