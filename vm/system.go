package vm

type IODestination int

const (
	STDIN IODestination = iota
	STDOUT
	STDERR
	FILE
)

type SystemCall int

const (
	READ SystemCall = iota
	WRITE
)
