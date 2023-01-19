package parse

import (
	"fmt"
	"testing"
)

func TestNewNode(t *testing.T) {
	n := NewNode(NdVarDecl, nil)
	fmt.Println(n)
}
