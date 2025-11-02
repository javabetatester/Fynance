package infrastructure

import (
	"Fynance/internal/domain/transaction"

	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

type TransactionRepository struct {
	DB *gorm.DB
}

func (r *TransactionRepository) Create(transaction *transaction.Transaction) error {
	return r.DB.Create(&transaction).Error
}

func (r *TransactionRepository) Update(transaction *transaction.Transaction) error {
	return r.DB.Save(&transaction).Error
}

func (r *TransactionRepository) Delete(transactionID ulid.ULID) error {
	return r.DB.Delete(&transaction.Transaction{}, transactionID.String()).Error
}

func (r *TransactionRepository) GetByID(transactionID ulid.ULID) (*transaction.Transaction, error) {
	var transaction transaction.Transaction
	err := r.DB.First(&transaction, transactionID.String()).Error
	return &transaction, err
}

func (r *TransactionRepository) GetAll(userID ulid.ULID) ([]*transaction.Transaction, error) {
	var transactions []*transaction.Transaction
	err := r.DB.Where("user_id = ?", userID.String()).Find(&transactions).Error
	return transactions, err
}

func (r *TransactionRepository) GetByAmount(amount float64) ([]*transaction.Transaction, error) {
	var transactions []*transaction.Transaction
	err := r.DB.Where("amount = ?", amount).Find(&transactions).Error
	return transactions, err
}

func (r *TransactionRepository) GetByName(name string) ([]*transaction.Transaction, error) {
	var transactions []*transaction.Transaction
	err := r.DB.Where("name LIKE ?", "%"+name+"%").Find(&transactions).Error
	return transactions, err
}

func (r *TransactionRepository) GetByCategory(userID ulid.ULID, categoryID ulid.ULID) ([]*transaction.Transaction, error) {
	var transactions []*transaction.Transaction
	err := r.DB.Where("user_id = ? AND category_id = ?", userID.String(), categoryID.String()).Find(&transactions).Error
	return transactions, err
}

func (r *TransactionRepository) GetNumberOfTransactions(userID ulid.ULID) (int64, error) {
	var count int64
	err := r.DB.Model(&transaction.Transaction{}).Where("user_id = ?", userID.String()).Count(&count).Error
	return count, err
}
