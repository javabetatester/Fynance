package goal

import "github.com/google/uuid"

type Repository interface {
	Create(goal *Goal) error
	List() ([]*Goal, error)
	Update(goal *Goal) error
	Delete(id uuid.UUID) error
	GetById(id uuid.UUID) (*Goal, error)
	GetByUserId(userId uuid.UUID) ([]*Goal, error)
}
