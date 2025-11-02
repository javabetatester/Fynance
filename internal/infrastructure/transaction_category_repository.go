package infrastructure

import (
	"Fynance/internal/domain/transaction"

	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

type TransactionCategoryRepository struct {
	DB *gorm.DB
}

func (r *TransactionCategoryRepository) Create(category *transaction.Category) error {
	return r.DB.Create(&category).Error
}

func (r *TransactionCategoryRepository) Update(category *transaction.Category) error {
	return r.DB.Save(&category).Error
}

func (r *TransactionCategoryRepository) Delete(categoryID ulid.ULID, userID ulid.ULID) error {
	return r.DB.Where("id = ? AND user_id = ?", categoryID, userID).Delete(&transaction.Category{}).Error
}

func (r *TransactionCategoryRepository) GetByID(categoryID ulid.ULID, userID ulid.ULID) (*transaction.Category, error) {
	var category transaction.Category
	err := r.DB.Where("id = ? AND user_id = ?", categoryID, userID).First(&category).Error
	return &category, err
}

func (r *TransactionCategoryRepository) GetAll(userID ulid.ULID) ([]*transaction.Category, error) {
	var categories []*transaction.Category
	err := r.DB.Where("user_id = ?", userID).Find(&categories).Error
	return categories, err
}

func (r *TransactionCategoryRepository) GetByUserID(userID ulid.ULID) ([]*transaction.Category, error) {
	var categories []*transaction.Category
	err := r.DB.Where("user_id = ?", userID).Find(&categories).Error
	return categories, err
}

func (r *TransactionCategoryRepository) GetByName(CategoryName string, userID ulid.ULID) (*transaction.Category, error) {
	var category transaction.Category
	err := r.DB.Where("name = ? AND user_id = ?", CategoryName, userID).First(&category).Error
	return &category, err
}

func (r *TransactionCategoryRepository) BelongsToUser(categoryID ulid.ULID, userID ulid.ULID) (bool, error) {
	var count int64
	err := r.DB.Model(&transaction.Category{}).Where("id = ? AND user_id = ?", categoryID, userID).Count(&count).Error
	return count > 0, err
}
