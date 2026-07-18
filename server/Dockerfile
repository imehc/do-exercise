# 1. 构建阶段
FROM golang:1.25-alpine AS builder

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
RUN go build -trimpath -ldflags="-s -w" -o server 

# 2. 运行阶段
FROM busybox:1.37.0-uclibc

# 设置必要的环境变量
ENV GIN_MODE=release \
    TZ=Asia/Shanghai 

WORKDIR /app

# 创建运行时需要的目录
RUN mkdir -p /etc/ssl/certs /usr/share/zoneinfo/Asia /app/logs

# 从构建阶段复制文件
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /usr/share/zoneinfo/Asia/Shanghai
COPY --from=builder /build/server ./
COPY --from=builder /build/config.yaml ./
COPY --from=builder /build/init.sql ./
COPY --from=builder /build/i18n ./i18n
COPY --from=builder /build/ip2region.xdb ./

RUN chown -R 10001:10001 /app/logs

# 暴露应用端口
EXPOSE 6020

# 设置用户权限
USER 10001:10001

# 健康检查
HEALTHCHECK --interval=30s --timeout=30s --start-period=5s --retries=3 \
    CMD wget -q -O /dev/null http://localhost:6020/health || exit 1

# 启动命令
CMD ["sh", "-c", "./server migrate --sql ./init.sql && ./server server"]
