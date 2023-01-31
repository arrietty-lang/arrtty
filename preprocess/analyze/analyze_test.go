package analyze_test

import (
	"github.com/arrietty-lang/arrtty/preprocess/analyze"
	"github.com/arrietty-lang/arrtty/preprocess/parse"
	"github.com/arrietty-lang/arrtty/preprocess/tokenize"
	"testing"
)

func TestAnalyze_1(t *testing.T) {
	code := `
	func sayHelloS(name string) string {
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
		} else if 1 == 2 {
			return z, true
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
	if _, err = analyze.Analyze(nodes); err != nil {
		t.Fatal(err)
	}
}
