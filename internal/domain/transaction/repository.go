package transaction

import "github.com/google/uuid"

type Repository interface {
	Create(transaction *Transaction) error
	Update(transaction *Transaction) error
	Delete(transactionID uuid.UUID) error
	GetByID(transactionID uuid.UUID) (*Transaction, error)
	GetAll(userID uuid.UUID) ([]*Transaction, error)
	GetByAmount(amount float64) ([]*Transaction, error)
	GetByName(name string) ([]*Transaction, error)
	GetByCategory(categoryID uuid.UUID, userID uuid.UUID) ([]*Transaction, error)
}

type CategoryRepository interface {
	Create(category *Category) error
	Update(category *Category) error
	Delete(categoryID uuid.UUID, userID uuid.UUID) error
	GetByID(categoryID uuid.UUID, userID uuid.UUID) (*Category, error)
	GetAll(userID uuid.UUID) ([]*Category, error)
	GetByUserID(userID uuid.UUID) ([]*Category, error)
	BelongsToUser(categoryID uuid.UUID, userID uuid.UUID) (bool, error)
	GetByName(categoryName string, userID uuid.UUID) (*Category, error)
}
