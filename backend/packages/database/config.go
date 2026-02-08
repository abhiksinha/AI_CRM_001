package database

import "fmt"

// DBConfig holds all configuration required for connecting to the database.
// The `mapstructure` tags ensure correct mapping from the TOML file.
type DBConfig struct {
	Dialect         string `mapstructure:"dialect"`
	Host            string `mapstructure:"host"`
	Port            int    `mapstructure:"port"`
	Username        string `mapstructure:"username"`
	Password        string `mapstructure:"password"`
	SslMode         string `mapstructure:"ssl_mode"`
	Name            string `mapstructure:"name"`
	MaxOpenConns    int    `mapstructure:"maxopenconns"`
	MaxIdleConns    int    `mapstructure:"maxidleconns"`
	ConnMaxLifetime int    `mapstructure:"connmaxlifetime"`
}

// DSN constructs the Data Source Name string for connecting to the database.
func (d DBConfig) DSN() string {
	// Example: "host=localhost user=postgres password=postgres dbname=crm port=5432 sslmode=disable"
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		d.Host, d.Username, d.Password, d.Name, d.Port, d.SslMode)
}
