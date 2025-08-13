package fastdfs

import (
	"fmt"
	"sync"
	"time"

	"fastdfs-migration-system/internal/models"
)

// ClusterManager 集群连接管理器
type ClusterManager struct {
	clusters map[string]*ClusterConnection
	mu       sync.RWMutex
}

// ClusterConnection 集群连接
type ClusterConnection struct {
	cluster    *models.Cluster
	pool       ConnectionPool
	client     *PooledClient
	lastCheck  time.Time
	isHealthy  bool
	mu         sync.RWMutex
}

// NewClusterManager 创建集群管理器
func NewClusterManager() *ClusterManager {
	return &ClusterManager{
		clusters: make(map[string]*ClusterConnection),
	}
}

// AddCluster 添加集群
func (cm *ClusterManager) AddCluster(cluster *models.Cluster) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	
	// 创建连接池
	pool := NewConnectionPool(cluster.TrackerAddr, cluster.TrackerPort, 10)
	
	// 创建带连接池的客户端
	client := NewPooledClient(pool)
	
	// 测试连接
	err := client.Ping()
	if err != nil {
		pool.Close()
		return fmt.Errorf("failed to connect to cluster %s: %w", cluster.Name, err)
	}
	
	connection := &ClusterConnection{
		cluster:   cluster,
		pool:      pool,
		client:    client,
		lastCheck: time.Now(),
		isHealthy: true,
	}
	
	cm.clusters[cluster.ID] = connection
	return nil
}

// RemoveCluster 移除集群
func (cm *ClusterManager) RemoveCluster(clusterID string) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	
	connection, exists := cm.clusters[clusterID]
	if !exists {
		return fmt.Errorf("cluster %s not found", clusterID)
	}
	
	// 关闭连接池
	connection.pool.Close()
	delete(cm.clusters, clusterID)
	
	return nil
}

// GetCluster 获取集群连接
func (cm *ClusterManager) GetCluster(clusterID string) (*ClusterConnection, error) {
	cm.mu.RLock()
	connection, exists := cm.clusters[clusterID]
	cm.mu.RUnlock()
	
	if !exists {
		return nil, fmt.Errorf("cluster %s not found", clusterID)
	}
	
	// 检查健康状态
	if err := connection.checkHealth(); err != nil {
		return nil, fmt.Errorf("cluster %s is unhealthy: %w", clusterID, err)
	}
	
	return connection, nil
}

// GetClient 获取集群客户端
func (cm *ClusterManager) GetClient(clusterID string) (*PooledClient, error) {
	connection, err := cm.GetCluster(clusterID)
	if err != nil {
		return nil, err
	}
	return connection.client, nil
}

// ListClusters 列出所有集群
func (cm *ClusterManager) ListClusters() []*models.Cluster {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	
	clusters := make([]*models.Cluster, 0, len(cm.clusters))
	for _, connection := range cm.clusters {
		clusters = append(clusters, connection.cluster)
	}
	
	return clusters
}

// HealthCheck 健康检查
func (cm *ClusterManager) HealthCheck() map[string]error {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	
	results := make(map[string]error)
	
	for clusterID, connection := range cm.clusters {
		err := connection.checkHealth()
		if err != nil {
			results[clusterID] = err
		}
	}
	
	return results
}

// Close 关闭所有连接
func (cm *ClusterManager) Close() error {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	
	for _, connection := range cm.clusters {
		connection.pool.Close()
	}
	
	cm.clusters = make(map[string]*ClusterConnection)
	return nil
}

// checkHealth 检查集群健康状态
func (cc *ClusterConnection) checkHealth() error {
	cc.mu.Lock()
	defer cc.mu.Unlock()
	
	// 如果最近检查过且健康，直接返回
	if time.Since(cc.lastCheck) < 30*time.Second && cc.isHealthy {
		return nil
	}
	
	// 执行健康检查
	err := cc.client.Ping()
	cc.lastCheck = time.Now()
	cc.isHealthy = (err == nil)
	
	return err
}

// GetCluster 获取集群信息
func (cc *ClusterConnection) GetCluster() *models.Cluster {
	cc.mu.RLock()
	defer cc.mu.RUnlock()
	return cc.cluster
}

// GetClient 获取客户端
func (cc *ClusterConnection) GetClient() *PooledClient {
	return cc.client
}

// IsHealthy 检查是否健康
func (cc *ClusterConnection) IsHealthy() bool {
	cc.mu.RLock()
	defer cc.mu.RUnlock()
	return cc.isHealthy
}

// GetConnectionStats 获取连接统计信息
func (cc *ClusterConnection) GetConnectionStats() map[string]interface{} {
	cc.mu.RLock()
	defer cc.mu.RUnlock()
	
	return map[string]interface{}{
		"cluster_id":     cc.cluster.ID,
		"cluster_name":   cc.cluster.Name,
		"tracker_addr":   cc.cluster.TrackerAddr,
		"tracker_port":   cc.cluster.TrackerPort,
		"is_healthy":     cc.isHealthy,
		"last_check":     cc.lastCheck,
		"pool_size":      cc.pool.Size(),
		"available_conn": cc.pool.Available(),
	}
}

// TestConnection 测试集群连接
func TestConnection(cluster *models.Cluster) error {
	client := NewClient(cluster.TrackerAddr, cluster.TrackerPort)
	err := client.Connect()
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	defer client.Close()
	
	err = client.Ping()
	if err != nil {
		return fmt.Errorf("ping failed: %w", err)
	}
	
	return nil
}

// GetClusterInfo 获取集群详细信息
func GetClusterInfo(cluster *models.Cluster) (*GroupInfo, error) {
	client := NewClient(cluster.TrackerAddr, cluster.TrackerPort)
	err := client.Connect()
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}
	defer client.Close()
	
	// 这里可以实现获取集群详细信息的逻辑
	// 由于FastDFS协议复杂，这里返回一个基本的信息结构
	return &GroupInfo{
		GroupName:      "group1", // 默认组名，实际应该从服务器获取
		TotalMB:        0,
		FreeMB:         0,
		StorageCount:   1,
		StoragePort:    23000,
		ActiveCount:    1,
		StorePathCount: 1,
	}, nil
}