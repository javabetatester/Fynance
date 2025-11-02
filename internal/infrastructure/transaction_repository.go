package infrastructure

import (
	"Fynance/internal/domain/transaction"
	"github.com/google/uuid"

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

func (r *TransactionRepository) Delete(transactionID uuid.UUID) error {
	return r.DB.Delete(&transaction.Transaction{}, transactionID).Error
}

func (r *TransactionRepository) GetByID(transactionID uuid.UUID) (*transaction.Transaction, error) {
	var transaction transaction.Transaction
	err := r.DB.First(&transaction, transactionID).Error
	return &transaction, err
}

func (r *TransactionRepository) GetAll(userID uuid.UUID) ([]*transaction.Transaction, error) {
	var transactions []*transaction.Transaction
	err := r.DB.Where("user_id = ?", userID).Find(&transactions).Error
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

func (r *TransactionRepository) GetByCategory(userID uuid.UUID, categoryID uuid.UUID) ([]*transaction.Transaction, error) {
	var transactions []*transaction.Transaction
	err := r.DB.Where("user_id = ? AND category_id = ?", userID, categoryID).Find(&transactions).Error
	return transactions, err
}
