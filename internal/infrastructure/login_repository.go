package infrastructure

import (
	"Fynance/internal/domain/login"
	"Fynance/internal/domain/user"

	"gorm.io/gorm"
)

type LoginRepository struct {
	DB *gorm.DB
}

func (r *LoginRepository) Login(login login.Login) (user.User, error) {
	var user user.User
	if err := r.DB.Where("email = ?", login.Email).First(&user).Error; err != nil {
		return user, err
	}

	if err := r.DB.Where("password = ?", login.Password).First(&user).Error; err != nil {
		return user, err
	}

	return user, nil
}

func (r *LoginRepository) GetByEmail(email string) (*user.User, error) {
	var user user.User
	if err := r.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
