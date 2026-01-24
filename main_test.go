package main

import (
	"fmt"
	"net/http"
	"os"
	"syscall"
	"testing"
	"time"
)

func TestMain_Integration(t *testing.T) {
	server := &http.Server{
		Addr: ":8081",
	}

	setupRoutes()

	errChan := make(chan error, 1)
	go func() {
		if err := startServer(server); err != nil {
			errChan <- err
		}
	}()

	time.Sleep(100 * time.Millisecond)

	resp, err := http.Get("http://localhost:8081/health")
	if err != nil {
		t.Fatalf("Failed to connect to server: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	if err := stopServer(server, ShutdownTimeout); err != nil {
		t.Errorf("Failed to stop server: %v", err)
	}

	select {
	case err := <-errChan:
		if err != http.ErrServerClosed {
			t.Errorf("Unexpected server error: %v", err)
		}
	default:
	}
}

func TestMain_SignalHandling(t *testing.T) {
	http.DefaultServeMux = new(http.ServeMux)

	server := &http.Server{
		Addr: ":8082",
	}

	setupRoutes()

	quit := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	go func() {
		startServer(server)
	}()

	time.Sleep(100 * time.Millisecond)

	go func() {
		<-quit
		fmt.Println(MsgShuttingDown)
		if err := stopServer(server, ShutdownTimeout); err != nil {
			fmt.Println(MsgShutdownTimeout)
		} else {
			fmt.Println(MsgShutdownComplete)
		}
		done <- true
	}()

	quit <- syscall.SIGTERM

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Error("Shutdown took too long")
	}
}
