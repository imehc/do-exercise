## RustFS 自动初始化

`docker compose up -d` 会同时启动 `rustfs-init` 一次性容器，自动完成以下操作：

- 等待 RustFS 服务可连接
- 创建 `.env` 中的 `RUSTFS_BUCKET_NAME`
- 设置 bucket 匿名只读访问
- 创建 server 使用的应用 Access Key（`rc admin service-account create`）

### 环境变量

部署前确认 `deploy/docker/.env` 包含以下配置：

#### 数据库

| 配置项 | 用途 |
| --- | --- |
| `POSTGRES_HOST` | PostgreSQL Docker 网络内主机名 |
| `POSTGRES_PORT` | PostgreSQL Docker 网络内端口 |
| `POSTGRES_DB` | PostgreSQL 数据库名 |
| `POSTGRES_USER` | PostgreSQL 用户名 |
| `POSTGRES_PASSWORD` | PostgreSQL 密码 |

#### RustFS

| 配置项 | 用途 |
| --- | --- |
| `RUSTFS_ACCESS_KEY` | RustFS root Access Key，用于初始化（勿用默认值 `rustfsadmin`） |
| `RUSTFS_SECRET_KEY` | RustFS root Secret Key，用于初始化 |
| `RUSTFS_RPC_SECRET` | 节点间 RPC 认证密钥，强随机且 ≠ `RUSTFS_SECRET_KEY` |
| `RUSTFS_HOST` | RustFS Docker 网络内主机名 |
| `RUSTFS_PORT` | RustFS Docker 网络内端口 |
| `RUSTFS_BUCKET_NAME` | 自动创建的 bucket 名称 |
| `RUSTFS_APP_ACCESS_KEY` | server 上传文件使用的 Access Key |
| `RUSTFS_APP_SECRET_KEY` | server 上传文件使用的 Secret Key |
| `RUSTFS_PRESIGNED_HOST` | 浏览器访问 RustFS 的公开反代路径 |

server 会从 `.env` 读取 RustFS 地址、Bucket、Access Key、Secret Key 和公开上传路径，覆盖 `server/config.yaml` 中的同名配置。

### 配置分工

Docker 部署时，`.env` 是运行环境入口，`server/config.yaml` 只保留应用默认参数。

| 配置来源 | 放置内容 |
| --- | --- |
| `deploy/docker/.env` | 数据库地址/端口/库名/账号/密码、RustFS 地址/端口/Bucket/Access Key/Secret Key/公开上传路径 |
| `server/config.yaml` | 数据库连接池、Redis 连接池与超时、日志、验证码、RustFS 预签名有效期、系统默认参数 |

如果同一个配置项同时存在于 `.env` 和 `server/config.yaml`，server 启动时以 `.env` 为准。
`docker-compose.yml` 会按服务显式映射所需变量，不会把整份 `.env` 注入每个容器。

### 开发模式

如果只用 Docker 启动数据库和 RustFS，后端在本机运行，推荐使用专用开发环境 compose（已把 PostgreSQL/Redis/RustFS 的端口映射到宿主机，并自动初始化 RustFS）：

```bash
# 一键开发：首次自动生成 .env.dev 并启动依赖容器 + 后端热重载
make dev
# 说明：后端热重载依赖 air，未安装时先执行 make dev-install

# 仅启动依赖容器（不带后端）
docker compose -f deploy/docker/docker-compose.dev.yml up -d

# 停止依赖容器
make dev-down
```

开发依赖端口（默认）：PostgreSQL `5432`、Redis `6379`、RustFS API `9000` / 控制台 `9001`。

`.env.dev` 与生产 `.env` 的差异在于主机名指向宿主机（`POSTGRES_HOST/REDIS_HOST/RUSTFS_HOST=localhost`），`rustfs-init` 容器内部会自动改用容器服务名访问 RustFS。

也可继续用旧的临时方式：取消 `docker-compose.yml` 中 `postgres` 的 `5432:5432` 和 `rustfs` 的 `9000:9000` 端口注释，然后只启动依赖服务：

```bash
docker compose -f deploy/docker/docker-compose.yml up -d postgres redis rustfs rustfs-init
```

开发模式下后端需要以下环境变量：

```bash
export POSTGRES_HOST=127.0.0.1
export POSTGRES_PORT=5432
export POSTGRES_USER=admin
export POSTGRES_PASSWORD=admin2025
export POSTGRES_DB=do-exercise

export RUSTFS_HOST=127.0.0.1
export RUSTFS_PORT=9000
export RUSTFS_BUCKET_NAME=do-exercise
export RUSTFS_APP_ACCESS_KEY=your-access-key
export RUSTFS_APP_SECRET_KEY=your-secret-key
export RUSTFS_PRESIGNED_HOST=/oss
```

`RUSTFS_PRESIGNED_HOST=/oss` 适配前端 Vite 开发代理；前端请求 `/oss` 时会转发到本机 RustFS API `127.0.0.1:9000`。

### 启动与验证

```bash
docker compose up -d --build
docker logs -f rustfs_init
docker logs -f server_data
```

`rustfs_init` 正常退出后，server 才会启动。首次部署不再需要手动登录 RustFS Console 创建 Bucket、Anonymous Access 或 Access Keys。

### 常见问题

**Q: server 仍然提示 RustFS 连接失败？**  
A: 先看 `docker logs rustfs_init`，通常是 `.env` 缺少 `RUSTFS_BUCKET_NAME`、`RUSTFS_APP_ACCESS_KEY` 或 `RUSTFS_APP_SECRET_KEY`。

**Q: 修改了 Access Key 后仍然无法上传？**  
A: 如果同一个 bucket/volume 已经初始化过旧凭据，需要删除旧的 Service Account 后重启 `rustfs-init`，或清理 RustFS 数据卷重新初始化。

**Q: 上传地址仍然返回 `rustfs:9000`？**  
A: 确认 `server/config.yaml` 中 `rustfs.presigned_host` 为 `/oss`，并重新构建 server 镜像。
