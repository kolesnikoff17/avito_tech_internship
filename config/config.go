package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type (
	// Config is a struct with ENV variables
	Config struct {
		HTTP
		PG
		Logger
	}
	// HTTP -.
	HTTP struct {
		Port            string
		ReadTimeout     time.Duration
		WriteTimeout    time.Duration
		ShutdownTimeout time.Duration
	}
	// PG -.
	PG struct {
		Host    string
		Port    string
		User    string
		Pwd     string
		Name    string
		MaxConn int
	}
	// Logger -.
	Logger struct {
		Level string
	}
)

// NewConfig gets values from ENV
func NewConfig() *Config {
	cfg := &Config{}
	cfg.HTTP.Port = os.Getenv("PORT")
	cfg.HTTP.ReadTimeout, _ = time.ParseDuration(os.Getenv("SERVER_READ_TIMEOUT") + "s")
	cfg.HTTP.WriteTimeout, _ = time.ParseDuration(os.Getenv("SERVER_WRITE_TIMEOUT") + "s")
	cfg.HTTP.ShutdownTimeout, _ = time.ParseDuration(os.Getenv("SERVER_SHUTDOWN_TIMEOUT") + "s")
	cfg.PG.Port = os.Getenv("DB_PORT")
	cfg.PG.Host = os.Getenv("DB_HOST")
	cfg.PG.User = os.Getenv("DB_USER")
	cfg.PG.Pwd = os.Getenv("DB_PWD")
	cfg.PG.Name = os.Getenv("DB_NAME")
	cfg.PG.MaxConn, _ = strconv.Atoi(os.Getenv("DB_MAXCONNS"))
	cfg.Logger.Level = os.Getenv("LOG_LVL")
	return cfg
}

// DbParams formats connection string from config
func DbParams(cfg *Config) string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable connect_timeout=5",
		cfg.PG.Host,
		cfg.PG.Port,
		cfg.PG.User,
		cfg.PG.Pwd,
		cfg.PG.Name)
}
