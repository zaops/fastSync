package database

import (
	"os"
	"testing"
	"time"

	"fastdfs-migration-system/internal/models"
)

func TestInitialize(t *testing.T) {
	// 使用内存数据库进行测试
	config := Config{
		Type:     "sqlite",
		DSN:      ":memory:",
		LogLevel: "silent",
	}

	err := Initialize(config)
	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}

	if DB == nil {
		t.Fatal("Database instance should not be nil after initialization")
	}

	// 测试数据库连接
	sqlDB, err := DB.DB()
	if err != nil {
		t.Fatalf("Failed to get underlying sql.DB: %v", err)
	}

	err = sqlDB.Ping()
	if err != nil {
		t.Fatalf("Failed to ping database: %v", err)
	}
}

func TestAutoMigrate(t *testing.T) {
	// 先初始化数据库
	config := Config{
		Type:     "sqlite",
		DSN:      ":memory:",
		LogLevel: "silent",
	}

	err := Initialize(config)
	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}

	// 测试自动迁移
	err = AutoMigrate()
	if err != nil {
		t.Fatalf("Failed to auto migrate: %v", err)
	}

	// 验证表是否创建成功
	tables := []string{"migrations", "clusters", "task_logs", "scheduled_tasks", "transfer_states"}
	for _, table := range tables {
		var count int64
		err := DB.Raw("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name=?", table).Scan(&count).Error
		if err != nil {
			t.Fatalf("Failed to check table %s: %v", table, err)
		}
		if count == 0 {
			t.Errorf("Table %s was not created", table)
		}
	}
}

func TestDatabaseOperations(t *testing.T) {
	// 初始化测试数据库
	config := Config{
		Type:     "sqlite",
		DSN:      ":memory:",
		LogLevel: "silent",
	}

	err := Initialize(config)
	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}

	err = AutoMigrate()
	if err != nil {
		t.Fatalf("Failed to auto migrate: %v", err)
	}

	// 测试创建迁移任务
	migration := &models.Migration{
		Name:            "Test Migration",
		SourceClusterID: "source-cluster-1",
		TargetClusterID: "target-cluster-1",
		Status:          models.MigrationStatusPending,
		Config: models.MigrationConfig{
			ConcurrentWorkers:   5,
			IncrementalSync:     true,
			VerificationEnabled: true,
		},
	}

	err = DB.Create(migration).Error
	if err != nil {
		t.Fatalf("Failed to create migration: %v", err)
	}

	if migration.ID == "" {
		t.Error("Migration ID should be generated automatically")
	}

	// 测试查询迁移任务
	var foundMigration models.Migration
	err = DB.Where("id = ?", migration.ID).First(&foundMigration).Error
	if err != nil {
		t.Fatalf("Failed to find migration: %v", err)
	}

	if foundMigration.Name != migration.Name {
		t.Errorf("Expected migration name %s, got %s", migration.Name, foundMigration.Name)
	}

	// 测试创建集群
	cluster := &models.Cluster{
		Name:        "Test Cluster",
		Version:     "5.0.7",
		TrackerAddr: "192.168.1.100",
		TrackerPort: 22122,
		Status:      models.ClusterStatusActive,
	}

	err = DB.Create(cluster).Error
	if err != nil {
		t.Fatalf("Failed to create cluster: %v", err)
	}

	// 测试创建任务日志
	taskLog := &models.TaskLog{
		TaskID:   migration.ID,
		TaskType: models.TaskTypeMigration,
		Level:    models.LogLevelInfo,
		Message:  "Test log message",
		Details: models.LogDetails{
			"key": "value",
		},
	}

	err = DB.Create(taskLog).Error
	if err != nil {
		t.Fatalf("Failed to create task log: %v", err)
	}

	// 测试创建定时任务
	scheduledTask := &models.ScheduledTask{
		Name:     "Test Scheduled Task",
		CronExpr: "0 0 * * *",
		Status:   models.ScheduleStatusActive,
		TaskConfig: models.TaskConfig{
			ConcurrentWorkers: 3,
			IncrementalSync:   false,
		},
	}

	err = DB.Create(scheduledTask).Error
	if err != nil {
		t.Fatalf("Failed to create scheduled task: %v", err)
	}

	// 测试创建传输状态
	transferState := &models.TransferState{
		TaskID:          migration.ID,
		FileID:          "test-file-id",
		FilePath:        "/test/file/path",
		TotalSize:       1024,
		TransferredSize: 512,
		ChunkSize:       256,
		Status:          models.TransferStatusRunning,
		ChunkStates: models.ChunkStates{
			{Index: 0, Offset: 0, Size: 256, Completed: true},
			{Index: 1, Offset: 256, Size: 256, Completed: false},
		},
	}

	err = DB.Create(transferState).Error
	if err != nil {
		t.Fatalf("Failed to create transfer state: %v", err)
	}

	// 验证所有记录都已创建
	var migrationCount, clusterCount, logCount, taskCount, stateCount int64
	
	DB.Model(&models.Migration{}).Count(&migrationCount)
	DB.Model(&models.Cluster{}).Count(&clusterCount)
	DB.Model(&models.TaskLog{}).Count(&logCount)
	DB.Model(&models.ScheduledTask{}).Count(&taskCount)
	DB.Model(&models.TransferState{}).Count(&stateCount)

	if migrationCount != 1 {
		t.Errorf("Expected 1 migration, got %d", migrationCount)
	}
	if clusterCount != 1 {
		t.Errorf("Expected 1 cluster, got %d", clusterCount)
	}
	if logCount != 1 {
		t.Errorf("Expected 1 log, got %d", logCount)
	}
	if taskCount != 1 {
		t.Errorf("Expected 1 scheduled task, got %d", taskCount)
	}
	if stateCount != 1 {
		t.Errorf("Expected 1 transfer state, got %d", stateCount)
	}
}

func TestClose(t *testing.T) {
	// 初始化数据库
	config := Config{
		Type:     "sqlite",
		DSN:      ":memory:",
		LogLevel: "silent",
	}

	err := Initialize(config)
	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}

	// 测试关闭数据库
	err = Close()
	if err != nil {
		t.Fatalf("Failed to close database: %v", err)
	}
}

func TestGetDB(t *testing.T) {
	// 初始化数据库
	config := Config{
		Type:     "sqlite",
		DSN:      ":memory:",
		LogLevel: "silent",
	}

	err := Initialize(config)
	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}

	db := GetDB()
	if db == nil {
		t.Error("GetDB should return non-nil database instance")
	}

	if db != DB {
		t.Error("GetDB should return the same instance as DB")
	}
}

func TestUnsupportedDatabaseType(t *testing.T) {
	config := Config{
		Type:     "mysql",
		DSN:      "test",
		LogLevel: "silent",
	}

	err := Initialize(config)
	if err == nil {
		t.Error("Should return error for unsupported database type")
	}

	expectedError := "unsupported database type: mysql"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}

func TestAutoMigrateWithoutInitialization(t *testing.T) {
	// 重置DB为nil
	originalDB := DB
	DB = nil
	defer func() {
		DB = originalDB
	}()

	err := AutoMigrate()
	if err == nil {
		t.Error("Should return error when database is not initialized")
	}

	expectedError := "database not initialized"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}

// TestMain 设置和清理测试环境
func TestMain(m *testing.M) {
	// 运行测试
	code := m.Run()
	
	// 清理
	if DB != nil {
		Close()
	}
	
	os.Exit(code)
}