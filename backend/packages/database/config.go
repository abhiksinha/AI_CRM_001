package database

// Config holds the configuration required for connecting to the database.
type Config struct {
	DSN string // Data Source Name, e.g., "host=localhost user=user password=pass dbname=crm port=5432"
}
