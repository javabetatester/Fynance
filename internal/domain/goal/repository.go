package goal

import (
	"context"

	"github.com/oklog/ulid/v2"
)

type Repository interface {
	Create(ctx context.Context, goal *Goal) error
	List(ctx context.Context) ([]*Goal, error)
	Update(ctx context.Context, goal *Goal) error
	UpdateFields(ctx context.Context, id ulid.ULID, fields map[string]interface{}) error
	Delete(ctx context.Context, id ulid.ULID) error
	GetById(ctx context.Context, id ulid.ULID) (*Goal, error)
	GetByUserId(ctx context.Context, userId ulid.ULID) ([]*Goal, error)
	CheckGoalBelongsToUser(ctx context.Context, goalID ulid.ULID, userID ulid.ULID) (bool, error)
}
