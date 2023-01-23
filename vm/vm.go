package vm

import "fmt"

type Vm struct {
	pc int
	bp int
	sp int
	zf int

	program   []*Data // 入力全体
	stack     []*Data
	registers []*Data
}

func (v *Vm) Run() {}

func (v *Vm) isProgramEof() bool {
	return v.pc >= len(v.program)
}

func (v *Vm) currentProgram() *Data {
	return v.program[v.pc]
}

func (v *Vm) getArgs(n int) []*Data {
	var args []*Data
	for i := 0; i < n; i++ {
		// なぜ+1しているのかは不明
		// 参考: func operands()
		// https://github.com/x0y14/volume/blob/main/src/vvm/vvm.go
		args = append(args, v.program[v.pc+i+1])
	}
	return args
}

func (v *Vm) addSPSafe(positiveDiff int) error {
	if v.sp+positiveDiff < 0 || len(v.stack)-1 < v.sp+positiveDiff {
		return fmt.Errorf("access error")
	}
	v.sp += positiveDiff
	return nil
}

func (v *Vm) subSPSafe(positiveDiff int) error {
	if v.sp-positiveDiff < 0 || len(v.stack)-1 < v.sp-positiveDiff {
		return fmt.Errorf("access error")
	}
	v.sp -= positiveDiff
	return nil
}

func (v *Vm) execute() {

}
