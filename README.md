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
- Docker 部署时，数据库和 MinIO 的连接地址、账号、密码、Bucket 等运行环境配置放在 `deploy/docker/.env`
- `server/config.yaml` 只保留应用默认参数，如连接池、日志、验证码、MinIO 预签名有效期等
- **本地开发时，如不使用 Docker，请通过环境变量或本地配置填写实际可访问的主机地址（如 `127.0.0.1` 或本机局域网 IP），而不是 Docker 容器名称**

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
- 后端开发：先设置 `POSTGRES_*`、`MINIO_*` 等环境变量，再执行 `cd server && go run main.go` 或使用 IDE 调试

#### 开发模式数据库和 MinIO

后端本机运行时不会自动读取 `deploy/docker/.env`，需要在当前 shell、IDE 运行配置或本地 dotenv 工具中设置环境变量。`server/config.yaml` 不保存数据库和 MinIO 的地址、账号、密码。

```bash
export POSTGRES_HOST=127.0.0.1
export POSTGRES_PORT=5432
export POSTGRES_USER=admin
export POSTGRES_PASSWORD=admin2025
export POSTGRES_DB=do-exercise

export MINIO_HOST=127.0.0.1
export MINIO_PORT=9000
export MINIO_BUCKET_NAME=do-exercise
export MINIO_APP_ACCESS_KEY=your-access-key
export MINIO_APP_SECRET_KEY=your-secret-key
export MINIO_PRESIGNED_HOST=/oss
```

- `MINIO_HOST` / `MINIO_PORT` 是后端连接 MinIO API 的地址，本地开发通常是 `127.0.0.1:9000`
- `MINIO_PRESIGNED_HOST=/oss` 配合 Vite 开发代理使用，前端上传会通过 `/oss` 转发到 `127.0.0.1:9000`
- 本地 MinIO 仍需提前创建 `MINIO_BUCKET_NAME` 对应 bucket、匿名只读规则和应用 Access Key
- 如只用 Docker 跑依赖服务，需要把 PostgreSQL `5432` 和 MinIO API `9000` 暴露到宿主机，并避免同时启动占用 `9000` 的 `web` 服务

#### 使用开发环境 env 文件

也可以为本机开发准备单独的环境变量文件，例如：

- `deploy/docker/.env`：Docker 部署使用
- `server/.env.development`：本机开发使用，保存真实本地配置，不提交仓库
- `server/.env.example`：本机开发示例，可提交仓库

`server/.env.development` 示例：

```env
POSTGRES_HOST=127.0.0.1
POSTGRES_PORT=5432
POSTGRES_USER=admin
POSTGRES_PASSWORD=admin2025
POSTGRES_DB=do-exercise

MINIO_HOST=127.0.0.1
MINIO_PORT=9000
MINIO_BUCKET_NAME=do-exercise
MINIO_APP_ACCESS_KEY=your-access-key
MINIO_APP_SECRET_KEY=your-secret-key
MINIO_PRESIGNED_HOST=/oss
```

启动后端前手动加载：

```bash
set -a
source server/.env.development
set +a

cd server
go run main.go
```

当前后端读取的是进程环境变量，不会自动加载 `server/.env.development`；如果需要自动加载，需要额外加入 dotenv 加载逻辑。

## 配置说明

- `deploy/docker/.env` 优先级高于 `server/config.yaml`
- 数据库与 MinIO 的部署差异和敏感信息放在 `.env`，不要写入 `config.yaml`
- MinIO Bucket、Access Key、Secret Key 会由 `minio-init` 容器自动初始化
- 邮箱配置用于找回密码、通知等功能

## 默认管理员账户

- 用户名：`admin`
- 密码：`@admin2025`

请在首次登录后及时修改默认管理员密码以确保安全。

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
