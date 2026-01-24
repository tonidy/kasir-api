package main

import (
	"os"
	"time"
)

// Server configuration
const (
	ContentTypeJSON = "application/json"
	InitialIDOffset = 1
	ShutdownTimeout = 10 * time.Second
)

// GetServerHost returns the server host from env or default
func GetServerHost() string {
	if host := os.Getenv("SERVER_HOST"); host != "" {
		return host
	}
	// Default to 0.0.0.0 if PORT env var exists (indicates cloud/container environment)
	if os.Getenv("PORT") != "" {
		return "0.0.0.0"
	}
	return "localhost"
}

// GetServerPort returns the server port from env or default
func GetServerPort() string {
	// Fallback to SERVER_PORT
	if port := os.Getenv("SERVER_PORT"); port != "" {
		return port
	}
	return ":8300"
}

// Response messages
const (
	MsgInvalidID        = "Invalid ID"
	MsgInvalidRequest   = "Invalid request"
	MsgNotFound         = "Data not found"
	MsgDeleteSuccess    = "Data successfully deleted"
	MsgAPIRunning       = "API is running 🔥"
	MsgServerFailed     = "Failed to start server"
	MsgShuttingDown     = "Received shutdown signal, stopping server..."
	MsgShutdownComplete = "Server successfully stopped"
	MsgShutdownTimeout  = "Shutdown timeout, forcing server to stop"
)

// GetServerRunningMsg returns the server running message with host and port
func GetServerRunningMsg() string {
	host := GetServerHost()
	// In production, bind to 0.0.0.0 but show localhost in message
	if host == "0.0.0.0" {
		host = "localhost"
	}
	return "Server running at http://" + host + GetServerPort()
}
