package vm

type STDKind int

const (
	STDIN STDKind = iota
	STDOUT
	STDERR
)

type SystemCall int

const (
	READ SystemCall = iota
	WRITE
)
