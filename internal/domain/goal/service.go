package goal

import (
	"time"

	"github.com/google/uuid"
)

type Service struct {
	Repository Repository
}

func (s *Service) CreateGoal(goal *Goal) error {

	goal.Id = uuid.New()
	goal.CreatedAt = time.Now()
	goal.UpdatedAt = time.Now()
	return s.Repository.Create(goal)
}

func (s *Service) UpdateGoal(goal *Goal) error {
	return s.Repository.Update(goal)
}

func (s *Service) DeleteGoal(id uuid.UUID) error {
	return s.Repository.Delete(id)
}

func (s *Service) GetGoalByID(id uuid.UUID) (*Goal, error) {
	return s.Repository.GetById(id)
}

func (s *Service) GetGoalsByUserID(userID uuid.UUID) ([]*Goal, error) {
	return s.Repository.GetByUserId(userID)
}

func (s *Service) ListGoals() ([]*Goal, error) {
	return s.Repository.List()
}
