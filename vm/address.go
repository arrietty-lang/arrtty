package vm

type Address struct {
	Base     int
	Relative int
}

func NewAddress(base, relative int) *Address {
	return &Address{
		Base:     base,
		Relative: relative,
	}
}
