package db

import (
	"pocket/pkg/auditlog"
	"pocket/pkg/utils"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var environment string
var DB *gorm.DB
var dbFileName string

func connect() {
	var err error
	if environment == "dev" {
		DB, err = gorm.Open(sqlite.Open(utils.Basedir+"/dev.db"), &gorm.Config{})
		if err != nil {
			auditlog.Errorlogger.Error().Str("error", err.Error()).Msg("Error connecting to database")
			panic(err)
		}
	} else {
		DB, err = gorm.Open(sqlite.Open(dbFileName), &gorm.Config{})
		if err != nil {
			auditlog.Errorlogger.Error().Str("error", err.Error()).Msg("Error connecting to database")
			panic(err)
		}
	}
}

func init() {
	environment = "dev"
	connect()
}
