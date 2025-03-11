## 生成依赖包

```bash
go mod vendor   
# 使用 go mod 并安装go依赖包
go generate 
```
## swagger

<!-- date: 2025-01-16 -->
<!-- 已知存在问题： -->
<!-- github.com/swaggo/swag/v2 v2.0.0-rc3 版本 不支持bearer认证 -->
<!-- github.com/swaggo/swag/v2 v2.0.0-rc4 版本 支持bearer认证但是重命名丢失 -->
<!-- swagger页面显示报错提示找不到doc.json -->
<!-- 相关issue: https://github.com/swaggo/swag/issues/1909 -->
```bash
go get -u github.com/swaggo/swag  
swag init --v3.1 
```

## 数据库
### 使用docker
``` bash
docker run --name do_exercise \
  -e POSTGRES_USER=[用户名] \
  -e POSTGRES_PASSWORD=[密码] \
  -e POSTGRES_DB=[数据库] \
  -e TZ='Asia/Shanghai' \
  -e ALLOW_IP_RANGE=0.0.0.0/0 \
  # -v 数据卷挂载
  # -v /Users/tom/person/learning/postgresql-data:/var/lib/postgresql \
  -p 5432:5432 \
  --restart always \
  -d postgres 
```

### 初始化数据
在启动服务时添加 `-init-data` 参数来初始化数据库：
```bash
go run main.go -init-data
```
该命令会：
1. 自动创建必要的数据库表
2. 初始化基础数据（如系统配置、默认角色等）
3. 完成数据库的初始化设置