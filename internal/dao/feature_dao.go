package dao

import (
	"context"
	"errors"

	"github.com/wu-clan/lykn/internal/common"
	"github.com/wu-clan/lykn/internal/model"
	"gorm.io/gorm"
)

type FeatureDAO struct{ db *gorm.DB }

func NewFeatureDAO(db *gorm.DB) *FeatureDAO { return &FeatureDAO{db: db} }

func (d *FeatureDAO) List(ctx context.Context) ([]model.Feature, error) {
	var features []model.Feature
	return features, d.db.WithContext(ctx).Order("id desc").Find(&features).Error
}

func (d *FeatureDAO) Create(ctx context.Context, feature *model.Feature) error {
	return d.db.WithContext(ctx).Create(feature).Error
}

func (d *FeatureDAO) Update(ctx context.Context, feature *model.Feature) error {
	return d.db.WithContext(ctx).Save(feature).Error
}

func (d *FeatureDAO) Delete(ctx context.Context, feature *model.Feature) error {
	return d.db.WithContext(ctx).Delete(feature).Error
}

func (d *FeatureDAO) GetByID(ctx context.Context, id uint) (*model.Feature, error) {
	var feature model.Feature
	if err := d.db.WithContext(ctx).First(&feature, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.ErrNotFound
		}
		return nil, err
	}
	return &feature, nil
}

func (d *FeatureDAO) GetByCode(ctx context.Context, code string) (*model.Feature, error) {
	var feature model.Feature
	if err := d.db.WithContext(ctx).Where("code = ?", code).First(&feature).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.ErrNotFound
		}
		return nil, err
	}
	return &feature, nil
}

func (d *FeatureDAO) CountPlansByFeatureID(ctx context.Context, featureID uint) (int64, error) {
	var count int64
	err := d.db.WithContext(ctx).Model(&model.PlanFeature{}).Where("feature_id = ?", featureID).Count(&count).Error
	return count, err
}
