package vm

type Address struct {
	Original Pointer
	Relative int
}

func NewAddress(original Pointer, relative int) *Address {
	return &Address{
		Original: original,
		Relative: relative,
	}
}
