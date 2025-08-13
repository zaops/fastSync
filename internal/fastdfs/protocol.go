package fastdfs

// FastDFS协议常量
const (
	// 协议头长度
	FDFS_PROTO_PKG_LEN_SIZE = 8
	
	// 组名最大长度
	FDFS_GROUP_NAME_MAX_LEN = 16
	
	// IP地址长度
	IP_ADDRESS_SIZE = 16
	
	// 文件名最大长度
	FDFS_FILE_NAME_MAX_LEN = 128
	
	// 文件扩展名最大长度
	FDFS_FILE_EXT_NAME_MAX_LEN = 6
	
	// Tracker查询存储服务器响应体长度
	TRACKER_QUERY_STORAGE_STORE_BODY_LEN = FDFS_GROUP_NAME_MAX_LEN + IP_ADDRESS_SIZE + FDFS_PROTO_PKG_LEN_SIZE
)

// 协议命令
const (
	// Tracker协议命令
	TRACKER_PROTO_CMD_SERVICE_QUERY_STORE_WITHOUT_GROUP_ONE = 101
	TRACKER_PROTO_CMD_SERVICE_QUERY_FETCH_ONE               = 102
	TRACKER_PROTO_CMD_SERVICE_QUERY_UPDATE                  = 103
	TRACKER_PROTO_CMD_SERVICE_QUERY_STORE_WITH_GROUP_ONE    = 104
	TRACKER_PROTO_CMD_SERVICE_QUERY_FETCH_ALL               = 105
	TRACKER_PROTO_CMD_SERVICE_QUERY_STORE_WITHOUT_GROUP_ALL = 106
	TRACKER_PROTO_CMD_SERVICE_QUERY_STORE_WITH_GROUP_ALL    = 107
	
	// Storage协议命令
	STORAGE_PROTO_CMD_UPLOAD_FILE                = 11
	STORAGE_PROTO_CMD_DELETE_FILE                = 12
	STORAGE_PROTO_CMD_SET_METADATA               = 13
	STORAGE_PROTO_CMD_DOWNLOAD_FILE              = 14
	STORAGE_PROTO_CMD_GET_METADATA               = 15
	STORAGE_PROTO_CMD_UPLOAD_SLAVE_FILE          = 21
	STORAGE_PROTO_CMD_QUERY_FILE_INFO            = 22
	STORAGE_PROTO_CMD_UPLOAD_APPENDER_FILE       = 23
	STORAGE_PROTO_CMD_APPEND_FILE                = 24
	STORAGE_PROTO_CMD_MODIFY_FILE                = 34
	STORAGE_PROTO_CMD_TRUNCATE_FILE              = 36
	STORAGE_PROTO_CMD_LIST_ONE_GROUP             = 40
	STORAGE_PROTO_CMD_LIST_ALL_GROUPS            = 41
	
	// 通用协议命令
	FDFS_PROTO_CMD_QUIT        = 82
	FDFS_PROTO_CMD_ACTIVE_TEST = 111
	FDFS_PROTO_CMD_RESP        = 100
)

// 协议状态码
const (
	FDFS_PROTO_STATUS_SUCCESS = 0
	FDFS_PROTO_STATUS_ERROR   = 1
)

// Header 协议头结构
type Header struct {
	Length  int64 // 包体长度
	Command byte  // 命令
	Status  byte  // 状态
}

// StorageServer 存储服务器信息
type StorageServer struct {
	GroupName string // 组名
	IPAddr    string // IP地址
	Port      int    // 端口
	StorePathIndex byte // 存储路径索引
}

// FileInfo 文件信息
type FileInfo struct {
	GroupName    string    // 组名
	FileName     string    // 文件名
	FileSize     int64     // 文件大小
	CreateTime   int64     // 创建时间
	CRC32        uint32    // CRC32校验值
	SourceIPAddr string    // 源IP地址
}

// UploadResponse 上传响应
type UploadResponse struct {
	GroupName string // 组名
	FileName  string // 文件名
}

// GroupInfo 组信息
type GroupInfo struct {
	GroupName      string // 组名
	TotalMB        int64  // 总容量(MB)
	FreeMB         int64  // 剩余容量(MB)
	TrunkFreeMB    int64  // Trunk剩余容量(MB)
	StorageCount   int    // 存储服务器数量
	StoragePort    int    // 存储服务器端口
	StorageHTTPPort int   // 存储服务器HTTP端口
	ActiveCount    int    // 活跃存储服务器数量
	CurrentWriteServer int // 当前写入服务器索引
	StorePathCount int    // 存储路径数量
	SubdirCountPerPath int // 每个路径的子目录数量
	CurrentTrunkFileID int // 当前Trunk文件ID
}

// ConnectionPool 连接池接口
type ConnectionPool interface {
	Get() (*Client, error)
	Put(*Client) error
	Close() error
	Size() int
	Available() int
}