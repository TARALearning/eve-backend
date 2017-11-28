package eve

import (
	"errors"
	"fmt"
	"strings"
	"testing"
)

func Test_EVLogPrintln(t *testing.T) {
	EVLogPrintln("test message")
	if !strings.Contains(EVLogBuffer.String(), "test message") {
		t.Error("EVLogPrintln does not work as expected")
	}
}

func Test_EVLogFatal(t *testing.T) {
	err := errors.New("testerror")
	defer func() {
		if r := recover(); r != nil {
			var ok bool
			err, ok = r.(error)
			if !ok {
				err = fmt.Errorf("pkg: %v", r)
			}
			if !strings.Contains(err.Error(), "testerror") {
				t.Error("EVLogFatal does not work as expected")
			}
		}
	}()
	EVLogFatal(err)
}
