package adapt

import (
	"github.com/arrietty-lang/arrtty/vm"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAdapt1(t *testing.T) {
	tests := []struct {
		name        string
		fragments   []vm.Fragment
		howManyLine int
	}{
		{
			"1",
			[]vm.Fragment{
				vm.NewDefLabelFragment("a"),
				vm.NewOpcodeFragment(vm.ADD),
				vm.NewLiteralFragment(vm.NewInt(1)),
				vm.NewLiteralFragment(vm.NewInt(1)),
			},
			2,
		},
		{
			"2",
			[]vm.Fragment{
				vm.NewDefLabelFragment("a"),
				vm.NewOpcodeFragment(vm.ADD),
				vm.NewLiteralFragment(vm.NewInt(1)),
				vm.NewLiteralFragment(vm.NewInt(1)),
				vm.NewOpcodeFragment(vm.SUB),
				vm.NewLiteralFragment(vm.NewInt(1)),
				vm.NewLiteralFragment(vm.NewInt(1)),
				vm.NewOpcodeFragment(vm.EXIT),
			},
			4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lines, err := Adapt(tt.fragments)
			if err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, len(lines), tt.howManyLine)
		})
	}
}
