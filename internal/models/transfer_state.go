package models

import (
	"encoding/json"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// TransferState 传输状态模型，用于断点续传
type TransferState struct {
	ID               string       `gorm:"primaryKey" json:"id"`
	TaskID           string       `gorm:"index" json:"task_id"`
	FileID           string       `gorm:"not null" json:"file_id"`
	FilePath         string       `gorm:"not null" json:"file_path"`
	TotalSize        int64        `gorm:"not null" json:"total_size"`
	TransferredSize  int64        `gorm:"default:0" json:"transferred_size"`
	ChunkSize        int64        `gorm:"not null" json:"chunk_size"`
	ChunkStates      ChunkStates `gorm:"type:json" json:"chunk_states"`
	Status           string       `gorm:"default:'pending'" json:"status"`
	Checksum         string       `json:"checksum,omitempty"`
	LastUpdate       time.Time    `json:"last_update"`
	CreatedAt        time.Time    `json:"created_at"`
}

// ChunkState 分块状态
type ChunkState struct {
	Index     int    `json:"index"`
	Offset    int64  `json:"offset"`
	Size      int64  `json:"size"`
	Completed bool   `json:"completed"`
	Checksum  string `json:"checksum,omitempty"`
}

// BeforeCreate GORM钩子，创建前生成ID
func (ts *TransferState) BeforeCreate(tx *gorm.DB) error {
	if ts.ID == "" {
		ts.ID = generateID()
	}
	return nil
}

// BeforeUpdate GORM钩子，更新前设置LastUpdate时间
func (ts *TransferState) BeforeUpdate(tx *gorm.DB) error {
	ts.LastUpdate = time.Now()
	return nil
}

// ChunkStates 分块状态列表类型
type ChunkStates []ChunkState

// 实现GORM的Valuer和Scanner接口，用于JSON字段的序列化
func (cs ChunkStates) Value() (interface{}, error) {
	return json.Marshal(cs)
}

func (cs *ChunkStates) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	
	return json.Unmarshal(bytes, cs)
}

// TransferStatus 传输状态常量
const (
	TransferStatusPending    = "pending"
	TransferStatusRunning    = "running"
	TransferStatusPaused     = "paused"
	TransferStatusCompleted  = "completed"
	TransferStatusFailed     = "failed"
)

// GetProgress 计算传输进度百分比
func (ts *TransferState) GetProgress() float64 {
	if ts.TotalSize == 0 {
		return 0
	}
	return float64(ts.TransferredSize) / float64(ts.TotalSize) * 100
}

// GetCompletedChunks 获取已完成的分块数量
func (ts *TransferState) GetCompletedChunks() int {
	completed := 0
	for _, chunk := range []ChunkState(ts.ChunkStates) {
		if chunk.Completed {
			completed++
		}
	}
	return completed
}

// GetTotalChunks 获取总分块数量
func (ts *TransferState) GetTotalChunks() int {
	return len(ts.ChunkStates)
}

// GetRemainingSize 获取剩余传输大小
func (ts *TransferState) GetRemainingSize() int64 {
	return ts.TotalSize - ts.TransferredSize
}

// IsCompleted 检查传输是否完成
func (ts *TransferState) IsCompleted() bool {
	return ts.Status == TransferStatusCompleted
}

// IsFailed 检查传输是否失败
func (ts *TransferState) IsFailed() bool {
	return ts.Status == TransferStatusFailed
}

// IsRunning 检查传输是否正在进行
func (ts *TransferState) IsRunning() bool {
	return ts.Status == TransferStatusRunning
}

// CanResume 检查是否可以恢复传输
func (ts *TransferState) CanResume() bool {
	return ts.Status == TransferStatusPaused || ts.Status == TransferStatusFailed
}

// GetNextIncompleteChunk 获取下一个未完成的分块
func (ts *TransferState) GetNextIncompleteChunk() *ChunkState {
	for i, chunk := range ts.ChunkStates {
		if !chunk.Completed {
			return &ts.ChunkStates[i]
		}
	}
	return nil
}

// UpdateChunkState 更新指定分块的状态
func (ts *TransferState) UpdateChunkState(index int, completed bool, checksum string) {
	if index >= 0 && index < len(ts.ChunkStates) {
		ts.ChunkStates[index].Completed = completed
		if checksum != "" {
			ts.ChunkStates[index].Checksum = checksum
		}
	}
}

// GetProgressString 获取进度字符串
func (ts *TransferState) GetProgressString() string {
	return fmt.Sprintf("%.2f%% (%d/%d chunks)", 
		ts.GetProgress(), ts.GetCompletedChunks(), ts.GetTotalChunks())
}

// GetTransferSpeed 计算传输速度（字节/秒）
func (ts *TransferState) GetTransferSpeed() float64 {
	duration := time.Since(ts.CreatedAt).Seconds()
	if duration <= 0 {
		return 0
	}
	return float64(ts.TransferredSize) / duration
}

// GetEstimatedTimeRemaining 估算剩余时间
func (ts *TransferState) GetEstimatedTimeRemaining() time.Duration {
	speed := ts.GetTransferSpeed()
	if speed <= 0 {
		return 0
	}
	remainingSeconds := float64(ts.GetRemainingSize()) / speed
	return time.Duration(remainingSeconds) * time.Second
}