package service

import (
	"fmt"
	"sync"
	"time"

	"fastdfs-migration-system/internal/fastdfs"
	"fastdfs-migration-system/internal/models"
	"fastdfs-migration-system/internal/repository"
	"github.com/sirupsen/logrus"
)

// FastDFSService FastDFS服务
type FastDFSService struct {
	clusterManager *fastdfs.ClusterManager
	repo           repository.Repository
	logger         *logrus.Logger
	mu             sync.RWMutex
}

// NewFastDFSService 创建FastDFS服务
func NewFastDFSService(repo repository.Repository, logger *logrus.Logger) *FastDFSService {
	return &FastDFSService{
		clusterManager: fastdfs.NewClusterManager(),
		repo:           repo,
		logger:         logger,
	}
}

// InitializeClusters 初始化集群连接
func (s *FastDFSService) InitializeClusters() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	// 从数据库加载所有活跃的集群
	clusters, err := s.repo.Cluster().GetByStatus(models.ClusterStatusActive)
	if err != nil {
		return fmt.Errorf("failed to load clusters: %w", err)
	}
	
	s.logger.Infof("Initializing %d clusters", len(clusters))
	
	for _, cluster := range clusters {
		err := s.clusterManager.AddCluster(cluster)
		if err != nil {
			s.logger.Errorf("Failed to add cluster %s: %v", cluster.Name, err)
			// 更新集群状态为错误
			s.repo.Cluster().UpdateStatus(cluster.ID, models.ClusterStatusError)
		} else {
			s.logger.Infof("Successfully added cluster %s", cluster.Name)
		}
	}
	
	return nil
}

// AddCluster 添加集群
func (s *FastDFSService) AddCluster(cluster *models.Cluster) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	// 测试连接
	err := fastdfs.TestConnection(cluster)
	if err != nil {
		return fmt.Errorf("connection test failed: %w", err)
	}
	
	// 保存到数据库
	err = s.repo.Cluster().Create(cluster)
	if err != nil {
		return fmt.Errorf("failed to save cluster: %w", err)
	}
	
	// 添加到集群管理器
	err = s.clusterManager.AddCluster(cluster)
	if err != nil {
		// 如果添加失败，从数据库删除
		s.repo.Cluster().Delete(cluster.ID)
		return fmt.Errorf("failed to add cluster to manager: %w", err)
	}
	
	s.logger.Infof("Added cluster %s (%s:%d)", cluster.Name, cluster.TrackerAddr, cluster.TrackerPort)
	return nil
}

// RemoveCluster 移除集群
func (s *FastDFSService) RemoveCluster(clusterID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	// 检查是否有关联的活跃任务
	migrations, err := s.repo.Migration().GetByStatus(models.MigrationStatusRunning)
	if err != nil {
		return fmt.Errorf("failed to check active migrations: %w", err)
	}
	
	for _, migration := range migrations {
		if migration.SourceClusterID == clusterID || migration.TargetClusterID == clusterID {
			return fmt.Errorf("cannot remove cluster: active migration %s is using this cluster", migration.Name)
		}
	}
	
	// 从集群管理器移除
	err = s.clusterManager.RemoveCluster(clusterID)
	if err != nil {
		s.logger.Warnf("Failed to remove cluster from manager: %v", err)
	}
	
	// 从数据库删除
	err = s.repo.Cluster().Delete(clusterID)
	if err != nil {
		return fmt.Errorf("failed to delete cluster from database: %w", err)
	}
	
	s.logger.Infof("Removed cluster %s", clusterID)
	return nil
}

// GetCluster 获取集群信息
func (s *FastDFSService) GetCluster(clusterID string) (*models.Cluster, error) {
	return s.repo.Cluster().GetByID(clusterID)
}

// ListClusters 列出所有集群
func (s *FastDFSService) ListClusters(pagination *models.Pagination) ([]*models.Cluster, error) {
	return s.repo.Cluster().GetAll(pagination)
}

// TestClusterConnection 测试集群连接
func (s *FastDFSService) TestClusterConnection(clusterID string) error {
	cluster, err := s.repo.Cluster().GetByID(clusterID)
	if err != nil {
		return fmt.Errorf("cluster not found: %w", err)
	}
	
	return fastdfs.TestConnection(cluster)
}

// GetClusterClient 获取集群客户端
func (s *FastDFSService) GetClusterClient(clusterID string) (*fastdfs.PooledClient, error) {
	return s.clusterManager.GetClient(clusterID)
}

// ListFiles 列出文件
func (s *FastDFSService) ListFiles(clusterID string, groupName string, startFileName string, limit int) ([]*fastdfs.FileInfo, error) {
	client, err := s.clusterManager.GetClient(clusterID)
	if err != nil {
		return nil, fmt.Errorf("failed to get cluster client: %w", err)
	}
	
	files, err := client.ListFiles(groupName, startFileName, limit)
	if err != nil {
		s.logger.Errorf("Failed to list files from cluster %s: %v", clusterID, err)
		return nil, err
	}
	
	s.logger.Debugf("Listed %d files from cluster %s, group %s", len(files), clusterID, groupName)
	return files, nil
}

// DownloadFile 下载文件
func (s *FastDFSService) DownloadFile(clusterID string, fileID string) ([]byte, error) {
	client, err := s.clusterManager.GetClient(clusterID)
	if err != nil {
		return nil, fmt.Errorf("failed to get cluster client: %w", err)
	}
	
	data, err := client.DownloadFile(fileID)
	if err != nil {
		s.logger.Errorf("Failed to download file %s from cluster %s: %v", fileID, clusterID, err)
		return nil, err
	}
	
	s.logger.Debugf("Downloaded file %s from cluster %s, size: %d bytes", fileID, clusterID, len(data))
	return data, nil
}

// UploadFile 上传文件
func (s *FastDFSService) UploadFile(clusterID string, groupName string, fileName string, data []byte) (string, error) {
	client, err := s.clusterManager.GetClient(clusterID)
	if err != nil {
		return "", fmt.Errorf("failed to get cluster client: %w", err)
	}
	
	fileID, err := client.UploadFile(groupName, fileName, data)
	if err != nil {
		s.logger.Errorf("Failed to upload file %s to cluster %s: %v", fileName, clusterID, err)
		return "", err
	}
	
	s.logger.Debugf("Uploaded file %s to cluster %s, fileID: %s", fileName, clusterID, fileID)
	return fileID, nil
}

// DeleteFile 删除文件
func (s *FastDFSService) DeleteFile(clusterID string, fileID string) error {
	client, err := s.clusterManager.GetClient(clusterID)
	if err != nil {
		return fmt.Errorf("failed to get cluster client: %w", err)
	}
	
	err = client.DeleteFile(fileID)
	if err != nil {
		s.logger.Errorf("Failed to delete file %s from cluster %s: %v", fileID, clusterID, err)
		return err
	}
	
	s.logger.Debugf("Deleted file %s from cluster %s", fileID, clusterID)
	return nil
}

// GetFileInfo 获取文件信息
func (s *FastDFSService) GetFileInfo(clusterID string, fileID string) (*fastdfs.FileInfo, error) {
	client, err := s.clusterManager.GetClient(clusterID)
	if err != nil {
		return nil, fmt.Errorf("failed to get cluster client: %w", err)
	}
	
	fileInfo, err := client.GetFileInfo(fileID)
	if err != nil {
		s.logger.Errorf("Failed to get file info %s from cluster %s: %v", fileID, clusterID, err)
		return nil, err
	}
	
	s.logger.Debugf("Got file info %s from cluster %s", fileID, clusterID)
	return fileInfo, nil
}

// HealthCheck 健康检查
func (s *FastDFSService) HealthCheck() map[string]error {
	return s.clusterManager.HealthCheck()
}

// GetConnectionStats 获取连接统计信息
func (s *FastDFSService) GetConnectionStats() map[string]interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	stats := make(map[string]interface{})
	clusters := s.clusterManager.ListClusters()
	
	for _, cluster := range clusters {
		connection, err := s.clusterManager.GetCluster(cluster.ID)
		if err != nil {
			continue
		}
		stats[cluster.ID] = connection.GetConnectionStats()
	}
	
	return stats
}

// StartHealthCheckRoutine 启动健康检查例程
func (s *FastDFSService) StartHealthCheckRoutine() {
	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()
		
		for range ticker.C {
			results := s.HealthCheck()
			for clusterID, err := range results {
				if err != nil {
					s.logger.Warnf("Cluster %s health check failed: %v", clusterID, err)
					// 更新数据库中的集群状态
					s.repo.Cluster().UpdateStatus(clusterID, models.ClusterStatusError)
				} else {
					// 恢复健康状态
					s.repo.Cluster().UpdateStatus(clusterID, models.ClusterStatusActive)
				}
			}
		}
	}()
}

// Close 关闭服务
func (s *FastDFSService) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	return s.clusterManager.Close()
}