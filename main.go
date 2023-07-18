package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/gcrtnst/sw-term-server/internal/xpty"
)

func main() {
	port := flag.Int("port", 0, "listen port")
	row := flag.Int("row", 27, "terminal rows")
	col := flag.Int("col", 58, "terminal columns")
	shell := flag.String("shell", defaultShell(), "shell")
	flag.Parse()

	if *row <= 0 {
		fmt.Fprintln(os.Stderr, "invalid row")
		os.Exit(1)
	}
	if *col <= 0 {
		fmt.Fprintln(os.Stderr, "invalid col")
		os.Exit(1)
	}
	if *shell == "" {
		fmt.Fprintln(os.Stderr, "shell not specified")
		os.Exit(1)
	}

	cfg := MainConfig{
		Port: *port,
		TermConfig: TermConfig{
			Open: xpty.Open,
			Row:  *row,
			Col:  *col,
			Cmd: xpty.Cmd{
				Path: *shell,
				Args: []string{*shell},
			},
		},
		LogWriter: os.Stdout,
	}
	code := Run(cfg)
	os.Exit(code)
}

func defaultShell() string {
	if runtime.GOOS == "windows" {
		comspec := os.Getenv("COMSPEC")
		if comspec != "" {
			return comspec
		}

		systemroot := os.Getenv("SYSTEMROOT")
		if systemroot != "" {
			return systemroot + `\system32\cmd.exe`
		}

		path, err := exec.LookPath("cmd.exe")
		if err == nil {
			return path
		}

		return ""
	}

	shell := os.Getenv("SHELL")
	if shell != "" {
		return shell
	}

	path, err := exec.LookPath("sh")
	if err == nil {
		return path
	}

	return ""
}
