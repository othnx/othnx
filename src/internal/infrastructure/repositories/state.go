package repositories

import (
	"aegis/internal/domain/entities"
	"aegis/internal/domain/ports/secondary"

	"gorm.io/gorm"
)

type StateRepository struct {
	db *gorm.DB
}

var _ secondary.StateRepository = (*StateRepository)(nil)

func NewStateRepository(db *gorm.DB) *StateRepository {
	return &StateRepository{db: db}
}

func (r *StateRepository) CreateState(state entities.State) error {
	return r.db.Create(&state).Error
}

func (r *StateRepository) GetAndDeleteState(value string) (entities.State, error) {
	var state entities.State
	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("value = ?", value).First(&state).Error; err != nil {
			return err
		}
		if err := tx.Where("value = ?", value).Delete(&entities.State{}).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return entities.State{}, err
	}
	return state, nil
}
