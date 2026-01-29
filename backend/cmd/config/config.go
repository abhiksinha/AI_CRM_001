package config

import "CRM/packages/database"

// Config is the top-level configuration struct that holds all application config.
type Config struct {
	App      AppConfig
	Database database.DBConfig
}

// AppConfig holds application-specific configuration.
type AppConfig struct {
	Env         string
	ServiceName string
	Host        string
	Port        string
}
