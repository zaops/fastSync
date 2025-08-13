package repository

import (
	"testing"

	"fastdfs-migration-system/internal/models"
)

// 这个测试验证Repository接口的完整性，不依赖具体的数据库实现
func TestRepositoryInterfaces(t *testing.T) {
	// 测试所有接口方法是否正确定义
	var repo Repository
	
	// 验证Repository接口包含所有子仓库
	if repo != nil {
		_ = repo.Migration()
		_ = repo.Cluster()
		_ = repo.TaskLog()
		_ = repo.ScheduledTask()
		_ = repo.TransferState()
	}
}

func TestMigrationRepositoryInterface(t *testing.T) {
	// 验证MigrationRepository接口的完整性
	var repo MigrationRepository
	
	if repo != nil {
		migration := &models.Migration{}
		pagination := &models.Pagination{}
		
		_ = repo.Create(migration)
		_, _ = repo.GetByID("test-id")
		_, _ = repo.GetAll(pagination)
		_ = repo.Update(migration)
		_ = repo.Delete("test-id")
		_, _ = repo.GetByStatus("pending")
		_ = repo.UpdateStatus("test-id", "running")
		_ = repo.UpdateProgress("test-id", 50.0, 100, 1024)
	}
}

func TestClusterRepositoryInterface(t *testing.T) {
	// 验证ClusterRepository接口的完整性
	var repo ClusterRepository
	
	if repo != nil {
		cluster := &models.Cluster{}
		pagination := &models.Pagination{}
		
		_ = repo.Create(cluster)
		_, _ = repo.GetByID("test-id")
		_, _ = repo.GetAll(pagination)
		_ = repo.Update(cluster)
		_ = repo.Delete("test-id")
		_, _ = repo.GetByStatus("active")
		_ = repo.UpdateStatus("test-id", "inactive")
	}
}

func TestTaskLogRepositoryInterface(t *testing.T) {
	// 验证TaskLogRepository接口的完整性
	var repo TaskLogRepository
	
	if repo != nil {
		log := &models.TaskLog{}
		pagination := &models.Pagination{}
		
		_ = repo.Create(log)
		_, _ = repo.GetByTaskID("task-id", pagination)
		_, _ = repo.GetByLevel("info", pagination)
		_, _ = repo.GetAll(pagination)
		_ = repo.Delete("log-id")
		_ = repo.DeleteByTaskID("task-id")
		_, _ = repo.Search("query", pagination)
	}
}

func TestScheduledTaskRepositoryInterface(t *testing.T) {
	// 验证ScheduledTaskRepository接口的完整性
	var repo ScheduledTaskRepository
	
	if repo != nil {
		task := &models.ScheduledTask{}
		pagination := &models.Pagination{}
		
		_ = repo.Create(task)
		_, _ = repo.GetByID("test-id")
		_, _ = repo.GetAll(pagination)
		_ = repo.Update(task)
		_ = repo.Delete("test-id")
		_, _ = repo.GetByStatus("active")
		_ = repo.UpdateStatus("test-id", "inactive")
		_ = repo.UpdateLastRun("test-id", "success")
	}
}

func TestTransferStateRepositoryInterface(t *testing.T) {
	// 验证TransferStateRepository接口的完整性
	var repo TransferStateRepository
	
	if repo != nil {
		state := &models.TransferState{}
		chunkStates := []models.ChunkState{}
		
		_ = repo.Create(state)
		_, _ = repo.GetByID("test-id")
		_, _ = repo.GetByTaskID("task-id")
		_, _ = repo.GetByFileID("file-id")
		_ = repo.Update(state)
		_ = repo.Delete("test-id")
		_ = repo.DeleteByTaskID("task-id")
		_, _ = repo.GetByStatus("pending")
		_ = repo.UpdateProgress("test-id", 1024, chunkStates)
	}
}