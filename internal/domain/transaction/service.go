package transaction

import (
	"context"
	"errors"

	"Fynance/internal/domain/user"
	appErrors "Fynance/internal/errors"
	"Fynance/internal/utils"
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

func (s *Service) CreateTransaction(ctx context.Context, transaction *Transaction) error {
	if err := s.ensureUserExists(ctx, transaction.UserId); err != nil {
		return err
	}

	err := s.CategoryValidation(ctx, transaction.CategoryId, transaction.UserId)
	if err != nil {
		return err
	}

	TransactionCreateStruct(transaction)

	if err := s.Repository.Create(ctx, transaction); err != nil {
		return appErrors.NewDatabaseError(err)
	}

	return nil
}

func (s *Service) UpdateTransaction(ctx context.Context, transaction *Transaction) error {
	if err := s.ensureUserExists(ctx, transaction.UserId); err != nil {
		return err
	}

	storedTransaction, err := s.GetTransactionByID(ctx, transaction.Id, transaction.UserId)
	if err != nil {
		return err
	}

	transaction.UpdatedAt = time.Now()

	err = s.UpdateTransactionValidation(ctx, transaction)
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

	return s.Repository.Update(ctx, storedTransaction)
}

func (s *Service) DeleteTransaction(ctx context.Context, transactionID ulid.ULID, userID ulid.ULID) error {
	if err := s.TransactionExists(ctx, transactionID, userID); err != nil {
		return err
	}
	return s.Repository.Delete(ctx, transactionID)
}

func (s *Service) GetTransactionByID(ctx context.Context, transactionID ulid.ULID, userID ulid.ULID) (*Transaction, error) {
	transaction, err := s.Repository.GetByID(ctx, transactionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, appErrors.ErrTransactionNotFound
		}
		return nil, appErrors.NewDatabaseError(err)
	}
	if transaction.UserId != userID {
		return nil, appErrors.ErrResourceNotOwned
	}
	return transaction, nil
}

func (s *Service) GetAllTransactions(ctx context.Context, userID ulid.ULID) ([]*Transaction, error) {
	transactions, err := s.Repository.GetAll(ctx, userID)
	if err != nil {
		return nil, appErrors.NewDatabaseError(err)
	}
	return transactions, nil
}

func (s *Service) GetTransactionsByAmount(ctx context.Context, amount float64) ([]*Transaction, error) {
	transactions, err := s.Repository.GetByAmount(ctx, amount)
	if err != nil {
		return nil, appErrors.NewDatabaseError(err)
	}
	return transactions, nil
}

func (s *Service) GetTransactionsByName(ctx context.Context, name string) ([]*Transaction, error) {
	transactions, err := s.Repository.GetByName(ctx, name)
	if err != nil {
		return nil, appErrors.NewDatabaseError(err)
	}
	return transactions, nil
}

func (s *Service) GetTransactionsByCategory(ctx context.Context, categoryID ulid.ULID, userID ulid.ULID) ([]*Transaction, error) {
	transactions, err := s.Repository.GetByCategory(ctx, categoryID, userID)
	if err != nil {
		return nil, appErrors.NewDatabaseError(err)
	}
	return transactions, nil
}

func (s *Service) CreateCategory(ctx context.Context, category *Category) error {
	if err := s.ensureUserExists(ctx, category.UserId); err != nil {
		return err
	}

	category.Name = strings.TrimSpace(category.Name)
	if category.Name == "" {
		return appErrors.NewValidationError("name", "é obrigatório")
	}

	if err := s.CategoryExists(ctx, category.Name, category.UserId); err != nil {
		return err
	}

	CategoryCreateStruct(category)

	if err := s.CategoryRepository.Create(ctx, category); err != nil {
		return appErrors.NewDatabaseError(err)
	}

	return nil
}

func (s *Service) UpdateCategory(ctx context.Context, category *Category) error {
	if err := s.ensureUserExists(ctx, category.UserId); err != nil {
		return err
	}

	existingCategory, err := s.CategoryRepository.GetByID(ctx, category.Id, category.UserId)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return appErrors.ErrCategoryNotFound
	}
	if err != nil {
		return appErrors.NewDatabaseError(err)
	}

	category.Name = strings.TrimSpace(category.Name)
	if category.Name == "" {
		return appErrors.NewValidationError("name", "é obrigatório")
	}

	if !strings.EqualFold(existingCategory.Name, category.Name) {
		if err := s.CategoryExists(ctx, category.Name, category.UserId); err != nil {
			return err
		}
	}

	existingCategory.Name = category.Name
	existingCategory.Icon = category.Icon
	existingCategory.UpdatedAt = time.Now()

	return s.CategoryRepository.Update(ctx, existingCategory)
}

func (s *Service) DeleteCategory(ctx context.Context, categoryID ulid.ULID, userID ulid.ULID) error {
	if err := s.ensureUserExists(ctx, userID); err != nil {
		return err
	}

	if _, err := s.CategoryRepository.GetByID(ctx, categoryID, userID); errors.Is(err, gorm.ErrRecordNotFound) {
		return appErrors.ErrCategoryNotFound
	} else if err != nil {
		return appErrors.NewDatabaseError(err)
	}
	return s.CategoryRepository.Delete(ctx, categoryID, userID)
}

func (s *Service) GetCategoryByID(ctx context.Context, categoryID ulid.ULID, userID ulid.ULID) (*Category, error) {
	if err := s.ensureUserExists(ctx, userID); err != nil {
		return nil, err
	}

	category, err := s.CategoryRepository.GetByID(ctx, categoryID, userID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, appErrors.ErrCategoryNotFound
	}
	if err != nil {
		return nil, appErrors.NewDatabaseError(err)
	}

	return category, nil
}

func (s *Service) GetAllCategories(ctx context.Context, userID ulid.ULID) ([]*Category, error) {
	if err := s.ensureUserExists(ctx, userID); err != nil {
		return nil, err
	}
	categories, err := s.CategoryRepository.GetAll(ctx, userID)
	if err != nil {
		return nil, appErrors.NewDatabaseError(err)
	}
	return categories, nil
}

func (s *Service) CategoryExists(ctx context.Context, categoryName string, userID ulid.ULID) error {
	trimmedName := strings.TrimSpace(categoryName)
	if trimmedName == "" {
		return appErrors.NewValidationError("name", "é obrigatório")
	}

	_, err := s.CategoryRepository.GetByName(ctx, trimmedName, userID)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}

	if err != nil {
		return appErrors.NewDatabaseError(err)
	}

	return appErrors.NewConflictError("categoria")
}

func (s *Service) CategoryValidation(ctx context.Context, categoryId ulid.ULID, userID ulid.ULID) error {
	_, err := s.CategoryRepository.GetByID(ctx, categoryId, userID)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return appErrors.ErrCategoryNotFound
	}

	if err != nil {
		return appErrors.NewDatabaseError(err)
	}

	return nil
}

func (s *Service) GetNumberOfTransactions(ctx context.Context, userID ulid.ULID) (int64, error) {
	count, err := s.Repository.GetNumberOfTransactions(ctx, userID)
	if err != nil {
		return 0, appErrors.NewDatabaseError(err)
	}
	return count, nil
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

func (s *Service) UpdateTransactionValidation(ctx context.Context, transaction *Transaction) error {
	if transaction.Amount < 0 {
		return appErrors.NewValidationError("amount", "deve ser maior que zero")
	}

	if _, err := s.GetCategoryByID(ctx, transaction.CategoryId, transaction.UserId); err != nil {
		return err
	}

	return nil
}

func (s *Service) TransactionExists(ctx context.Context, transactionID ulid.ULID, userID ulid.ULID) error {
	_, err := s.GetTransactionByID(ctx, transactionID, userID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return appErrors.ErrTransactionNotFound
	}
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) ensureUserExists(ctx context.Context, userID ulid.ULID) error {
	if s.UserService == nil {
		return appErrors.ErrInternalServer.WithError(errors.New("user service not configured"))
	}
	_, err := s.UserService.GetByID(ctx, userID.String())
	if err != nil {
		return appErrors.ErrUserNotFound
	}
	return nil
}
