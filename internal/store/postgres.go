package store

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func GetDSN(user, password, host, db string) string {
	return fmt.Sprintf("user=%s password=%s host=%s dbname=%s", user, password, host, db)
}

func GetDB(dsn string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		// @todo log + propogate
		return nil
	}

	return db
}
