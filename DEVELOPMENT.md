# FastDFS迁移系统 - 开发指南

## 项目概述

这是一个用于将FastDFS 5.0.7集群数据迁移到6.0.6集群的Web管理系统。项目采用Go语言开发，支持断点续传、定时任务、增量同步等功能。

## 开发环境要求

- Go 1.21+
- Redis (可选，用于任务队列)
- Git

## 快速开始

### 1. 克隆项目
```bash
git clone <repository-url>
cd fastdfs-migration-system
```

### 2. 安装依赖
```bash
make deps
```

### 3. 初始化项目
```bash
make init
```

### 4. 运行开发服务器
```bash
make run
```

### 5. 运行测试
```bash
make test
```

## 项目结构

```
fastdfs-migration-system/
├── .kiro/specs/                    # 项目规格文档
│   └── fastdfs-migration-system/
│       ├── requirements.md         # 需求文档
│       ├── design.md              # 设计文档
│       └── tasks.md               # 任务计划
├── cmd/server/                     # 应用入口
│   └── main.go
├── internal/                       # 内部包
│   ├── config/                    # 配置管理
│   ├── logger/                    # 日志系统
│   └── server/                    # HTTP服务器
├── web/                           # 前端资源
│   └── templates/                 # HTML模板
├── bin/                           # 构建输出
├── logs/                          # 日志文件
├── config.yaml                    # 配置文件
├── go.mod                         # Go模块定义
├── Makefile                       # 构建脚本
└── README.md                      # 项目说明
```

## 开发状态

### ✅ 已完成
- [x] 项目初始化和基础架构搭建
  - Go模块配置和依赖管理
  - 配置管理系统(viper)
  - 日志系统(logrus)
  - HTTP服务器框架(Gin)
  - 基础测试框架

### 🚧 进行中
- [ ] 数据库模型和存储层实现
- [ ] FastDFS客户端集成
- [ ] 核心迁移引擎开发

### 📋 待开发
查看 `.kiro/specs/fastdfs-migration-system/tasks.md` 获取完整的任务列表

## 开发工作流

### 1. 选择任务
打开 `.kiro/specs/fastdfs-migration-system/tasks.md` 文件，选择下一个要开发的任务。

### 2. 阅读相关文档
- **需求文档**: `.kiro/specs/fastdfs-migration-system/requirements.md`
- **设计文档**: `.kiro/specs/fastdfs-migration-system/design.md`

### 3. 开发流程
1. 创建功能分支: `git checkout -b feature/task-name`
2. 实现功能代码
3. 编写单元测试
4. 运行测试: `make test`
5. 提交代码: `git commit -m "feat: implement task description"`
6. 推送分支: `git push origin feature/task-name`

### 4. 代码规范
- 遵循Go官方代码规范
- 每个公共函数都要有注释
- 单元测试覆盖率保持在80%以上
- 使用有意义的变量和函数名

## 核心技术栈

- **后端框架**: Gin (HTTP路由和中间件)
- **配置管理**: Viper (支持YAML、环境变量)
- **日志系统**: Logrus (结构化日志)
- **数据库**: GORM + SQLite/PostgreSQL
- **缓存队列**: Redis
- **任务调度**: Cron
- **测试框架**: Go标准testing + testify

## API设计

### 当前可用端点
- `GET /health` - 健康检查
- `GET /api/v1/ping` - API连通性测试
- `GET /` - Web界面首页

### 计划中的API端点
```
POST   /api/v1/migrations          # 创建迁移任务
GET    /api/v1/migrations          # 获取任务列表
GET    /api/v1/migrations/:id      # 获取任务详情
PUT    /api/v1/migrations/:id      # 更新任务配置
DELETE /api/v1/migrations/:id      # 删除任务
POST   /api/v1/migrations/:id/start   # 启动任务
POST   /api/v1/migrations/:id/pause   # 暂停任务
POST   /api/v1/migrations/:id/resume  # 恢复任务

POST   /api/v1/schedules          # 创建定时任务
GET    /api/v1/schedules          # 获取定时任务列表
PUT    /api/v1/schedules/:id      # 更新定时任务
DELETE /api/v1/schedules/:id      # 删除定时任务

POST   /api/v1/clusters           # 添加集群配置
GET    /api/v1/clusters           # 获取集群列表
PUT    /api/v1/clusters/:id       # 更新集群配置
DELETE /api/v1/clusters/:id       # 删除集群配置
POST   /api/v1/clusters/:id/test  # 测试集群连接
```

## 配置说明

编辑 `config.yaml` 文件进行配置：

```yaml
server:
  port: "8080"          # 服务端口
  host: "0.0.0.0"       # 监听地址

database:
  type: "sqlite"        # 数据库类型: sqlite/postgresql
  dsn: "./migration.db" # 数据库连接字符串

redis:
  addr: "localhost:6379"  # Redis地址
  password: ""            # Redis密码
  db: 0                   # Redis数据库编号

migration:
  default_workers: 5      # 默认并发工作线程数
  chunk_size: 1048576     # 分块大小(字节)
  max_retry: 3            # 最大重试次数
  retry_interval: "30s"   # 重试间隔

logging:
  level: "info"                    # 日志级别: debug/info/warn/error
  file: "./logs/migration.log"     # 日志文件路径
  max_size: 100                    # 日志文件最大大小(MB)
  max_backups: 5                   # 保留的日志文件数量
```

## 测试指南

### 运行所有测试
```bash
make test
```

### 运行特定包的测试
```bash
go test ./internal/server -v
```

### 运行测试并查看覆盖率
```bash
go test -cover ./...
```

### 生成测试覆盖率报告
```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## 构建和部署

### 开发构建
```bash
make build
```

### 生产构建
```bash
make build-prod
```

### 清理构建文件
```bash
make clean
```

## 调试技巧

### 1. 启用调试日志
在 `config.yaml` 中设置:
```yaml
logging:
  level: "debug"
```

### 2. 使用Delve调试器
```bash
go install github.com/go-delve/delve/cmd/dlv@latest
dlv debug ./cmd/server
```

### 3. 查看实时日志
```bash
tail -f logs/migration.log
```

## 常见问题

### Q: 如何添加新的API端点？
A: 在 `internal/server/server.go` 的 `setupRoutes()` 函数中添加路由定义。

### Q: 如何修改数据库模型？
A: 在 `internal/models/` 目录下定义新的结构体，使用GORM标签。

### Q: 如何添加新的配置项？
A: 在 `internal/config/config.go` 中添加字段，并在 `setDefaults()` 函数中设置默认值。

### Q: 测试失败怎么办？
A: 检查测试环境是否正确设置，确保没有端口冲突，查看详细错误信息。

## 联系方式

如有问题，请联系项目负责人或在项目仓库中创建Issue。

## 下一步开发建议

1. **优先任务**: 实现数据库模型和存储层 (任务2)
2. **关键功能**: FastDFS客户端集成 (任务3)
3. **核心逻辑**: 迁移引擎开发 (任务4)

建议按照 `tasks.md` 中的顺序逐步实现，每完成一个任务就进行测试和代码审查。