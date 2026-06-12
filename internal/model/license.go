package model

import (
	"time"

	"gorm.io/datatypes"
)

type License struct {
	ID           uint           `gorm:"primaryKey;comment:License ID"`
	UUID         string         `gorm:"size:100;uniqueIndex;not null;comment:License UUID"`
	ProjectID    uint           `gorm:"index;not null;comment:Project ID"`
	SubjectName  string         `gorm:"size:255;not null;comment:Licensed subject name"`
	SubjectEmail string         `gorm:"size:255;comment:Licensed subject email"`
	SubjectOrg   string         `gorm:"size:255;comment:Licensed subject organization"`
	PlanID       *uint          `gorm:"index;comment:Plan ID snapshot source"`
	PlanName     string         `gorm:"size:255;comment:License plan name snapshot"`
	Plan         string         `gorm:"size:100;comment:License plan"`
	NotBefore    time.Time      `gorm:"not null;comment:Valid from time"`
	NotAfter     time.Time      `gorm:"not null;comment:Valid until time"`
	Hardware     datatypes.JSON `gorm:"type:json;comment:Hardware binding JSON"`
	Features     datatypes.JSON `gorm:"type:json;not null;comment:Licensed features JSON"`
	Limits       datatypes.JSON `gorm:"type:json;comment:License limits JSON"`
	Metadata     datatypes.JSON `gorm:"type:json;not null;comment:License metadata JSON"`
	LicContent   string         `gorm:"type:text;not null;comment:Signed license content"`
	CreatedAt    time.Time      `gorm:"comment:Creation time"`
	Project      Project        `gorm:"foreignKey:ProjectID"`
}

func (License) TableName() string {
	return "licenses"
}
