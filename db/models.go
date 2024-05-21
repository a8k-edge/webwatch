package db

import (
	"time"

	"gorm.io/gorm"
)

type Target struct {
	gorm.Model
	ID       uint `gorm:"primaryKey;autoIncrement"`
	URL      string
	Name     string `gorm:"not null;default:''"`
	IsActive bool   `gorm:"default:true"`

	History []History `gorm:"foreignKey:TargetID"`
}

type History struct {
	ID         uint `gorm:"primaryKey;autoIncrement"`
	TargetID   uint
	Hash       string
	IsChanged  bool   `gorm:"default:false"`
	StatusCode int    `gorm:"default:0"`
	Diff       string `gorm:"default:''"`

	CreatedAt time.Time
}
