# Makefile for Go project
# 
# 注意事项：
# 1. 确保已安装必要的 Go 工具：
#    go install mvdan.cc/gofumpt@latest
#    go install golang.org/x/tools/cmd/goimports@latest
#    go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest
#
# 2. 如果遇到 "No such file or directory" 错误，请确保 Go bin 目录在 PATH 中：
#    export PATH=$PATH:$(go env GOPATH)/bin
#
# 3. 建议将 Go bin 目录添加到 ~/.zshrc 或 ~/.bashrc 中：
#    echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.zshrc
#
# 4. 格式化前建议先提交当前更改，格式化会修改文件内容

.PHONY: format import help lint lint-fix install-lint pprof dev dev-install gen migrate

# 代码生成器参数（详见 cmd/gen/README.md）
# 用法: make gen NAME=SysDevice DESC=设备 ROUTE=devices FIELDS="code:string:设备编码 name:string:设备名称"
NAME ?=
DESC ?=
ROUTE ?=
FIELDS ?=

# 显示帮助信息
help:
	@echo "可用的 make 命令："
	@echo "  make format     - 格式化 Go 代码 (使用 gofumpt)"
	@echo "  make import     - 自动整理导入包 (使用 goimports)"
	@echo "  make lint       - 运行代码检查 (使用 golangci-lint)"
	@echo "  make lint-fix   - 运行代码检查并自动修复可修复的问题"
	@echo "  make install-lint - 安装 golangci-lint"
	@echo "  make dev        - 使用 air 启动开发服务器（热重载）"
	@echo "  make dev-install - 安装 air 开发工具"
	@echo "  make pprof      - 启动 pprof 性能分析服务并显示常用分析命令"
	@echo "  make gen        - 使用代码生成器生成实体 CRUD 六件套"
	@echo "  make migrate    - 初始化数据库（建表 + 种子数据，已有数据自动跳过）"
	@echo "  make help       - 显示此帮助信息"
	@echo ""
	@echo "注意事项："
	@echo "  - 格式化会修改文件内容，建议先提交当前更改"
	@echo "  - 确保 Go 工具已安装且 PATH 配置正确"
	@echo "  - 开发时建议使用 make dev 命令启动热重载服务"

# 安装 golangci-lint
install-lint:
	@echo "正在安装 golangci-lint..."
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest
	@echo "golangci-lint 安装完成！"

# 格式化代码
format:
	@echo "正在格式化 Go 代码..."
	gofumpt -w .
	@echo "格式化完成！"

# 自动导入包
import:
	@echo "正在整理导入包..."
	goimports -w .
	@echo "导入包整理完成！"

# 运行代码检查
lint:
	@echo "正在运行代码检查..."
	golangci-lint run
	@echo "代码检查完成！"

# 运行代码检查并自动修复
lint-fix:
	@echo "正在运行代码检查并自动修复..."
	golangci-lint run --fix
	@echo "代码检查和修复完成！"

# 安装 air 开发工具（注意：air 已迁移到 github.com/air-verse/air）
dev-install:
	@echo "正在安装 air 开发工具..."
	go install github.com/air-verse/air@latest
	@echo "air 安装完成！"
	@echo "现在可以使用 make dev 命令启动开发服务器"

# 使用 air 启动开发服务器（热重载）
dev:
	@echo "启动开发服务器（热重载模式）..."
	@echo "服务器将自动监听文件变化并重新编译运行"
	@echo "按 Ctrl+C 停止服务器"
	air

# 初始化数据库（建表 + 种子数据；幂等，已有数据自动跳过）
migrate:
	@echo "正在初始化数据库..."
	@go run main.go migrate --sql=init.sql
	@echo "数据库初始化完成！"

# 启动 pprof 性能分析服务（监听 6060 端口）
pprof:
	@echo "启动 pprof 性能分析服务，监听 http://localhost:6060 ..."
	@go run main.go &
	@sleep 2
	@echo "你可以使用如下命令进行分析："
	@echo "  CPU分析: go tool pprof http://localhost:6060/debug/pprof/profile"
	@echo "  内存分析: go tool pprof http://localhost:6060/debug/pprof/heap"
	@echo "  协程分析: go tool pprof http://localhost:6060/debug/pprof/goroutine"
	@echo "访问 http://localhost:6060/debug/pprof/ 查看更多分析入口"

# 代码生成器（详见 cmd/gen/README.md）
gen:
	@if [ -z "$(NAME)" ]; then \
		echo "用法: make gen NAME=SysDevice DESC=设备 ROUTE=devices FIELDS='code:string:设备编码 name:string:设备名称'"; \
		echo "更多说明见 cmd/gen/README.md"; \
		exit 1; \
	fi
	@echo "正在为 $(NAME) 生成实体代码..."
	@go run ./cmd/gen -name="$(NAME)" -desc="$(DESC)" -path="$(ROUTE)" $(foreach f,$(FIELDS),-field="$(f)")
	@echo "生成完成，请执行 go build ./... 验证编译"