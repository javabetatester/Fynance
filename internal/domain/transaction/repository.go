package transaction

type Repository interface {
	Create(transaction Transaction) error
	Update(transaction Transaction) error
	Delete(transactionID int) error
	GetByID(transactionID int) (Transaction, error)
	GetAll() ([]Transaction, error)
	GetByAmount(amount float64) ([]Transaction, error)
	GetByName(name string) ([]Transaction, error)
	GetByCategory(categoryID int) ([]Transaction, error)
}

type CategoryRepositoiry interface {
	Create(category Category) error
	Update(category Category) error
	Delete(categoryID int) error
	GetByID(categoryID int) (Category, error)
	GetAll() ([]Category, error)
}
