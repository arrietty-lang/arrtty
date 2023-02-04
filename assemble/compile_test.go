package assemble_test

import (
	"fmt"
	"github.com/arrietty-lang/arrtty/assemble"
	"github.com/arrietty-lang/arrtty/preprocess/analyze"
	"github.com/arrietty-lang/arrtty/preprocess/parse"
	"github.com/arrietty-lang/arrtty/preprocess/tokenize"
	"github.com/arrietty-lang/arrtty/vm"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCompile_1(t *testing.T) {
	code := `
	func main() int {
		return 1000
	}
	`
	token, err := tokenize.Tokenize(code)
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
	fragments, err := assemble.Compile(obj.SemanticsNode)
	if err != nil {
		t.Fatal(err)
	}

	v := vm.NewVm(fragments)
	fmt.Println(v.Export())
	err = v.Execute()
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 1000, v.ExitCode())
}

func TestCompile_2(t *testing.T) {
	code := `
	func main() int {
		return 1 + 10
	}
	`
	token, err := tokenize.Tokenize(code)
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
	fragments, err := assemble.Compile(obj.SemanticsNode)
	if err != nil {
		t.Fatal(err)
	}

	v := vm.NewVm(fragments)
	fmt.Println(v.Export())
	err = v.Execute()
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 11, v.ExitCode())
}

func TestCompile_3(t *testing.T) {
	code := `
	func main() int {
		return 1 - 10
	}
	`
	token, err := tokenize.Tokenize(code)
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
	fragments, err := assemble.Compile(obj.SemanticsNode)
	if err != nil {
		t.Fatal(err)
	}

	v := vm.NewVm(fragments)
	fmt.Println(v.Export())
	err = v.Execute()
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, -9, v.ExitCode())
}

func TestCompile_4(t *testing.T) {
	code := `
	func main() int {
		return 1 * 10
	}
	`
	token, err := tokenize.Tokenize(code)
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
	fragments, err := assemble.Compile(obj.SemanticsNode)
	if err != nil {
		t.Fatal(err)
	}

	v := vm.NewVm(fragments)
	fmt.Println(v.Export())
	err = v.Execute()
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 10, v.ExitCode())
}
func TestCompile_Math(t *testing.T) {
	tests := []struct {
		name   string
		code   string
		expect int
	}{
		{
			"*+",
			"func main() int { return 1 + ( 3 - 2 ) * 4 }",
			5,
		},
	}

	for _, tt := range tests {
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
		fragments, err := assemble.Compile(obj.SemanticsNode)
		if err != nil {
			t.Fatal(err)
		}

		v := vm.NewVm(fragments)
		fmt.Println(v.Export())
		err = v.Execute()
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, tt.expect, v.ExitCode())
	}
}

func TestCompile_CALL(t *testing.T) {
	tests := []struct {
		name   string
		code   string
		expect int
	}{
		{
			"+-",
			`
				func sub(a int, b int) int { return a-b }
				func add(x int, y int) int {
					return x+y
				}
				func main() int {
					return sub(add(2, 1), 2)
				}`,
			1,
		}, {
			"*+",
			`
				func mul(a int, b int) int { return a*b }
				func add(x int, y int) int {
					return x+y
				}
				func main() int {
					return mul(3, add(2, 3))
				}`,
			15,
		},
		{
			"g",
			`
				func main() int {
					var a int = 1
					var b int  = a + 1
					return b
				}
				`,
			2,
		},
	}

	for _, tt := range tests {
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
		fragments, err := assemble.Compile(obj.SemanticsNode)
		if err != nil {
			t.Fatal(err)
		}

		v := vm.NewVm(fragments)
		fmt.Println(v.Export())
		err = v.Execute()
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, tt.expect, v.ExitCode())
	}
}
