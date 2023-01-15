//go:build linux

package main

import (
	"os"

	"golang.org/x/sys/unix"
)

var signals = []os.Signal{
	unix.SIGHUP,
	unix.SIGINT,
	unix.SIGTERM,
}
