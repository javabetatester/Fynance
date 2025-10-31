package infrastructure

import (
	"Fynance/internal/domain/user"

	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func (r *UserRepository) Create(user *user.User) error {
	return r.DB.Create(user).Error
}

func (r *UserRepository) Update(user *user.User) error {
	return r.DB.Save(user).Error
}

func (r *UserRepository) Delete(id string) error {
	return r.DB.Where("id = ?", id).Delete(&user.User{}).Error
}

func (r *UserRepository) GetById(id string) (*user.User, error) {
	var user user.User
	if err := r.DB.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByEmail(email string) (*user.User, error) {
	var user user.User
	if err := r.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
