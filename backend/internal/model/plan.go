package model

import "time"

type Plan struct {
	ID          uint      `gorm:"primaryKey;comment:Plan ID"`
	Code        string    `gorm:"size:100;uniqueIndex;not null;comment:Plan code"`
	Name        string    `gorm:"size:255;not null;comment:Plan name"`
	Description string    `gorm:"type:text;comment:Plan description"`
	MaxUsers    int       `gorm:"not null;default:0;comment:Max user count, 0 means unlimited"`
	MaxDevices  int       `gorm:"not null;default:1;comment:Max device count"`
	Enabled     bool      `gorm:"not null;default:true;comment:Whether plan is enabled"`
	Features    []Feature `gorm:"many2many:plan_features;constraint:OnDelete:CASCADE;"`
	CreatedAt   time.Time `gorm:"comment:Creation time"`
	UpdatedAt   time.Time `gorm:"comment:Last update time"`
}

func (Plan) TableName() string {
	return "plans"
}

type PlanFeature struct {
	PlanID    uint      `gorm:"primaryKey;comment:Plan ID"`
	FeatureID uint      `gorm:"primaryKey;comment:Feature ID"`
	CreatedAt time.Time `gorm:"comment:Creation time"`
}

func (PlanFeature) TableName() string {
	return "plan_features"
}
