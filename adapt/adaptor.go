package adapt

import "github.com/arrietty-lang/arrtty/vm"

var fragmentPerLine [][]vm.Fragment
var line []vm.Fragment

func addLine() {
	fragmentPerLine = append(fragmentPerLine, line)
	refreshLine()
}

func addFragmentToLine(fragment vm.Fragment) {
	line = append(line, fragment)
}
func refreshLine() {
	line = []vm.Fragment{}
}

func Adapt(fragments []vm.Fragment) ([][]vm.Fragment, error) {
	var pos int
	for pos < len(fragments) {
		curt := fragments[pos]
		switch curt.Kind {
		case vm.LABEL:
			refreshLine()
			addFragmentToLine(curt)
			addLine()
			pos++
		case vm.OPCODE:
			refreshLine()
			addFragmentToLine(curt)
			for i := 0; i < curt.CountOfOperand(); i++ {
				addFragmentToLine(fragments[pos+i+1])
			}
			pos += 1 + curt.CountOfOperand()
			addLine()
		}
	}
	return fragmentPerLine, nil
}
