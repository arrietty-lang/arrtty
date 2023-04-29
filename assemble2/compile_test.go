package assemble2

import (
	"github.com/arrietty-lang/arrtty/assemble"
	"github.com/arrietty-lang/arrtty/preprocess/analyze"
	"github.com/arrietty-lang/arrtty/preprocess/parse"
	"github.com/arrietty-lang/arrtty/preprocess/tokenize"
	"github.com/arrietty-lang/arrtty/vm3"
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

			obj, err := assemble.Link([]*assemble.Object{
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

			virtualMachine := vm3.NewVm(program, 100)
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
