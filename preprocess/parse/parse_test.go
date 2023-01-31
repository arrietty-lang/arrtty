package parse_test

import (
	"fmt"
	"github.com/arrietty-lang/arrtty/preprocess/parse"
	"github.com/arrietty-lang/arrtty/preprocess/tokenize"
	"testing"
)

func TestNewNode(t *testing.T) {
	n := parse.NewNode(parse.NdVarDecl, nil)
	fmt.Println(n)
}

func TestParseWork(t *testing.T) {
	code := `
	func sayHelloS(name sting) string {
		var hello string = "hello, "
		return hello + name + "!"
	}
	func sayHello(name string) {
		fmt.printf(sayHelloS(name))
	}
	func sub(x int, y int) (int, bool) {
		var isMinus bool
		z := x - y
		if z < 0 {
			isMinus = true
		} else {
			isMinus = false
		}
		return z, isMinus
	}
	`
	head, err := tokenize.Tokenize(code)
	if err != nil {
		t.Fatal(err)
	}
	nodes, err := parse.Parse(head)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(nodes)
}
