package models

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

// Cluster FastDFS集群配置模型
type Cluster struct {
	ID          string    `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"not null" json:"name"`
	Version     string    `gorm:"not null" json:"version"`
	TrackerAddr string    `gorm:"not null" json:"tracker_addr"`
	TrackerPort int       `gorm:"not null" json:"tracker_port"`
	Username    string    `json:"username,omitempty"`
	Password    string    `json:"password,omitempty"`
	Status      string    `gorm:"default:'active'" json:"status"`
	Description string    `gorm:"type:text" json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// BeforeCreate GORM钩子，创建前生成ID
func (c *Cluster) BeforeCreate(tx *gorm.DB) error {
	if c.ID == "" {
		c.ID = generateID()
	}
	return nil
}

// ClusterStatus 集群状态常量
const (
	ClusterStatusActive   = "active"
	ClusterStatusInactive = "inactive"
	ClusterStatusError    = "error"
)

// GetConnectionString 获取连接字符串
func (c *Cluster) GetConnectionString() string {
	return fmt.Sprintf("%s:%d", c.TrackerAddr, c.TrackerPort)
}

// IsActive 检查集群是否处于活跃状态
func (c *Cluster) IsActive() bool {
	return c.Status == ClusterStatusActive
}

// Validate 验证集群配置
func (c *Cluster) Validate() error {
	if c.Name == "" {
		return fmt.Errorf("cluster name is required")
	}
	if c.TrackerAddr == "" {
		return fmt.Errorf("tracker address is required")
	}
	if c.TrackerPort <= 0 || c.TrackerPort > 65535 {
		return fmt.Errorf("tracker port must be between 1 and 65535")
	}
	if c.Version == "" {
		return fmt.Errorf("cluster version is required")
	}
	return nil
}