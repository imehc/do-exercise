## MinIO 自动初始化

`docker compose up -d` 会同时启动 `minio-init` 一次性容器，自动完成以下操作：

- 等待 MinIO 服务可连接
- 创建 `.env` 中的 `MINIO_BUCKET_NAME`
- 设置 bucket 匿名只读访问
- 创建 server 使用的应用 Access Key

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

#### MinIO

| 配置项 | 用途 |
| --- | --- |
| `MINIO_ROOT_USER` | MinIO root 用户，用于初始化 |
| `MINIO_ROOT_PASSWORD` | MinIO root 密码，用于初始化 |
| `MINIO_HOST` | MinIO Docker 网络内主机名 |
| `MINIO_PORT` | MinIO Docker 网络内端口 |
| `MINIO_BUCKET_NAME` | 自动创建的 bucket 名称 |
| `MINIO_APP_ACCESS_KEY` | server 上传文件使用的 Access Key |
| `MINIO_APP_SECRET_KEY` | server 上传文件使用的 Secret Key |
| `MINIO_PRESIGNED_HOST` | 浏览器访问 MinIO 的公开反代路径 |

server 会从 `.env` 读取 MinIO 地址、Bucket、Access Key、Secret Key 和公开上传路径，覆盖 `server/config.yaml` 中的同名配置。

### 配置分工

Docker 部署时，`.env` 是运行环境入口，`server/config.yaml` 只保留应用默认参数。

| 配置来源 | 放置内容 |
| --- | --- |
| `deploy/docker/.env` | 数据库地址/端口/库名/账号/密码、MinIO 地址/端口/Bucket/Access Key/Secret Key/公开上传路径 |
| `server/config.yaml` | 数据库连接池、Redis 连接池与超时、日志、验证码、MinIO 预签名有效期、系统默认参数 |

如果同一个配置项同时存在于 `.env` 和 `server/config.yaml`，server 启动时以 `.env` 为准。
`docker-compose.yml` 会按服务显式映射所需变量，不会把整份 `.env` 注入每个容器。

### 开发模式

如果只用 Docker 启动数据库和 MinIO，后端在本机运行，需要让本机能直连依赖服务：

- PostgreSQL 暴露到宿主机 `5432`
- MinIO API 暴露到宿主机 `9000`
- 不启动 `web` 服务，避免它占用宿主机 `9000`

可临时取消 `docker-compose.yml` 中 `postgres` 的 `5432:5432` 和 `minio` 的 `9000:9000` 端口注释，然后只启动依赖服务：

```bash
docker compose -f deploy/docker/docker-compose.yml up -d postgres redis minio minio-init
```

开发模式下后端需要以下环境变量：

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

`MINIO_PRESIGNED_HOST=/oss` 适配前端 Vite 开发代理；前端请求 `/oss` 时会转发到本机 MinIO API `127.0.0.1:9000`。

### 启动与验证

```bash
docker compose up -d --build
docker logs -f minio_init
docker logs -f server_data
```

`minio_init` 正常退出后，server 才会启动。首次部署不再需要手动登录 MinIO Console 创建 Bucket、Anonymous Access 或 Access Keys。

### 常见问题

**Q: server 仍然提示 MinIO 连接失败？**  
A: 先看 `docker logs minio_init`，通常是 `.env` 缺少 `MINIO_BUCKET_NAME`、`MINIO_APP_ACCESS_KEY` 或 `MINIO_APP_SECRET_KEY`。

**Q: 修改了 Access Key 后仍然无法上传？**  
A: 如果同一个 bucket/volume 已经初始化过旧凭据，需要在 MinIO Console 删除旧 Access Key 后重启 `minio-init`，或清理 MinIO 数据卷重新初始化。

**Q: 上传地址仍然返回 `minio:9000`？**  
A: 确认 `server/config.yaml` 中 `minio.presigned_host` 为 `/oss`，并重新构建 server 镜像。
