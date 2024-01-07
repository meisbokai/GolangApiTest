package seeders

import (
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/meisbokai/GolangApiTest/internal/constants"
	"github.com/meisbokai/GolangApiTest/internal/datasources/records"
	"github.com/meisbokai/GolangApiTest/pkg/logger"
	"github.com/sirupsen/logrus"
)

type Seeder interface {
	UserSeeder(userData []records.Users) (err error)
}

type seeder struct {
	db *sqlx.DB
}

func NewSeeder(db *sqlx.DB) Seeder {
	return &seeder{db: db}
}

func (s *seeder) UserSeeder(userData []records.Users) (err error) {
	if len(userData) == 0 {
		return errors.New("users data is empty")
	}

	logger.Info("inserting users data...", logrus.Fields{constants.LoggerCategory: constants.LoggerCategorySeeder})
	for _, user := range userData {
		user.CreatedAt = time.Now().In(time.FixedZone("GMT+8", 8*60*60))
		if _, err = s.db.NamedQuery(`INSERT INTO users(id, username, email, password, role_id, created_at) VALUES (uuid_generate_v4(), :username, :email, :password, :role_id, :created_at)`, user); err != nil {
			return err
		}
	}
	logger.Info("users data inserted successfully", logrus.Fields{constants.LoggerCategory: constants.LoggerCategorySeeder})

	return
}
