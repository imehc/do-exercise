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