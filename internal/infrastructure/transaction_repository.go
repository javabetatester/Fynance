package infrastructure

import (
	"Fynance/internal/domain/transaction"

	"gorm.io/gorm"
)

type TransactionRepository struct {
	DB *gorm.DB
}

func (r *TransactionRepository) Create(transaction transaction.Transaction) error {
	return r.DB.Create(&transaction).Error
}

func (r *TransactionRepository) Update(transaction transaction.Transaction) error {
	return r.DB.Save(&transaction).Error
}

func (r *TransactionRepository) Delete(transactionID int) error {
	return r.DB.Delete(&transaction.Transaction{}, transactionID).Error
}

func (r *TransactionRepository) GetByID(transactionID int) (transaction.Transaction, error) {
	var transaction transaction.Transaction
	err := r.DB.First(&transaction, transactionID).Error
	return transaction, err
}

func (r *TransactionRepository) GetAll() ([]transaction.Transaction, error) {
	var transactions []transaction.Transaction
	err := r.DB.Find(&transactions).Error
	return transactions, err
}

func (r *TransactionRepository) GetByAmount(amount float64) ([]transaction.Transaction, error) {
	var transactions []transaction.Transaction
	err := r.DB.Where("amount = ?", amount).Find(&transactions).Error
	return transactions, err
}

func (r *TransactionRepository) GetByName(name string) ([]transaction.Transaction, error) {
	var transactions []transaction.Transaction
	err := r.DB.Where("name LIKE ?", "%"+name+"%").Find(&transactions).Error
	return transactions, err
}

func (r *TransactionRepository) GetByCategory(categoryID int) ([]transaction.Transaction, error) {
	var transactions []transaction.Transaction
	err := r.DB.Where("category_id = ?", categoryID).Find(&transactions).Error
	return transactions, err
}
