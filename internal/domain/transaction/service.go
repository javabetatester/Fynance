package transaction

import (
	"Fynance/internal/domain/user"
	"Fynance/internal/utils"
	"errors"
	"strings"
	"time"

	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

type Service struct {
	Repository         Repository
	CategoryRepository CategoryRepository
	UserService        *user.Service
}

func (s *Service) CreateTransaction(transaction *Transaction) error {
	if err := s.ensureUserExists(transaction.UserId); err != nil {
		return err
	}

	err := s.CategoryValidation(transaction.CategoryId, transaction.UserId)
	if err != nil {
		return err
	}

	TransactionCreateStruct(transaction)

	if err := s.Repository.Create(transaction); err != nil {
		return errors.New("failed to create transaction")
	}

	return nil
}

func (s *Service) UpdateTransaction(transaction *Transaction) error {
	if err := s.ensureUserExists(transaction.UserId); err != nil {
		return err
	}

	storedTransaction, err := s.GetTransactionByID(transaction.Id, transaction.UserId)
	if err != nil {
		return err
	}

	transaction.UpdatedAt = time.Now()

	err = s.UpdateTransactionValidation(transaction)
	if err != nil {
		return err
	}

	storedTransaction.CategoryId = transaction.CategoryId
	storedTransaction.Amount = transaction.Amount
	storedTransaction.Description = transaction.Description
	storedTransaction.Type = transaction.Type
	if !transaction.Date.IsZero() {
		storedTransaction.Date = transaction.Date
	}
	storedTransaction.UpdatedAt = transaction.UpdatedAt

	return s.Repository.Update(storedTransaction)
}

func (s *Service) DeleteTransaction(transactionID ulid.ULID, userID ulid.ULID) error {
	if err := s.TransactionExists(transactionID, userID); err != nil {
		return err
	}
	return s.Repository.Delete(transactionID)
}

func (s *Service) GetTransactionByID(transactionID ulid.ULID, userID ulid.ULID) (*Transaction, error) {
	transaction, err := s.Repository.GetByID(transactionID)
	if err != nil {
		return nil, err
	}
	if transaction.UserId != userID {
		return nil, errors.New("transaction does not belong to user")
	}
	return transaction, nil
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

func (s *Service) CreateCategory(category *Category) error {
	if err := s.ensureUserExists(category.UserId); err != nil {
		return err
	}

	category.Name = strings.TrimSpace(category.Name)
	if category.Name == "" {
		return errors.New("name is required")
	}

	if err := s.CategoryExists(category.Name, category.UserId); err != nil {
		return err
	}

	CategoryCreateStruct(category)

	return s.CategoryRepository.Create(category)
}

func (s *Service) UpdateCategory(category *Category) error {
	if err := s.ensureUserExists(category.UserId); err != nil {
		return err
	}

	existingCategory, err := s.CategoryRepository.GetByID(category.Id, category.UserId)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("category not found")
	}
	if err != nil {
		return err
	}

	category.Name = strings.TrimSpace(category.Name)
	if category.Name == "" {
		return errors.New("name is required")
	}

	if !strings.EqualFold(existingCategory.Name, category.Name) {
		if err := s.CategoryExists(category.Name, category.UserId); err != nil {
			return err
		}
	}

	existingCategory.Name = category.Name
	existingCategory.Icon = category.Icon
	existingCategory.UpdatedAt = time.Now()

	return s.CategoryRepository.Update(existingCategory)
}

func (s *Service) DeleteCategory(categoryID ulid.ULID, userID ulid.ULID) error {
	if err := s.ensureUserExists(userID); err != nil {
		return err
	}

	if _, err := s.CategoryRepository.GetByID(categoryID, userID); errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("category not found")
	} else if err != nil {
		return err
	}
	return s.CategoryRepository.Delete(categoryID, userID)
}

func (s *Service) GetCategoryByID(categoryID ulid.ULID, userID ulid.ULID) (*Category, error) {
	if err := s.ensureUserExists(userID); err != nil {
		return nil, err
	}

	category, err := s.CategoryRepository.GetByID(categoryID, userID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("category not found")
	}
	if err != nil {
		return nil, err
	}

	return category, nil
}

func (s *Service) GetAllCategories(userID ulid.ULID) ([]*Category, error) {
	if err := s.ensureUserExists(userID); err != nil {
		return nil, err
	}
	return s.CategoryRepository.GetAll(userID)
}

func (s *Service) CategoryExists(categoryName string, userID ulid.ULID) error {
	trimmedName := strings.TrimSpace(categoryName)
	if trimmedName == "" {
		return errors.New("name is required")
	}

	_, err := s.CategoryRepository.GetByName(trimmedName, userID)

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

func TransactionCreateStruct(transaction *Transaction) {
	transaction.Id = utils.GenerateULIDObject()
	now := utils.SetTimestamps()
	transaction.CreatedAt = now
	transaction.UpdatedAt = now
}

func CategoryCreateStruct(category *Category) {
	category.Id = utils.GenerateULIDObject()
	category.CreatedAt = time.Now()
	category.UpdatedAt = time.Now()
}

func (s *Service) UpdateTransactionValidation(transaction *Transaction) error {

	if transaction.Amount < 0 {
		return errors.New("amount must be greater than 0")
	}

	if _, err := s.GetCategoryByID(transaction.CategoryId, transaction.UserId); err != nil {
		return err
	}

	return nil
}

func (s *Service) TransactionExists(transactionID ulid.ULID, userID ulid.ULID) error {
	_, err := s.GetTransactionByID(transactionID, userID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("transaction does not exist")
	}
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) ensureUserExists(userID ulid.ULID) error {
	if s.UserService == nil {
		return errors.New("user service not configured")
	}
	_, err := s.UserService.GetByID(userID.String())
	if err != nil {
		return errors.New("user not found")
	}
	return nil
}
