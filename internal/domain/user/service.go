package user

import "github.com/google/uuid"

type Service struct {
	Repository Repository
}



func (s *Service) Create(user *User) error {
	return s.Repository.Create(user)
}

func (s *Service) Update(user *User) error {
	return s.Repository.Update(user)
}

func (s *Service) Delete(id string) error {
	return s.Repository.Delete(id)
}

func (s *Service) FindByID(id uuid.UUID) (*User, error) {
	return s.Repository.FindByID(id)
}

func (s *Service) FindByEmail(email string) (*User, error) {
	return s.Repository.FindByEmail(email)
}
