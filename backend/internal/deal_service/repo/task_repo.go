package repo

import (
	"CRM/internal/deal_service/model"

	"gorm.io/gorm"
)

type TaskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) CreateTask(task *model.Task) error {
	return r.db.Create(task).Error
}

func (r *TaskRepository) ListByDealID(dealID string) ([]model.Task, error) {
	var tasks []model.Task
	err := r.db.Where("deal_id = ?", dealID).Order("created_at desc").Find(&tasks).Error
	return tasks, err
}
