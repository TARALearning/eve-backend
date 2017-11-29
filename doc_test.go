package eve

import "testing"

func Test_SetDebug(t *testing.T) {
	SetDebug(true)
	if !debug {
		t.Error("SetDebug does not work as expected")
	}
}

func Test_SetDebugProcesses(t *testing.T) {
	SetDebugProcesses(true)
	if !debugProcesses {
		t.Error("SetDebugProcesses does not work as expected")
	}
}
