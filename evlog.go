package eve

import (
	"bytes"
	"log"
	"os"
)

var (
	EVLogBuffer = bytes.NewBuffer(nil)
	EVLogger    = log.New(EVLogBuffer, "", log.LstdFlags)
)

func EVLogPrintln(v ...interface{}) {
	EVLogger.Println(v...)
}

func EVLogFatal(v ...interface{}) {
	EVLogPrintln(v...)
	os.Exit(1)
}
