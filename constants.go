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
	return "localhost"
}

// GetServerPort returns the server port from env or default
func GetServerPort() string {
	if port := os.Getenv("SERVER_PORT"); port != "" {
		return port
	}
	return ":8300"
}

// Response messages in Indonesian
const (
	MsgInvalidID        = "ID tidak valid"
	MsgInvalidRequest   = "Request tidak valid"
	MsgNotFound         = "Data tidak ditemukan"
	MsgDeleteSuccess    = "Data berhasil dihapus"
	MsgAPIRunning       = "API berjalan dengan baik"
	MsgServerFailed     = "Gagal menjalankan server"
	MsgShuttingDown     = "Menerima sinyal shutdown, menghentikan server..."
	MsgShutdownComplete = "Server berhasil dihentikan"
	MsgShutdownTimeout  = "Shutdown timeout, memaksa server berhenti"
)

// GetServerRunningMsg returns the server running message with host and port
func GetServerRunningMsg() string {
	return "Server berjalan di http://" + GetServerHost() + GetServerPort()
}
