.PHONY: build worker cli test clean

# 构建所有二进制文件
build: build-worker build-cli

# 构建Worker
build-worker:
	go build -o bin/worker cmd/worker/main.go

# 构建CLI
build-cli:
	go build -o bin/cli cmd/cli/main.go

# 运行Worker
worker: build-worker
	./bin/worker

# 运行CLI
cli: build-cli
	./bin/cli

# 运行测试
test:
	go test -v ./...

# 清理构建文件
clean:
	rm -rf bin/
	rm -rf tmp/

# 整理依赖
deps:
	go mod tidy

# 查看依赖
deps-graph:
	go mod graph

# 格式化代码
fmt:
	gofmt -w .

# 检查代码
lint:
	golangci-lint run

# 启动本地服务（使用Docker）
up:
	docker-compose up -d

# 停止本地服务
down:
	docker-compose down

# 查看日志
logs:
	docker-compose logs -f

# 帮助
help:
	@echo "Available commands:"
	@echo "  build       - Build all binaries"
	@echo "  worker      - Build and run worker"
	@echo "  cli         - Build CLI"
	@echo "  test        - Run tests"
	@echo "  clean       - Clean build files"
	@echo "  deps        - Tidy dependencies"
	@echo "  fmt         - Format code"
	@echo "  lint        - Run linter"
	@echo "  up          - Start local services"
	@echo "  down        - Stop local services"
	@echo "  logs        - View logs"
