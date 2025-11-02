package transaction

import "github.com/oklog/ulid/v2"

type Repository interface {
	Create(transaction *Transaction) error
	Update(transaction *Transaction) error
	Delete(transactionID ulid.ULID) error
	GetByID(transactionID ulid.ULID) (*Transaction, error)
	GetAll(userID ulid.ULID) ([]*Transaction, error)
	GetByAmount(amount float64) ([]*Transaction, error)
	GetByName(name string) ([]*Transaction, error)
	GetByCategory(categoryID ulid.ULID, userID ulid.ULID) ([]*Transaction, error)
	GetNumberOfTransactions(userID ulid.ULID) (int64, error)
}

type CategoryRepository interface {
	Create(category *Category) error
	Update(category *Category) error
	Delete(categoryID ulid.ULID, userID ulid.ULID) error
	GetByID(categoryID ulid.ULID, userID ulid.ULID) (*Category, error)
	GetAll(userID ulid.ULID) ([]*Category, error)
	GetByUserID(userID ulid.ULID) ([]*Category, error)
	BelongsToUser(categoryID ulid.ULID, userID ulid.ULID) (bool, error)
	GetByName(categoryName string, userID ulid.ULID) (*Category, error)
}
