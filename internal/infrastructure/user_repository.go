package infrastructure

import (
	"context"
	"errors"

	"Fynance/internal/domain/user"
	appErrors "Fynance/internal/errors"

	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func (r *UserRepository) Create(ctx context.Context, user *user.User) error {
	if err := r.DB.WithContext(ctx).Create(user).Error; err != nil {
		return appErrors.NewDatabaseError(err)
	}
	return nil
}

func (r *UserRepository) Update(ctx context.Context, user *user.User) error {
	if err := r.DB.WithContext(ctx).Save(user).Error; err != nil {
		return appErrors.NewDatabaseError(err)
	}
	return nil
}

func (r *UserRepository) Delete(ctx context.Context, id string) error {
	result := r.DB.WithContext(ctx).Where("id = ?", id).Delete(&user.User{})
	if result.Error != nil {
		return appErrors.NewDatabaseError(result.Error)
	}
	if result.RowsAffected == 0 {
		return appErrors.ErrUserNotFound
	}
	return nil
}

func (r *UserRepository) GetById(ctx context.Context, id string) (*user.User, error) {
	var entity user.User
	if err := r.DB.WithContext(ctx).Where("id = ?", id).First(&entity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, appErrors.ErrUserNotFound.WithError(err)
		}
		return nil, appErrors.NewDatabaseError(err)
	}
	return &entity, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	var entity user.User
	if err := r.DB.WithContext(ctx).Where("email = ?", email).First(&entity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, appErrors.ErrUserNotFound.WithError(err)
		}
		return nil, appErrors.NewDatabaseError(err)
	}
	return &entity, nil
}

func (r *UserRepository) GetPlan(ctx context.Context, id ulid.ULID) (user.Plan, error) {
	var entity user.User
	if err := r.DB.WithContext(ctx).Select("plan").Where("id = ?", id.String()).First(&entity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", appErrors.ErrUserNotFound.WithError(err)
		}
		return "", appErrors.NewDatabaseError(err)
	}
	return entity.Plan, nil
}
