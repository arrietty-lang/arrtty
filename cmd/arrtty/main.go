package main

import (
	"github.com/arrietty-lang/arrtty/assemble"
	"github.com/arrietty-lang/arrtty/preprocess/analyze"
	"github.com/arrietty-lang/arrtty/preprocess/parse"
	"github.com/arrietty-lang/arrtty/preprocess/tokenize"
	"github.com/arrietty-lang/arrtty/vm"
	"log"
	"os"
)

func main() {
	if len(os.Args) <= 1 {
		log.Fatalf("bin <filepath>")
	}

	bytes, err := os.ReadFile(os.Args[1])
	if err != nil {
		log.Fatalf("failed to read: %s", os.Args[1])
	}

	token, err := tokenize.Tokenize(string(bytes))
	if err != nil {
		log.Fatalf("failed to tokenize: %s", err)
	}

	nodes, err := parse.Parse(token)
	if err != nil {
		log.Fatalf("failed to parse: %s", err)
	}

	sem, err := analyze.Analyze(nodes)
	if err != nil {
		log.Fatalf("failed to analyze: %s", err)
	}

	obj, err := assemble.Link([]*assemble.Object{
		{
			"", sem,
		},
	})
	if err != nil {
		log.Fatalf("failed to link: %s", err)
	}

	program, err := assemble.Compile(obj.SemanticsNode)
	if err != nil {
		log.Fatalf("failed to compile: %s", err)
	}

	virtualMachine := vm.NewVm(program, 100)
	err = virtualMachine.Execute()
	if err != nil {
		log.Fatalf("failed to run: %s", err)
	}

	exitCode, err := virtualMachine.ExitCode()
	if err != nil {
		log.Fatalf("failed to get exitCode: %s", err)
	}
	os.Exit(exitCode)
	return
}
