package fastdfs

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"strings"
	"time"
)

// Client FastDFS客户端
type Client struct {
	trackerAddr string
	trackerPort int
	timeout     time.Duration
	conn        net.Conn
}

// NewClient 创建新的FastDFS客户端
func NewClient(trackerAddr string, trackerPort int) *Client {
	return &Client{
		trackerAddr: trackerAddr,
		trackerPort: trackerPort,
		timeout:     30 * time.Second,
	}
}

// Connect 连接到tracker服务器
func (c *Client) Connect() error {
	addr := fmt.Sprintf("%s:%d", c.trackerAddr, c.trackerPort)
	conn, err := net.DialTimeout("tcp", addr, c.timeout)
	if err != nil {
		return fmt.Errorf("failed to connect to tracker %s: %w", addr, err)
	}
	c.conn = conn
	return nil
}

// Close 关闭连接
func (c *Client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// IsConnected 检查连接状态
func (c *Client) IsConnected() bool {
	return c.conn != nil
}

// Ping 测试连接
func (c *Client) Ping() error {
	if !c.IsConnected() {
		return fmt.Errorf("client not connected")
	}
	
	// 发送ping命令
	header := &Header{
		Length: 0,
		Command: FDFS_PROTO_CMD_ACTIVE_TEST,
		Status: 0,
	}
	
	err := c.sendHeader(header)
	if err != nil {
		return fmt.Errorf("failed to send ping: %w", err)
	}
	
	// 接收响应
	respHeader, err := c.receiveHeader()
	if err != nil {
		return fmt.Errorf("failed to receive ping response: %w", err)
	}
	
	if respHeader.Status != 0 {
		return fmt.Errorf("ping failed with status: %d", respHeader.Status)
	}
	
	return nil
}

// GetStorageServer 获取存储服务器信息
func (c *Client) GetStorageServer(groupName string) (*StorageServer, error) {
	if !c.IsConnected() {
		return nil, fmt.Errorf("client not connected")
	}
	
	// 构建请求数据
	data := make([]byte, FDFS_GROUP_NAME_MAX_LEN)
	copy(data, []byte(groupName))
	
	header := &Header{
		Length: int64(len(data)),
		Command: TRACKER_PROTO_CMD_SERVICE_QUERY_STORE_WITHOUT_GROUP_ONE,
		Status: 0,
	}
	
	err := c.sendHeader(header)
	if err != nil {
		return nil, fmt.Errorf("failed to send get storage request: %w", err)
	}
	
	err = c.sendData(data)
	if err != nil {
		return nil, fmt.Errorf("failed to send get storage data: %w", err)
	}
	
	// 接收响应
	respHeader, err := c.receiveHeader()
	if err != nil {
		return nil, fmt.Errorf("failed to receive get storage response: %w", err)
	}
	
	if respHeader.Status != 0 {
		return nil, fmt.Errorf("get storage failed with status: %d", respHeader.Status)
	}
	
	if respHeader.Length < TRACKER_QUERY_STORAGE_STORE_BODY_LEN {
		return nil, fmt.Errorf("invalid response length: %d", respHeader.Length)
	}
	
	respData := make([]byte, respHeader.Length)
	err = c.receiveData(respData)
	if err != nil {
		return nil, fmt.Errorf("failed to receive storage server data: %w", err)
	}
	
	return parseStorageServer(respData)
}

// ListFiles 列出指定组的文件
func (c *Client) ListFiles(groupName string, startFileName string, limit int) ([]*FileInfo, error) {
	storageServer, err := c.GetStorageServer(groupName)
	if err != nil {
		return nil, fmt.Errorf("failed to get storage server: %w", err)
	}
	
	// 连接到存储服务器
	storageClient := NewClient(storageServer.IPAddr, storageServer.Port)
	err = storageClient.Connect()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to storage server: %w", err)
	}
	defer storageClient.Close()
	
	return storageClient.listFilesFromStorage(groupName, startFileName, limit)
}

// DownloadFile 下载文件
func (c *Client) DownloadFile(fileID string) ([]byte, error) {
	groupName, fileName, err := parseFileID(fileID)
	if err != nil {
		return nil, fmt.Errorf("invalid file ID: %w", err)
	}
	
	storageServer, err := c.GetStorageServer(groupName)
	if err != nil {
		return nil, fmt.Errorf("failed to get storage server: %w", err)
	}
	
	// 连接到存储服务器
	storageClient := NewClient(storageServer.IPAddr, storageServer.Port)
	err = storageClient.Connect()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to storage server: %w", err)
	}
	defer storageClient.Close()
	
	return storageClient.downloadFromStorage(groupName, fileName)
}

// UploadFile 上传文件
func (c *Client) UploadFile(groupName string, fileName string, data []byte) (string, error) {
	storageServer, err := c.GetStorageServer(groupName)
	if err != nil {
		return "", fmt.Errorf("failed to get storage server: %w", err)
	}
	
	// 连接到存储服务器
	storageClient := NewClient(storageServer.IPAddr, storageServer.Port)
	err = storageClient.Connect()
	if err != nil {
		return "", fmt.Errorf("failed to connect to storage server: %w", err)
	}
	defer storageClient.Close()
	
	return storageClient.uploadToStorage(groupName, fileName, data)
}

// DeleteFile 删除文件
func (c *Client) DeleteFile(fileID string) error {
	groupName, fileName, err := parseFileID(fileID)
	if err != nil {
		return fmt.Errorf("invalid file ID: %w", err)
	}
	
	storageServer, err := c.GetStorageServer(groupName)
	if err != nil {
		return fmt.Errorf("failed to get storage server: %w", err)
	}
	
	// 连接到存储服务器
	storageClient := NewClient(storageServer.IPAddr, storageServer.Port)
	err = storageClient.Connect()
	if err != nil {
		return fmt.Errorf("failed to connect to storage server: %w", err)
	}
	defer storageClient.Close()
	
	return storageClient.deleteFromStorage(groupName, fileName)
}

// GetFileInfo 获取文件信息
func (c *Client) GetFileInfo(fileID string) (*FileInfo, error) {
	groupName, fileName, err := parseFileID(fileID)
	if err != nil {
		return nil, fmt.Errorf("invalid file ID: %w", err)
	}
	
	storageServer, err := c.GetStorageServer(groupName)
	if err != nil {
		return nil, fmt.Errorf("failed to get storage server: %w", err)
	}
	
	// 连接到存储服务器
	storageClient := NewClient(storageServer.IPAddr, storageServer.Port)
	err = storageClient.Connect()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to storage server: %w", err)
	}
	defer storageClient.Close()
	
	return storageClient.getFileInfoFromStorage(groupName, fileName)
}

// sendHeader 发送协议头
func (c *Client) sendHeader(header *Header) error {
	buf := make([]byte, FDFS_PROTO_PKG_LEN_SIZE+2)
	binary.BigEndian.PutUint64(buf[0:8], uint64(header.Length))
	buf[8] = header.Command
	buf[9] = header.Status
	
	_, err := c.conn.Write(buf)
	return err
}

// receiveHeader 接收协议头
func (c *Client) receiveHeader() (*Header, error) {
	buf := make([]byte, FDFS_PROTO_PKG_LEN_SIZE+2)
	_, err := io.ReadFull(c.conn, buf)
	if err != nil {
		return nil, err
	}
	
	header := &Header{
		Length:  int64(binary.BigEndian.Uint64(buf[0:8])),
		Command: buf[8],
		Status:  buf[9],
	}
	
	return header, nil
}

// sendData 发送数据
func (c *Client) sendData(data []byte) error {
	_, err := c.conn.Write(data)
	return err
}

// receiveData 接收数据
func (c *Client) receiveData(data []byte) error {
	_, err := io.ReadFull(c.conn, data)
	return err
}

// parseFileID 解析文件ID
func parseFileID(fileID string) (string, string, error) {
	parts := strings.SplitN(fileID, "/", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid file ID format: %s", fileID)
	}
	return parts[0], parts[1], nil
}

// parseStorageServer 解析存储服务器响应
func parseStorageServer(data []byte) (*StorageServer, error) {
	if len(data) < TRACKER_QUERY_STORAGE_STORE_BODY_LEN {
		return nil, fmt.Errorf("invalid storage server response length")
	}
	
	server := &StorageServer{
		GroupName: strings.TrimRight(string(data[0:FDFS_GROUP_NAME_MAX_LEN]), "\x00"),
		IPAddr:    strings.TrimRight(string(data[FDFS_GROUP_NAME_MAX_LEN:FDFS_GROUP_NAME_MAX_LEN+IP_ADDRESS_SIZE-1]), "\x00"),
		Port:      int(binary.BigEndian.Uint64(data[FDFS_GROUP_NAME_MAX_LEN+IP_ADDRESS_SIZE-1:])),
	}
	
	return server, nil
}