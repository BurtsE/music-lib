package models

import (
	"time"
)

type Song struct {
	ID          uint `gorm:"type:serial;primaryKey;"`
	Name        string
	GroupName   string
	Text        string
	Link        string
	ReleaseDate time.Time
}
