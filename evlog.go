package eve

import (
	"bytes"
	"log"
	"os"
)

var (
	// EVLogBuffer is the buffer where all msg will be stored
	EVLogBuffer = bytes.NewBuffer(nil)
	// EVLogger logs into the buffer
	EVLogger = log.New(EVLogBuffer, "", log.LstdFlags)
)

// EVLogPrintln writes the messages into the buffer with the logger
func EVLogPrintln(v ...interface{}) {
	EVLogger.Println(v...)
}

// EVLogFatal writes the message into the buffer and quits the execution with an error 1 message
func EVLogFatal(v ...interface{}) {
	EVLogPrintln(v...)
	os.Exit(1)
}
