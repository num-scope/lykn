package model

import "time"

type Project struct {
	ID          uint      `gorm:"primaryKey;comment:Project ID"`
	Name        string    `gorm:"size:255;not null;comment:Project name"`
	Description string    `gorm:"type:text;comment:Project description"`
	PrivateKey  string    `gorm:"type:text;not null;comment:Encrypted private key"`
	PublicKey   string    `gorm:"type:text;not null;comment:Public key"`
	KeyBits     int       `gorm:"not null;default:2048;comment:RSA key size"`
	CreatedAt   time.Time `gorm:"comment:Creation time"`
	UpdatedAt   time.Time `gorm:"comment:Last update time"`
	Licenses    []License `gorm:"foreignKey:ProjectID"`
}

func (Project) TableName() string {
	return "projects"
}
