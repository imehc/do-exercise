# Docker Compose 相关命令
.PHONY: up down build logs ps restart clean certs dev dev-down dev-logs

# 默认使用 deploy/docker 目录下的 docker-compose.yml
DOCKER_COMPOSE_FILE := deploy/docker/docker-compose.yml
# 开发环境依赖（数据库/Redis/MinIO）
DEV_COMPOSE_FILE := deploy/docker/docker-compose.dev.yml

# SSL证书相关配置
NGINX_DIR := deploy/docker/nginx
SSL_CERT_DIR := $(NGINX_DIR)/certs

# 生成SSL证书
certs:
	@mkdir -p $(SSL_CERT_DIR)
	openssl req -x509 -newkey rsa:2048 -sha256 -days 365 \
		-nodes \
		-keyout $(SSL_CERT_DIR)/server.key \
		-out $(SSL_CERT_DIR)/server.crt \
		-config $(NGINX_DIR)/openssl.cnf

# 启动所有服务
up:
	docker compose -f $(DOCKER_COMPOSE_FILE) up -d

# 停止并移除所有服务
down:
	docker compose -f $(DOCKER_COMPOSE_FILE) down

# 构建服务
build:
	docker compose -f $(DOCKER_COMPOSE_FILE) build

# 查看服务日志
logs:
	docker compose -f $(DOCKER_COMPOSE_FILE) logs -f

# 查看服务状态
ps:
	docker compose -f $(DOCKER_COMPOSE_FILE) ps

# 重启服务
restart:
	docker compose -f $(DOCKER_COMPOSE_FILE) restart

# 清理所有容器和卷
clean:
	docker compose -f $(DOCKER_COMPOSE_FILE) down -v --remove-orphans

# 一键启动开发环境：依赖容器 + 数据库初始化 + 后端热重载
dev:
	@if [ ! -f deploy/docker/.env.dev ]; then \
		cp deploy/docker/.env.dev.example deploy/docker/.env.dev; \
		echo "已创建 deploy/docker/.env.dev（按需修改）"; \
	fi
	docker compose -f $(DEV_COMPOSE_FILE) up -d
	@echo "等待数据库就绪..."
	@i=0; until docker compose -f $(DEV_COMPOSE_FILE) exec -T postgres pg_isready >/dev/null 2>&1; do \
		i=$$((i+1)); \
		if [ $$i -gt 30 ]; then echo "等待数据库就绪超时"; exit 1; fi; \
		sleep 2; \
	done; echo "数据库已就绪"
	@if command -v air >/dev/null 2>&1 || [ -x "$$(go env GOPATH)/bin/air" ]; then \
		echo "初始化数据库（已有数据自动跳过）..."; \
		export PATH="$$(go env GOPATH)/bin:$$PATH"; \
		cd server && set -a && . ../deploy/docker/.env.dev && set +a && $(MAKE) migrate && $(MAKE) dev; \
	else \
		echo "依赖已就绪，但未检测到 air（后端热重载工具）。"; \
		echo "请先安装：cd server && make dev-install，再执行 make dev。"; \
	fi

# 停止开发环境依赖容器
dev-down:
	docker compose -f $(DEV_COMPOSE_FILE) down

# 查看开发环境依赖日志
dev-logs:
	docker compose -f $(DEV_COMPOSE_FILE) logs -f

# 帮助信息
help:
	@echo "Docker Compose 管理命令："
	@echo "make up        - 启动所有服务"
	@echo "make down      - 停止并移除所有服务"
	@echo "make build     - 构建服务"
	@echo "make logs      - 查看服务日志"
	@echo "make ps        - 查看服务状态"
	@echo "make restart   - 重启服务"
	@echo "make clean     - 清理所有容器和卷"
	@echo "make certs     - 生成SSL证书"
	@echo "make dev       - 一键开发：启动依赖容器（数据库/Redis/MinIO）+ 后端热重载"
	@echo "make dev-down  - 停止开发环境依赖容器"
	@echo "make dev-logs  - 查看开发环境依赖容器日志"