## 注意事项

> **提示**：如果启动容器或服务连接失败，请优先检查各服务的配置文件（如 .env、config.yaml 等）是否填写正确，端口、数据库、MinIO、Redis 等连接信息需与实际环境一致。

### 配置文件准备
在执行 docker-compose 之前，请确保完成以下配置文件的准备：

1. 环境变量配置
   - 将 `deploy/docker/.env.example` 重命名为 `.env`
   - 根据实际需求填写 `.env` 文件中的配置项

2. 服务配置
   - 将 `server/config.example.yaml` 重命名为 `config.yaml`
   - 根据实际需求修改 `config.yaml` 中的配置项

### MinIO 配置相关
如果执行 docker-compose 后容器无法正常启动，请检查以下事项：

1. 检查 MinIO Access Keys 是否已配置
   - 如果 MinIO 服务无法正常启动，很可能是因为未配置 Access Keys
   - 请确保已正确配置 MinIO 的 Access Key 和 Secret Key
   - 可以通过 MinIO 控制台进行配置

如需配置 MinIO Access Keys：
1. 访问 MinIO 控制台
2. 登录后进入设置页面
3. 在 Access Keys 部分添加新的访问密钥
4. 保存配置后重启容器
