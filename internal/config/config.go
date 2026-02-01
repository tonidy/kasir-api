package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/knadh/koanf/parsers/dotenv"
	"github.com/knadh/koanf/providers/env/v2"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
}

type ServerConfig struct {
	Host         string
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
	MaxConns int
	MinConns int
}

func Load() (*Config, error) {
	k := koanf.New(".")

	// Load from .env file if it exists
	if _, err := os.Stat(".env"); err == nil {
		if err := k.Load(file.Provider(".env"), dotenv.ParserEnv("APP_", ".", func(s string) string {
			return strings.Replace(strings.ToLower(
				strings.TrimPrefix(s, "APP_")), "_", ".", -1)
		})); err != nil {
			return nil, fmt.Errorf("failed to load .env file: %w", err)
		}
	}

	// Load from environment variables (overrides .env)
	if err := k.Load(env.Provider("", env.Opt{
		Prefix: "APP_",
		TransformFunc: func(key, value string) (string, any) {
			key = strings.ToLower(strings.TrimPrefix(key, "APP_"))
			key = strings.ReplaceAll(key, "_", ".")
			return key, value
		},
	}), nil); err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	cfg := &Config{
		Server: ServerConfig{
			Host:         k.String("server.host"),
			Port:         k.String("server.port"),
			ReadTimeout:  k.Duration("server.readtimeout"),
			WriteTimeout: k.Duration("server.writetimeout"),
		},
		Database: DatabaseConfig{
			Host:     k.String("database.host"),
			Port:     k.Int("database.port"),
			User:     k.String("database.user"),
			Password: k.String("database.password"),
			DBName:   k.String("database.dbname"),
			SSLMode:  k.String("database.sslmode"),
			MaxConns: k.Int("database.maxconns"),
			MinConns: k.Int("database.minconns"),
		},
	}

	setDefaults(cfg)
	return cfg, nil
}

func setDefaults(cfg *Config) {
	if cfg.Server.Host == "" {
		cfg.Server.Host = "localhost"
	}
	if cfg.Server.Port == "" {
		cfg.Server.Port = ":8300"
	}
	if cfg.Server.ReadTimeout == 0 {
		cfg.Server.ReadTimeout = 10 * time.Second
	}
	if cfg.Server.WriteTimeout == 0 {
		cfg.Server.WriteTimeout = 10 * time.Second
	}
	if cfg.Database.Port == 0 {
		cfg.Database.Port = 5432
	}
	if cfg.Database.SSLMode == "" {
		cfg.Database.SSLMode = "require"
	}
	if cfg.Database.MaxConns == 0 {
		cfg.Database.MaxConns = 25
	}
	if cfg.Database.MinConns == 0 {
		cfg.Database.MinConns = 5
	}
}
