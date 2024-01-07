package util

import (
	"time"

	"github.com/jmoiron/sqlx"
	config "github.com/meisbokai/GolangApiTest/internal/configs"
	"github.com/meisbokai/GolangApiTest/internal/datasources/drivers"
)

func SetupPostgresConnection() (*sqlx.DB, error) {
	// Setup sqlx config of postgreSQL
	config := drivers.SQLXConfig{
		DriverName:     config.AppConfig.DBPostgreDriver,
		DataSourceName: config.AppConfig.DBPostgreDsn,
		MaxOpenConns:   100,
		MaxIdleConns:   10,
		MaxLifetime:    15 * time.Minute,
	}

	// Initialize postgreSQL connection with sqlx
	conn, err := config.InitializeSQLXDatabase()
	if err != nil {
		return nil, err
	}

	return conn, nil
}

// To view the database
// docker exec -it -u postgres some-postgres bash
// psql -U postgres
// \c postgres
// select * from users;
