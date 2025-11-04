package goal

import (
	"Fynance/internal/domain/user"
	"Fynance/internal/utils"
	"fmt"
	"time"

	"github.com/oklog/ulid/v2"
)

type Service struct {
	Repository  Repository
	UserService user.Service
}

func (s *Service) CreateGoal(goal *GoalCreateRequest) error {
	err := Validate(*goal)
	if err != nil {
		return err
	}

	goalEntity := &Goal{
		Id:            utils.GenerateULIDObject(),
		UserId:        goal.UserId,
		Name:          goal.Name,
		TargetAmount:  goal.Target,
		CurrentAmount: goal.Target,
		StartedAt:     time.Now(),
		EndedAt:       goal.EndedAt,
		Status:        Active,
	}

	return s.Repository.Create(goalEntity)
}

func (s *Service) UpdateGoal(goal *GoalUpdateRequest) error {
	err := ValidateUpdateGoal(*goal)
	if err != nil {
		return err
	}

	err = s.CheckGoalBelongsToUser(goal.Id, goal.UserId)
	if err != nil {
		return err
	}

	goalEntity := &Goal{
		Id:            goal.Id,
		UserId:        goal.UserId,
		Name:          goal.Name,
		TargetAmount:  goal.Target,
		EndedAt:       goal.EndedAt,
		Status:        Active,
		UpdatedAt:     time.Now(),
	}

	return s.Repository.Update(goalEntity)
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

func (s *Service) CheckGoalBelongsToUser(goalID ulid.ULID, userID ulid.ULID) error {
	userBelongs, err := s.Repository.CheckGoalBelongsToUser(goalID, userID)
	if err != nil {
		return err
	}
	if !userBelongs {
		return fmt.Errorf("goal does not belong to user")
	}

	return nil
}

func Validate(goal GoalCreateRequest) error {
	if goal.Name == "" {
		return fmt.Errorf("name is required")
	}
	if goal.Target <= 0 {
		return fmt.Errorf("target must be greater than 0")
	}
	if goal.EndedAt != nil && goal.EndedAt.Before(time.Now()) {
		return fmt.Errorf("ended_at must be in the future")
	}

	return nil
}

func ValidateUpdateGoal(goal GoalUpdateRequest) error {
	if goal.Name == "" {
		return fmt.Errorf("name is required")
	}
	if goal.Target == 0 {
		return fmt.Errorf("target must be greater than 0")
	}
	if goal.EndedAt != nil && goal.EndedAt.Before(time.Now()) {
		return fmt.Errorf("ended_at must be in the future")
	}
	return nil
}
