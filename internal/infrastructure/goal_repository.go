package infrastructure

import (
	"Fynance/internal/domain/goal"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GoalRepository struct {
	DB *gorm.DB
}

func (r *GoalRepository) Create(g *goal.Goal) error {
	return r.DB.Create(&g).Error
}

func (r *GoalRepository) Delete(id uuid.UUID) error {
	return r.DB.Delete(&goal.Goal{}, id).Error
}

func (r *GoalRepository) GetById(id uuid.UUID) (*goal.Goal, error) {
	var g goal.Goal
	if err := r.DB.First(&g, id).Error; err != nil {
		return nil, err
	}
	return &g, nil
}

func (r *GoalRepository) GetByUserId(userID uuid.UUID) ([]*goal.Goal, error) {
	var goals []*goal.Goal
	if err := r.DB.Where("user_id = ?", userID).Find(&goals).Error; err != nil {
		return nil, err
	}
	return goals, nil
}

func (r *GoalRepository) List() ([]*goal.Goal, error) {
	var goals []*goal.Goal
	if err := r.DB.Find(&goals).Error; err != nil {
		return nil, err
	}
	return goals, nil
}

func (r *GoalRepository) Update(g *goal.Goal) error {
	return r.DB.Save(&g).Error
}
