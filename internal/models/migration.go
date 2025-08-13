package models

import (
	"encoding/json"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// Migration 迁移任务模型
type Migration struct {
	ID              string          `gorm:"primaryKey" json:"id"`
	Name            string          `gorm:"not null" json:"name"`
	SourceClusterID string          `gorm:"not null" json:"source_cluster_id"`
	TargetClusterID string          `gorm:"not null" json:"target_cluster_id"`
	Config          MigrationConfig `gorm:"type:json" json:"config"`
	Status          string          `gorm:"default:'pending'" json:"status"`
	Progress        float64         `gorm:"default:0" json:"progress"`
	TotalFiles      int64           `gorm:"default:0" json:"total_files"`
	ProcessedFiles  int64           `gorm:"default:0" json:"processed_files"`
	TotalSize       int64           `gorm:"default:0" json:"total_size"`
	ProcessedSize   int64           `gorm:"default:0" json:"processed_size"`
	ErrorMessage    string          `gorm:"type:text" json:"error_message,omitempty"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
	CompletedAt     *time.Time      `json:"completed_at,omitempty"`
}

// MigrationConfig 迁移配置
type MigrationConfig struct {
	// 时间过滤
	TimeFilter *TimeFilter `json:"time_filter,omitempty"`
	
	// 文件类型过滤
	FileTypeFilter *FileTypeFilter `json:"file_type_filter,omitempty"`
	
	// 增量同步配置
	IncrementalSync bool `json:"incremental_sync"`
	
	// 并发配置
	ConcurrentWorkers int `json:"concurrent_workers"`
	
	// 重试配置
	RetryConfig *RetryConfig `json:"retry_config"`
	
	// 验证配置
	VerificationEnabled bool `json:"verification_enabled"`
}

// TimeFilter 时间过滤器
type TimeFilter struct {
	StartTime *time.Time `json:"start_time,omitempty"`
	EndTime   *time.Time `json:"end_time,omitempty"`
}

// FileTypeFilter 文件类型过滤器
type FileTypeFilter struct {
	IncludeExtensions []string `json:"include_extensions,omitempty"`
	ExcludeExtensions []string `json:"exclude_extensions,omitempty"`
	IncludeMimeTypes  []string `json:"include_mime_types,omitempty"`
	ExcludeMimeTypes  []string `json:"exclude_mime_types,omitempty"`
}

// RetryConfig 重试配置
type RetryConfig struct {
	MaxRetries    int           `json:"max_retries"`
	RetryInterval time.Duration `json:"retry_interval"`
	BackoffFactor float64       `json:"backoff_factor"`
}

// 实现GORM的Valuer和Scanner接口，用于JSON字段的序列化
func (mc MigrationConfig) Value() (interface{}, error) {
	return json.Marshal(mc)
}

func (mc *MigrationConfig) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	
	return json.Unmarshal(bytes, mc)
}

// BeforeCreate GORM钩子，创建前生成ID
func (m *Migration) BeforeCreate(tx *gorm.DB) error {
	if m.ID == "" {
		m.ID = generateID()
	}
	return nil
}

// MigrationStatus 迁移状态常量
const (
	MigrationStatusPending    = "pending"
	MigrationStatusRunning    = "running"
	MigrationStatusPaused     = "paused"
	MigrationStatusCompleted  = "completed"
	MigrationStatusFailed     = "failed"
	MigrationStatusCancelled  = "cancelled"
)

// IsRunning 检查迁移是否正在运行
func (m *Migration) IsRunning() bool {
	return m.Status == MigrationStatusRunning
}

// IsCompleted 检查迁移是否已完成
func (m *Migration) IsCompleted() bool {
	return m.Status == MigrationStatusCompleted
}

// IsFailed 检查迁移是否失败
func (m *Migration) IsFailed() bool {
	return m.Status == MigrationStatusFailed
}

// CanStart 检查迁移是否可以启动
func (m *Migration) CanStart() bool {
	return m.Status == MigrationStatusPending || m.Status == MigrationStatusPaused
}

// CanPause 检查迁移是否可以暂停
func (m *Migration) CanPause() bool {
	return m.Status == MigrationStatusRunning
}

// CanResume 检查迁移是否可以恢复
func (m *Migration) CanResume() bool {
	return m.Status == MigrationStatusPaused
}

// GetProgressPercentage 获取进度百分比字符串
func (m *Migration) GetProgressPercentage() string {
	return fmt.Sprintf("%.2f%%", m.Progress)
}

// GetProcessedFilesRatio 获取已处理文件比例
func (m *Migration) GetProcessedFilesRatio() float64 {
	if m.TotalFiles == 0 {
		return 0
	}
	return float64(m.ProcessedFiles) / float64(m.TotalFiles)
}

// GetProcessedSizeRatio 获取已处理大小比例
func (m *Migration) GetProcessedSizeRatio() float64 {
	if m.TotalSize == 0 {
		return 0
	}
	return float64(m.ProcessedSize) / float64(m.TotalSize)
}

// Validate 验证迁移配置
func (m *Migration) Validate() error {
	if m.Name == "" {
		return fmt.Errorf("migration name is required")
	}
	if m.SourceClusterID == "" {
		return fmt.Errorf("source cluster ID is required")
	}
	if m.TargetClusterID == "" {
		return fmt.Errorf("target cluster ID is required")
	}
	if m.SourceClusterID == m.TargetClusterID {
		return fmt.Errorf("source and target cluster cannot be the same")
	}
	return nil
}