# 1. 构建阶段
FROM golang:1.24-alpine AS builder

# 设置必要的环境变量
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build

# 安装基础工具
RUN apk add --no-cache git ca-certificates tzdata

# 首先复制依赖文件
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 编译可执行文件
RUN go build -ldflags="-s -w" -o server 

# 2. 运行阶段
FROM alpine:latest

# 设置时区
ENV TZ=Asia/Shanghai

# 创建非 root 用户
RUN adduser -D -g '' appuser && \
    apk add --no-cache ca-certificates tzdata

WORKDIR /app

# 从构建阶段复制文件
COPY --from=builder /build/server ./
COPY --from=builder /build/config.yaml ./
COPY --from=builder /build/init.sql ./
COPY --from=builder /build/i18n ./i18n
COPY --from=builder /build/ip2region.xdb ./

# 暴露应用端口
EXPOSE 6020

# 设置用户权限
RUN chown -R appuser:appuser /app
USER appuser

# 健康检查
HEALTHCHECK --interval=30s --timeout=30s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:6020/health || exit 1

# 启动命令
CMD ["sh", "-c", "./server migrate --sql ./init.sql && ./server"]
