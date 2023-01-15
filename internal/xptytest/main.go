package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"net"
	"net/rpc"
	"os"
	"strconv"
	"time"

	"github.com/gcrtnst/sw-term-server/internal/xpty"
	"golang.org/x/term"
)

func main() {
	port := flag.Int("port", 0, "")
	flag.Parse()

	if *port == 0 {
		parent()
	} else {
		child(*port)
	}
}

func parent() {
	exe, err := os.Executable()
	if err != nil {
		panic(err)
	}

	lis, err := net.ListenTCP("tcp", &net.TCPAddr{IP: net.IPv4(127, 0, 0, 1)})
	if err != nil {
		panic(err)
	}
	defer lis.Close() // error ignored
	port := lis.Addr().(*net.TCPAddr).Port

	pty, err := xpty.Open()
	if err != nil {
		panic(err)
	}

	errCopy := make(chan error)
	go func() {
		_, err := io.Copy(io.Discard, pty)
		errCopy <- err
	}()

	sess, err := pty.Session(xpty.Size{Row: 30, Col: 120})
	if err != nil {
		_ = pty.Close()
		<-errCopy
		panic(err)
	}

	proc, err := sess.StartProcess(xpty.Cmd{
		Path: exe,
		Args: []string{exe, "-port", strconv.Itoa(port)},
	})
	if err != nil {
		_ = sess.Close()
		_ = pty.Close()
		<-errCopy
		panic(err)
	}

	conn, err := lis.AcceptTCP()
	if err != nil {
		_ = sess.Close()
		_ = pty.Close()
		<-errCopy
		panic(err)
	}

	cli := rpc.NewClient(conn)
	defer cli.Close() // error ignored

	size, err := sess.GetSize()
	if err != nil {
		_ = sess.Close()
		_ = pty.Close()
		<-errCopy
		panic(err)
	}
	if size.Row != 30 || size.Col != 120 {
		_ = sess.Close()
		_ = pty.Close()
		<-errCopy
		panic(fmt.Errorf("%d, %d", size.Row, size.Col))
	}

	size = xpty.Size{}
	err = cli.Call("API.GetSize", struct{}{}, &size)
	if err != nil {
		_ = sess.Close()
		_ = pty.Close()
		<-errCopy
		panic(err)
	}
	if size.Row != 30 || size.Col != 120 {
		_ = sess.Close()
		_ = pty.Close()
		<-errCopy
		panic(fmt.Errorf("%d, %d", size.Row, size.Col))
	}

	err = sess.SetSize(xpty.Size{Row: 40, Col: 80})
	if err != nil {
		_ = sess.Close()
		_ = pty.Close()
		<-errCopy
		panic(err)
	}

	size, err = sess.GetSize()
	if err != nil {
		_ = sess.Close()
		_ = pty.Close()
		<-errCopy
		panic(err)
	}
	if size.Row != 40 || size.Col != 80 {
		_ = sess.Close()
		_ = pty.Close()
		<-errCopy
		panic(fmt.Errorf("%d, %d", size.Row, size.Col))
	}

	size = xpty.Size{}
	err = cli.Call("API.GetSize", struct{}{}, &size)
	if err != nil {
		_ = sess.Close()
		_ = pty.Close()
		<-errCopy
		panic(err)
	}
	if size.Row != 40 || size.Col != 80 {
		_ = sess.Close()
		_ = pty.Close()
		<-errCopy
		panic(fmt.Errorf("%d, %d", size.Row, size.Col))
	}

	err = sess.Close()
	if err != nil {
		_ = pty.Close()
		<-errCopy
		panic(err)
	}

	select {
	case err = <-errCopy:
		_ = pty.Close()
		panic(err)
	default:
	}

	err = pty.Close()
	if err != nil {
		<-errCopy
		panic(err)
	}

	_, err = proc.Wait()
	if err != nil {
		<-errCopy
		panic(err)
	}

	err = <-errCopy
	if !errors.Is(err, os.ErrClosed) {
		panic(err)
	}
}

func child(port int) {
	srv := rpc.NewServer()
	err := srv.Register(API{})
	if err != nil {
		panic(err)
	}

	conn, err := net.DialTCP("tcp", nil, &net.TCPAddr{IP: net.IPv4(127, 0, 0, 1), Port: port})
	if err != nil {
		panic(err)
	}
	srv.ServeConn(conn)

	// expect force termination by SIGHUP
	for {
		time.Sleep(time.Second)
	}
}

type API struct{}

func (API) GetSize(in struct{}, out *xpty.Size) error {
	size, err := getSize()
	*out = size
	return err
}

func getSize() (xpty.Size, error) {
	raw, errRaw := os.Stdout.SyscallConn()
	if errRaw != nil {
		return xpty.Size{}, errRaw
	}

	var width, height int
	var errSize error
	errCtrl := raw.Control(func(fd uintptr) {
		width, height, errSize = term.GetSize(int(fd))
	})
	if errCtrl != nil {
		return xpty.Size{}, errCtrl
	}
	if errSize != nil {
		return xpty.Size{}, errSize
	}

	size := xpty.Size{
		Row: height,
		Col: width,
	}
	return size, nil
}
