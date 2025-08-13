package fastdfs

import (
	"encoding/binary"
	"fmt"
	"strings"
	"time"
)

// listFilesFromStorage 从存储服务器列出文件
func (c *Client) listFilesFromStorage(groupName string, startFileName string, limit int) ([]*FileInfo, error) {
	// 构建请求数据
	data := make([]byte, FDFS_GROUP_NAME_MAX_LEN+FDFS_FILE_NAME_MAX_LEN+8)
	copy(data[0:FDFS_GROUP_NAME_MAX_LEN], []byte(groupName))
	copy(data[FDFS_GROUP_NAME_MAX_LEN:FDFS_GROUP_NAME_MAX_LEN+FDFS_FILE_NAME_MAX_LEN], []byte(startFileName))
	binary.BigEndian.PutUint64(data[FDFS_GROUP_NAME_MAX_LEN+FDFS_FILE_NAME_MAX_LEN:], uint64(limit))
	
	header := &Header{
		Length:  int64(len(data)),
		Command: STORAGE_PROTO_CMD_LIST_ONE_GROUP,
		Status:  0,
	}
	
	err := c.sendHeader(header)
	if err != nil {
		return nil, fmt.Errorf("failed to send list files request: %w", err)
	}
	
	err = c.sendData(data)
	if err != nil {
		return nil, fmt.Errorf("failed to send list files data: %w", err)
	}
	
	// 接收响应
	respHeader, err := c.receiveHeader()
	if err != nil {
		return nil, fmt.Errorf("failed to receive list files response: %w", err)
	}
	
	if respHeader.Status != 0 {
		return nil, fmt.Errorf("list files failed with status: %d", respHeader.Status)
	}
	
	if respHeader.Length == 0 {
		return []*FileInfo{}, nil
	}
	
	respData := make([]byte, respHeader.Length)
	err = c.receiveData(respData)
	if err != nil {
		return nil, fmt.Errorf("failed to receive file list data: %w", err)
	}
	
	return parseFileList(respData)
}

// downloadFromStorage 从存储服务器下载文件
func (c *Client) downloadFromStorage(groupName string, fileName string) ([]byte, error) {
	// 构建请求数据
	data := make([]byte, 16+FDFS_GROUP_NAME_MAX_LEN+len(fileName))
	binary.BigEndian.PutUint64(data[0:8], 0)  // 文件偏移量
	binary.BigEndian.PutUint64(data[8:16], 0) // 下载字节数，0表示下载整个文件
	copy(data[16:16+FDFS_GROUP_NAME_MAX_LEN], []byte(groupName))
	copy(data[16+FDFS_GROUP_NAME_MAX_LEN:], []byte(fileName))
	
	header := &Header{
		Length:  int64(len(data)),
		Command: STORAGE_PROTO_CMD_DOWNLOAD_FILE,
		Status:  0,
	}
	
	err := c.sendHeader(header)
	if err != nil {
		return nil, fmt.Errorf("failed to send download request: %w", err)
	}
	
	err = c.sendData(data)
	if err != nil {
		return nil, fmt.Errorf("failed to send download data: %w", err)
	}
	
	// 接收响应
	respHeader, err := c.receiveHeader()
	if err != nil {
		return nil, fmt.Errorf("failed to receive download response: %w", err)
	}
	
	if respHeader.Status != 0 {
		return nil, fmt.Errorf("download failed with status: %d", respHeader.Status)
	}
	
	if respHeader.Length == 0 {
		return []byte{}, nil
	}
	
	fileData := make([]byte, respHeader.Length)
	err = c.receiveData(fileData)
	if err != nil {
		return nil, fmt.Errorf("failed to receive file data: %w", err)
	}
	
	return fileData, nil
}

// uploadToStorage 上传文件到存储服务器
func (c *Client) uploadToStorage(groupName string, fileName string, data []byte) (string, error) {
	// 构建请求数据
	extName := getFileExtension(fileName)
	requestData := make([]byte, 1+FDFS_FILE_EXT_NAME_MAX_LEN+len(data))
	requestData[0] = 0 // store_path_index
	copy(requestData[1:1+FDFS_FILE_EXT_NAME_MAX_LEN], []byte(extName))
	copy(requestData[1+FDFS_FILE_EXT_NAME_MAX_LEN:], data)
	
	header := &Header{
		Length:  int64(len(requestData)),
		Command: STORAGE_PROTO_CMD_UPLOAD_FILE,
		Status:  0,
	}
	
	err := c.sendHeader(header)
	if err != nil {
		return "", fmt.Errorf("failed to send upload request: %w", err)
	}
	
	err = c.sendData(requestData)
	if err != nil {
		return "", fmt.Errorf("failed to send upload data: %w", err)
	}
	
	// 接收响应
	respHeader, err := c.receiveHeader()
	if err != nil {
		return "", fmt.Errorf("failed to receive upload response: %w", err)
	}
	
	if respHeader.Status != 0 {
		return "", fmt.Errorf("upload failed with status: %d", respHeader.Status)
	}
	
	if respHeader.Length < FDFS_GROUP_NAME_MAX_LEN {
		return "", fmt.Errorf("invalid upload response length: %d", respHeader.Length)
	}
	
	respData := make([]byte, respHeader.Length)
	err = c.receiveData(respData)
	if err != nil {
		return "", fmt.Errorf("failed to receive upload response data: %w", err)
	}
	
	uploadResp := parseUploadResponse(respData)
	return fmt.Sprintf("%s/%s", uploadResp.GroupName, uploadResp.FileName), nil
}

// deleteFromStorage 从存储服务器删除文件
func (c *Client) deleteFromStorage(groupName string, fileName string) error {
	// 构建请求数据
	data := make([]byte, FDFS_GROUP_NAME_MAX_LEN+len(fileName))
	copy(data[0:FDFS_GROUP_NAME_MAX_LEN], []byte(groupName))
	copy(data[FDFS_GROUP_NAME_MAX_LEN:], []byte(fileName))
	
	header := &Header{
		Length:  int64(len(data)),
		Command: STORAGE_PROTO_CMD_DELETE_FILE,
		Status:  0,
	}
	
	err := c.sendHeader(header)
	if err != nil {
		return fmt.Errorf("failed to send delete request: %w", err)
	}
	
	err = c.sendData(data)
	if err != nil {
		return fmt.Errorf("failed to send delete data: %w", err)
	}
	
	// 接收响应
	respHeader, err := c.receiveHeader()
	if err != nil {
		return fmt.Errorf("failed to receive delete response: %w", err)
	}
	
	if respHeader.Status != 0 {
		return fmt.Errorf("delete failed with status: %d", respHeader.Status)
	}
	
	return nil
}

// getFileInfoFromStorage 从存储服务器获取文件信息
func (c *Client) getFileInfoFromStorage(groupName string, fileName string) (*FileInfo, error) {
	// 构建请求数据
	data := make([]byte, FDFS_GROUP_NAME_MAX_LEN+len(fileName))
	copy(data[0:FDFS_GROUP_NAME_MAX_LEN], []byte(groupName))
	copy(data[FDFS_GROUP_NAME_MAX_LEN:], []byte(fileName))
	
	header := &Header{
		Length:  int64(len(data)),
		Command: STORAGE_PROTO_CMD_QUERY_FILE_INFO,
		Status:  0,
	}
	
	err := c.sendHeader(header)
	if err != nil {
		return nil, fmt.Errorf("failed to send file info request: %w", err)
	}
	
	err = c.sendData(data)
	if err != nil {
		return nil, fmt.Errorf("failed to send file info data: %w", err)
	}
	
	// 接收响应
	respHeader, err := c.receiveHeader()
	if err != nil {
		return nil, fmt.Errorf("failed to receive file info response: %w", err)
	}
	
	if respHeader.Status != 0 {
		return nil, fmt.Errorf("get file info failed with status: %d", respHeader.Status)
	}
	
	if respHeader.Length < 3*8+IP_ADDRESS_SIZE {
		return nil, fmt.Errorf("invalid file info response length: %d", respHeader.Length)
	}
	
	respData := make([]byte, respHeader.Length)
	err = c.receiveData(respData)
	if err != nil {
		return nil, fmt.Errorf("failed to receive file info data: %w", err)
	}
	
	return parseFileInfo(groupName, fileName, respData)
}

// parseFileList 解析文件列表响应
func parseFileList(data []byte) ([]*FileInfo, error) {
	var files []*FileInfo
	
	// 每个文件信息的固定长度
	const fileInfoLen = FDFS_GROUP_NAME_MAX_LEN + FDFS_FILE_NAME_MAX_LEN + 3*8 + 4 + IP_ADDRESS_SIZE
	
	for i := 0; i < len(data); i += fileInfoLen {
		if i+fileInfoLen > len(data) {
			break
		}
		
		fileData := data[i : i+fileInfoLen]
		fileInfo := &FileInfo{
			GroupName:    strings.TrimRight(string(fileData[0:FDFS_GROUP_NAME_MAX_LEN]), "\x00"),
			FileName:     strings.TrimRight(string(fileData[FDFS_GROUP_NAME_MAX_LEN:FDFS_GROUP_NAME_MAX_LEN+FDFS_FILE_NAME_MAX_LEN]), "\x00"),
			FileSize:     int64(binary.BigEndian.Uint64(fileData[FDFS_GROUP_NAME_MAX_LEN+FDFS_FILE_NAME_MAX_LEN:])),
			CreateTime:   int64(binary.BigEndian.Uint64(fileData[FDFS_GROUP_NAME_MAX_LEN+FDFS_FILE_NAME_MAX_LEN+8:])),
			CRC32:        binary.BigEndian.Uint32(fileData[FDFS_GROUP_NAME_MAX_LEN+FDFS_FILE_NAME_MAX_LEN+24:]),
			SourceIPAddr: strings.TrimRight(string(fileData[FDFS_GROUP_NAME_MAX_LEN+FDFS_FILE_NAME_MAX_LEN+28:]), "\x00"),
		}
		
		files = append(files, fileInfo)
	}
	
	return files, nil
}

// parseUploadResponse 解析上传响应
func parseUploadResponse(data []byte) *UploadResponse {
	return &UploadResponse{
		GroupName: strings.TrimRight(string(data[0:FDFS_GROUP_NAME_MAX_LEN]), "\x00"),
		FileName:  strings.TrimRight(string(data[FDFS_GROUP_NAME_MAX_LEN:]), "\x00"),
	}
}

// parseFileInfo 解析文件信息响应
func parseFileInfo(groupName, fileName string, data []byte) (*FileInfo, error) {
	if len(data) < 3*8+4+IP_ADDRESS_SIZE {
		return nil, fmt.Errorf("invalid file info data length")
	}
	
	return &FileInfo{
		GroupName:    groupName,
		FileName:     fileName,
		FileSize:     int64(binary.BigEndian.Uint64(data[0:8])),
		CreateTime:   int64(binary.BigEndian.Uint64(data[8:16])),
		CRC32:        binary.BigEndian.Uint32(data[24:28]),
		SourceIPAddr: strings.TrimRight(string(data[28:28+IP_ADDRESS_SIZE]), "\x00"),
	}, nil
}

// getFileExtension 获取文件扩展名
func getFileExtension(fileName string) string {
	parts := strings.Split(fileName, ".")
	if len(parts) > 1 {
		ext := parts[len(parts)-1]
		if len(ext) > FDFS_FILE_EXT_NAME_MAX_LEN {
			ext = ext[:FDFS_FILE_EXT_NAME_MAX_LEN]
		}
		return ext
	}
	return ""
}

// GetCreateTime 获取文件创建时间
func (f *FileInfo) GetCreateTime() time.Time {
	return time.Unix(f.CreateTime, 0)
}

// GetFileID 获取文件ID
func (f *FileInfo) GetFileID() string {
	return fmt.Sprintf("%s/%s", f.GroupName, f.FileName)
}