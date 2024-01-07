package drivers

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// SQLXConfig holds the configuration for the database instance
type SQLXConfig struct {
	DriverName     string
	DataSourceName string
	MaxOpenConns   int
	MaxIdleConns   int
	MaxLifetime    time.Duration
}

// InitializeSQLXDatabase returns a new DBInstance
func (config *SQLXConfig) InitializeSQLXDatabase() (*sqlx.DB, error) {
	db, err := sqlx.Open(config.DriverName, config.DataSourceName)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}

	// set maximum number of open connections to database
	db.SetMaxOpenConns(config.MaxOpenConns)

	// set maximum number of idle connections in the pool
	db.SetMaxIdleConns(config.MaxIdleConns)

	// set maximum time to wait for new connection
	db.SetConnMaxLifetime(config.MaxLifetime)

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("error pinging database: %v", err)
	}

	return db, nil
}
