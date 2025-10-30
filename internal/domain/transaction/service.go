package transaction

type Service struct {
	Repository          Repository
	CategoryRepositoiry CategoryRepositoiry
}

func (s *Service) CreateTransaction(transaction Transaction) error {
	return s.Repository.Create(transaction)
}

func (s *Service) UpdateTransaction(transaction Transaction) error {
	return s.Repository.Update(transaction)
}

func (s *Service) DeleteTransaction(transactionID int) error {
	return s.Repository.Delete(transactionID)
}

func (s *Service) GetTransactionByID(transactionID int) (Transaction, error) {
	return s.Repository.GetByID(transactionID)
}

func (s *Service) GetAllTransactions() ([]Transaction, error) {
	return s.Repository.GetAll()
}

func (s *Service) GetTransactionsByAmount(amount float64) ([]Transaction, error) {
	return s.Repository.GetByAmount(amount)
}

func (s *Service) GetTransactionsByName(name string) ([]Transaction, error) {
	return s.Repository.GetByName(name)
}

func (s *Service) GetTransactionsByCategory(categoryID int) ([]Transaction, error) {
	return s.Repository.GetByCategory(categoryID)
}

//CATEGORYS

func (s *Service) CreateCategory(category Category) error {
	return s.CategoryRepositoiry.Create(category)
}

func (s *Service) UpdateCategory(category Category) error {
	return s.CategoryRepositoiry.Update(category)
}

func (s *Service) DeleteCategory(categoryID int) error {
	return s.CategoryRepositoiry.Delete(categoryID)
}

func (s *Service) GetCategoryByID(categoryID int) (Category, error) {
	return s.CategoryRepositoiry.GetByID(categoryID)
}

func (s *Service) GetAllCategories() ([]Category, error) {
	return s.CategoryRepositoiry.GetAll()
}
