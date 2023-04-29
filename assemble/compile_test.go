package assemble

import (
	"github.com/arrietty-lang/arrtty/preprocess/analyze"
	"github.com/arrietty-lang/arrtty/preprocess/parse"
	"github.com/arrietty-lang/arrtty/preprocess/tokenize"
	"github.com/arrietty-lang/arrtty/vm"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCompile_Call(t *testing.T) {
	tests := []struct {
		name   string
		code   string
		expect int
	}{
		{
			"+-",
			`
func add(a int, b int) int {
	return a+b
}
func sub(a int, b int) int {
	return a-b
}
func main() int {
	return sub(add(2, 4), 3)
}
		`,
			3,
		},
		{
			"f",
			`
func f(n int) int {
	var x int = n + 1
	return x + n
}
func main() int {
	return f(f(1))
}
				`,
			7,
		},
		{
			"1",
			`
func f(n int) int {
	if n < 10 {
		return f(n+1)
	}
	return n
}

func main() int {
	return f(0)
}
				`,
			10,
		},
		{
			"fib",
			`
func fib(n int) int {
	if n < 2 {
		return n
	}
	return fib(n-1) + fib(n-2)
}
func main() int {
	return fib(10)
}
				`,
			55,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := tokenize.Tokenize(tt.code)
			if err != nil {
				t.Fatal(err)
			}

			nodes, err := parse.Parse(token)
			if err != nil {
				t.Fatal(err)
			}

			sem, err := analyze.Analyze(nodes)
			if err != nil {
				t.Fatal(err)
			}

			obj, err := Link([]*Object{
				{
					Identifier:    "",
					SemanticsNode: sem,
				},
			})
			if err != nil {
				t.Fatal(err)
			}

			program, err := Compile(obj.SemanticsNode)
			if err != nil {
				t.Fatal(err)
			}

			virtualMachine := vm.NewVm(program, 100)
			err = virtualMachine.Execute()
			if err != nil {
				t.Fatal(err)
			}
			ec, err := virtualMachine.ExitCode()
			if err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, tt.expect, ec)
		})
	}
}
