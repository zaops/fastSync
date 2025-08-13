package repository

import (
	"fastdfs-migration-system/internal/models"
	"time"

	"gorm.io/gorm"
)

// scheduledTaskRepository 定时任务仓库实现
type scheduledTaskRepository struct {
	db *gorm.DB
}

// NewScheduledTaskRepository 创建定时任务仓库
func NewScheduledTaskRepository(db *gorm.DB) ScheduledTaskRepository {
	return &scheduledTaskRepository{db: db}
}

// Create 创建定时任务
func (r *scheduledTaskRepository) Create(task *models.ScheduledTask) error {
	return r.db.Create(task).Error
}

// GetByID 根据ID获取定时任务
func (r *scheduledTaskRepository) GetByID(id string) (*models.ScheduledTask, error) {
	var task models.ScheduledTask
	err := r.db.Where("id = ?", id).First(&task).Error
	if err != nil {
		return nil, err
	}
	return &task, nil
}

// GetAll 获取所有定时任务
func (r *scheduledTaskRepository) GetAll(pagination *models.Pagination) ([]*models.ScheduledTask, error) {
	var tasks []*models.ScheduledTask
	var total int64
	
	// 计算总数
	if err := r.db.Model(&models.ScheduledTask{}).Count(&total).Error; err != nil {
		return nil, err
	}
	pagination.Total = total
	
	// 分页查询
	err := r.db.Offset(pagination.GetOffset()).
		Limit(pagination.GetLimit()).
		Order("created_at DESC").
		Find(&tasks).Error
	
	return tasks, err
}

// Update 更新定时任务
func (r *scheduledTaskRepository) Update(task *models.ScheduledTask) error {
	return r.db.Save(task).Error
}

// Delete 删除定时任务
func (r *scheduledTaskRepository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&models.ScheduledTask{}).Error
}

// GetByStatus 根据状态获取定时任务
func (r *scheduledTaskRepository) GetByStatus(status string) ([]*models.ScheduledTask, error) {
	var tasks []*models.ScheduledTask
	err := r.db.Where("status = ?", status).
		Order("created_at DESC").
		Find(&tasks).Error
	return tasks, err
}

// UpdateStatus 更新定时任务状态
func (r *scheduledTaskRepository) UpdateStatus(id string, status string) error {
	return r.db.Model(&models.ScheduledTask{}).
		Where("id = ?", id).
		Update("status", status).Error
}

// UpdateLastRun 更新最后运行时间和结果
func (r *scheduledTaskRepository) UpdateLastRun(id string, result string) error {
	now := time.Now()
	return r.db.Model(&models.ScheduledTask{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"last_run":    &now,
			"last_result": result,
		}).Error
}