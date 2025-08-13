package repository

import (
	"fastdfs-migration-system/internal/models"
)

// MigrationRepository 迁移任务仓库接口
type MigrationRepository interface {
	Create(migration *models.Migration) error
	GetByID(id string) (*models.Migration, error)
	GetAll(pagination *models.Pagination) ([]*models.Migration, error)
	Update(migration *models.Migration) error
	Delete(id string) error
	GetByStatus(status string) ([]*models.Migration, error)
	UpdateStatus(id string, status string) error
	UpdateProgress(id string, progress float64, processedFiles, processedSize int64) error
}

// ClusterRepository 集群仓库接口
type ClusterRepository interface {
	Create(cluster *models.Cluster) error
	GetByID(id string) (*models.Cluster, error)
	GetAll(pagination *models.Pagination) ([]*models.Cluster, error)
	Update(cluster *models.Cluster) error
	Delete(id string) error
	GetByStatus(status string) ([]*models.Cluster, error)
	UpdateStatus(id string, status string) error
}

// TaskLogRepository 任务日志仓库接口
type TaskLogRepository interface {
	Create(log *models.TaskLog) error
	GetByTaskID(taskID string, pagination *models.Pagination) ([]*models.TaskLog, error)
	GetByLevel(level string, pagination *models.Pagination) ([]*models.TaskLog, error)
	GetAll(pagination *models.Pagination) ([]*models.TaskLog, error)
	Delete(id string) error
	DeleteByTaskID(taskID string) error
	Search(query string, pagination *models.Pagination) ([]*models.TaskLog, error)
}

// ScheduledTaskRepository 定时任务仓库接口
type ScheduledTaskRepository interface {
	Create(task *models.ScheduledTask) error
	GetByID(id string) (*models.ScheduledTask, error)
	GetAll(pagination *models.Pagination) ([]*models.ScheduledTask, error)
	Update(task *models.ScheduledTask) error
	Delete(id string) error
	GetByStatus(status string) ([]*models.ScheduledTask, error)
	UpdateStatus(id string, status string) error
	UpdateLastRun(id string, result string) error
}

// TransferStateRepository 传输状态仓库接口
type TransferStateRepository interface {
	Create(state *models.TransferState) error
	GetByID(id string) (*models.TransferState, error)
	GetByTaskID(taskID string) ([]*models.TransferState, error)
	GetByFileID(fileID string) (*models.TransferState, error)
	Update(state *models.TransferState) error
	Delete(id string) error
	DeleteByTaskID(taskID string) error
	GetByStatus(status string) ([]*models.TransferState, error)
	UpdateProgress(id string, transferredSize int64, chunkStates []models.ChunkState) error
}

// Repository 仓库集合接口
type Repository interface {
	Migration() MigrationRepository
	Cluster() ClusterRepository
	TaskLog() TaskLogRepository
	ScheduledTask() ScheduledTaskRepository
	TransferState() TransferStateRepository
}