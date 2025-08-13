package repository

import (
	"os"
	"testing"
	"time"

	"fastdfs-migration-system/internal/database"
	"fastdfs-migration-system/internal/models"
)

var testRepo Repository

func TestMain(m *testing.M) {
	// 设置测试数据库
	config := database.Config{
		Type:     "sqlite",
		DSN:      ":memory:",
		LogLevel: "silent",
	}

	err := database.Initialize(config)
	if err != nil {
		panic("Failed to initialize test database: " + err.Error())
	}

	err = database.AutoMigrate()
	if err != nil {
		panic("Failed to migrate test database: " + err.Error())
	}

	testRepo = NewRepository(database.GetDB())

	// 运行测试
	code := m.Run()

	// 清理
	database.Close()
	os.Exit(code)
}

func TestMigrationRepository_CRUD(t *testing.T) {
	repo := testRepo.Migration()

	// 测试创建
	migration := &models.Migration{
		Name:            "Test Migration",
		SourceClusterID: "source-1",
		TargetClusterID: "target-1",
		Status:          models.MigrationStatusPending,
		Config: models.MigrationConfig{
			ConcurrentWorkers:   5,
			IncrementalSync:     true,
			VerificationEnabled: true,
		},
	}

	err := repo.Create(migration)
	if err != nil {
		t.Fatalf("Failed to create migration: %v", err)
	}

	if migration.ID == "" {
		t.Error("Migration ID should be generated")
	}

	// 测试根据ID查询
	found, err := repo.GetByID(migration.ID)
	if err != nil {
		t.Fatalf("Failed to get migration by ID: %v", err)
	}

	if found.Name != migration.Name {
		t.Errorf("Expected name %s, got %s", migration.Name, found.Name)
	}

	// 测试更新
	migration.Status = models.MigrationStatusRunning
	migration.Progress = 50.0
	err = repo.Update(migration)
	if err != nil {
		t.Fatalf("Failed to update migration: %v", err)
	}

	updated, err := repo.GetByID(migration.ID)
	if err != nil {
		t.Fatalf("Failed to get updated migration: %v", err)
	}

	if updated.Status != models.MigrationStatusRunning {
		t.Errorf("Expected status %s, got %s", models.MigrationStatusRunning, updated.Status)
	}

	if updated.Progress != 50.0 {
		t.Errorf("Expected progress 50.0, got %f", updated.Progress)
	}

	// 测试根据状态查询
	runningMigrations, err := repo.GetByStatus(models.MigrationStatusRunning)
	if err != nil {
		t.Fatalf("Failed to get migrations by status: %v", err)
	}

	if len(runningMigrations) != 1 {
		t.Errorf("Expected 1 running migration, got %d", len(runningMigrations))
	}

	// 测试更新状态
	err = repo.UpdateStatus(migration.ID, models.MigrationStatusCompleted)
	if err != nil {
		t.Fatalf("Failed to update migration status: %v", err)
	}

	// 测试更新进度
	err = repo.UpdateProgress(migration.ID, 100.0, 1000, 1024*1024)
	if err != nil {
		t.Fatalf("Failed to update migration progress: %v", err)
	}

	final, err := repo.GetByID(migration.ID)
	if err != nil {
		t.Fatalf("Failed to get final migration: %v", err)
	}

	if final.Status != models.MigrationStatusCompleted {
		t.Errorf("Expected status %s, got %s", models.MigrationStatusCompleted, final.Status)
	}

	if final.Progress != 100.0 {
		t.Errorf("Expected progress 100.0, got %f", final.Progress)
	}

	if final.ProcessedFiles != 1000 {
		t.Errorf("Expected processed files 1000, got %d", final.ProcessedFiles)
	}

	// 测试分页查询
	pagination := &models.Pagination{Page: 1, PageSize: 10}
	migrations, err := repo.GetAll(pagination)
	if err != nil {
		t.Fatalf("Failed to get all migrations: %v", err)
	}

	if len(migrations) != 1 {
		t.Errorf("Expected 1 migration, got %d", len(migrations))
	}

	if pagination.Total != 1 {
		t.Errorf("Expected total 1, got %d", pagination.Total)
	}

	// 测试删除
	err = repo.Delete(migration.ID)
	if err != nil {
		t.Fatalf("Failed to delete migration: %v", err)
	}

	_, err = repo.GetByID(migration.ID)
	if err == nil {
		t.Error("Should return error when getting deleted migration")
	}
}

func TestClusterRepository_CRUD(t *testing.T) {
	repo := testRepo.Cluster()

	// 测试创建
	cluster := &models.Cluster{
		Name:        "Test Cluster",
		Version:     "5.0.7",
		TrackerAddr: "192.168.1.100",
		TrackerPort: 22122,
		Status:      models.ClusterStatusActive,
		Description: "Test cluster description",
	}

	err := repo.Create(cluster)
	if err != nil {
		t.Fatalf("Failed to create cluster: %v", err)
	}

	// 测试查询
	found, err := repo.GetByID(cluster.ID)
	if err != nil {
		t.Fatalf("Failed to get cluster by ID: %v", err)
	}

	if found.Name != cluster.Name {
		t.Errorf("Expected name %s, got %s", cluster.Name, found.Name)
	}

	// 测试更新
	cluster.Status = models.ClusterStatusInactive
	err = repo.Update(cluster)
	if err != nil {
		t.Fatalf("Failed to update cluster: %v", err)
	}

	// 测试根据状态查询
	inactiveClusters, err := repo.GetByStatus(models.ClusterStatusInactive)
	if err != nil {
		t.Fatalf("Failed to get clusters by status: %v", err)
	}

	if len(inactiveClusters) != 1 {
		t.Errorf("Expected 1 inactive cluster, got %d", len(inactiveClusters))
	}

	// 测试更新状态
	err = repo.UpdateStatus(cluster.ID, models.ClusterStatusActive)
	if err != nil {
		t.Fatalf("Failed to update cluster status: %v", err)
	}

	// 测试分页查询
	pagination := &models.Pagination{Page: 1, PageSize: 10}
	clusters, err := repo.GetAll(pagination)
	if err != nil {
		t.Fatalf("Failed to get all clusters: %v", err)
	}

	if len(clusters) != 1 {
		t.Errorf("Expected 1 cluster, got %d", len(clusters))
	}

	// 测试删除
	err = repo.Delete(cluster.ID)
	if err != nil {
		t.Fatalf("Failed to delete cluster: %v", err)
	}
}

func TestTaskLogRepository_CRUD(t *testing.T) {
	repo := testRepo.TaskLog()

	// 创建测试任务ID
	taskID := "test-task-id"

	// 测试创建多个日志
	logs := []*models.TaskLog{
		{
			TaskID:   taskID,
			TaskType: models.TaskTypeMigration,
			Level:    models.LogLevelInfo,
			Message:  "Info message",
			Details:  models.LogDetails{"key": "value"},
		},
		{
			TaskID:   taskID,
			TaskType: models.TaskTypeMigration,
			Level:    models.LogLevelError,
			Message:  "Error message",
			Details:  models.LogDetails{"error": "test error"},
		},
		{
			TaskID:   "other-task",
			TaskType: models.TaskTypeSystem,
			Level:    models.LogLevelWarn,
			Message:  "Warning message",
		},
	}

	for _, log := range logs {
		err := repo.Create(log)
		if err != nil {
			t.Fatalf("Failed to create task log: %v", err)
		}
	}

	// 测试根据任务ID查询
	pagination := &models.Pagination{Page: 1, PageSize: 10}
	taskLogs, err := repo.GetByTaskID(taskID, pagination)
	if err != nil {
		t.Fatalf("Failed to get logs by task ID: %v", err)
	}

	if len(taskLogs) != 2 {
		t.Errorf("Expected 2 logs for task ID, got %d", len(taskLogs))
	}

	// 测试根据级别查询
	pagination = &models.Pagination{Page: 1, PageSize: 10}
	errorLogs, err := repo.GetByLevel(models.LogLevelError, pagination)
	if err != nil {
		t.Fatalf("Failed to get logs by level: %v", err)
	}

	if len(errorLogs) != 1 {
		t.Errorf("Expected 1 error log, got %d", len(errorLogs))
	}

	// 测试搜索
	pagination = &models.Pagination{Page: 1, PageSize: 10}
	searchResults, err := repo.Search("Info", pagination)
	if err != nil {
		t.Fatalf("Failed to search logs: %v", err)
	}

	if len(searchResults) != 1 {
		t.Errorf("Expected 1 search result, got %d", len(searchResults))
	}

	// 测试获取所有日志
	pagination = &models.Pagination{Page: 1, PageSize: 10}
	allLogs, err := repo.GetAll(pagination)
	if err != nil {
		t.Fatalf("Failed to get all logs: %v", err)
	}

	if len(allLogs) != 3 {
		t.Errorf("Expected 3 logs, got %d", len(allLogs))
	}

	// 测试根据任务ID删除
	err = repo.DeleteByTaskID(taskID)
	if err != nil {
		t.Fatalf("Failed to delete logs by task ID: %v", err)
	}

	// 验证删除结果
	pagination = &models.Pagination{Page: 1, PageSize: 10}
	remainingLogs, err := repo.GetAll(pagination)
	if err != nil {
		t.Fatalf("Failed to get remaining logs: %v", err)
	}

	if len(remainingLogs) != 1 {
		t.Errorf("Expected 1 remaining log, got %d", len(remainingLogs))
	}

	// 测试删除单个日志
	err = repo.Delete(remainingLogs[0].ID)
	if err != nil {
		t.Fatalf("Failed to delete single log: %v", err)
	}
}

func TestScheduledTaskRepository_CRUD(t *testing.T) {
	repo := testRepo.ScheduledTask()

	// 测试创建
	task := &models.ScheduledTask{
		Name:        "Test Scheduled Task",
		CronExpr:    "0 0 * * *",
		Status:      models.ScheduleStatusActive,
		Description: "Test description",
		TaskConfig: models.TaskConfig{
			ConcurrentWorkers: 3,
			IncrementalSync:   false,
		},
	}

	err := repo.Create(task)
	if err != nil {
		t.Fatalf("Failed to create scheduled task: %v", err)
	}

	// 测试查询
	found, err := repo.GetByID(task.ID)
	if err != nil {
		t.Fatalf("Failed to get scheduled task by ID: %v", err)
	}

	if found.Name != task.Name {
		t.Errorf("Expected name %s, got %s", task.Name, found.Name)
	}

	// 测试更新
	task.Status = models.ScheduleStatusInactive
	err = repo.Update(task)
	if err != nil {
		t.Fatalf("Failed to update scheduled task: %v", err)
	}

	// 测试根据状态查询
	inactiveTasks, err := repo.GetByStatus(models.ScheduleStatusInactive)
	if err != nil {
		t.Fatalf("Failed to get tasks by status: %v", err)
	}

	if len(inactiveTasks) != 1 {
		t.Errorf("Expected 1 inactive task, got %d", len(inactiveTasks))
	}

	// 测试更新状态
	err = repo.UpdateStatus(task.ID, models.ScheduleStatusActive)
	if err != nil {
		t.Fatalf("Failed to update task status: %v", err)
	}

	// 测试更新最后运行时间
	err = repo.UpdateLastRun(task.ID, models.ScheduleResultSuccess)
	if err != nil {
		t.Fatalf("Failed to update last run: %v", err)
	}

	updated, err := repo.GetByID(task.ID)
	if err != nil {
		t.Fatalf("Failed to get updated task: %v", err)
	}

	if updated.LastResult != models.ScheduleResultSuccess {
		t.Errorf("Expected last result %s, got %s", models.ScheduleResultSuccess, updated.LastResult)
	}

	if updated.LastRun == nil {
		t.Error("Last run time should be set")
	}

	// 测试分页查询
	pagination := &models.Pagination{Page: 1, PageSize: 10}
	tasks, err := repo.GetAll(pagination)
	if err != nil {
		t.Fatalf("Failed to get all tasks: %v", err)
	}

	if len(tasks) != 1 {
		t.Errorf("Expected 1 task, got %d", len(tasks))
	}

	// 测试删除
	err = repo.Delete(task.ID)
	if err != nil {
		t.Fatalf("Failed to delete task: %v", err)
	}
}

func TestTransferStateRepository_CRUD(t *testing.T) {
	repo := testRepo.TransferState()

	// 测试创建
	state := &models.TransferState{
		TaskID:          "test-task-id",
		FileID:          "test-file-id",
		FilePath:        "/test/file/path",
		TotalSize:       1024,
		TransferredSize: 512,
		ChunkSize:       256,
		Status:          models.TransferStatusRunning,
		ChunkStates: models.ChunkStates{
			{Index: 0, Offset: 0, Size: 256, Completed: true, Checksum: "abc123"},
			{Index: 1, Offset: 256, Size: 256, Completed: false},
		},
		Checksum: "def456",
	}

	err := repo.Create(state)
	if err != nil {
		t.Fatalf("Failed to create transfer state: %v", err)
	}

	// 测试查询
	found, err := repo.GetByID(state.ID)
	if err != nil {
		t.Fatalf("Failed to get transfer state by ID: %v", err)
	}

	if found.FileID != state.FileID {
		t.Errorf("Expected file ID %s, got %s", state.FileID, found.FileID)
	}

	// 测试根据任务ID查询
	taskStates, err := repo.GetByTaskID(state.TaskID)
	if err != nil {
		t.Fatalf("Failed to get states by task ID: %v", err)
	}

	if len(taskStates) != 1 {
		t.Errorf("Expected 1 state for task, got %d", len(taskStates))
	}

	// 测试根据文件ID查询
	fileState, err := repo.GetByFileID(state.FileID)
	if err != nil {
		t.Fatalf("Failed to get state by file ID: %v", err)
	}

	if fileState.ID != state.ID {
		t.Error("Should return the same state")
	}

	// 测试更新
	state.Status = models.TransferStatusCompleted
	state.TransferredSize = 1024
	err = repo.Update(state)
	if err != nil {
		t.Fatalf("Failed to update transfer state: %v", err)
	}

	// 测试根据状态查询
	completedStates, err := repo.GetByStatus(models.TransferStatusCompleted)
	if err != nil {
		t.Fatalf("Failed to get states by status: %v", err)
	}

	if len(completedStates) != 1 {
		t.Errorf("Expected 1 completed state, got %d", len(completedStates))
	}

	// 测试更新进度
	newChunkStates := []models.ChunkState{
		{Index: 0, Offset: 0, Size: 256, Completed: true, Checksum: "abc123"},
		{Index: 1, Offset: 256, Size: 256, Completed: true, Checksum: "def456"},
	}

	err = repo.UpdateProgress(state.ID, 1024, newChunkStates)
	if err != nil {
		t.Fatalf("Failed to update progress: %v", err)
	}

	updated, err := repo.GetByID(state.ID)
	if err != nil {
		t.Fatalf("Failed to get updated state: %v", err)
	}

	if updated.TransferredSize != 1024 {
		t.Errorf("Expected transferred size 1024, got %d", updated.TransferredSize)
	}

	if len(updated.ChunkStates) != 2 {
		t.Errorf("Expected 2 chunk states, got %d", len(updated.ChunkStates))
	}

	// 测试根据任务ID删除
	err = repo.DeleteByTaskID(state.TaskID)
	if err != nil {
		t.Fatalf("Failed to delete states by task ID: %v", err)
	}

	// 验证删除结果
	_, err = repo.GetByID(state.ID)
	if err == nil {
		t.Error("Should return error when getting deleted state")
	}
}

func TestRepository_Integration(t *testing.T) {
	// 测试仓库集合的完整性
	if testRepo.Migration() == nil {
		t.Error("Migration repository should not be nil")
	}

	if testRepo.Cluster() == nil {
		t.Error("Cluster repository should not be nil")
	}

	if testRepo.TaskLog() == nil {
		t.Error("TaskLog repository should not be nil")
	}

	if testRepo.ScheduledTask() == nil {
		t.Error("ScheduledTask repository should not be nil")
	}

	if testRepo.TransferState() == nil {
		t.Error("TransferState repository should not be nil")
	}
}