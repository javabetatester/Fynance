package goal

import "github.com/google/uuid"

type Service struct {
	Repository Repository
}

func (s *Service) CreateGoal(goal *Goal) error {
	return s.Repository.Create(goal)
}

func (s *Service) UpdateGoal(goal *Goal) error {
	return s.Repository.Update(goal)
}

func (s *Service) DeleteGoal(id int) error {
	return s.Repository.Delete(id)
}

func (s *Service) GetGoalByID(id int) (*Goal, error) {
	return s.Repository.GetById(id)
}

func (s *Service) GetGoalsByUserID(userID uuid.UUID) ([]*Goal, error) {
	return s.Repository.GetByUserId(userID)
}

func (s *Service) ListGoals() ([]*Goal, error) {
	return s.Repository.List()
}
