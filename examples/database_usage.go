package main

import (
	"fmt"
	"log"
	"time"

	"fastdfs-migration-system/internal/database"
	"fastdfs-migration-system/internal/models"
	"fastdfs-migration-system/internal/repository"
)

func main() {
	// 初始化数据库
	config := database.Config{
		Type:     "sqlite",
		DSN:      "example.db",
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

	// 示例：创建集群配置
	cluster := &models.Cluster{
		Name:        "FastDFS 5.0.7 源集群",
		Version:     "5.0.7",
		TrackerAddr: "192.168.1.100",
		TrackerPort: 22122,
		Status:      models.ClusterStatusActive,
		Description: "源集群配置",
	}

	err = repo.Cluster().Create(cluster)
	if err != nil {
		log.Printf("Failed to create cluster: %v", err)
	} else {
		fmt.Printf("Created cluster: %s (ID: %s)\n", cluster.Name, cluster.ID)
	}

	// 示例：创建迁移任务
	migration := &models.Migration{
		Name:            "测试迁移任务",
		SourceClusterID: cluster.ID,
		TargetClusterID: "target-cluster-id",
		Status:          models.MigrationStatusPending,
		Config: models.MigrationConfig{
			ConcurrentWorkers:   5,
			IncrementalSync:     false,
			VerificationEnabled: true,
			RetryConfig: &models.RetryConfig{
				MaxRetries:    3,
				RetryInterval: time.Minute,
				BackoffFactor: 2.0,
			},
		},
	}

	err = repo.Migration().Create(migration)
	if err != nil {
		log.Printf("Failed to create migration: %v", err)
	} else {
		fmt.Printf("Created migration: %s (ID: %s)\n", migration.Name, migration.ID)
	}

	// 示例：创建任务日志
	taskLog := &models.TaskLog{
		TaskID:   migration.ID,
		TaskType: models.TaskTypeMigration,
		Level:    models.LogLevelInfo,
		Message:  "迁移任务已创建",
		Details: models.LogDetails{
			"source_cluster": cluster.Name,
			"status":         migration.Status,
		},
	}

	err = repo.TaskLog().Create(taskLog)
	if err != nil {
		log.Printf("Failed to create task log: %v", err)
	} else {
		fmt.Printf("Created task log: %s (ID: %s)\n", taskLog.Message, taskLog.ID)
	}

	// 示例：查询数据
	pagination := &models.Pagination{
		Page:     1,
		PageSize: 10,
	}

	migrations, err := repo.Migration().GetAll(pagination)
	if err != nil {
		log.Printf("Failed to get migrations: %v", err)
	} else {
		fmt.Printf("Found %d migrations (total: %d)\n", len(migrations), pagination.Total)
		for _, m := range migrations {
			fmt.Printf("  - %s: %s\n", m.Name, m.Status)
		}
	}

	// 示例：更新迁移进度
	err = repo.Migration().UpdateProgress(migration.ID, 25.0, 100, 1024*1024)
	if err != nil {
		log.Printf("Failed to update migration progress: %v", err)
	} else {
		fmt.Printf("Updated migration progress to 25%%\n")
	}

	fmt.Println("Database usage example completed successfully!")
}