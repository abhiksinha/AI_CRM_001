package database

import "fmt"

// DBConfig holds all configuration required for connecting to the database.
type DBConfig struct {
	Dialect         string
	Host            string
	Port            int
	Username        string
	Password        string
	SslMode         string
	Name            string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime int
}

// DSN constructs the Data Source Name string for connecting to the database.
func (d DBConfig) DSN() string {
	// Example: "host=localhost user=postgres password=postgres dbname=crm port=5432 sslmode=disable"
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		d.Host, d.Username, d.Password, d.Name, d.Port, d.SslMode)
}
