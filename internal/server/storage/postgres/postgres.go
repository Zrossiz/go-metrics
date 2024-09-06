package postgres

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var PgConn *gorm.DB

func InitConnect(connStr string) error {
	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		return err
	}

	err = Ping(db)
	if err != nil {
		return err
	}

	PgConn = db

	return nil
}

func Ping(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	err = sqlDB.Ping()
	if err != nil {
		return err
	}

	return nil
}
