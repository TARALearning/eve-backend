package eve

import "testing"

func Test_Command(t *testing.T) {
	cmd := NewEVServiceDefaultCommandFlags()
	t.Log(cmd)
}
