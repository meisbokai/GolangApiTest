package util

import (
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/meisbokai/GolangApiTest/internal/datasources/drivers"
)

func SetupPostgresConnection() (*sqlx.DB, error) {
	// Setup sqlx config of postgreSQL
	config := drivers.SQLXConfig{
		DriverName:     "postgres",
		DataSourceName: "user=postgres password=golangapitest host=localhost port=5432 dbname=postgres sslmode=disable timezone=Asia/Singapore",
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
