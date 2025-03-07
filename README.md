# do-exercise

## 项目概述

这是一个基于前后端分离架构的现代化Web应用项目。项目采用Go语言开发后端API服务，Next.js构建前端用户界面，提供了完整的用户认证、国际化支持等功能。

## 技术栈

### 后端 (server)
- Go 语言
- Gin Web框架
- GORM ORM框架
- PostgreSQL 数据库
- Redis 缓存
- Swagger API文档
- JWT认证

### 前端 (web)
- Next.js
- TypeScript
- Tailwind CSS
- Shadcn UI组件库
- next-intl 国际化

## 环境要求

### 后端环境
- Go 1.21+
- PostgreSQL
- Redis

### 前端环境
- Node.js 18+
- pnpm

## 项目结构

```
├── server/          # 后端服务
│   ├── api/         # API接口
│   ├── config/      # 配置文件
│   ├── core/        # 核心代码
│   ├── docs/        # Swagger文档
│   ├── middleware/  # 中间件
│   ├── model/       # 数据模型
│   ├── router/      # 路由
│   └── service/     # 业务逻辑
│
└── web/            # 前端应用
    ├── app/         # 页面组件
    ├── components/  # 通用组件
    ├── helper/      # 工具函数
    ├── i18n/        # 国际化配置
    └── public/      # 静态资源
```

## 快速开始

### 后端服务启动

1. 安装依赖
```bash
cd server
go mod vendor   # 安装依赖包
go generate     # 生成所需文件
```

2. 配置环境
- 复制 `.env.example` 为 `.env`
- 修改数据库、Redis等配置

3. 生成Swagger文档
```bash
go get -u github.com/swaggo/swag
swag init --v3.1
```

4. 启动服务
```bash
go run main.go
```

### 前端应用启动

1. 安装依赖
```bash
cd web
pnpm install
```

2. 启动开发服务器
```bash
pnpm dev
```

## 开发指南

### 添加新的UI组件
```bash
pnpm dlx shadcn@latest add [组件名称]
```

### 国际化支持
项目使用next-intl实现国际化，详细文档请参考：[next-intl文档](https://next-intl.dev/)

## 部署

### 后端部署
1. 编译
```bash
go build -o app
```

2. 运行
```bash
./app
```

### 前端部署
1. 构建
```bash
pnpm build
```

2. 运行
```bash
pnpm start
```

## 贡献指南

1. Fork 项目
2. 创建功能分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 提交Pull Request

## 许可证

本项目基于 MIT 许可证开源 - 查看 [LICENSE](LICENSE) 文件了解更多细节