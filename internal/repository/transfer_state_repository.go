package repository

import (
	"fastdfs-migration-system/internal/models"
	"gorm.io/gorm"
)

// transferStateRepository 传输状态仓库实现
type transferStateRepository struct {
	db *gorm.DB
}

// NewTransferStateRepository 创建传输状态仓库
func NewTransferStateRepository(db *gorm.DB) TransferStateRepository {
	return &transferStateRepository{db: db}
}

// Create 创建传输状态
func (r *transferStateRepository) Create(state *models.TransferState) error {
	return r.db.Create(state).Error
}

// GetByID 根据ID获取传输状态
func (r *transferStateRepository) GetByID(id string) (*models.TransferState, error) {
	var state models.TransferState
	err := r.db.Where("id = ?", id).First(&state).Error
	if err != nil {
		return nil, err
	}
	return &state, nil
}

// GetByTaskID 根据任务ID获取传输状态
func (r *transferStateRepository) GetByTaskID(taskID string) ([]*models.TransferState, error) {
	var states []*models.TransferState
	err := r.db.Where("task_id = ?", taskID).
		Order("created_at DESC").
		Find(&states).Error
	return states, err
}

// GetByFileID 根据文件ID获取传输状态
func (r *transferStateRepository) GetByFileID(fileID string) (*models.TransferState, error) {
	var state models.TransferState
	err := r.db.Where("file_id = ?", fileID).First(&state).Error
	if err != nil {
		return nil, err
	}
	return &state, nil
}

// Update 更新传输状态
func (r *transferStateRepository) Update(state *models.TransferState) error {
	return r.db.Save(state).Error
}

// Delete 删除传输状态
func (r *transferStateRepository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&models.TransferState{}).Error
}

// DeleteByTaskID 根据任务ID删除传输状态
func (r *transferStateRepository) DeleteByTaskID(taskID string) error {
	return r.db.Where("task_id = ?", taskID).Delete(&models.TransferState{}).Error
}

// GetByStatus 根据状态获取传输状态
func (r *transferStateRepository) GetByStatus(status string) ([]*models.TransferState, error) {
	var states []*models.TransferState
	err := r.db.Where("status = ?", status).
		Order("created_at DESC").
		Find(&states).Error
	return states, err
}

// UpdateProgress 更新传输进度
func (r *transferStateRepository) UpdateProgress(id string, transferredSize int64, chunkStates []models.ChunkState) error {
	return r.db.Model(&models.TransferState{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"transferred_size": transferredSize,
			"chunk_states":     models.ChunkStates(chunkStates),
		}).Error
}