package user

import (
	"Fynance/internal/utils"
	"errors"

	"github.com/oklog/ulid/v2"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	Repository Repository
}

func (s *Service) Create(user *User) error {
	user.Id = utils.GenerateULID()

	now := utils.SetTimestamps()
	user.CreatedAt = now
	user.UpdatedAt = now

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)

	return s.Repository.Create(user)
}

func (s *Service) Update(user *User) error {
	return s.Repository.Update(user)
}

func (s *Service) Delete(id string) error {
	return s.Repository.Delete(id)
}

func (s *Service) GetByID(id string) (*User, error) {
	user, err := s.Repository.GetById(id)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (s *Service) GetByEmail(email string) (*User, error) {
	return s.Repository.GetByEmail(email)
}

func (s *Service) GetPlan(id ulid.ULID) (Plan, error) {
	plan, err := s.Repository.GetPlan(id)
	if err != nil {
		return "", err
	}
	return plan, nil
}
