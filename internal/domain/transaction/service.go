package transaction

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Service struct {
	Repository         Repository
	CategoryRepository CategoryRepository
}

func (s *Service) CreateTransaction(transaction *Transaction) error {
	transaction.Id = uuid.New()

	_, err := s.CategoryRepository.GetByID(transaction.CategoryId, transaction.UserId)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("category does not exist")
	}

	if err != nil {
		return err
	}

	return s.Repository.Create(transaction)
}

func (s *Service) UpdateTransaction(transaction *Transaction) error {
	return s.Repository.Update(transaction)
}

func (s *Service) DeleteTransaction(transactionID uuid.UUID) error {
	return s.Repository.Delete(transactionID)
}

func (s *Service) GetTransactionByID(transactionID uuid.UUID) (*Transaction, error) {
	return s.Repository.GetByID(transactionID)
}

func (s *Service) GetAllTransactions(userID uuid.UUID) ([]*Transaction, error) {
	return s.Repository.GetAll(userID)
}

func (s *Service) GetTransactionsByAmount(amount float64) ([]*Transaction, error) {
	return s.Repository.GetByAmount(amount)
}

func (s *Service) GetTransactionsByName(name string) ([]*Transaction, error) {
	return s.Repository.GetByName(name)
}

func (s *Service) GetTransactionsByCategory(categoryID uuid.UUID, userID uuid.UUID) ([]*Transaction, error) {
	return s.Repository.GetByCategory(categoryID, userID)
}

//CATEGORYS

func (s *Service) CreateCategory(category *Category) error {
	category.Id = uuid.New()
	if err := s.CategoryExists(category.Name, category.UserId); err != nil {
		return err
	}
	return s.CategoryRepository.Create(category)
}

func (s *Service) UpdateCategory(category *Category) error {
	return s.CategoryRepository.Update(category)
}

func (s *Service) DeleteCategory(categoryID uuid.UUID, userID uuid.UUID) error {
	return s.CategoryRepository.Delete(categoryID, userID)
}

func (s *Service) GetCategoryByID(categoryID uuid.UUID, userID uuid.UUID) (*Category, error) {
	return s.CategoryRepository.GetByID(categoryID, userID)
}

func (s *Service) GetAllCategories(userID uuid.UUID) ([]*Category, error) {
	return s.CategoryRepository.GetAll(userID)
}

func (s *Service) CategoryExists(categoryName string, userID uuid.UUID) error {
	_, err := s.CategoryRepository.GetByName(categoryName, userID)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}

	if err != nil {
		return err
	}

	return errors.New("category already exists")
}
