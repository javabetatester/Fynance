package infrastructure

import (
	"Fynance/internal/domain/transaction"
	"Fynance/internal/utils"

	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

type TransactionCategoryRepository struct {
	DB *gorm.DB
}

type categoryDB struct {
	UserId string `gorm:"type:varchar(26);index;not null"`
	Id     string `gorm:"type:varchar(26);primaryKey"`
	Name   string `gorm:"size:100;not null"`
	Icon   string `gorm:"size:50"`
}

func toDomainCategory(cdb *categoryDB) (*transaction.Category, error) {
	uid, err := utils.ParseULID(cdb.UserId)
	if err != nil {
		return nil, err
	}
	id, err := utils.ParseULID(cdb.Id)
	if err != nil {
		return nil, err
	}
	return &transaction.Category{
		UserId: uid,
		Id:     id,
		Name:   cdb.Name,
		Icon:   cdb.Icon,
	}, nil
}

func toDBCategory(c *transaction.Category) *categoryDB {
	return &categoryDB{
		UserId: c.UserId.String(),
		Id:     c.Id.String(),
		Name:   c.Name,
		Icon:   c.Icon,
	}
}

func (r *TransactionCategoryRepository) Create(category *transaction.Category) error {
	cdb := toDBCategory(category)
	return r.DB.Table("categories").Create(&cdb).Error
}

func (r *TransactionCategoryRepository) Update(category *transaction.Category) error {
	cdb := toDBCategory(category)
	return r.DB.Table("categories").Where("id = ?", cdb.Id).Updates(&cdb).Error
}

func (r *TransactionCategoryRepository) Delete(categoryID ulid.ULID, userID ulid.ULID) error {
	return r.DB.Table("categories").Where("id = ? AND user_id = ?", categoryID.String(), userID.String()).Delete(&categoryDB{}).Error
}

func (r *TransactionCategoryRepository) GetByID(categoryID ulid.ULID, userID ulid.ULID) (*transaction.Category, error) {
	var row categoryDB
	err := r.DB.Table("categories").Where("id = ? AND user_id = ?", categoryID.String(), userID.String()).First(&row).Error
	if err != nil {
		return nil, err
	}
	return toDomainCategory(&row)
}

func (r *TransactionCategoryRepository) GetAll(userID ulid.ULID) ([]*transaction.Category, error) {
	var rows []categoryDB
	err := r.DB.Table("categories").Where("user_id = ?", userID.String()).Find(&rows).Error
	if err != nil {
		return nil, err
	}
	out := make([]*transaction.Category, 0, len(rows))
	for i := range rows {
		c, err := toDomainCategory(&rows[i])
		if err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, nil
}

func (r *TransactionCategoryRepository) GetByUserID(userID ulid.ULID) ([]*transaction.Category, error) {
	var rows []categoryDB
	err := r.DB.Table("categories").Where("user_id = ?", userID.String()).Find(&rows).Error
	if err != nil {
		return nil, err
	}
	out := make([]*transaction.Category, 0, len(rows))
	for i := range rows {
		c, err := toDomainCategory(&rows[i])
		if err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, nil
}

func (r *TransactionCategoryRepository) GetByName(CategoryName string, userID ulid.ULID) (*transaction.Category, error) {
	var row categoryDB
	err := r.DB.Table("categories").Where("name = ? AND user_id = ?", CategoryName, userID.String()).First(&row).Error
	if err != nil {
		return nil, err
	}
	return toDomainCategory(&row)
}

func (r *TransactionCategoryRepository) BelongsToUser(categoryID ulid.ULID, userID ulid.ULID) (bool, error) {
	var count int64
	err := r.DB.Table("categories").Where("id = ? AND user_id = ?", categoryID.String(), userID.String()).Count(&count).Error
	return count > 0, err
}
