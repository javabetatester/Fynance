package login

import "Fynance/internal/domain/user"

type Repository interface {
	GetByEmail(email string) (*user.User, error)
}
