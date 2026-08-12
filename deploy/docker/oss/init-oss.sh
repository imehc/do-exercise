#!/bin/sh
set -eu

: "${OSS_ROOT_ACCESS_KEY:?OSS_ROOT_ACCESS_KEY is required}"
: "${OSS_ROOT_SECRET_KEY:?OSS_ROOT_SECRET_KEY is required}"
: "${OSS_HOST:?OSS_HOST is required}"
: "${OSS_PORT:?OSS_PORT is required}"
: "${OSS_BUCKET_NAME:?OSS_BUCKET_NAME is required}"
: "${OSS_APP_ACCESS_KEY:?OSS_APP_ACCESS_KEY is required}"
: "${OSS_APP_SECRET_KEY:?OSS_APP_SECRET_KEY is required}"

until rc alias set local "http://${OSS_HOST}:${OSS_PORT}" "${OSS_ROOT_ACCESS_KEY}" "${OSS_ROOT_SECRET_KEY}"; do
  sleep 2
done

rc mb --ignore-existing "local/${OSS_BUCKET_NAME}"
rc anonymous set download "local/${OSS_BUCKET_NAME}"
# 应用专用凭证：RustFS 管理面与 MinIO Admin API 不兼容，改用官方 rc CLI（见 docs/minio-to-rustfs-migration.md R2）。
# 幂等处理：首次运行创建 service account；卷已初始化过时（access key 已被占用）
# 视为成功，避免 init 容器在重启/重建后反复失败，阻塞 server 的依赖条件。
if ! output=$(rc admin service-account create local/ "${OSS_APP_ACCESS_KEY}" "${OSS_APP_SECRET_KEY}" \
  --name do-exercise-app 2>&1); then
  case "$output" in
    *"access key is already in use"*)
      echo "Service account '${OSS_APP_ACCESS_KEY}' already exists, skipping." ;;
    *)
      echo "Failed to create service account: $output" >&2
      exit 1
      ;;
  esac
fi