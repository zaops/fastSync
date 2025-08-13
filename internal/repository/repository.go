package repository

import (
	"gorm.io/gorm"
)

// repository 仓库集合实现
type repository struct {
	db                  *gorm.DB
	migrationRepo       MigrationRepository
	clusterRepo         ClusterRepository
	taskLogRepo         TaskLogRepository
	scheduledTaskRepo   ScheduledTaskRepository
	transferStateRepo   TransferStateRepository
}

// NewRepository 创建仓库集合
func NewRepository(db *gorm.DB) Repository {
	return &repository{
		db:                db,
		migrationRepo:     NewMigrationRepository(db),
		clusterRepo:       NewClusterRepository(db),
		taskLogRepo:       NewTaskLogRepository(db),
		scheduledTaskRepo: NewScheduledTaskRepository(db),
		transferStateRepo: NewTransferStateRepository(db),
	}
}

// Migration 获取迁移任务仓库
func (r *repository) Migration() MigrationRepository {
	return r.migrationRepo
}

// Cluster 获取集群仓库
func (r *repository) Cluster() ClusterRepository {
	return r.clusterRepo
}

// TaskLog 获取任务日志仓库
func (r *repository) TaskLog() TaskLogRepository {
	return r.taskLogRepo
}

// ScheduledTask 获取定时任务仓库
func (r *repository) ScheduledTask() ScheduledTaskRepository {
	return r.scheduledTaskRepo
}

// TransferState 获取传输状态仓库
func (r *repository) TransferState() TransferStateRepository {
	return r.transferStateRepo
}