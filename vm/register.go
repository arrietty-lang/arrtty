package vm

type Register int

const (
	R1 Register = iota
	R2
	R3
	R4
	R5

	R10 // 関数の戻り値にも使える
	R11

	EC // ループなどのカウンタ
	ED // 読み込むデータサイズ | 書き込むデータサイズ?
	EM // 書き込み先の読み込み先の種類(STDIN, STDOUT, STDERR, FILE)
	EP // FILE PATH
	EW // 書き込むデータの格納
	ER // 読み込んだデータの格納
)

var regs = [...]string{
	R1:  "R1",
	R2:  "R2",
	R3:  "R3",
	R4:  "R4",
	R5:  "R5",
	R10: "R10",
	R11: "R11",
	EC:  "EC",
	ED:  "ED",
	EM:  "EM",
	EP:  "EP",
	EW:  "EW",
	ER:  "ER",
}

func (r Register) String() string {
	return regs[r]
}
