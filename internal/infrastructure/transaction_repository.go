package infrastructure

import (
	"Fynance/internal/domain/transaction"
	"Fynance/internal/utils"
	"time"

	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

type TransactionRepository struct {
	DB *gorm.DB
}

type transactionDB struct {
	Id           string    `gorm:"type:varchar(26);primaryKey"`
	UserId       string    `gorm:"type:varchar(26);index;not null"`
	Type         string    `gorm:"type:varchar(15);not null"`
	CategoryId   string    `gorm:"type:varchar(26);index"`
	InvestmentId *string   `gorm:"type:varchar(26);index"`
	Amount       float64   `gorm:"not null"`
	Description  string    `gorm:"size:255"`
	Date         time.Time `gorm:"not null"`
	CreatedAt    time.Time `gorm:"not null"`
	UpdatedAt    time.Time `gorm:"not null"`
}

func toDomainTransaction(tdb *transactionDB) (*transaction.Transaction, error) {
	id, err := utils.ParseULID(tdb.Id)
	if err != nil {
		return nil, err
	}
	uid, err := utils.ParseULID(tdb.UserId)
	if err != nil {
		return nil, err
	}
	cid, err := utils.ParseULID(tdb.CategoryId)
	if err != nil {
		return nil, err
	}

	var invID *ulid.ULID
	if tdb.InvestmentId != nil && *tdb.InvestmentId != "" {
		parsed, err := utils.ParseULID(*tdb.InvestmentId)
		if err != nil {
			return nil, err
		}
		invID = &parsed
	}

	return &transaction.Transaction{
		Id:           id,
		UserId:       uid,
		Type:         transaction.Types(tdb.Type),
		CategoryId:   cid,
		InvestmentId: invID,
		Amount:       tdb.Amount,
		Description:  tdb.Description,
		Date:         tdb.Date,
		CreatedAt:    tdb.CreatedAt,
		UpdatedAt:    tdb.UpdatedAt,
	}, nil
}

func toDBTransaction(t *transaction.Transaction) *transactionDB {
	var invID *string
	if t.InvestmentId != nil {
		s := t.InvestmentId.String()
		invID = &s
	}
	return &transactionDB{
		Id:           t.Id.String(),
		UserId:       t.UserId.String(),
		Type:         string(t.Type),
		CategoryId:   t.CategoryId.String(),
		InvestmentId: invID,
		Amount:       t.Amount,
		Description:  t.Description,
		Date:         t.Date,
		CreatedAt:    t.CreatedAt,
		UpdatedAt:    t.UpdatedAt,
	}
}

func (r *TransactionRepository) Create(t *transaction.Transaction) error {
	tdb := toDBTransaction(t)
	return r.DB.Table("transactions").Create(tdb).Error
}

func (r *TransactionRepository) Update(t *transaction.Transaction) error {
	tdb := toDBTransaction(t)
	return r.DB.Table("transactions").Where("id = ?", tdb.Id).Updates(tdb).Error
}

func (r *TransactionRepository) Delete(transactionID ulid.ULID) error {
	return r.DB.Delete(&transaction.Transaction{}, transactionID.String()).Error
}

func (r *TransactionRepository) GetByID(transactionID ulid.ULID) (*transaction.Transaction, error) {
	var tdb transactionDB
	err := r.DB.Table("transactions").Where("id = ?", transactionID.String()).First(&tdb).Error
	if err != nil {
		return nil, err
	}
	return toDomainTransaction(&tdb)
}

func (r *TransactionRepository) GetAll(userID ulid.ULID) ([]*transaction.Transaction, error) {
	var rows []transactionDB
	err := r.DB.Table("transactions").Where("user_id = ?", userID.String()).Find(&rows).Error
	if err != nil {
		return nil, err
	}
	out := make([]*transaction.Transaction, 0, len(rows))
	for i := range rows {
		t, err := toDomainTransaction(&rows[i])
		if err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	return out, nil
}

func (r *TransactionRepository) GetByAmount(amount float64) ([]*transaction.Transaction, error) {
	var rows []transactionDB
	err := r.DB.Table("transactions").Where("amount = ?", amount).Find(&rows).Error
	if err != nil {
		return nil, err
	}
	out := make([]*transaction.Transaction, 0, len(rows))
	for i := range rows {
		t, err := toDomainTransaction(&rows[i])
		if err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	return out, nil
}

func (r *TransactionRepository) GetByName(name string) ([]*transaction.Transaction, error) {
	var rows []transactionDB
	err := r.DB.Table("transactions").Where("name LIKE ?", "%"+name+"%").Find(&rows).Error
	if err != nil {
		return nil, err
	}
	out := make([]*transaction.Transaction, 0, len(rows))
	for i := range rows {
		t, err := toDomainTransaction(&rows[i])
		if err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	return out, nil
}

func (r *TransactionRepository) GetByCategory(userID ulid.ULID, categoryID ulid.ULID) ([]*transaction.Transaction, error) {
	var rows []transactionDB
	err := r.DB.Table("transactions").Where("user_id = ? AND category_id = ?", userID.String(), categoryID.String()).Find(&rows).Error
	if err != nil {
		return nil, err
	}
	out := make([]*transaction.Transaction, 0, len(rows))
	for i := range rows {
		t, err := toDomainTransaction(&rows[i])
		if err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	return out, nil
}

func (r *TransactionRepository) GetNumberOfTransactions(userID ulid.ULID) (int64, error) {
	var count int64
	err := r.DB.Model(&transaction.Transaction{}).Where("user_id = ?", userID.String()).Count(&count).Error
	return count, err
}

func (r *TransactionRepository) GetByInvestmentId(investmentID ulid.ULID, userID ulid.ULID) ([]*transaction.Transaction, error) {
	var rows []transactionDB
	err := r.DB.Table("transactions").Where("investment_id = ? AND user_id = ?", investmentID.String(), userID.String()).Find(&rows).Error
	if err != nil {
		return nil, err
	}
	out := make([]*transaction.Transaction, 0, len(rows))
	for i := range rows {
		t, err := toDomainTransaction(&rows[i])
		if err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	return out, nil
}
