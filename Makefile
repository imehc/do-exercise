# Docker Compose 相关命令
.PHONY: up down build logs ps restart clean certs

# 默认使用 deploy/docker 目录下的 docker-compose.yml
DOCKER_COMPOSE_FILE := deploy/docker/docker-compose.yml

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