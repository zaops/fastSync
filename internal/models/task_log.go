package models

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

// TaskLog 任务日志模型
type TaskLog struct {
	ID        string                 `gorm:"primaryKey" json:"id"`
	TaskID    string                 `gorm:"index" json:"task_id"`
	TaskType  string                 `gorm:"not null" json:"task_type"` // migration, schedule
	Level     string                 `gorm:"not null" json:"level"`
	Message   string                 `gorm:"type:text" json:"message"`
	Details   LogDetails `gorm:"type:json" json:"details,omitempty"`
	CreatedAt time.Time              `json:"created_at"`
}

// BeforeCreate GORM钩子，创建前生成ID
func (tl *TaskLog) BeforeCreate(tx *gorm.DB) error {
	if tl.ID == "" {
		tl.ID = generateID()
	}
	return nil
}

// LogDetails 日志详情类型
type LogDetails map[string]interface{}

// 实现GORM的Valuer和Scanner接口，用于JSON字段的序列化
func (ld LogDetails) Value() (interface{}, error) {
	return json.Marshal(ld)
}

func (ld *LogDetails) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	
	return json.Unmarshal(bytes, ld)
}

// LogLevel 日志级别常量
const (
	LogLevelDebug = "debug"
	LogLevelInfo  = "info"
	LogLevelWarn  = "warn"
	LogLevelError = "error"
	LogLevelFatal = "fatal"
)

// TaskType 任务类型常量
const (
	TaskTypeMigration = "migration"
	TaskTypeSchedule  = "schedule"
	TaskTypeSystem    = "system"
)

// IsError 检查是否为错误日志
func (tl *TaskLog) IsError() bool {
	return tl.Level == LogLevelError || tl.Level == LogLevelFatal
}

// IsWarning 检查是否为警告日志
func (tl *TaskLog) IsWarning() bool {
	return tl.Level == LogLevelWarn
}

// GetFormattedTime 获取格式化的时间字符串
func (tl *TaskLog) GetFormattedTime() string {
	return tl.CreatedAt.Format("2006-01-02 15:04:05")
}

// GetLevelColor 获取日志级别对应的颜色（用于前端显示）
func (tl *TaskLog) GetLevelColor() string {
	switch tl.Level {
	case LogLevelError, LogLevelFatal:
		return "red"
	case LogLevelWarn:
		return "orange"
	case LogLevelInfo:
		return "blue"
	case LogLevelDebug:
		return "gray"
	default:
		return "black"
	}
}

// AddDetail 添加详情信息
func (tl *TaskLog) AddDetail(key string, value interface{}) {
	if tl.Details == nil {
		tl.Details = make(LogDetails)
	}
	tl.Details[key] = value
}

// GetDetail 获取详情信息
func (tl *TaskLog) GetDetail(key string) (interface{}, bool) {
	if tl.Details == nil {
		return nil, false
	}
	value, exists := tl.Details[key]
	return value, exists
}