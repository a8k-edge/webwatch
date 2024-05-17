package main

import (
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

type Target struct {
	gorm.Model
	ID       uint `gorm:"primaryKey;autoIncrement"`
	URL      string
	Name     string `gorm:"not null;default:''"`
	IsActive bool   `gorm:"default:true"`

	History []History `gorm:"foreignKey:TargetID"`
}

type History struct {
	ID        uint `gorm:"primaryKey;autoIncrement"`
	TargetID  uint
	Hash      string
	CreatedAt time.Time
}

func initializeDB(dbPath string) {
	var err error
	db, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&Target{})
	db.AutoMigrate(&History{})
}
