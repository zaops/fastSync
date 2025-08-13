package fastdfs

import (
	"fmt"
	"sync"
)

// connectionPool 连接池实现
type connectionPool struct {
	trackerAddr    string
	trackerPort    int
	maxConnections int
	connections    chan *Client
	mu             sync.RWMutex
	closed         bool
}

// NewConnectionPool 创建新的连接池
func NewConnectionPool(trackerAddr string, trackerPort int, maxConnections int) ConnectionPool {
	pool := &connectionPool{
		trackerAddr:    trackerAddr,
		trackerPort:    trackerPort,
		maxConnections: maxConnections,
		connections:    make(chan *Client, maxConnections),
		closed:         false,
	}
	
	// 预创建一些连接
	for i := 0; i < maxConnections/2; i++ {
		client := NewClient(trackerAddr, trackerPort)
		if err := client.Connect(); err == nil {
			pool.connections <- client
		}
	}
	
	return pool
}

// Get 从连接池获取连接
func (p *connectionPool) Get() (*Client, error) {
	p.mu.RLock()
	if p.closed {
		p.mu.RUnlock()
		return nil, fmt.Errorf("connection pool is closed")
	}
	p.mu.RUnlock()
	
	select {
	case client := <-p.connections:
		// 检查连接是否有效
		if client.IsConnected() {
			if err := client.Ping(); err == nil {
				return client, nil
			}
		}
		// 连接无效，重新创建
		client.Close()
		return p.createNewConnection()
	default:
		// 没有可用连接，创建新连接
		return p.createNewConnection()
	}
}

// Put 将连接放回连接池
func (p *connectionPool) Put(client *Client) error {
	if client == nil {
		return nil
	}
	
	p.mu.RLock()
	if p.closed {
		p.mu.RUnlock()
		return client.Close()
	}
	p.mu.RUnlock()
	
	// 检查连接是否有效
	if !client.IsConnected() || client.Ping() != nil {
		client.Close()
		return nil
	}
	
	select {
	case p.connections <- client:
		return nil
	default:
		// 连接池已满，关闭连接
		return client.Close()
	}
}

// Close 关闭连接池
func (p *connectionPool) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	
	if p.closed {
		return nil
	}
	
	p.closed = true
	close(p.connections)
	
	// 关闭所有连接
	for client := range p.connections {
		client.Close()
	}
	
	return nil
}

// Size 获取连接池大小
func (p *connectionPool) Size() int {
	return p.maxConnections
}

// Available 获取可用连接数
func (p *connectionPool) Available() int {
	return len(p.connections)
}

// createNewConnection 创建新连接
func (p *connectionPool) createNewConnection() (*Client, error) {
	client := NewClient(p.trackerAddr, p.trackerPort)
	err := client.Connect()
	if err != nil {
		return nil, fmt.Errorf("failed to create new connection: %w", err)
	}
	return client, nil
}

// PooledClient 带连接池的客户端
type PooledClient struct {
	pool   ConnectionPool
	client *Client
}

// NewPooledClient 创建带连接池的客户端
func NewPooledClient(pool ConnectionPool) *PooledClient {
	return &PooledClient{
		pool: pool,
	}
}

// getClient 获取客户端连接
func (pc *PooledClient) getClient() (*Client, error) {
	if pc.client == nil {
		client, err := pc.pool.Get()
		if err != nil {
			return nil, err
		}
		pc.client = client
	}
	return pc.client, nil
}

// releaseClient 释放客户端连接
func (pc *PooledClient) releaseClient() {
	if pc.client != nil {
		pc.pool.Put(pc.client)
		pc.client = nil
	}
}

// Ping 测试连接
func (pc *PooledClient) Ping() error {
	client, err := pc.getClient()
	if err != nil {
		return err
	}
	defer pc.releaseClient()
	
	return client.Ping()
}

// GetStorageServer 获取存储服务器信息
func (pc *PooledClient) GetStorageServer(groupName string) (*StorageServer, error) {
	client, err := pc.getClient()
	if err != nil {
		return nil, err
	}
	defer pc.releaseClient()
	
	return client.GetStorageServer(groupName)
}

// ListFiles 列出文件
func (pc *PooledClient) ListFiles(groupName string, startFileName string, limit int) ([]*FileInfo, error) {
	client, err := pc.getClient()
	if err != nil {
		return nil, err
	}
	defer pc.releaseClient()
	
	return client.ListFiles(groupName, startFileName, limit)
}

// DownloadFile 下载文件
func (pc *PooledClient) DownloadFile(fileID string) ([]byte, error) {
	client, err := pc.getClient()
	if err != nil {
		return nil, err
	}
	defer pc.releaseClient()
	
	return client.DownloadFile(fileID)
}

// UploadFile 上传文件
func (pc *PooledClient) UploadFile(groupName string, fileName string, data []byte) (string, error) {
	client, err := pc.getClient()
	if err != nil {
		return "", err
	}
	defer pc.releaseClient()
	
	return client.UploadFile(groupName, fileName, data)
}

// DeleteFile 删除文件
func (pc *PooledClient) DeleteFile(fileID string) error {
	client, err := pc.getClient()
	if err != nil {
		return err
	}
	defer pc.releaseClient()
	
	return client.DeleteFile(fileID)
}

// GetFileInfo 获取文件信息
func (pc *PooledClient) GetFileInfo(fileID string) (*FileInfo, error) {
	client, err := pc.getClient()
	if err != nil {
		return nil, err
	}
	defer pc.releaseClient()
	
	return client.GetFileInfo(fileID)
}