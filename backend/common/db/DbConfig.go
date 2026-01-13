package db

// DbConfig holds the parameters for connecting to the database.
type DbConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}
