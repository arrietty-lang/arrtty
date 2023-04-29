package vm

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLiteral_String(t *testing.T) {
	l := NewLiteral("hello")
	assert.Equal(t, "hello", l.GetString())
	l.SetString("world")
	assert.Equal(t, "world", l.GetString())
	assert.Equal(t, KString, l.GetKind())
}
