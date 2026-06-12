package dao

import (
	"context"
	"errors"
	"time"

	"github.com/wu-clan/lykn/backend/internal/common"
	"github.com/wu-clan/lykn/backend/internal/model"
	"gorm.io/gorm"
)

type LicenseDAO struct{ db *gorm.DB }

func NewLicenseDAO(db *gorm.DB) *LicenseDAO { return &LicenseDAO{db: db} }

func (d *LicenseDAO) Create(ctx context.Context, license *model.License) error {
	return d.db.WithContext(ctx).Create(license).Error
}

func (d *LicenseDAO) GetByID(ctx context.Context, id uint) (*model.License, error) {
	var license model.License
	if err := d.db.WithContext(ctx).First(&license, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.ErrNotFound
		}
		return nil, err
	}
	return &license, nil
}

func (d *LicenseDAO) ListByProjectID(ctx context.Context, projectID uint) ([]model.License, error) {
	var licenses []model.License
	return licenses, d.db.WithContext(ctx).Where("project_id = ?", projectID).Order("id desc").Find(&licenses).Error
}

func (d *LicenseDAO) CountByProjectID(ctx context.Context, projectID uint) (int64, error) {
	var count int64
	if err := d.db.WithContext(ctx).Model(&model.License{}).Where("project_id = ?", projectID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (d *LicenseDAO) Count(ctx context.Context) (int64, error) {
	var count int64
	if err := d.db.WithContext(ctx).Model(&model.License{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (d *LicenseDAO) CountActive(ctx context.Context, now time.Time) (int64, error) {
	var count int64
	if err := d.db.WithContext(ctx).Model(&model.License{}).
		Where("not_before <= ? AND not_after >= ?", now, now).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (d *LicenseDAO) CountExpired(ctx context.Context, now time.Time) (int64, error) {
	var count int64
	if err := d.db.WithContext(ctx).Model(&model.License{}).Where("not_after < ?", now).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
