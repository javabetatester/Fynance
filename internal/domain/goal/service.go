package goal

import (
	"context"

	"Fynance/internal/domain/user"
	"Fynance/internal/pkg"
	"fmt"
	"time"

	"github.com/oklog/ulid/v2"
)

type Service struct {
	Repository  Repository
	UserService user.Service
}

func (s *Service) CreateGoal(ctx context.Context, goal *GoalCreateRequest) error {
	err := Validate(*goal)
	if err != nil {
		return err
	}

	if _, err := s.UserService.GetByID(ctx, goal.UserId.String()); err != nil {
		return fmt.Errorf("user not found")
	}

	now := time.Now()

	goalEntity := &Goal{
		Id:            pkg.GenerateULIDObject(),
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

	return s.Repository.Create(ctx, goalEntity)
}

func (s *Service) UpdateGoal(ctx context.Context, goal *GoalUpdateRequest) error {
	err := ValidateUpdateGoal(*goal)
	if err != nil {
		return err
	}

	err = s.CheckGoalBelongsToUser(ctx, goal.Id, goal.UserId)
	if err != nil {
		return err
	}

	currentGoal, err := s.Repository.GetById(ctx, goal.Id)
	if err != nil {
		return err
	}

	currentGoal.Name = goal.Name
	currentGoal.TargetAmount = goal.Target
	currentGoal.EndedAt = goal.EndedAt
	currentGoal.UpdatedAt = time.Now()

	return s.Repository.Update(ctx, currentGoal)
}

func (s *Service) DeleteGoal(ctx context.Context, goalID ulid.ULID, userID ulid.ULID) error {
	if err := s.CheckGoalBelongsToUser(ctx, goalID, userID); err != nil {
		return err
	}

	return s.Repository.Delete(ctx, goalID)
}

func (s *Service) GetGoalByID(ctx context.Context, goalID ulid.ULID, userID ulid.ULID) (*Goal, error) {
	goal, err := s.Repository.GetById(ctx, goalID)
	if err != nil {
		return nil, err
	}

	if goal.UserId != userID {
		return nil, fmt.Errorf("goal does not belong to user")
	}

	return goal, nil
}

func (s *Service) GetGoalsByUserID(ctx context.Context, userID ulid.ULID) ([]*Goal, error) {
	return s.Repository.GetByUserId(ctx, userID)
}

func (s *Service) ListGoals(ctx context.Context) ([]*Goal, error) {
	return s.Repository.List(ctx)
}

func (s *Service) CheckGoalBelongsToUser(ctx context.Context, goalID ulid.ULID, userID ulid.ULID) error {
	userBelongs, err := s.Repository.CheckGoalBelongsToUser(ctx, goalID, userID)
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
