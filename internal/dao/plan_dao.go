package dao

import (
	"context"
	"errors"

	"github.com/wu-clan/lykn/internal/common"
	"github.com/wu-clan/lykn/internal/model"
	"gorm.io/gorm"
)

type PlanDAO struct{ db *gorm.DB }

func NewPlanDAO(db *gorm.DB) *PlanDAO { return &PlanDAO{db: db} }

func (d *PlanDAO) List(ctx context.Context) ([]model.Plan, error) {
	var plans []model.Plan
	return plans, d.db.WithContext(ctx).Preload("Features").Order("id desc").Find(&plans).Error
}

func (d *PlanDAO) Create(ctx context.Context, plan *model.Plan, features []model.Feature) error {
	return d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(plan).Error; err != nil {
			return err
		}
		if len(features) == 0 {
			return nil
		}
		return tx.Model(plan).Association("Features").Replace(features)
	})
}

func (d *PlanDAO) Update(ctx context.Context, plan *model.Plan, features []model.Feature) error {
	return d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(plan).Error; err != nil {
			return err
		}
		return tx.Model(plan).Association("Features").Replace(features)
	})
}

func (d *PlanDAO) Delete(ctx context.Context, plan *model.Plan) error {
	return d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(plan).Association("Features").Clear(); err != nil {
			return err
		}
		return tx.Delete(plan).Error
	})
}

func (d *PlanDAO) GetByID(ctx context.Context, id uint) (*model.Plan, error) {
	var plan model.Plan
	if err := d.db.WithContext(ctx).Preload("Features").First(&plan, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.ErrNotFound
		}
		return nil, err
	}
	return &plan, nil
}

func (d *PlanDAO) GetByCode(ctx context.Context, code string) (*model.Plan, error) {
	var plan model.Plan
	if err := d.db.WithContext(ctx).Preload("Features").Where("code = ?", code).First(&plan).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.ErrNotFound
		}
		return nil, err
	}
	return &plan, nil
}
