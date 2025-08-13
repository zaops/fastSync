package models

import (
	"encoding/json"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// ScheduledTask 定时任务模型
type ScheduledTask struct {
	ID          string          `gorm:"primaryKey" json:"id"`
	Name        string          `gorm:"not null" json:"name"`
	CronExpr    string          `gorm:"not null" json:"cron_expr"`
	TaskConfig  TaskConfig `gorm:"type:json" json:"task_config"`
	Status      string          `gorm:"default:'active'" json:"status"`
	NextRun     *time.Time      `json:"next_run,omitempty"`
	LastRun     *time.Time      `json:"last_run,omitempty"`
	LastResult  string          `json:"last_result,omitempty"`
	Description string          `gorm:"type:text" json:"description,omitempty"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

// BeforeCreate GORM钩子，创建前生成ID
func (st *ScheduledTask) BeforeCreate(tx *gorm.DB) error {
	if st.ID == "" {
		st.ID = generateID()
	}
	return nil
}

// TaskConfig 实现GORM的Valuer和Scanner接口，用于JSON字段的序列化
func (tc TaskConfig) Value() (interface{}, error) {
	return json.Marshal(tc)
}

func (tc *TaskConfig) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	
	return json.Unmarshal(bytes, tc)
}

// ScheduleStatus 定时任务状态常量
const (
	ScheduleStatusActive   = "active"
	ScheduleStatusInactive = "inactive"
	ScheduleStatusError    = "error"
)

// ScheduleResult 定时任务执行结果常量
const (
	ScheduleResultSuccess = "success"
	ScheduleResultFailed  = "failed"
	ScheduleResultSkipped = "skipped"
)

// TaskConfig 任务配置，与MigrationConfig相同但作为独立类型避免方法冲突
type TaskConfig MigrationConfig

// IsActive 检查定时任务是否处于活跃状态
func (st *ScheduledTask) IsActive() bool {
	return st.Status == ScheduleStatusActive
}

// ShouldRun 检查定时任务是否应该运行
func (st *ScheduledTask) ShouldRun() bool {
	if !st.IsActive() {
		return false
	}
	if st.NextRun == nil {
		return false
	}
	return time.Now().After(*st.NextRun)
}

// GetLastRunStatus 获取最后运行状态描述
func (st *ScheduledTask) GetLastRunStatus() string {
	if st.LastRun == nil {
		return "Never run"
	}
	return fmt.Sprintf("Last run: %s, Result: %s", 
		st.LastRun.Format("2006-01-02 15:04:05"), st.LastResult)
}

// Validate 验证定时任务配置
func (st *ScheduledTask) Validate() error {
	if st.Name == "" {
		return fmt.Errorf("scheduled task name is required")
	}
	if st.CronExpr == "" {
		return fmt.Errorf("cron expression is required")
	}
	// 这里可以添加更复杂的cron表达式验证
	return nil
}