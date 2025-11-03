package goal

import "github.com/oklog/ulid/v2"

type Repository interface {
	Create(goal *Goal) error
	List() ([]*Goal, error)
	Update(goal *Goal) error
	UpdateFields(id ulid.ULID, fields map[string]interface{}) error
	Delete(id ulid.ULID) error
	GetById(id ulid.ULID) (*Goal, error)
	GetByUserId(userId ulid.ULID) ([]*Goal, error)
	CheckGoalBelongsToUser(goalID ulid.ULID, userID ulid.ULID) (bool, error)
}
