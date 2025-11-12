package goal

import (
	"context"
	"time"

	domaincontracts "Fynance/internal/domain/contracts"
	"Fynance/internal/domain/user"
	appErrors "Fynance/internal/errors"
	"Fynance/internal/pkg"

	"github.com/oklog/ulid/v2"
)

type Service struct {
	Repository  Repository
	UserService user.Service
}

func (s *Service) CreateGoal(ctx context.Context, request *domaincontracts.GoalCreateRequest) error {
	if err := Validate(*request); err != nil {
		return err
	}

	if _, err := s.UserService.GetByID(ctx, request.UserId.String()); err != nil {
		return appErrors.ErrUserNotFound.WithError(err)
	}

	now := time.Now()
	entity := &Goal{
		Id:            pkg.GenerateULIDObject(),
		UserId:        request.UserId,
		Name:          request.Name,
		TargetAmount:  request.Target,
		CurrentAmount: 0,
		StartedAt:     now,
		EndedAt:       request.EndedAt,
		Status:        Active,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	return s.Repository.Create(ctx, entity)
}

func (s *Service) UpdateGoal(ctx context.Context, request *domaincontracts.GoalUpdateRequest) error {
	if err := ValidateUpdateGoal(*request); err != nil {
		return err
	}

	if err := s.CheckGoalBelongsToUser(ctx, request.Id, request.UserId); err != nil {
		return err
	}

	current, err := s.Repository.GetById(ctx, request.Id)
	if err != nil {
		return err
	}

	current.Name = request.Name
	current.TargetAmount = request.Target
	current.EndedAt = request.EndedAt
	current.UpdatedAt = time.Now()

	return s.Repository.Update(ctx, current)
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
		return nil, appErrors.ErrResourceNotOwned
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
		return appErrors.ErrResourceNotOwned
	}
	return nil
}

func Validate(request domaincontracts.GoalCreateRequest) error {
	if request.Name == "" {
		return appErrors.NewValidationError("name", "é obrigatório")
	}
	if request.Target <= 0 {
		return appErrors.NewValidationError("target", "deve ser maior que zero")
	}
	if request.EndedAt != nil && request.EndedAt.Before(time.Now()) {
		return appErrors.NewValidationError("ended_at", "deve ser uma data futura")
	}
	return nil
}

func ValidateUpdateGoal(request domaincontracts.GoalUpdateRequest) error {
	if request.Name == "" {
		return appErrors.NewValidationError("name", "é obrigatório")
	}
	if request.Target <= 0 {
		return appErrors.NewValidationError("target", "deve ser maior que zero")
	}
	if request.EndedAt != nil && request.EndedAt.Before(time.Now()) {
		return appErrors.NewValidationError("ended_at", "deve ser uma data futura")
	}
	return nil
}
