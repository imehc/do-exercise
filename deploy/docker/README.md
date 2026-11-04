# Docker 部署说明

## 配置

将 `.env.example` 复制为 `.env`，填写实际环境配置；将 `server/config.example.yaml` 复制为 `server/config.yaml`。以下命令在项目根目录执行。

```bash
# 已有配置文件时保留原文件
cp -n deploy/docker/.env.example deploy/docker/.env
cp -n server/config.example.yaml server/config.yaml
```

Compose 从 `.env` 获取变量，再按服务注入容器。后端读取进程环境变量，本地直接运行时需要手动加载环境文件或通过 IDE 注入。环境变量优先于 `server/config.yaml` 中的对应配置。

| 配置项 | 用途 |
| --- | --- |
| `POSTGRES_HOST`、`POSTGRES_PORT` | PostgreSQL 地址及端口，容器内通常为 `postgres:5432` |
| `POSTGRES_USER`、`POSTGRES_PASSWORD`、`POSTGRES_DB` | 数据库账号、密码及库名 |
| `REDIS_HOST`、`REDIS_PORT`、`REDIS_PASSWORD`、`REDIS_DATABASE` | Redis 地址、端口、密码及逻辑库 |
| `BOOTSTRAP_ADMIN_PASSWORD` | 初始管理员密码，12–72 字节，包含大小写字母、数字和符号 |
| `OSS_ROOT_ACCESS_KEY`、`OSS_ROOT_SECRET_KEY` | RustFS 初始化使用的 root 凭据 |
| `OSS_RPC_SECRET` | 节点间 RPC 密钥，使用强随机值且与 root 密码不同 |
| `OSS_HOST`、`OSS_PORT` | 对象存储地址及端口，容器内通常为 `oss:9000` |
| `OSS_BUCKET_NAME` | 存储桶名称 |
| `OSS_APP_ACCESS_KEY`、`OSS_APP_SECRET_KEY` | 后端使用的应用凭据 |
| `OSS_PRESIGNED_HOST` | 浏览器访问对象存储的地址，生产反代路径为 `/oss` |
| `EMAIL_HOST`、`EMAIL_PORT`、`EMAIL_USER`、`EMAIL_PASS` | 邮箱服务配置 |

应用连接池、超时和日志等参数放在 `server/config.yaml`。部署凭据放在环境变量中，不要写入镜像内的配置文件或提交到仓库。RustFS 容器内部的 `RUSTFS_*` 变量由 Compose 映射，应用配置统一使用 `OSS_*`。

默认 Compose 使用内部服务地址 `postgres:5432`、`redis:6379` 和 `oss:9000`，浏览器对象存储路径为 `/oss`。生产环境设置独立的强随机密码与密钥，并填写 `BOOTSTRAP_ADMIN_PASSWORD`。`REDIS_PASSWORD` 同时配置 Redis 服务与后端，Redis 启用 AOF 持久化。邮箱功能需要有效的 SMTP 配置。

构建后端需要访问 Go 模块源。默认 `GOPROXY=https://proxy.golang.org,direct`；该源不可达时可在 `.env` 中设置 `GOPROXY=https://goproxy.cn,direct`，模块校验仍遵循 `go.sum` 和 Go 校验数据库。

## 证书与启动

将可信证书和私钥放入 `deploy/docker/nginx/certs`，文件名为 `server.crt` 和 `server.key`。本地开发可以使用 `make certs` 生成自签名证书。

证书和私钥只在运行时挂载，不写入前端镜像。使用正式域名时，将 `deploy/docker/nginx/nginx.conf` 中的 `server_name localhost` 改为实际域名，并使用匹配该域名的证书。

```bash
docker compose --env-file deploy/docker/.env -f deploy/docker/docker-compose.yml up -d --build
docker compose --env-file deploy/docker/.env -f deploy/docker/docker-compose.yml ps
docker compose --env-file deploy/docker/.env -f deploy/docker/docker-compose.yml logs -f server oss-init
```

PostgreSQL、Redis 和 RustFS 就绪后，`oss-init` 创建存储桶及应用凭据，完成后后端才启动。后端镜像先执行数据库迁移，再启动服务；后端健康检查通过后启动 Web 入口。首次部署不需要手动导入 SQL，也不需要另装 Kafka 或 RabbitMQ。首次对象存储初始化可能扫描既有文件，应预留时间和存储空间。

同样的启动命令可以简写为 `make up`。初始平台管理员用户名为 `superAdmin`，登录时租户编码留空，密码为 `BOOTSTRAP_ADMIN_PASSWORD`，首次登录必须修改密码。

浏览器入口为 `https://localhost`，HTTP 80 端口跳转到 HTTPS。对象存储数据访问通过 HTTPS 下的 `/oss/` 路径进行。RustFS 控制台仅绑定宿主机 `127.0.0.1:9001`，远程管理使用可信隧道。

业务文件保持私有，只有公开头像目录允许匿名读取。文件上传、下载和账号权限的使用规则见 [多租户说明与注意事项](../../MULTI_TENANCY.md)。

## 本地开发

将 `.env.dev.example` 复制为 `.env.dev`，设置管理员初始密码和实际凭据。本地后端的 `POSTGRES_HOST`、`REDIS_HOST`、`OSS_HOST` 使用 `127.0.0.1` 或 `localhost`；`oss-init` 在容器内使用 `oss` 服务地址。

```bash
# 仅启动依赖服务
docker compose --env-file deploy/docker/.env.dev -f deploy/docker/docker-compose.dev.yml up -d

# 依赖服务就绪后启动后端热重载
make dev
```

后端热重载依赖 `air`，可以在 `server` 目录执行 `make dev-install` 安装。默认开发端口为 PostgreSQL `5432`、Redis `6379`、RustFS API `9000` 和控制台 `9001`。

直接运行后端时，先向进程注入环境变量，再执行数据库迁移。Vite 的 `/oss` 开发代理转发到 `http://127.0.0.1:9000`；浏览器必须能够访问返回的上传或下载地址。

## 注意事项

- 初始化后的应用 Access Key 需要与服务端凭据保持一致。变更凭据时应管理对应的服务账号，并更新后端环境变量。
- RustFS 应用 Access Key 为 3–128 字节，Secret Key 为 8–40 字节；建议使用 32 字符随机密钥。
- 修改镜像内的配置文件后需要重新构建镜像；仅修改环境变量时需要重新创建对应容器。
- 持久卷保存数据库、Redis 和对象存储数据。删除卷会丢失数据，清理前应确认备份。
- 生产环境使用强随机凭据、可信证书和受限管理网络，并为持久数据安排备份与恢复流程。

## 备份与整体恢复

`backup.sh` 生成 PostgreSQL custom dump、Redis 数据目录归档、对象存储数据目录归档，以及镜像标识和 SHA-256 清单。Redis 与对象存储必须把 `/data` 挂到持久化卷。对象目录归档适用于当前单节点部署，不能对持续写入的多节点对象集群直接使用。

开始前应关闭入口和所有后端实例，停止其他数据库、Redis、对象存储写入者。`--quiesced` 表示调用者已经完成写入隔离；工具会检查 `APP_CONTAINERS` 中的容器均已停止，并在采集期间停止 Redis 和对象存储，结束后恢复这两个服务。入口和后端由操作者在操作完成后重新开启。

```bash
docker compose --env-file .env -f docker-compose.yml stop web server
PG_CONTAINER=postgres_data REDIS_CONTAINER=redis_data OSS_CONTAINER=oss_data \
  DB_USER=admin DB_NAME=do-exercise APP_CONTAINERS=server_data \
  ./backup.sh --quiesced /srv/backups/do-exercise-20261103
```

设置了 Redis 密码时，另传 `REDIS_PASSWORD`。备份期间数据库保持启动，但禁止业务写入。备份目录权限为 0700，包含用户密码哈希、会话凭据和文件内容；应加密保存并复制到独立存储。应用环境变量、TLS 证书、对象存储管理凭据、Redis 启动参数和部署配置需另外保管，归档不包含这些配置。

恢复应使用独立的空数据卷和空数据库，保持入口与后端关闭。目标 PostgreSQL 先完成初始化并启动；目标 Redis 和对象存储容器先创建、保持停止，`/data` 必须为空。三个目标容器的镜像标识必须与备份中的标识相同，Redis 持久化启动参数和对象存储凭据也应一致。

```bash
PG_CONTAINER=postgres-restore REDIS_CONTAINER=redis-restore OSS_CONTAINER=oss-restore \
  DB_USER=admin DB_NAME=do-exercise \
  ./restore.sh --empty-targets /srv/backups/do-exercise-20261103
```

工具先核对校验和、镜像标识和目标空目录，再恢复数据库与归档。完成后 Redis 和对象存储启动，后端和入口保持关闭。中途失败时保留目标状态；重新操作应换用新的空目标，不能在部分恢复的数据上重复覆盖。

开放入口前确认数据库迁移版本、租户及成员关系、角色权限、存储占用、文件字节内容与私有访问规则。恢复的任务队列会在后端启动后继续处理，应先决定是否继续队列中的工作。恢复仅覆盖整个应用数据集，不提供按租户选择性恢复。

Redis 归档可能包含备份时仍有效的会话，备份之后的吊销状态不会自动出现。通常在恢复后、开放入口前清空本应用专用 Redis 数据库并要求重新登录；确认所选数据库没有其他应用数据后再操作。用户密码及角色变化同样只能恢复到备份时的状态。数据恢复点是最近一次完整备份，停写时长、数据量和复制速度共同决定恢复耗时；需要更短恢复点时另配置 PostgreSQL WAL 归档及对象存储版本/复制策略。
