//go:build windows

package main

import (
	"os"

	"golang.org/x/sys/windows"
)

var signals = []os.Signal{
	os.Interrupt,
	windows.SIGTERM,
}
