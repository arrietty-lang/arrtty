package assemble_test

import (
	"fmt"
	"github.com/arrietty-lang/arrtty/assemble"
	"github.com/arrietty-lang/arrtty/preprocess/analyze"
	"github.com/arrietty-lang/arrtty/preprocess/parse"
	"github.com/arrietty-lang/arrtty/preprocess/tokenize"
	"github.com/arrietty-lang/arrtty/vm"
	"github.com/stretchr/testify/assert"
	"log"
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
		//{
		//	"+-",
		//	`
		//		func sub(a int, b int) int { return a-b }
		//		func add(x int, y int) int {
		//			return x+y
		//		}
		//		func main() int {
		//			return sub(add(2, 1), 2)
		//		}`,
		//	1,
		//}, {
		//	"*+",
		//	`
		//		func mul(a int, b int) int { return a*b }
		//		func add(x int, y int) int {
		//			return x+y
		//		}
		//		func main() int {
		//			return mul(3, add(2, 3))
		//		}`,
		//	15,
		//},
		//{
		//	"g",
		//	`
		//		var a int = 1
		//		var x int
		//		func add(x int, y int) int {
		//			return x + y
		//		}
		//		func main() int {
		//			var b int  = a + 1
		//			var c int = add(2, b)
		//			a = c * 2
		//			var d int
		//			d = a
		//			x = 22
		//			return d * x
		//		}
		//		`,
		//	176,
		//},
		//{
		//	"dec",
		//	`
		//		func dec(i int) int {
		//			i = i-1
		//			return i
		//		}
		//		func main() int {
		//			return dec(5)
		//		}
		//		`,
		//	4,
		//},
		//{
		//	"if 1",
		//	`
		//		func isMinus(i int) int {
		//			if i < 0 {
		//				return 1
		//			}
		//			return 0
		//		}
		//		func main() int {
		//			return isMinus(1)
		//		}
		//		`,
		//	0,
		//},
		//{
		//	"if 2",
		//	`
		//		func isMinus(i int) int {
		//			if i < 0 {
		//				return 1
		//			}
		//			return 0
		//		}
		//		func main() int {
		//			return isMinus(-1)
		//		}
		//		`,
		//	1,
		//},
		{
			"f",
			`
				func f(i int) int {
					var x int = i + 1
					var y int = i + 1
					return x + y
					//var x int = i+1
					// var y int = i+1
					i + 1
					return i
				}
				func main() int {
					return f(1)
				}
				`,
			4,
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
		log.Println(tt.name)
		assert.Equal(t, tt.expect, v.ExitCode())
	}
}
