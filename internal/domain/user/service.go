package user

import (
	"context"

	"Fynance/internal/utils"

	"github.com/oklog/ulid/v2"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	Repository Repository
}

func (s *Service) Create(ctx context.Context, user *User) error {
	user.Id = utils.GenerateULID()

	now := utils.SetTimestamps()
	user.CreatedAt = now
	user.UpdatedAt = now

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)

	return s.Repository.Create(ctx, user)
}

func (s *Service) Update(ctx context.Context, user *User) error {
	return s.Repository.Update(ctx, user)
}

func (s *Service) Delete(ctx context.Context, id string) error {
	return s.Repository.Delete(ctx, id)
}

func (s *Service) GetByID(ctx context.Context, id string) (*User, error) {
	user, err := s.Repository.GetById(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *Service) GetByEmail(ctx context.Context, email string) (*User, error) {
	return s.Repository.GetByEmail(ctx, email)
}

func (s *Service) GetPlan(ctx context.Context, id ulid.ULID) (Plan, error) {
	plan, err := s.Repository.GetPlan(ctx, id)
	if err != nil {
		return "", err
	}
	return plan, nil
}
