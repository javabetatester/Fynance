package goal

import (
	"Fynance/internal/utils"

	"github.com/oklog/ulid/v2"
)

type Service struct {
	Repository Repository
}

func (s *Service) CreateGoal(goal *Goal) error {
    goal.Id = utils.GenerateULIDObject()
    now := utils.SetTimestamps()
    goal.CreatedAt = now
    goal.UpdatedAt = now
    return s.Repository.Create(goal)
}

func (s *Service) UpdateGoal(goal *Goal) error {
	return s.Repository.Update(goal)
}

func (s *Service) DeleteGoal(id ulid.ULID) error {
	return s.Repository.Delete(id)
}

func (s *Service) GetGoalByID(id ulid.ULID) (*Goal, error) {
	return s.Repository.GetById(id)
}

func (s *Service) GetGoalsByUserID(userID ulid.ULID) ([]*Goal, error) {
	return s.Repository.GetByUserId(userID)
}

func (s *Service) ListGoals() ([]*Goal, error) {
	return s.Repository.List()
}
