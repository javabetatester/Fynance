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

	if _, err := s.UserService.GetByID(goal.UserId.String()); err != nil {
		return fmt.Errorf("user not found")
	}

	now := time.Now()

	goalEntity := &Goal{
		Id:            utils.GenerateULIDObject(),
		UserId:        goal.UserId,
		Name:          goal.Name,
		TargetAmount:  goal.Target,
		CurrentAmount: 0,
		StartedAt:     now,
		EndedAt:       goal.EndedAt,
		Status:        Active,
		CreatedAt:     now,
		UpdatedAt:     now,
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

	currentGoal, err := s.Repository.GetById(goal.Id)
	if err != nil {
		return err
	}

	currentGoal.Name = goal.Name
	currentGoal.TargetAmount = goal.Target
	currentGoal.EndedAt = goal.EndedAt
	currentGoal.UpdatedAt = time.Now()

	return s.Repository.Update(currentGoal)
}

func (s *Service) DeleteGoal(goalID ulid.ULID, userID ulid.ULID) error {
	if err := s.CheckGoalBelongsToUser(goalID, userID); err != nil {
		return err
	}

	return s.Repository.Delete(goalID)
}

func (s *Service) GetGoalByID(goalID ulid.ULID, userID ulid.ULID) (*Goal, error) {
	goal, err := s.Repository.GetById(goalID)
	if err != nil {
		return nil, err
	}

	if goal.UserId != userID {
		return nil, fmt.Errorf("goal does not belong to user")
	}

	return goal, nil
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
	if goal.Target <= 0 {
		return fmt.Errorf("target must be greater than 0")
	}
	if goal.EndedAt != nil && goal.EndedAt.Before(time.Now()) {
		return fmt.Errorf("ended_at must be in the future")
	}
	return nil
}
