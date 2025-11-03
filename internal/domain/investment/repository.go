package investment

import "github.com/oklog/ulid/v2"

type Repository interface {
	Create(investment *Investment) error
	List(userId ulid.ULID) ([]*Investment, error)
	Update(investment *Investment) error
	Delete(id ulid.ULID, userId ulid.ULID) error
	GetInvestmentById(id ulid.ULID, userId ulid.ULID) (*Investment, error)
	GetByUserId(userId ulid.ULID) ([]*Investment, error)
	GetTotalBalance(userId ulid.ULID) (float64, error)
	GetByType(userId ulid.ULID, investmentType Types) ([]*Investment, error)
}


