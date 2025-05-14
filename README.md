## 注意事项

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
