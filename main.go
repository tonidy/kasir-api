package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// @title Kasir API
// @version 1.0
// @description REST API for managing products and categories
// @servers http://localhost:8300

func main() {
	server := &http.Server{
		Addr: GetServerHost() + GetServerPort(),
	}

	setupRoutes()
	fmt.Println("Routes configured successfully")

	// Channel to capture server startup errors
	errChan := make(chan error, 1)

	go func() {
		fmt.Printf("Starting server on %s\n", server.Addr)
		fmt.Println(GetServerRunningMsg())
		if err := startServer(server); err != nil {
			errChan <- err
		}
	}()

	// Give server a moment to start listening
	time.Sleep(100 * time.Millisecond)

	// Wait for either shutdown signal or server error
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errChan:
		fmt.Println(MsgServerFailed, err)
		os.Exit(1)
	case sig := <-quit:
		fmt.Printf("Received signal: %v\n", sig)
		fmt.Println(MsgShuttingDown)
	}

	if err := stopServer(server, ShutdownTimeout); err != nil {
		fmt.Println(MsgShutdownTimeout)
	} else {
		fmt.Println(MsgShutdownComplete)
	}
}
