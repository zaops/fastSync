package repository

import (
	"fastdfs-migration-system/internal/models"
	"gorm.io/gorm"
)

// taskLogRepository 任务日志仓库实现
type taskLogRepository struct {
	db *gorm.DB
}

// NewTaskLogRepository 创建任务日志仓库
func NewTaskLogRepository(db *gorm.DB) TaskLogRepository {
	return &taskLogRepository{db: db}
}

// Create 创建任务日志
func (r *taskLogRepository) Create(log *models.TaskLog) error {
	return r.db.Create(log).Error
}

// GetByTaskID 根据任务ID获取日志
func (r *taskLogRepository) GetByTaskID(taskID string, pagination *models.Pagination) ([]*models.TaskLog, error) {
	var logs []*models.TaskLog
	var total int64
	
	query := r.db.Model(&models.TaskLog{}).Where("task_id = ?", taskID)
	
	// 计算总数
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}
	pagination.Total = total
	
	// 分页查询
	err := query.Offset(pagination.GetOffset()).
		Limit(pagination.GetLimit()).
		Order("created_at DESC").
		Find(&logs).Error
	
	return logs, err
}

// GetByLevel 根据日志级别获取日志
func (r *taskLogRepository) GetByLevel(level string, pagination *models.Pagination) ([]*models.TaskLog, error) {
	var logs []*models.TaskLog
	var total int64
	
	query := r.db.Model(&models.TaskLog{}).Where("level = ?", level)
	
	// 计算总数
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}
	pagination.Total = total
	
	// 分页查询
	err := query.Offset(pagination.GetOffset()).
		Limit(pagination.GetLimit()).
		Order("created_at DESC").
		Find(&logs).Error
	
	return logs, err
}

// GetAll 获取所有日志
func (r *taskLogRepository) GetAll(pagination *models.Pagination) ([]*models.TaskLog, error) {
	var logs []*models.TaskLog
	var total int64
	
	// 计算总数
	if err := r.db.Model(&models.TaskLog{}).Count(&total).Error; err != nil {
		return nil, err
	}
	pagination.Total = total
	
	// 分页查询
	err := r.db.Offset(pagination.GetOffset()).
		Limit(pagination.GetLimit()).
		Order("created_at DESC").
		Find(&logs).Error
	
	return logs, err
}

// Delete 删除日志
func (r *taskLogRepository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&models.TaskLog{}).Error
}

// DeleteByTaskID 根据任务ID删除日志
func (r *taskLogRepository) DeleteByTaskID(taskID string) error {
	return r.db.Where("task_id = ?", taskID).Delete(&models.TaskLog{}).Error
}

// Search 搜索日志
func (r *taskLogRepository) Search(query string, pagination *models.Pagination) ([]*models.TaskLog, error) {
	var logs []*models.TaskLog
	var total int64
	
	dbQuery := r.db.Model(&models.TaskLog{}).Where("message LIKE ?", "%"+query+"%")
	
	// 计算总数
	if err := dbQuery.Count(&total).Error; err != nil {
		return nil, err
	}
	pagination.Total = total
	
	// 分页查询
	err := dbQuery.Offset(pagination.GetOffset()).
		Limit(pagination.GetLimit()).
		Order("created_at DESC").
		Find(&logs).Error
	
	return logs, err
}