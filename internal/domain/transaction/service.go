package transaction

import (
	"Fynance/internal/utils"
	"errors"

	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

type Service struct {
	Repository         Repository
	CategoryRepository CategoryRepository
}

func (s *Service) CreateTransaction(transaction *Transaction) error {
	_, err := s.CategoryRepository.GetByID(transaction.CategoryId, transaction.UserId)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("category does not exist")
	}

	if err != nil {
		return err
	}

	transaction.Id = ulid.MustNew(ulid.Timestamp(utils.SetTimestamps()), ulid.DefaultEntropy())
	now := utils.SetTimestamps()
	transaction.CreatedAt = now
	transaction.UpdatedAt = now

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
	if err := s.CategoryExists(category.Name, category.UserId); err != nil {
		return err
	}

	category.Id = ulid.MustNew(ulid.Timestamp(utils.SetTimestamps()), ulid.DefaultEntropy())

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

func (s *Service) EnsureDefaultInvestmentCategory(userID ulid.ULID) (ulid.ULID, error) {
	const defaultName = "Investment"
	cat, err := s.CategoryRepository.GetByName(defaultName, userID)
	if err == nil {
		return cat.Id, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c := &Category{
			UserId: userID,
			Name:   defaultName,
		}
		c.Id = ulid.MustNew(ulid.Timestamp(utils.SetTimestamps()), ulid.DefaultEntropy())
		if errs := s.CategoryRepository.Create(c); errs != nil {
			return ulid.ULID{}, errs
		}
		return c.Id, nil
	}
	return ulid.ULID{}, err
}

func (s *Service) GetNumberOfTransactions(userID ulid.ULID) (int64, error) {
	return s.Repository.GetNumberOfTransactions(userID)
}
