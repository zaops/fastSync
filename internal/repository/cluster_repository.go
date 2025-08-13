package repository

import (
	"fastdfs-migration-system/internal/models"
	"gorm.io/gorm"
)

// clusterRepository 集群仓库实现
type clusterRepository struct {
	db *gorm.DB
}

// NewClusterRepository 创建集群仓库
func NewClusterRepository(db *gorm.DB) ClusterRepository {
	return &clusterRepository{db: db}
}

// Create 创建集群
func (r *clusterRepository) Create(cluster *models.Cluster) error {
	return r.db.Create(cluster).Error
}

// GetByID 根据ID获取集群
func (r *clusterRepository) GetByID(id string) (*models.Cluster, error) {
	var cluster models.Cluster
	err := r.db.Where("id = ?", id).First(&cluster).Error
	if err != nil {
		return nil, err
	}
	return &cluster, nil
}

// GetAll 获取所有集群
func (r *clusterRepository) GetAll(pagination *models.Pagination) ([]*models.Cluster, error) {
	var clusters []*models.Cluster
	var total int64
	
	// 计算总数
	if err := r.db.Model(&models.Cluster{}).Count(&total).Error; err != nil {
		return nil, err
	}
	pagination.Total = total
	
	// 分页查询
	err := r.db.Offset(pagination.GetOffset()).
		Limit(pagination.GetLimit()).
		Order("created_at DESC").
		Find(&clusters).Error
	
	return clusters, err
}

// Update 更新集群
func (r *clusterRepository) Update(cluster *models.Cluster) error {
	return r.db.Save(cluster).Error
}

// Delete 删除集群
func (r *clusterRepository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&models.Cluster{}).Error
}

// GetByStatus 根据状态获取集群
func (r *clusterRepository) GetByStatus(status string) ([]*models.Cluster, error) {
	var clusters []*models.Cluster
	err := r.db.Where("status = ?", status).
		Order("created_at DESC").
		Find(&clusters).Error
	return clusters, err
}

// UpdateStatus 更新集群状态
func (r *clusterRepository) UpdateStatus(id string, status string) error {
	return r.db.Model(&models.Cluster{}).
		Where("id = ?", id).
		Update("status", status).Error
}