//go:build !linux && !windows

package main

import "os"

var signals = []os.Signal{
	os.Interrupt,
}
