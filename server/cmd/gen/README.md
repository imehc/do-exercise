# cmd/gen 代码生成器

根据实体名与字段定义，一键生成一套完整的 CRUD 六件套，并自动完成注册，免去新表接管的样板代码。

## 用法

在 `server/` 目录下执行：

```bash
# 直接使用 go run
go run ./cmd/gen -name=SysDevice -desc=设备 -path=devices \
  -field=code:string:设备编码 \
  -field=name:string:设备名称 \
  -field=enabled:bool:是否启用 \
  -field=price:float:单价

# 或通过 Makefile（推荐）
make gen NAME=SysDevice DESC=设备 ROUTE=devices \
  FIELDS="code:string:设备编码 name:string:设备名称 enabled:bool:是否启用 price:float:单价"
```

### 参数说明

| 参数 | 必填 | 说明 |
| --- | --- | --- |
| `-name` | 是 | 实体结构体名，Pascal 风格，如 `SysDevice`、`Device` |
| `-desc` | 否 | 实体中文描述，写入注释，如 `设备` |
| `-path` | 否 | 路由分组路径，默认取实体名的短横线形式（`SysDevice` → `sys-device`） |
| `-field` | 否 | 字段定义，可多次指定，格式 `name:type[:注释]`；见下方「字段类型」 |
| `-dir` | 否 | server 根目录，默认当前目录 |

> 说明：`-path`（即 Makefile 的 `ROUTE`）是路由分组路径片段，如 `devices`，最终生成 `/system/devices`。Makefile 参数名刻意用 `ROUTE`，避免与 make 内置的 `PATH` 变量冲突。

### 字段类型

| type | Go 类型 | gorm 标签 |
| --- | --- | --- |
| `string` | `string` | `size:255` |
| `text` | `string` | `type:text` |
| `bool` | `bool` | `default:false` |
| `int` | `int` | — |
| `uint` | `uint` | — |
| `int64` | `int64` | — |
| `float` | `float64` | — |
| `time` | `time.Time` | — |
| `id_array` | `[]uint` | — |

## 生成内容

以 `-name=SysDevice -path=devices` 为例：

| 文件 | 内容 |
| --- | --- |
| `model/system/sys-device.go` | 模型（`IdWrapper` + 自定义字段 + 时间戳 + 软删除） |
| `model/system/request/sys-device.go` | `CreateSysDeviceReq` / `UpdateSysDeviceReq` |
| `model/system/response/sys-device.go` | `SysDeviceResp` |
| `service/system/sys-device.go` | `SysDeviceService`：Create / Update / Get / GetList / Delete（均以 `db *gorm.DB` 为首参） |
| `api/v1/system/sys-device.go` | `SysDeviceApi`：REST 处理器，使用 `util.DB(ctx)` |
| `router/system/sys-device.go` | `SysDeviceRouter`：注册以下路由 |

自动生成的路由（挂在 `/system` 受保护组下）：

```
POST   /system/devices        # 创建
PUT    /system/devices/:id    # 更新
GET    /system/devices/:id    # 详情
GET    /system/devices        # 分页列表（page / page_size）
DELETE /system/devices/:id    # 删除
```

同时自动注册（幂等，重复执行不会重复插入）：

- `service/system/enter.go`：ServiceGroup 增加 `SysDeviceService`
- `api/v1/system/enter.go`：ApiGroup 增加 `SysDeviceApi` + `sysDeviceService` 变量
- `router/system/enter.go`：RouterGroup 增加 `SysDeviceRouter` + `sysDeviceApi` 变量
- `router/router.go`：挂载 `system.InitSysDeviceRouter(protected)`
- `internal/gorm.go`：AutoMigrate 增加 `system.SysDevice{}`

## 注意事项

- 生成的文件会被**直接覆盖**（`已生成`），自定义改动会在重新生成时丢失；注册类文件为幂等追加（`已更新`）。
- 生成后需执行 `go build ./...` 验证编译，并按需补充业务逻辑（如字段校验、过滤条件、审计钩子）。
- 新表通过 `AutoMigrate` 自动建表；新路由默认受 casbin 保护，需在后台为角色分配对应权限后方可访问。
- 约定遵循现有代码风格：服务方法 `db *gorm.DB` 首参、API 层用 `util.DB(ctx)` 取请求级连接、错误信息走 `errors.New("xxxKey")` + i18n。
