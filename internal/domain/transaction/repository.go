package transaction

import (
	"context"

	"github.com/oklog/ulid/v2"
)

type Repository interface {
	Create(ctx context.Context, transaction *Transaction) error
	Update(ctx context.Context, transaction *Transaction) error
	Delete(ctx context.Context, transactionID ulid.ULID) error
	GetByID(ctx context.Context, transactionID ulid.ULID) (*Transaction, error)
	GetAll(ctx context.Context, userID ulid.ULID) ([]*Transaction, error)
	GetByAmount(ctx context.Context, amount float64) ([]*Transaction, error)
	GetByName(ctx context.Context, name string) ([]*Transaction, error)
	GetByCategory(ctx context.Context, categoryID ulid.ULID, userID ulid.ULID) ([]*Transaction, error)
	GetByInvestmentId(ctx context.Context, investmentID ulid.ULID, userID ulid.ULID) ([]*Transaction, error)
	GetNumberOfTransactions(ctx context.Context, userID ulid.ULID) (int64, error)
}

type CategoryRepository interface {
	Create(ctx context.Context, category *Category) error
	Update(ctx context.Context, category *Category) error
	Delete(ctx context.Context, categoryID ulid.ULID, userID ulid.ULID) error
	GetByID(ctx context.Context, categoryID ulid.ULID, userID ulid.ULID) (*Category, error)
	GetAll(ctx context.Context, userID ulid.ULID) ([]*Category, error)
	GetByUserID(ctx context.Context, userID ulid.ULID) ([]*Category, error)
	BelongsToUser(ctx context.Context, categoryID ulid.ULID, userID ulid.ULID) (bool, error)
	GetByName(ctx context.Context, categoryName string, userID ulid.ULID) (*Category, error)
}
