package login

import (
	"Fynance/internal/domain/user"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	Repository user.Repository
}

func (s *Service) Login(login Login) (*user.User, error) {

	user, err := s.Repository.GetByEmail(login.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(login.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	user.Password = ""
	return user, nil
}

func (s *Service) GetByEmail(email string) (*user.User, error) {
	user, err := s.Repository.GetByEmail(email)
	if err != nil {
		return nil, err
	}

	return user, nil
}
