package user

import "github.com/oklog/ulid/v2"

type Repository interface {
	Create(user *User) error
	Update(user *User) error
	Delete(id string) error
	GetById(id string) (*User, error)
	GetByEmail(email string) (*User, error)
	GetPlan(id ulid.ULID) (Plan, error)
}
