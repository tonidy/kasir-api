package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

// @title Kasir API
// @version 1.0
// @description REST API untuk mengelola produk dan kategori
// @servers http://localhost:8300

func main() {
	server := &http.Server{
		Addr: GetServerPort(),
	}

	setupRoutes()

	go func() {
		fmt.Println(GetServerRunningMsg())
		if err := startServer(server); err != nil {
			fmt.Println(MsgServerFailed)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println(MsgShuttingDown)

	if err := stopServer(server, ShutdownTimeout); err != nil {
		fmt.Println(MsgShutdownTimeout)
	} else {
		fmt.Println(MsgShutdownComplete)
	}
}
