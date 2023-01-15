package main

import (
	"context"
	"io"
	"log"
	"net"
	"net/http"
	"os/signal"
)

const logFlags = log.Ldate | log.Ltime | log.Lmsgprefix

type MainConfig struct {
	Port       int
	TermConfig TermConfig
	LogWriter  io.Writer
}

func Run(cfg MainConfig) int {
	logger := log.New(cfg.LogWriter, "main: ", logFlags)

	ctx, cancel := signal.NotifyContext(context.Background(), signals...)
	defer cancel()

	slot := NewTermSlot(cfg.TermConfig)
	defer slot.Stop()

	addr := &net.TCPAddr{
		IP:   net.IPv4(127, 0, 0, 1),
		Port: cfg.Port,
	}
	lis, err := net.ListenTCP("tcp", addr)
	if err != nil {
		logger.Printf("error: %s", err.Error())
		return 1
	}
	addr = lis.Addr().(*net.TCPAddr)
	logger.Printf("listening on %s", addr.String())

	server := BuildServer(slot, cfg.LogWriter)
	serverDone := make(chan error)
	go func() {
		err := server.Serve(lis)
		serverDone <- err
		close(serverDone)
	}()
	<-ctx.Done()

	code := 0
	err = server.Shutdown(context.Background())
	if err != nil {
		logger.Printf("error: %s", err.Error())
		code = 1
	}
	err = <-serverDone
	if err != http.ErrServerClosed {
		logger.Printf("error: %s", err.Error())
		code = 1
	}

	return code
}

func BuildServeMux(slot *TermSlot, logw io.Writer) *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/keyboard", &ServiceHandler{
		Service: &KeyboardService{
			TermSlot: slot,
			Logger:   log.New(logw, "keyboard: ", logFlags),
		},
	})
	mux.Handle("/screen", &ServiceHandler{
		Service: &ScreenService{
			TermSlot: slot,
			Logger:   log.New(logw, "screen: ", logFlags),
		},
	})
	mux.Handle("/stop", &ServiceHandler{
		Service: &StopService{
			TermSlot: slot,
		},
	})
	return mux
}

func BuildServer(slot *TermSlot, logw io.Writer) *http.Server {
	return &http.Server{
		Handler:  BuildServeMux(slot, logw),
		ErrorLog: log.New(logw, "server: ", logFlags),
	}
}
