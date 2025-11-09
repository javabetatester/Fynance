package infrastructure

import (
	"context"

	"Fynance/internal/domain/goal"
	"Fynance/internal/pkg"
	"time"

	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

type GoalRepository struct {
	DB *gorm.DB
}

type goalDB struct {
	Id            string  `gorm:"type:varchar(26);primaryKey"`
	UserId        string  `gorm:"type:varchar(26);index;not null"`
	Name          string  `gorm:"not null"`
	TargetAmount  float64 `gorm:"not null"`
	CurrentAmount float64 `gorm:"not null"`
	StartedAt     time.Time
	EndedAt       *time.Time
	Status        goal.GoalStatus `gorm:"not null"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func toDomainGoal(gdb *goalDB) (*goal.Goal, error) {
	id, err := pkg.ParseULID(gdb.Id)
	if err != nil {
		return nil, err
	}
	uid, err := pkg.ParseULID(gdb.UserId)
	if err != nil {
		return nil, err
	}
	return &goal.Goal{
		Id:            id,
		UserId:        uid,
		Name:          gdb.Name,
		TargetAmount:  gdb.TargetAmount,
		CurrentAmount: gdb.CurrentAmount,
		StartedAt:     gdb.StartedAt,
		EndedAt:       gdb.EndedAt,
		Status:        gdb.Status,
		CreatedAt:     gdb.CreatedAt,
		UpdatedAt:     gdb.UpdatedAt,
	}, nil
}

func toDBGoal(g *goal.Goal) *goalDB {
	return &goalDB{
		Id:            g.Id.String(),
		UserId:        g.UserId.String(),
		Name:          g.Name,
		TargetAmount:  g.TargetAmount,
		CurrentAmount: g.CurrentAmount,
		StartedAt:     g.StartedAt,
		EndedAt:       g.EndedAt,
		Status:        g.Status,
		CreatedAt:     g.CreatedAt,
		UpdatedAt:     g.UpdatedAt,
	}
}

func (r *GoalRepository) Create(ctx context.Context, g *goal.Goal) error {
	gdb := toDBGoal(g)
	return r.DB.WithContext(ctx).Table("goals").Create(&gdb).Error
}

func (r *GoalRepository) Delete(ctx context.Context, id ulid.ULID) error {
	return r.DB.WithContext(ctx).Table("goals").Where("id = ?", id.String()).Delete(&goalDB{}).Error
}

func (r *GoalRepository) GetById(ctx context.Context, id ulid.ULID) (*goal.Goal, error) {
	var gdb goalDB
	if err := r.DB.WithContext(ctx).Table("goals").Where("id = ?", id.String()).First(&gdb).Error; err != nil {
		return nil, err
	}
	return toDomainGoal(&gdb)
}

func (r *GoalRepository) GetByUserId(ctx context.Context, userID ulid.ULID) ([]*goal.Goal, error) {
	var rows []goalDB
	if err := r.DB.WithContext(ctx).Table("goals").Where("user_id = ?", userID.String()).Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]*goal.Goal, 0, len(rows))
	for i := range rows {
		g, err := toDomainGoal(&rows[i])
		if err != nil {
			return nil, err
		}
		out = append(out, g)
	}
	return out, nil
}

func (r *GoalRepository) List(ctx context.Context) ([]*goal.Goal, error) {
	var rows []goalDB
	if err := r.DB.WithContext(ctx).Table("goals").Find(&rows).Error; err != nil {
		return nil, err
	}

	out := make([]*goal.Goal, 0, len(rows))
	for i := range rows {
		g, err := toDomainGoal(&rows[i])
		if err != nil {
			return nil, err
		}
		out = append(out, g)
	}
	return out, nil
}

func (r *GoalRepository) Update(ctx context.Context, g *goal.Goal) error {
	gdb := toDBGoal(g)
	return r.DB.WithContext(ctx).Table("goals").Where("id = ?", gdb.Id).Updates(&gdb).Error
}

func (r *GoalRepository) UpdateFields(ctx context.Context, id ulid.ULID, fields map[string]interface{}) error {
	return r.DB.WithContext(ctx).Table("goals").Where("id = ?", id.String()).Updates(fields).Error
}

func (r *GoalRepository) CheckGoalBelongsToUser(ctx context.Context, goalID ulid.ULID, userID ulid.ULID) (bool, error) {
	var count int64
	err := r.DB.WithContext(ctx).Table("goals").Where("id = ? AND user_id = ?", goalID.String(), userID.String()).Count(&count).Error

	if err != nil {
		return false, err
	}
	return count > 0, nil
}
