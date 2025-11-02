package goal

import "github.com/oklog/ulid/v2"

type Repository interface {
	Create(goal *Goal) error
	List() ([]*Goal, error)
	Update(goal *Goal) error
	Delete(id ulid.ULID) error
	GetById(id ulid.ULID) (*Goal, error)
	GetByUserId(userId ulid.ULID) ([]*Goal, error)
}
