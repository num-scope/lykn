package dao

import (
	"context"
	"errors"

	"github.com/wu-clan/lykn/internal/common"
	"github.com/wu-clan/lykn/internal/model"
	"gorm.io/gorm"
)

type UserDAO struct{ db *gorm.DB }

func NewUserDAO(db *gorm.DB) *UserDAO { return &UserDAO{db: db} }

func (d *UserDAO) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	if err := d.db.WithContext(ctx).Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.ErrNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (d *UserDAO) GetByID(ctx context.Context, id uint) (*model.User, error) {
	var user model.User
	if err := d.db.WithContext(ctx).First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.ErrNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (d *UserDAO) Count(ctx context.Context) (int64, error) {
	var count int64
	if err := d.db.WithContext(ctx).Model(&model.User{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (d *UserDAO) Create(ctx context.Context, user *model.User) error {
	return d.db.WithContext(ctx).Create(user).Error
}
