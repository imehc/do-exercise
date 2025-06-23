# Build stage
FROM node:20-alpine AS builder
WORKDIR /app

# 安装 openjdk
RUN apk add --no-cache openjdk17-jre

# 复制依赖文件并安装依赖
COPY web/package.json web/pnpm-lock.yaml ./
RUN npm install -g pnpm && pnpm install --frozen-lockfile

# 复制所有web源码
COPY web/. .

# 复制openapi.yaml
ARG OPENAPI_YAML
COPY ${OPENAPI_YAML} ../server/openapi.yaml

# 生成API代码
RUN pnpm run gen:apis

# 构建前端
RUN pnpm run build

# Production stage
FROM nginx:alpine
COPY --from=builder /app/dist /usr/share/nginx/html
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"] 