// Copyright 2011 The EvAlgo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*

 */

package eve

const (
	// VERSION of the eve library and all its's services
	VERSION = "0.0.2"
)

var (
	// debug is used to set the debug flag for printing out information
	debug = false
	// debugProcesses is used to display all found processes by the schedule
	debugProcesses = false
)

// SetDebug sets the debug value
func SetDebug(value bool) {
	debug = value
}

// SetDebugProcesses sets the debug processes flag
func SetDebugProcesses(value bool) {
	debugProcesses = value
}
