package model

import "time"

type User struct {
	ID        uint      `gorm:"primaryKey;comment:User ID"`
	Username  string    `gorm:"size:100;uniqueIndex;not null;comment:Login username"`
	Password  string    `gorm:"not null;comment:Password hash"`
	CreatedAt time.Time `gorm:"comment:Creation time"`
	UpdatedAt time.Time `gorm:"comment:Last update time"`
}

func (User) TableName() string {
	return "users"
}
