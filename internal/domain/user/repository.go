package user

type Repository interface {
	Create(user *User) error
	Update(user *User) error
	Delete(id string) error
	GetById(id string) (*User, error)
	GetByEmail(email string) (*User, error)
}
