package infrastructure

import (
	"Fynance/internal/domain/transaction"

	"gorm.io/gorm"
)

type TransactionCategoryRepository struct {
	DB *gorm.DB
}

func (r *TransactionCategoryRepository) Create(category transaction.Category) error {
	return r.DB.Create(&category).Error
}

func (r *TransactionCategoryRepository) Update(category transaction.Category) error {
	return r.DB.Save(&category).Error
}

func (r *TransactionCategoryRepository) Delete(categoryID int) error {
	return r.DB.Delete(&transaction.Category{}, categoryID).Error
}

func (r *TransactionCategoryRepository) GetByID(categoryID int) (transaction.Category, error) {
	var category transaction.Category
	err := r.DB.First(&category, categoryID).Error
	return category, err
}

func (r *TransactionCategoryRepository) GetAll() ([]transaction.Category, error) {
	var categories []transaction.Category
	err := r.DB.Find(&categories).Error
	return categories, err
}
