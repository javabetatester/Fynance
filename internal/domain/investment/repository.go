package investment

import "github.com/oklog/ulid/v2"

type Repository interface {
	Create(investment *Investment) error
	List() ([]*Investment, error)
	Update(investment *Investment) error
	Delete(id int) error
	GetById(id int) (*Investment, error)
	GetByUserId(userId ulid.ULID) ([]*Investment, error)
}
