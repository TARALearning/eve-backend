package eve

import "testing"

func Test_Command(t *testing.T) {
	cmd := NewEVServiceDefaultCommandFlags()
	if cmd.Name != "help" {
		t.Error("NewEVServiceDefaultCommandFlags does not work as expected")
	}
}
