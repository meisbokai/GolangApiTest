package seeders

import (
	"github.com/meisbokai/GolangApiTest/internal/constants"
	"github.com/meisbokai/GolangApiTest/internal/datasources/records"
	"github.com/meisbokai/GolangApiTest/internal/util"
	"github.com/meisbokai/GolangApiTest/pkg/logger"

	"github.com/sirupsen/logrus"
)

var pass string
var UserData []records.Users

func init() {
	var err error
	pass, err = util.GenerateHash("12345")
	if err != nil {
		logger.Panic(err.Error(), logrus.Fields{constants.LoggerCategory: constants.LoggerCategorySeeder})
	}

	UserData = []records.Users{
		{
			Username: "test user 1",
			Email:    "test1@example.com",
			Password: pass,
			RoleId:   1,
		},
		{
			Username: "test user 2",
			Email:    "test2@example.com",
			Password: pass,
			RoleId:   2,
		},
		{
			Username: "test user 3",
			Email:    "test3@example.com",
			Password: pass,
			RoleId:   2,
		},
	}
}
