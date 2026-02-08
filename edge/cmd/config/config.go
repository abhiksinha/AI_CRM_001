package config

import "edge/packages/redis"

// Config is the top-level configuration struct that holds all application config.
type Config struct {
	App             AppConfig
	Redis           redis.Config
	BackendServices map[string]string
}

// AppConfig holds application-specific configuration.
type AppConfig struct {
	Env         string
	ServiceName string
	Host        string
	Port        string
}
