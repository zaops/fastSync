package fastdfs

import (
	"testing"
	"time"

	"fastdfs-migration-system/internal/models"
)

func TestNewClient(t *testing.T) {
	client := NewClient("127.0.0.1", 22122)
	
	if client.trackerAddr != "127.0.0.1" {
		t.Errorf("Expected tracker addr 127.0.0.1, got %s", client.trackerAddr)
	}
	
	if client.trackerPort != 22122 {
		t.Errorf("Expected tracker port 22122, got %d", client.trackerPort)
	}
	
	if client.timeout != 30*time.Second {
		t.Errorf("Expected timeout 30s, got %v", client.timeout)
	}
}

func TestParseFileID(t *testing.T) {
	tests := []struct {
		fileID    string
		groupName string
		fileName  string
		hasError  bool
	}{
		{"group1/M00/00/00/test.jpg", "group1", "M00/00/00/test.jpg", false},
		{"group2/file.txt", "group2", "file.txt", false},
		{"invalid", "", "", true},
		{"", "", "", true},
	}
	
	for _, test := range tests {
		groupName, fileName, err := parseFileID(test.fileID)
		
		if test.hasError {
			if err == nil {
				t.Errorf("Expected error for fileID %s", test.fileID)
			}
		} else {
			if err != nil {
				t.Errorf("Unexpected error for fileID %s: %v", test.fileID, err)
			}
			
			if groupName != test.groupName {
				t.Errorf("Expected group name %s, got %s", test.groupName, groupName)
			}
			
			if fileName != test.fileName {
				t.Errorf("Expected file name %s, got %s", test.fileName, fileName)
			}
		}
	}
}

func TestGetFileExtension(t *testing.T) {
	tests := []struct {
		fileName string
		expected string
	}{
		{"test.jpg", "jpg"},
		{"document.pdf", "pdf"},
		{"archive.tar.gz", "gz"},
		{"noextension", ""},
		{"", ""},
		{"file.verylongextension", "verylo"}, // 测试长度限制
	}
	
	for _, test := range tests {
		result := getFileExtension(test.fileName)
		if result != test.expected {
			t.Errorf("getFileExtension(%s): expected %s, got %s", test.fileName, test.expected, result)
		}
	}
}

func TestFileInfo_GetFileID(t *testing.T) {
	fileInfo := &FileInfo{
		GroupName: "group1",
		FileName:  "M00/00/00/test.jpg",
	}
	
	expected := "group1/M00/00/00/test.jpg"
	result := fileInfo.GetFileID()
	
	if result != expected {
		t.Errorf("Expected file ID %s, got %s", expected, result)
	}
}

func TestFileInfo_GetCreateTime(t *testing.T) {
	timestamp := int64(1640995200) // 2022-01-01 00:00:00 UTC
	fileInfo := &FileInfo{
		CreateTime: timestamp,
	}
	
	expected := time.Unix(timestamp, 0)
	result := fileInfo.GetCreateTime()
	
	if !result.Equal(expected) {
		t.Errorf("Expected create time %v, got %v", expected, result)
	}
}

func TestConnectionPool(t *testing.T) {
	// 测试连接池创建
	pool := NewConnectionPool("127.0.0.1", 22122, 5)
	
	if pool.Size() != 5 {
		t.Errorf("Expected pool size 5, got %d", pool.Size())
	}
	
	// 测试关闭连接池
	err := pool.Close()
	if err != nil {
		t.Errorf("Failed to close connection pool: %v", err)
	}
}

func TestClusterManager(t *testing.T) {
	manager := NewClusterManager()
	
	// 测试添加集群（这里会失败，因为没有真实的FastDFS服务器）
	cluster := &models.Cluster{
		ID:          "test-cluster",
		Name:        "Test Cluster",
		TrackerAddr: "127.0.0.1",
		TrackerPort: 22122,
	}
	
	err := manager.AddCluster(cluster)
	// 预期会失败，因为没有真实的服务器
	if err == nil {
		t.Log("Cluster added successfully (unexpected in test environment)")
		
		// 如果成功添加，测试获取集群
		_, err = manager.GetCluster("test-cluster")
		if err != nil {
			t.Errorf("Failed to get cluster: %v", err)
		}
		
		// 测试移除集群
		err = manager.RemoveCluster("test-cluster")
		if err != nil {
			t.Errorf("Failed to remove cluster: %v", err)
		}
	} else {
		t.Logf("Expected connection failure in test environment: %v", err)
	}
	
	// 测试关闭管理器
	err = manager.Close()
	if err != nil {
		t.Errorf("Failed to close cluster manager: %v", err)
	}
}

func TestTestConnection(t *testing.T) {
	cluster := &models.Cluster{
		ID:          "test-cluster",
		Name:        "Test Cluster",
		TrackerAddr: "127.0.0.1",
		TrackerPort: 22122,
	}
	
	err := TestConnection(cluster)
	// 预期会失败，因为没有真实的FastDFS服务器
	if err == nil {
		t.Log("Connection test passed (unexpected in test environment)")
	} else {
		t.Logf("Expected connection failure in test environment: %v", err)
	}
}

func TestProtocolConstants(t *testing.T) {
	// 测试协议常量
	if FDFS_PROTO_PKG_LEN_SIZE != 8 {
		t.Errorf("Expected FDFS_PROTO_PKG_LEN_SIZE to be 8, got %d", FDFS_PROTO_PKG_LEN_SIZE)
	}
	
	if FDFS_GROUP_NAME_MAX_LEN != 16 {
		t.Errorf("Expected FDFS_GROUP_NAME_MAX_LEN to be 16, got %d", FDFS_GROUP_NAME_MAX_LEN)
	}
	
	if IP_ADDRESS_SIZE != 16 {
		t.Errorf("Expected IP_ADDRESS_SIZE to be 16, got %d", IP_ADDRESS_SIZE)
	}
}

func TestHeader(t *testing.T) {
	header := &Header{
		Length:  1024,
		Command: FDFS_PROTO_CMD_ACTIVE_TEST,
		Status:  FDFS_PROTO_STATUS_SUCCESS,
	}
	
	if header.Length != 1024 {
		t.Errorf("Expected length 1024, got %d", header.Length)
	}
	
	if header.Command != FDFS_PROTO_CMD_ACTIVE_TEST {
		t.Errorf("Expected command %d, got %d", FDFS_PROTO_CMD_ACTIVE_TEST, header.Command)
	}
	
	if header.Status != FDFS_PROTO_STATUS_SUCCESS {
		t.Errorf("Expected status %d, got %d", FDFS_PROTO_STATUS_SUCCESS, header.Status)
	}
}