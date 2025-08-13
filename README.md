# FastDFS Migration System

FastDFS 5.0.7 到 6.0.6 集群迁移管理系统

## 功能特性

- 🚀 高性能文件迁移，支持断点续传
- 🕒 定时任务调度，支持增量同步
- 🎯 多种过滤条件（时间、文件类型）
- 📊 实时监控和进度跟踪
- 🌐 Web管理界面
- 📦 单一二进制部署

## 快速开始

### 环境要求

- Go 1.21+
- Redis (可选，用于任务队列)

### 安装和运行

1. 克隆项目
```bash
git clone <repository-url>
cd fastdfs-migration-system
```

2. 初始化项目
```bash
make init
make deps
```

3. 运行测试
```bash
make test
```

4. 运行开发服务器
```bash
make run
```

5. 访问Web界面
```
http://localhost:8080
```

### 项目交接

如果你是新接手这个项目的开发者，请先阅读：
- **交接文档**: `HANDOVER.md` - 项目状态和交接清单
- **开发指南**: `DEVELOPMENT.md` - 详细的开发说明
- **规格文档**: `.kiro/specs/fastdfs-migration-system/` - 完整的需求、设计和任务计划

### 构建生产版本

```bash
make build-prod
```

## 配置说明

编辑 `config.yaml` 文件进行配置：

```yaml
server:
  port: "8080"
  host: "0.0.0.0"

database:
  type: "sqlite"
  dsn: "./migration.db"

migration:
  default_workers: 5
  chunk_size: 1048576  # 1MB
  max_retry: 3
```

## 开发状态

- [x] 项目初始化和基础架构
- [x] 数据库模型和存储层
- [x] FastDFS客户端集成
- [ ] 核心迁移引擎 ⬅️ **进行中**
- [ ] 断点续传功能
- [ ] Web管理界面

## 许可证

MIT License