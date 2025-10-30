package goal

import "github.com/google/uuid"

type Repository interface {
	Create(goal *Goal) error
	List() ([]*Goal, error)
	Update(goal *Goal) error
	Delete(id int) error
	GetById(id int) (*Goal, error)
	GetByUserId(userId uuid.UUID) ([]*Goal, error)
}
