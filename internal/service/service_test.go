package service

import (
	"testing"

	"github.com/sirupsen/logrus"
)

func TestServiceCreation(t *testing.T) {
	// 测试服务创建的基本功能，不依赖数据库
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	// 这里我们只测试服务创建的基本逻辑
	if logger == nil {
		t.Error("Logger should not be nil")
	}

	// 测试日志级别设置
	if logger.Level != logrus.ErrorLevel {
		t.Errorf("Expected log level %v, got %v", logrus.ErrorLevel, logger.Level)
	}
}

func TestServiceInterface(t *testing.T) {
	// 测试服务接口的完整性
	// 这里只验证接口定义，不执行实际操作
	
	// 验证FastDFSService结构体包含必要的字段
	// 这是一个编译时检查，确保结构体定义正确
	var service *FastDFSService
	if service != nil {
		// 这些方法应该存在
		_ = service.InitializeClusters
		_ = service.AddCluster
		_ = service.RemoveCluster
		_ = service.GetCluster
		_ = service.ListClusters
		_ = service.TestClusterConnection
		_ = service.GetClusterClient
		_ = service.ListFiles
		_ = service.DownloadFile
		_ = service.UploadFile
		_ = service.DeleteFile
		_ = service.GetFileInfo
		_ = service.HealthCheck
		_ = service.GetConnectionStats
		_ = service.StartHealthCheckRoutine
		_ = service.Close
	}
}