package model

import "time"

type Feature struct {
	ID          uint      `gorm:"primaryKey;comment:Feature ID"`
	Code        string    `gorm:"size:100;uniqueIndex;not null;comment:Feature code"`
	Name        string    `gorm:"size:255;not null;comment:Feature name"`
	Description string    `gorm:"type:text;comment:Feature description"`
	Enabled     bool      `gorm:"not null;default:true;comment:Whether feature is enabled"`
	CreatedAt   time.Time `gorm:"comment:Creation time"`
	UpdatedAt   time.Time `gorm:"comment:Last update time"`
}

func (Feature) TableName() string {
	return "features"
}
