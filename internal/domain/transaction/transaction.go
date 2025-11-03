package transaction

import (
	"time"

	"github.com/oklog/ulid/v2"
)

type Transaction struct {
	Id           ulid.ULID  `gorm:"type:varchar(26);primaryKey" json:"id"`
	UserId       ulid.ULID  `gorm:"type:varchar(26);index:idx_transactions_user_id,priority:1;index:idx_transactions_user_date;not null" json:"user_id"`
	Type         Types      `gorm:"type:varchar(10);not null;index:idx_transactions_type" json:"type"`
	CategoryId   ulid.ULID  `gorm:"type:varchar(26);index:idx_transactions_category_id" json:"category_id"`
	InvestmentId *ulid.ULID `gorm:"type:varchar(26);index:idx_transactions_investment_id" json:"investment_id"`
	Amount       float64    `gorm:"type:decimal(15,2);not null" json:"amount"`
	Description  string     `gorm:"type:varchar(255)" json:"description"`
	Date         time.Time  `gorm:"type:date;not null;index:idx_transactions_user_date,priority:2;index:idx_transactions_date" json:"date"`
	CreatedAt    time.Time  `gorm:"autoCreateTime;not null" json:"created_at"`
	UpdatedAt    time.Time  `gorm:"autoUpdateTime;not null" json:"updated_at"`
}

func (Transaction) TableName() string {
	return "transactions"
}

type Category struct {
	Id        ulid.ULID `gorm:"type:varchar(26);primaryKey" json:"id"`
	UserId    ulid.ULID `gorm:"type:varchar(26);index:idx_categories_user_id;not null" json:"user_id"`
	Name      string    `gorm:"type:varchar(100);not null;index:idx_categories_user_name,unique" json:"name"`
	Icon      string    `gorm:"type:varchar(50)" json:"icon"`
	CreatedAt time.Time `gorm:"autoCreateTime;not null" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime;not null" json:"updated_at"`
}

func (Category) TableName() string {
	return "categories"
}
