# do-exercise

## 项目简介

do-exercise 是一个基于 Go + React + Vite + Tailwind CSS 的全栈练习项目，支持用户认证、权限管理、菜单管理、操作日志、系统信息、Token 管理等常见后台功能。

## 主要功能

- 用户注册、登录、找回密码、验证码校验
- 权限管理（基于 Casbin）
- 菜单与角色管理
- 操作日志记录与查询
- Token 管理
- 系统信息展示
- 国际化支持（中英文）
- 文件上传（MinIO 对象存储）
- 邮件通知
- 响应式前端 UI，支持暗黑模式

## 技术栈

- **后端**：Go 1.24、Gin、GORM、Casbin、Zap、MinIO、Redis、PostgreSQL
- **前端**：React 19、Vite、Tailwind CSS、Radix UI、TanStack Router & Query、Zod、Lingui（i18n）
- **容器化/部署**：Docker Compose、Nginx

## 快速开始

### 1. 克隆项目

```bash
git clone https://github.com/imehc/do-exercise.git
cd do-exercise
```

### 2. 配置文件准备

- 将 `deploy/docker/.env.example` 重命名为 `.env`，并根据实际需求填写配置项
- 将 `server/config.example.yaml` 重命名为 `config.yaml`，并根据实际需求修改配置项
- **本地开发时，`.env` 和 `config.yaml` 中的 host 配置应填写实际可访问的主机地址（如 `127.0.0.1` 或本机局域网 IP），而不是 Docker 容器名称**

### 3. 证书准备

- 如运行 `docker-compose` 时提示 `certs` 目录不存在，请执行 `make certs` 生成自签名证书，或将你自己的证书手动放到 `deploy/docker/nginx/certs` 目录下

### 4. 启动服务

```bash
# 构建并启动所有服务
make up

# 查看服务状态
make ps

# 查看服务日志
make logs
```

前端默认访问：https://localhost  

### 5. 本地开发

- 前端开发：`cd web && pnpm install && pnpm dev`
- 后端开发：`cd server && go run main.go` 或使用 IDE 调试

## 配置说明

- 数据库、Redis、MinIO、邮箱等服务的连接信息需与实际环境一致
- MinIO Access Key/Secret Key 必须配置，Bucket 名需与 `config.yaml` 保持一致
- 邮箱配置用于找回密码、通知等功能

## 常见问题

- **MinIO 启动失败**：请检查 Access Key/Secret Key 是否配置，Bucket 名是否一致
- **证书目录不存在**：请执行 `make certs` 或手动放置证书
- **服务无法访问**：本地开发请确保 host 配置为实际可访问地址

## 目录结构

```
do-exercise/
  ├── server/      # Go 后端服务
  ├── web/         # React 前端项目
  ├── deploy/      # Docker、Nginx 配置
  └── Makefile     # 常用命令
```

## 贡献指南

欢迎提交 Issue 或 PR，建议先阅读代码结构和配置说明。

## License

MIT
