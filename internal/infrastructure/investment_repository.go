package infrastructure

import (
	"Fynance/internal/domain/investment"
	"Fynance/internal/utils"
	"time"

	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

type InvestmentRepository struct {
	DB *gorm.DB
}

type investmentDB struct {
	Id              string    `gorm:"type:varchar(26);primaryKey"`
	UserId          string    `gorm:"type:varchar(26);index;not null"`
	Type            string    `gorm:"type:varchar(20);not null"`
	Name            string    `gorm:"size:100;not null"`
	CurrentBalance  float64   `gorm:"not null;default:0"`
	ReturnBalance   float64   `gorm:"not null;default:0"`
	ReturnRate      float64   `gorm:"default:0"`
	ApplicationDate time.Time `gorm:"not null"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func toDomainInvestment(idb *investmentDB) (*investment.Investment, error) {
	id, err := utils.ParseULID(idb.Id)
	if err != nil {
		return nil, err
	}
	uid, err := utils.ParseULID(idb.UserId)
	if err != nil {
		return nil, err
	}
	return &investment.Investment{
		Id:              id,
		UserId:          uid,
		Type:            investment.Types(idb.Type),
		Name:            idb.Name,
		CurrentBalance:  idb.CurrentBalance,
		ReturnBalance:   idb.ReturnBalance,
		ReturnRate:      idb.ReturnRate,
		ApplicationDate: idb.ApplicationDate,
		CreatedAt:       idb.CreatedAt,
		UpdatedAt:       idb.UpdatedAt,
	}, nil
}

func toDBInvestment(inv *investment.Investment) *investmentDB {
	return &investmentDB{
		Id:              inv.Id.String(),
		UserId:          inv.UserId.String(),
		Type:            string(inv.Type),
		Name:            inv.Name,
		CurrentBalance:  inv.CurrentBalance,
		ReturnBalance:   inv.ReturnBalance,
		ReturnRate:      inv.ReturnRate,
		ApplicationDate: inv.ApplicationDate,
		CreatedAt:       inv.CreatedAt,
		UpdatedAt:       inv.UpdatedAt,
	}
}

func (r *InvestmentRepository) Create(inv *investment.Investment) error {
	idb := toDBInvestment(inv)
	return r.DB.Table("investments").Create(idb).Error
}

func (r *InvestmentRepository) List(userId ulid.ULID) ([]*investment.Investment, error) {
	var rows []investmentDB
	err := r.DB.Table("investments").Where("user_id = ?", userId.String()).
		Order("application_date DESC").
		Find(&rows).Error
	if err != nil {
		return nil, err
	}
	out := make([]*investment.Investment, 0, len(rows))
	for i := range rows {
		inv, err := toDomainInvestment(&rows[i])
		if err != nil {
			return nil, err
		}
		out = append(out, inv)
	}
	return out, nil
}

func (r *InvestmentRepository) Update(inv *investment.Investment) error {
	idb := toDBInvestment(inv)
	return r.DB.Table("investments").Where("id = ?", idb.Id).Updates(idb).Error
}

func (r *InvestmentRepository) Delete(id ulid.ULID, userId ulid.ULID) error {
	return r.DB.Table("investments").Where("id = ? AND user_id = ?", id.String(), userId.String()).
		Delete(&investmentDB{}).Error
}

func (r *InvestmentRepository) GetInvestmentById(id ulid.ULID, userId ulid.ULID) (*investment.Investment, error) {
	var row investmentDB
	err := r.DB.Table("investments").Where("id = ? AND user_id = ?", id.String(), userId.String()).
		First(&row).Error
	if err != nil {
		return nil, err
	}
	return toDomainInvestment(&row)
}

func (r *InvestmentRepository) GetByUserId(userId ulid.ULID) ([]*investment.Investment, error) {
	var rows []investmentDB
	err := r.DB.Table("investments").Where("user_id = ?", userId.String()).
		Order("application_date DESC").
		Find(&rows).Error
	if err != nil {
		return nil, err
	}
	out := make([]*investment.Investment, 0, len(rows))
	for i := range rows {
		inv, err := toDomainInvestment(&rows[i])
		if err != nil {
			return nil, err
		}
		out = append(out, inv)
	}
	return out, nil
}

func (r *InvestmentRepository) GetTotalBalance(userId ulid.ULID) (float64, error) {
	var total float64
	err := r.DB.Table("investments").
		Where("user_id = ?", userId.String()).
		Select("COALESCE(SUM(current_balance), 0)").
		Scan(&total).Error
	return total, err
}

func (r *InvestmentRepository) GetByType(userId ulid.ULID, investmentType investment.Types) ([]*investment.Investment, error) {
	var rows []investmentDB
	err := r.DB.Table("investments").Where("user_id = ? AND type = ?", userId.String(), string(investmentType)).
		Order("application_date DESC").
		Find(&rows).Error
	if err != nil {
		return nil, err
	}
	out := make([]*investment.Investment, 0, len(rows))
	for i := range rows {
		inv, err := toDomainInvestment(&rows[i])
		if err != nil {
			return nil, err
		}
		out = append(out, inv)
	}
	return out, nil
}
