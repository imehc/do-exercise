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
- 文件上传（RustFS 对象存储）
- 邮件通知
- 响应式前端 UI，支持暗黑模式

## 技术栈

- **后端**：Go 1.24、Gin、GORM、Casbin、Zap、RustFS、Redis、PostgreSQL
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
- Docker 部署时，数据库和 RustFS 的连接地址、账号、密码、Bucket 等运行环境配置放在 `deploy/docker/.env`
- `server/config.yaml` 只保留应用默认参数，如连接池、日志、验证码、RustFS 预签名有效期等
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

推荐使用 `make dev` 一键启动开发环境：它会自动复制开发环境变量、启动依赖容器（PostgreSQL / Redis / RustFS）、等待数据库就绪，并初始化数据库后通过 `air` 启动后端热重载。

#### 前置准备（仅首次）

1. 安装后端热重载工具 `air`：

   ```bash
   cd server && make dev-install
   ```

2. 确保 Go 的 bin 目录（`$(go env GOPATH)/bin`）已加入 `PATH`，否则 `air` 无法被找到。

#### 一键启动（推荐）

在项目根目录执行：

```bash
make dev
```

该命令会依次完成：

- 若 `deploy/docker/.env.dev` 不存在，则从 `.env.dev.example` 复制一份（按需修改）
- 启动 `deploy/docker/docker-compose.dev.yml` 中的依赖容器（PostgreSQL / Redis / RustFS，端口映射到宿主机 5432/6379/9000/9001）
- 等待数据库就绪后，自动加载 `.env.dev` 并执行数据库初始化（`make migrate`，幂等，已有数据自动跳过）
- 启动后端 `air` 热重载服务（监听文件变化自动重新编译运行）

> 注意：首次执行若未检测到 `air`，`make dev` 会提示你先执行 `cd server && make dev-install`，安装后再重新运行 `make dev`。

#### 前端开发

另开一个终端，单独启动前端：

```bash
cd web
pnpm install
pnpm dev
```

前端默认通过 Vite 开发服务器访问（配合 `/oss` 代理转发到本地 RustFS `127.0.0.1:9000` 处理上传）。

#### 开发环境依赖容器管理

```bash
make dev-down   # 停止并移除开发依赖容器
make dev-logs   # 查看开发依赖容器日志
```

#### 开发环境变量说明

开发环境变量统一放在 `deploy/docker/.env.dev`（`make dev` 自动生成，不提交仓库），由后端进程和 `docker-compose.dev.yml` 共用：

- `POSTGRES_HOST` / `REDIS_HOST` / `RUSTFS_HOST` 在本地均指向 `localhost`，便于本机运行的 Go 进程直接连接
- 容器内部的 `rustfs-init` 会自动将 `RUSTFS_HOST` 覆盖为服务名 `rustfs` 完成 bucket 初始化
- `RUSTFS_PRESIGNED_HOST` 在本地为 `http://localhost:9000`，生产环境则改为走 Nginx 反代路径（如 `/oss`）

如需要手动调整开发配置，直接编辑 `deploy/docker/.env.dev` 即可，无需改动 `server/config.yaml`。

#### 仅手动启动后端（不使用 make dev）

如果你只想手动启动后端，可先加载开发环境变量再运行：

```bash
set -a
source deploy/docker/.env.dev
set +a

cd server
make migrate   # 初始化数据库（首次或重置时）
make dev       # air 热重载；或用 go run main.go 直接运行
```

> 后端读取的是进程环境变量，不会自动加载 `.env.dev`，需手动 `source` 或借助 IDE 运行配置注入。

## 配置说明

- `deploy/docker/.env` 优先级高于 `server/config.yaml`
- 数据库与 RustFS 的部署差异和敏感信息放在 `.env`，不要写入 `config.yaml`
- RustFS Bucket、Access Key、Secret Key 会由 `rustfs-init` 容器自动初始化
- 邮箱配置用于找回密码、通知等功能

## 默认管理员账户

- 用户名：`superAdmin`
- 密码：`@admin2025`

请在首次登录后及时修改默认管理员密码以确保安全。

## 常见问题

- **RustFS 启动失败**：请检查 Access Key/Secret Key 是否配置，Bucket 名是否一致
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
