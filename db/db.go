package db

import (
	"webwatch/config"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func Init() {
	var err error
	db, err = gorm.Open(sqlite.Open(config.Cfg.Database.DBPath), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&Target{})
	db.AutoMigrate(&History{})
}

func GetDB() *gorm.DB {
	return db
}
