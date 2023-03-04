package main

import (
	"flag"
	"fmt"
	"net"
	"net/rpc"
	"os"
	"strconv"

	"github.com/gcrtnst/sw-term-server/internal/xpty"
	"golang.org/x/term"
)

func main() {
	port := flag.Int("port", 0, "")
	flag.Parse()

	if *port != 0 {
		child(*port)
	} else {
		parent()
	}
}

func parent() {
	exe, err := os.Executable()
	if err != nil {
		panic(err)
	}

	lis, err := net.ListenTCP("tcp", &net.TCPAddr{
		IP: net.IPv4(127, 0, 0, 1),
	})
	if err != nil {
		panic(err)
	}
	defer func() {
		err := lis.Close()
		if err != nil {
			panic(err)
		}
	}()
	port := lis.Addr().(*net.TCPAddr).Port

	pty, err := xpty.Open()
	if err != nil {
		panic(err)
	}
	defer func() {
		if pty == nil {
			return
		}

		err := pty.Close()
		if err != nil {
			panic(err)
		}
	}()

	proc, err := pty.Start(&xpty.Cmd{
		Path: exe,
		Args: []string{exe, "-port", strconv.Itoa(port)},
	})
	if err != nil {
		panic(err)
	}
	defer func() {
		if proc == nil {
			return
		}

		err := proc.Release()
		if err != nil {
			panic(err)
		}
	}()

	conn, err := lis.AcceptTCP()
	if err != nil {
		panic(err)
	}

	c := rpc.NewClient(conn)
	defer func() {
		if c == nil {
			return
		}

		err := c.Close()
		if err != nil {
			panic(err)
		}
	}()

	row, col, err := pty.GetSize()
	if err != nil {
		panic(err)
	}
	if row != 0 || col != 0 {
		panic(fmt.Errorf("row=%d, col=%d", row, col))
	}

	size := &Size{}
	err = c.Call("ServerAPI.GetSize", struct{}{}, &size)
	if err != nil {
		panic(err)
	}
	if size.Row != 0 || size.Col != 0 {
		panic(fmt.Errorf("%#v", size))
	}

	err = pty.SetSize(120, 30)
	if err != nil {
		panic(err)
	}

	row, col, err = pty.GetSize()
	if err != nil {
		panic(err)
	}
	if row != 120 || col != 30 {
		panic(fmt.Errorf("row=%d, col=%d", row, col))
	}

	size = &Size{}
	err = c.Call("ServerAPI.GetSize", struct{}{}, &size)
	if err != nil {
		panic(err)
	}
	if size.Row != 120 || size.Col != 30 {
		panic(fmt.Errorf("%#v", size))
	}
	err = c.Close()
	c = nil
	if err != nil {
		panic(err)
	}

	err = pty.Close()
	pty = nil
	if err != nil {
		panic(err)
	}

	_, err = proc.Wait()
	proc = nil
	if err != nil {
		panic(err)
	}
}

func child(port int) {
	s := rpc.NewServer()
	err := s.Register(&ServerAPI{})
	if err != nil {
		panic(err)
	}

	conn, err := net.DialTCP("tcp", nil, &net.TCPAddr{
		IP:   net.IPv4(127, 0, 0, 1),
		Port: port,
	})
	if err != nil {
		panic(err)
	}
	s.ServeConn(conn)
}

type ServerAPI struct{}

func (api *ServerAPI) GetSize(arg struct{}, ret *Size) error {
	col, row, err := term.GetSize(0)
	if err != nil {
		panic(err)
	}

	ret.Row = row
	ret.Col = col
	return nil
}

type Size struct {
	Row, Col int
}
