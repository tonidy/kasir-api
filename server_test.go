package main

import (
	"net/http"
	"testing"
	"time"
)

func TestStopServer(t *testing.T) {
	server := &http.Server{
		Addr: ":9999",
	}

	go func() {
		startServer(server)
	}()

	err := stopServer(server, 1*time.Second)
	if err != nil {
		t.Errorf("Expected successful shutdown, got error: %v", err)
	}
}

func TestStopServer_Timeout(t *testing.T) {
	server := &http.Server{
		Addr: ":9998",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(3 * time.Second)
		}),
	}

	go func() {
		startServer(server)
	}()

	time.Sleep(100 * time.Millisecond)

	go func() {
		http.Get("http://localhost:9998")
	}()

	time.Sleep(100 * time.Millisecond)

	err := stopServer(server, 500*time.Millisecond)
	if err == nil {
		t.Error("Expected timeout error, got nil")
	}
}
