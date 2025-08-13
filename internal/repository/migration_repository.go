package repository

import (
	"fastdfs-migration-system/internal/models"
	"gorm.io/gorm"
)

// migrationRepository 迁移任务仓库实现
type migrationRepository struct {
	db *gorm.DB
}

// NewMigrationRepository 创建迁移任务仓库
func NewMigrationRepository(db *gorm.DB) MigrationRepository {
	return &migrationRepository{db: db}
}

// Create 创建迁移任务
func (r *migrationRepository) Create(migration *models.Migration) error {
	return r.db.Create(migration).Error
}

// GetByID 根据ID获取迁移任务
func (r *migrationRepository) GetByID(id string) (*models.Migration, error) {
	var migration models.Migration
	err := r.db.Where("id = ?", id).First(&migration).Error
	if err != nil {
		return nil, err
	}
	return &migration, nil
}

// GetAll 获取所有迁移任务
func (r *migrationRepository) GetAll(pagination *models.Pagination) ([]*models.Migration, error) {
	var migrations []*models.Migration
	var total int64
	
	// 计算总数
	if err := r.db.Model(&models.Migration{}).Count(&total).Error; err != nil {
		return nil, err
	}
	pagination.Total = total
	
	// 分页查询
	err := r.db.Offset(pagination.GetOffset()).
		Limit(pagination.GetLimit()).
		Order("created_at DESC").
		Find(&migrations).Error
	
	return migrations, err
}

// Update 更新迁移任务
func (r *migrationRepository) Update(migration *models.Migration) error {
	return r.db.Save(migration).Error
}

// Delete 删除迁移任务
func (r *migrationRepository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&models.Migration{}).Error
}

// GetByStatus 根据状态获取迁移任务
func (r *migrationRepository) GetByStatus(status string) ([]*models.Migration, error) {
	var migrations []*models.Migration
	err := r.db.Where("status = ?", status).
		Order("created_at DESC").
		Find(&migrations).Error
	return migrations, err
}

// UpdateStatus 更新迁移任务状态
func (r *migrationRepository) UpdateStatus(id string, status string) error {
	return r.db.Model(&models.Migration{}).
		Where("id = ?", id).
		Update("status", status).Error
}

// UpdateProgress 更新迁移进度
func (r *migrationRepository) UpdateProgress(id string, progress float64, processedFiles, processedSize int64) error {
	return r.db.Model(&models.Migration{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"progress":        progress,
			"processed_files": processedFiles,
			"processed_size":  processedSize,
		}).Error
}