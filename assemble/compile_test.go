package assemble_test

import (
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
	func main() {
		//return 1
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
	err = v.Execute()
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 0, v.ExitCode())
}
