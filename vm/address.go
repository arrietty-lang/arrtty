package vm

import "fmt"

type Address struct {
	Original Pointer
	Relative int
}

func (a *Address) String() string {
	var rel string
	if a.Relative >= 0 {
		rel = fmt.Sprintf("+%d", a.Relative)
	} else {
		rel = fmt.Sprintf("%d", a.Relative)
	}
	return fmt.Sprintf("[%s%s]", a.Original.String(), rel)
}

func NewAddress(original Pointer, relative int) *Address {
	return &Address{
		Original: original,
		Relative: relative,
	}
}
