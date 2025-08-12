.PHONY: build run clean test deps

# 构建二进制文件
build:
	go build -o bin/migration-system ./cmd/server

# 运行开发服务器
run:
	go run ./cmd/server

# 清理构建文件
clean:
	rm -rf bin/
	rm -rf logs/
	rm -f migration.db

# 运行测试
test:
	go test -v ./...

# 下载依赖
deps:
	go mod download
	go mod tidy

# 构建生产版本
build-prod:
	CGO_ENABLED=1 GOOS=linux go build -ldflags="-w -s" -o bin/migration-system ./cmd/server

# 初始化项目
init:
	go mod download
	mkdir -p logs
	mkdir -p bin

# 验证环境设置
verify:
	go run scripts/verify-setup.go