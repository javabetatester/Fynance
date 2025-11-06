package investment

import (
	"context"

	"github.com/oklog/ulid/v2"
)

type Repository interface {
	Create(ctx context.Context, investment *Investment) error
	List(ctx context.Context, userId ulid.ULID) ([]*Investment, error)
	Update(ctx context.Context, investment *Investment) error
	Delete(ctx context.Context, id ulid.ULID, userId ulid.ULID) error
	GetInvestmentById(ctx context.Context, id ulid.ULID, userId ulid.ULID) (*Investment, error)
	GetByUserId(ctx context.Context, userId ulid.ULID) ([]*Investment, error)
	GetTotalBalance(ctx context.Context, userId ulid.ULID) (float64, error)
	GetByType(ctx context.Context, userId ulid.ULID, investmentType Types) ([]*Investment, error)
}
