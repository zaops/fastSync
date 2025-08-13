package models

import (
	"testing"
	"time"
)

func TestMigrationConfig_JSON(t *testing.T) {
	config := MigrationConfig{
		ConcurrentWorkers:   5,
		IncrementalSync:     true,
		VerificationEnabled: true,
		TimeFilter: &TimeFilter{
			StartTime: &time.Time{},
			EndTime:   &time.Time{},
		},
		RetryConfig: &RetryConfig{
			MaxRetries:    3,
			RetryInterval: time.Minute,
			BackoffFactor: 2.0,
		},
	}

	// 测试序列化
	data, err := config.Value()
	if err != nil {
		t.Errorf("Failed to serialize MigrationConfig: %v", err)
	}

	// 测试反序列化
	var newConfig MigrationConfig
	err = newConfig.Scan(data)
	if err != nil {
		t.Errorf("Failed to deserialize MigrationConfig: %v", err)
	}

	if newConfig.ConcurrentWorkers != config.ConcurrentWorkers {
		t.Errorf("Expected ConcurrentWorkers %d, got %d", config.ConcurrentWorkers, newConfig.ConcurrentWorkers)
	}

	if newConfig.IncrementalSync != config.IncrementalSync {
		t.Errorf("Expected IncrementalSync %v, got %v", config.IncrementalSync, newConfig.IncrementalSync)
	}
}

func TestTransferState_GetProgress(t *testing.T) {
	state := &TransferState{
		TotalSize:       1000,
		TransferredSize: 500,
	}

	progress := state.GetProgress()
	expected := 50.0

	if progress != expected {
		t.Errorf("Expected progress %f, got %f", expected, progress)
	}
}

func TestTransferState_GetCompletedChunks(t *testing.T) {
	state := &TransferState{
		ChunkStates: ChunkStates{
			{Index: 0, Completed: true},
			{Index: 1, Completed: false},
			{Index: 2, Completed: true},
		},
	}

	completed := state.GetCompletedChunks()
	expected := 2

	if completed != expected {
		t.Errorf("Expected completed chunks %d, got %d", expected, completed)
	}
}

func TestPagination_GetOffset(t *testing.T) {
	tests := []struct {
		page     int
		pageSize int
		expected int
	}{
		{1, 10, 0},
		{2, 10, 10},
		{3, 5, 10},
		{0, 10, 0}, // 测试边界情况
		{-1, 10, 0}, // 测试边界情况
	}

	for _, test := range tests {
		p := &Pagination{
			Page:     test.page,
			PageSize: test.pageSize,
		}

		offset := p.GetOffset()
		if offset != test.expected {
			t.Errorf("Page %d, PageSize %d: expected offset %d, got %d", 
				test.page, test.pageSize, test.expected, offset)
		}
	}
}

func TestPagination_GetLimit(t *testing.T) {
	tests := []struct {
		pageSize int
		expected int
	}{
		{10, 10},
		{50, 50},
		{150, 100}, // 测试最大限制
		{0, 10},    // 测试默认值
		{-1, 10},   // 测试默认值
	}

	for _, test := range tests {
		p := &Pagination{
			PageSize: test.pageSize,
		}

		limit := p.GetLimit()
		if limit != test.expected {
			t.Errorf("PageSize %d: expected limit %d, got %d", 
				test.pageSize, test.expected, limit)
		}
	}
}

func TestGenerateID(t *testing.T) {
	id1 := generateID()
	id2 := generateID()

	if id1 == "" {
		t.Error("Generated ID should not be empty")
	}

	if id1 == id2 {
		t.Error("Generated IDs should be unique")
	}

	// ID应该包含时间戳和随机部分
	if len(id1) < 10 {
		t.Error("Generated ID should be at least 10 characters long")
	}
}

func TestNewSuccessResponse(t *testing.T) {
	data := map[string]string{"key": "value"}
	response := NewSuccessResponse(data)

	if response.Code != 200 {
		t.Errorf("Expected code 200, got %d", response.Code)
	}

	if response.Message != "success" {
		t.Errorf("Expected message 'success', got %s", response.Message)
	}

	if response.Data == nil {
		t.Error("Response data should not be nil")
	}
}

func TestNewErrorResponse(t *testing.T) {
	response := NewErrorResponse(400, "Bad Request")

	if response.Code != 400 {
		t.Errorf("Expected code 400, got %d", response.Code)
	}

	if response.Message != "Bad Request" {
		t.Errorf("Expected message 'Bad Request', got %s", response.Message)
	}

	if response.Data != nil {
		t.Error("Error response data should be nil")
	}
}

// 测试Cluster模型的新方法
func TestCluster_IsActive(t *testing.T) {
	cluster := &Cluster{Status: ClusterStatusActive}
	if !cluster.IsActive() {
		t.Error("Cluster should be active")
	}

	cluster.Status = ClusterStatusInactive
	if cluster.IsActive() {
		t.Error("Cluster should not be active")
	}
}

func TestCluster_GetConnectionString(t *testing.T) {
	cluster := &Cluster{
		TrackerAddr: "192.168.1.100",
		TrackerPort: 22122,
	}

	expected := "192.168.1.100:22122"
	result := cluster.GetConnectionString()

	if result != expected {
		t.Errorf("Expected connection string %s, got %s", expected, result)
	}
}

func TestCluster_Validate(t *testing.T) {
	tests := []struct {
		cluster   *Cluster
		shouldErr bool
		errMsg    string
	}{
		{
			cluster: &Cluster{
				Name:        "Test Cluster",
				TrackerAddr: "192.168.1.100",
				TrackerPort: 22122,
				Version:     "5.0.7",
			},
			shouldErr: false,
		},
		{
			cluster:   &Cluster{},
			shouldErr: true,
			errMsg:    "cluster name is required",
		},
		{
			cluster: &Cluster{
				Name:        "Test",
				TrackerPort: 22122,
				Version:     "5.0.7",
			},
			shouldErr: true,
			errMsg:    "tracker address is required",
		},
		{
			cluster: &Cluster{
				Name:        "Test",
				TrackerAddr: "192.168.1.100",
				TrackerPort: 70000,
				Version:     "5.0.7",
			},
			shouldErr: true,
			errMsg:    "tracker port must be between 1 and 65535",
		},
	}

	for i, test := range tests {
		err := test.cluster.Validate()
		if test.shouldErr {
			if err == nil {
				t.Errorf("Test %d: expected error but got none", i)
			} else if err.Error() != test.errMsg {
				t.Errorf("Test %d: expected error '%s', got '%s'", i, test.errMsg, err.Error())
			}
		} else {
			if err != nil {
				t.Errorf("Test %d: expected no error but got: %v", i, err)
			}
		}
	}
}

// 测试Migration模型的新方法
func TestMigration_StatusChecks(t *testing.T) {
	migration := &Migration{}

	// 测试IsRunning
	migration.Status = MigrationStatusRunning
	if !migration.IsRunning() {
		t.Error("Migration should be running")
	}

	// 测试IsCompleted
	migration.Status = MigrationStatusCompleted
	if !migration.IsCompleted() {
		t.Error("Migration should be completed")
	}

	// 测试IsFailed
	migration.Status = MigrationStatusFailed
	if !migration.IsFailed() {
		t.Error("Migration should be failed")
	}

	// 测试CanStart
	migration.Status = MigrationStatusPending
	if !migration.CanStart() {
		t.Error("Migration should be able to start")
	}

	migration.Status = MigrationStatusPaused
	if !migration.CanStart() {
		t.Error("Migration should be able to start from paused")
	}

	// 测试CanPause
	migration.Status = MigrationStatusRunning
	if !migration.CanPause() {
		t.Error("Migration should be able to pause")
	}

	// 测试CanResume
	migration.Status = MigrationStatusPaused
	if !migration.CanResume() {
		t.Error("Migration should be able to resume")
	}
}

func TestMigration_GetProgressPercentage(t *testing.T) {
	migration := &Migration{Progress: 75.5}
	expected := "75.50%"
	result := migration.GetProgressPercentage()

	if result != expected {
		t.Errorf("Expected progress percentage %s, got %s", expected, result)
	}
}

func TestMigration_GetProcessedRatios(t *testing.T) {
	migration := &Migration{
		TotalFiles:      1000,
		ProcessedFiles:  250,
		TotalSize:       1024 * 1024,
		ProcessedSize:   256 * 1024,
	}

	// 测试文件比例
	expectedFileRatio := 0.25
	fileRatio := migration.GetProcessedFilesRatio()
	if fileRatio != expectedFileRatio {
		t.Errorf("Expected file ratio %f, got %f", expectedFileRatio, fileRatio)
	}

	// 测试大小比例
	expectedSizeRatio := 0.25
	sizeRatio := migration.GetProcessedSizeRatio()
	if sizeRatio != expectedSizeRatio {
		t.Errorf("Expected size ratio %f, got %f", expectedSizeRatio, sizeRatio)
	}

	// 测试零除情况
	migration.TotalFiles = 0
	migration.TotalSize = 0
	if migration.GetProcessedFilesRatio() != 0 {
		t.Error("File ratio should be 0 when total files is 0")
	}
	if migration.GetProcessedSizeRatio() != 0 {
		t.Error("Size ratio should be 0 when total size is 0")
	}
}

func TestMigration_Validate(t *testing.T) {
	tests := []struct {
		migration *Migration
		shouldErr bool
		errMsg    string
	}{
		{
			migration: &Migration{
				Name:            "Test Migration",
				SourceClusterID: "source-1",
				TargetClusterID: "target-1",
			},
			shouldErr: false,
		},
		{
			migration: &Migration{},
			shouldErr: true,
			errMsg:    "migration name is required",
		},
		{
			migration: &Migration{
				Name:            "Test",
				TargetClusterID: "target-1",
			},
			shouldErr: true,
			errMsg:    "source cluster ID is required",
		},
		{
			migration: &Migration{
				Name:            "Test",
				SourceClusterID: "source-1",
			},
			shouldErr: true,
			errMsg:    "target cluster ID is required",
		},
		{
			migration: &Migration{
				Name:            "Test",
				SourceClusterID: "same-id",
				TargetClusterID: "same-id",
			},
			shouldErr: true,
			errMsg:    "source and target cluster cannot be the same",
		},
	}

	for i, test := range tests {
		err := test.migration.Validate()
		if test.shouldErr {
			if err == nil {
				t.Errorf("Test %d: expected error but got none", i)
			} else if err.Error() != test.errMsg {
				t.Errorf("Test %d: expected error '%s', got '%s'", i, test.errMsg, err.Error())
			}
		} else {
			if err != nil {
				t.Errorf("Test %d: expected no error but got: %v", i, err)
			}
		}
	}
}

// 测试ScheduledTask模型的新方法
func TestScheduledTask_IsActive(t *testing.T) {
	task := &ScheduledTask{Status: ScheduleStatusActive}
	if !task.IsActive() {
		t.Error("Scheduled task should be active")
	}

	task.Status = ScheduleStatusInactive
	if task.IsActive() {
		t.Error("Scheduled task should not be active")
	}
}

func TestScheduledTask_ShouldRun(t *testing.T) {
	now := time.Now()
	pastTime := now.Add(-time.Hour)
	futureTime := now.Add(time.Hour)

	tests := []struct {
		task       *ScheduledTask
		shouldRun  bool
		description string
	}{
		{
			task: &ScheduledTask{
				Status:  ScheduleStatusActive,
				NextRun: &pastTime,
			},
			shouldRun:   true,
			description: "active task with past next run time",
		},
		{
			task: &ScheduledTask{
				Status:  ScheduleStatusActive,
				NextRun: &futureTime,
			},
			shouldRun:   false,
			description: "active task with future next run time",
		},
		{
			task: &ScheduledTask{
				Status:  ScheduleStatusInactive,
				NextRun: &pastTime,
			},
			shouldRun:   false,
			description: "inactive task",
		},
		{
			task: &ScheduledTask{
				Status:  ScheduleStatusActive,
				NextRun: nil,
			},
			shouldRun:   false,
			description: "active task with no next run time",
		},
	}

	for _, test := range tests {
		result := test.task.ShouldRun()
		if result != test.shouldRun {
			t.Errorf("%s: expected %v, got %v", test.description, test.shouldRun, result)
		}
	}
}

func TestScheduledTask_GetLastRunStatus(t *testing.T) {
	task := &ScheduledTask{}

	// 测试从未运行
	result := task.GetLastRunStatus()
	if result != "Never run" {
		t.Errorf("Expected 'Never run', got %s", result)
	}

	// 测试有运行记录
	lastRun := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
	task.LastRun = &lastRun
	task.LastResult = ScheduleResultSuccess

	result = task.GetLastRunStatus()
	expected := "Last run: 2023-01-01 12:00:00, Result: success"
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

// 测试TaskLog模型的新方法
func TestTaskLog_LevelChecks(t *testing.T) {
	log := &TaskLog{}

	// 测试IsError
	log.Level = LogLevelError
	if !log.IsError() {
		t.Error("Log should be error level")
	}

	log.Level = LogLevelFatal
	if !log.IsError() {
		t.Error("Log should be error level (fatal)")
	}

	// 测试IsWarning
	log.Level = LogLevelWarn
	if !log.IsWarning() {
		t.Error("Log should be warning level")
	}
}

func TestTaskLog_GetLevelColor(t *testing.T) {
	tests := []struct {
		level    string
		expected string
	}{
		{LogLevelError, "red"},
		{LogLevelFatal, "red"},
		{LogLevelWarn, "orange"},
		{LogLevelInfo, "blue"},
		{LogLevelDebug, "gray"},
		{"unknown", "black"},
	}

	for _, test := range tests {
		log := &TaskLog{Level: test.level}
		result := log.GetLevelColor()
		if result != test.expected {
			t.Errorf("Level %s: expected color %s, got %s", test.level, test.expected, result)
		}
	}
}

func TestTaskLog_DetailOperations(t *testing.T) {
	log := &TaskLog{}

	// 测试添加详情
	log.AddDetail("key1", "value1")
	log.AddDetail("key2", 123)

	// 测试获取详情
	value1, exists1 := log.GetDetail("key1")
	if !exists1 {
		t.Error("Detail key1 should exist")
	}
	if value1 != "value1" {
		t.Errorf("Expected value1, got %v", value1)
	}

	value2, exists2 := log.GetDetail("key2")
	if !exists2 {
		t.Error("Detail key2 should exist")
	}
	if value2 != 123 {
		t.Errorf("Expected 123, got %v", value2)
	}

	// 测试不存在的键
	_, exists3 := log.GetDetail("nonexistent")
	if exists3 {
		t.Error("Nonexistent key should not exist")
	}
}

// 测试TransferState模型的新方法
func TestTransferState_AdditionalMethods(t *testing.T) {
	state := &TransferState{
		TotalSize:       1000,
		TransferredSize: 600,
		Status:          TransferStatusRunning,
		ChunkStates: ChunkStates{
			{Index: 0, Completed: true},
			{Index: 1, Completed: true},
			{Index: 2, Completed: false},
			{Index: 3, Completed: false},
		},
		CreatedAt: time.Now().Add(-time.Minute),
	}

	// 测试GetTotalChunks
	if state.GetTotalChunks() != 4 {
		t.Errorf("Expected 4 total chunks, got %d", state.GetTotalChunks())
	}

	// 测试GetRemainingSize
	if state.GetRemainingSize() != 400 {
		t.Errorf("Expected 400 remaining size, got %d", state.GetRemainingSize())
	}

	// 测试状态检查
	if !state.IsRunning() {
		t.Error("State should be running")
	}

	state.Status = TransferStatusCompleted
	if !state.IsCompleted() {
		t.Error("State should be completed")
	}

	state.Status = TransferStatusFailed
	if !state.IsFailed() {
		t.Error("State should be failed")
	}

	if !state.CanResume() {
		t.Error("Failed state should be resumable")
	}

	state.Status = TransferStatusPaused
	if !state.CanResume() {
		t.Error("Paused state should be resumable")
	}

	// 测试GetNextIncompleteChunk
	nextChunk := state.GetNextIncompleteChunk()
	if nextChunk == nil {
		t.Error("Should find next incomplete chunk")
	}
	if nextChunk.Index != 2 {
		t.Errorf("Expected next chunk index 2, got %d", nextChunk.Index)
	}

	// 测试UpdateChunkState
	state.UpdateChunkState(2, true, "checksum123")
	if !state.ChunkStates[2].Completed {
		t.Error("Chunk 2 should be completed after update")
	}
	if state.ChunkStates[2].Checksum != "checksum123" {
		t.Errorf("Expected checksum 'checksum123', got '%s'", state.ChunkStates[2].Checksum)
	}

	// 测试GetProgressString
	progressStr := state.GetProgressString()
	if progressStr == "" {
		t.Error("Progress string should not be empty")
	}

	// 测试GetTransferSpeed
	speed := state.GetTransferSpeed()
	if speed <= 0 {
		t.Error("Transfer speed should be positive")
	}

	// 测试GetEstimatedTimeRemaining
	eta := state.GetEstimatedTimeRemaining()
	if eta <= 0 {
		t.Error("ETA should be positive")
	}
}