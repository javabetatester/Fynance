package infrastructure

import (
	"context"

	"Fynance/internal/domain/user"

	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func (r *UserRepository) Create(ctx context.Context, user *user.User) error {
	return r.DB.WithContext(ctx).Create(user).Error
}

func (r *UserRepository) Update(ctx context.Context, user *user.User) error {
	return r.DB.WithContext(ctx).Save(user).Error
}

func (r *UserRepository) Delete(ctx context.Context, id string) error {
	return r.DB.WithContext(ctx).Where("id = ?", id).Delete(&user.User{}).Error
}

func (r *UserRepository) GetById(ctx context.Context, id string) (*user.User, error) {
	var user user.User
	if err := r.DB.WithContext(ctx).Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	var user user.User
	if err := r.DB.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetPlan(ctx context.Context, id ulid.ULID) (user.Plan, error) {
	var u user.User

	err := r.DB.WithContext(ctx).Select("plan").Where("id = ?", id.String()).First(&u).Error
	if err != nil {
		return "", err
	}

	return u.Plan, nil
}
