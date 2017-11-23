package eve

import (
	"strings"
	"testing"
)

func Test_EVLogPrintln(t *testing.T) {
	EVLogPrintln("test message")
	if !strings.Contains(EVLogBuffer.String(), "test message") {
		t.Error("EVLogPrintln does not work as expected")
	}
}
