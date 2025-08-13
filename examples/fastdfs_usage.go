package main

import (
	"fmt"
	"log"

	"fastdfs-migration-system/internal/database"
	"fastdfs-migration-system/internal/models"
	"fastdfs-migration-system/internal/repository"
	"fastdfs-migration-system/internal/service"
	"github.com/sirupsen/logrus"
)

func main() {
	// 初始化日志
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)

	// 初始化数据库
	config := database.Config{
		Type:     "sqlite",
		DSN:      "fastdfs_example.db",
		LogLevel: "info",
	}

	err := database.Initialize(config)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	// 自动迁移数据库表
	err = database.AutoMigrate()
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// 创建Repository实例
	repo := repository.NewRepository(database.GetDB())

	// 创建FastDFS服务
	fastdfsService := service.NewFastDFSService(repo, logger)
	defer fastdfsService.Close()

	// 示例：添加源集群配置
	sourceCluster := &models.Cluster{
		Name:        "FastDFS 5.0.7 源集群",
		Version:     "5.0.7",
		TrackerAddr: "192.168.1.100",
		TrackerPort: 22122,
		Status:      models.ClusterStatusActive,
		Description: "源集群，用于迁移数据",
	}

	err = fastdfsService.AddCluster(sourceCluster)
	if err != nil {
		log.Printf("Failed to add source cluster (expected in demo): %v", err)
	} else {
		fmt.Printf("Added source cluster: %s\n", sourceCluster.Name)
	}

	// 示例：添加目标集群配置
	targetCluster := &models.Cluster{
		Name:        "FastDFS 6.0.6 目标集群",
		Version:     "6.0.6",
		TrackerAddr: "192.168.1.200",
		TrackerPort: 22122,
		Status:      models.ClusterStatusActive,
		Description: "目标集群，用于接收迁移数据",
	}

	err = fastdfsService.AddCluster(targetCluster)
	if err != nil {
		log.Printf("Failed to add target cluster (expected in demo): %v", err)
	} else {
		fmt.Printf("Added target cluster: %s\n", targetCluster.Name)
	}

	// 示例：列出所有集群
	pagination := &models.Pagination{
		Page:     1,
		PageSize: 10,
	}

	clusters, err := fastdfsService.ListClusters(pagination)
	if err != nil {
		log.Printf("Failed to list clusters: %v", err)
	} else {
		fmt.Printf("Found %d clusters:\n", len(clusters))
		for _, cluster := range clusters {
			fmt.Printf("  - %s (%s:%d) - %s\n", 
				cluster.Name, cluster.TrackerAddr, cluster.TrackerPort, cluster.Status)
		}
	}

	// 示例：测试集群连接
	for _, cluster := range clusters {
		err = fastdfsService.TestClusterConnection(cluster.ID)
		if err != nil {
			fmt.Printf("Cluster %s connection test failed (expected): %v\n", cluster.Name, err)
		} else {
			fmt.Printf("Cluster %s connection test passed\n", cluster.Name)
		}
	}

	// 示例：健康检查
	healthResults := fastdfsService.HealthCheck()
	fmt.Printf("Health check results:\n")
	for clusterID, err := range healthResults {
		if err != nil {
			fmt.Printf("  - Cluster %s: UNHEALTHY (%v)\n", clusterID, err)
		} else {
			fmt.Printf("  - Cluster %s: HEALTHY\n", clusterID)
		}
	}

	// 示例：获取连接统计信息
	stats := fastdfsService.GetConnectionStats()
	fmt.Printf("Connection statistics:\n")
	for clusterID, stat := range stats {
		if statMap, ok := stat.(map[string]interface{}); ok {
			fmt.Printf("  - Cluster %s:\n", clusterID)
			fmt.Printf("    Name: %v\n", statMap["cluster_name"])
			fmt.Printf("    Address: %v:%v\n", statMap["tracker_addr"], statMap["tracker_port"])
			fmt.Printf("    Healthy: %v\n", statMap["is_healthy"])
			fmt.Printf("    Pool Size: %v\n", statMap["pool_size"])
			fmt.Printf("    Available Connections: %v\n", statMap["available_conn"])
		}
	}

	// 示例：模拟文件操作（这些操作会失败，因为没有真实的FastDFS服务器）
	if len(clusters) > 0 {
		clusterID := clusters[0].ID
		
		// 尝试列出文件
		files, err := fastdfsService.ListFiles(clusterID, "group1", "", 10)
		if err != nil {
			fmt.Printf("List files failed (expected): %v\n", err)
		} else {
			fmt.Printf("Found %d files\n", len(files))
		}
		
		// 尝试获取文件信息
		fileInfo, err := fastdfsService.GetFileInfo(clusterID, "group1/M00/00/00/test.jpg")
		if err != nil {
			fmt.Printf("Get file info failed (expected): %v\n", err)
		} else {
			fmt.Printf("File info: %+v\n", fileInfo)
		}
		
		// 尝试上传文件
		testData := []byte("Hello, FastDFS!")
		fileID, err := fastdfsService.UploadFile(clusterID, "group1", "test.txt", testData)
		if err != nil {
			fmt.Printf("Upload file failed (expected): %v\n", err)
		} else {
			fmt.Printf("Uploaded file: %s\n", fileID)
			
			// 尝试下载文件
			data, err := fastdfsService.DownloadFile(clusterID, fileID)
			if err != nil {
				fmt.Printf("Download file failed: %v\n", err)
			} else {
				fmt.Printf("Downloaded file data: %s\n", string(data))
			}
			
			// 尝试删除文件
			err = fastdfsService.DeleteFile(clusterID, fileID)
			if err != nil {
				fmt.Printf("Delete file failed: %v\n", err)
			} else {
				fmt.Printf("Deleted file: %s\n", fileID)
			}
		}
	}

	fmt.Println("FastDFS service usage example completed!")
}