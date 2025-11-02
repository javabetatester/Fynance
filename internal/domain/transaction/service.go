package transaction

import (
	"errors"
	"time"

	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

type Service struct {
	Repository         Repository
	CategoryRepository CategoryRepository
}

func (s *Service) CreateTransaction(transaction *Transaction) error {
	categoryID, err := ulid.Parse(transaction.CategoryId)
	if err != nil {
		return errors.New("invalid category ID format")
	}

	userID, err := ulid.Parse(transaction.UserId)
	if err != nil {
		return errors.New("invalid user ID format")
	}

	_, err = s.CategoryRepository.GetByID(categoryID, userID)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("category does not exist")
	}

	if err != nil {
		return err
	}

	entropy := ulid.DefaultEntropy()
	transaction.Id = ulid.MustNew(ulid.Timestamp(time.Now()), entropy).String()
	transaction.CreatedAt = time.Now()
	transaction.UpdatedAt = time.Now()

	return s.Repository.Create(transaction)
}

func (s *Service) UpdateTransaction(transaction *Transaction) error {
	return s.Repository.Update(transaction)
}

func (s *Service) DeleteTransaction(transactionID ulid.ULID) error {
	return s.Repository.Delete(transactionID)
}

func (s *Service) GetTransactionByID(transactionID ulid.ULID) (*Transaction, error) {
	return s.Repository.GetByID(transactionID)
}

func (s *Service) GetAllTransactions(userID ulid.ULID) ([]*Transaction, error) {
	return s.Repository.GetAll(userID)
}

func (s *Service) GetTransactionsByAmount(amount float64) ([]*Transaction, error) {
	return s.Repository.GetByAmount(amount)
}

func (s *Service) GetTransactionsByName(name string) ([]*Transaction, error) {
	return s.Repository.GetByName(name)
}

func (s *Service) GetTransactionsByCategory(categoryID ulid.ULID, userID ulid.ULID) ([]*Transaction, error) {
	return s.Repository.GetByCategory(categoryID, userID)
}

// CATEGORYS
func (s *Service) CreateCategory(category *Category) error {
	userID, err := ulid.Parse(category.UserId)
	if err != nil {
		return errors.New("invalid user ID format")
	}

	if err := s.CategoryExists(category.Name, userID); err != nil {
		return err
	}

	entropy := ulid.DefaultEntropy()
	category.Id = ulid.MustNew(ulid.Timestamp(time.Now()), entropy).String()

	return s.CategoryRepository.Create(category)
}

func (s *Service) UpdateCategory(category *Category) error {
	return s.CategoryRepository.Update(category)
}

func (s *Service) DeleteCategory(categoryID ulid.ULID, userID ulid.ULID) error {
	return s.CategoryRepository.Delete(categoryID, userID)
}

func (s *Service) GetCategoryByID(categoryID ulid.ULID, userID ulid.ULID) (*Category, error) {
	return s.CategoryRepository.GetByID(categoryID, userID)
}

func (s *Service) GetAllCategories(userID ulid.ULID) ([]*Category, error) {
	return s.CategoryRepository.GetAll(userID)
}

func (s *Service) CategoryExists(categoryName string, userID ulid.ULID) error {
	_, err := s.CategoryRepository.GetByName(categoryName, userID)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}

	if err != nil {
		return err
	}

	return errors.New("category already exists")
}

func (s *Service) CategoryValidation(categoryId ulid.ULID, userID ulid.ULID) error {
	_, err := s.CategoryRepository.GetByID(categoryId, userID)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("category does not exist")
	}

	if err != nil {
		return err
	}

	return nil
}

func (s *Service) GetNumberOfTransactions(userID ulid.ULID) (int64, error) {
	return s.Repository.GetNumberOfTransactions(userID)
}
