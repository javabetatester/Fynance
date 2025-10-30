package investment

import "github.com/google/uuid"

type Repositort interface {
	Create(investment *Investment) error
	List() ([]*Investment, error)
	Update(investment *Investment) error
	Delete(id int) error
	GetById(id int) (*Investment, error)
	GetByUserId(userId uuid.UUID) ([]*Investment, error)
}
