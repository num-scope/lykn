package dao

import (
	"context"
	"errors"

	"github.com/wu-clan/lykn/internal/common"
	"github.com/wu-clan/lykn/internal/model"
	"gorm.io/gorm"
)

type ProjectDAO struct{ db *gorm.DB }

func NewProjectDAO(db *gorm.DB) *ProjectDAO { return &ProjectDAO{db: db} }

func (d *ProjectDAO) List(ctx context.Context) ([]model.Project, error) {
	var projects []model.Project
	return projects, d.db.WithContext(ctx).Order("id desc").Find(&projects).Error
}

func (d *ProjectDAO) Count(ctx context.Context) (int64, error) {
	var count int64
	if err := d.db.WithContext(ctx).Model(&model.Project{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (d *ProjectDAO) Create(ctx context.Context, project *model.Project) error {
	return d.db.WithContext(ctx).Create(project).Error
}

func (d *ProjectDAO) Update(ctx context.Context, project *model.Project) error {
	return d.db.WithContext(ctx).Save(project).Error
}

func (d *ProjectDAO) Delete(ctx context.Context, project *model.Project) error {
	return d.db.WithContext(ctx).Delete(project).Error
}

func (d *ProjectDAO) GetByID(ctx context.Context, id uint) (*model.Project, error) {
	var project model.Project
	if err := d.db.WithContext(ctx).First(&project, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.ErrNotFound
		}
		return nil, err
	}
	return &project, nil
}
