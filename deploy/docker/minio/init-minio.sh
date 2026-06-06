#!/bin/sh
set -eu

: "${MINIO_ROOT_USER:?MINIO_ROOT_USER is required}"
: "${MINIO_ROOT_PASSWORD:?MINIO_ROOT_PASSWORD is required}"
: "${MINIO_HOST:?MINIO_HOST is required}"
: "${MINIO_PORT:?MINIO_PORT is required}"
: "${MINIO_BUCKET_NAME:?MINIO_BUCKET_NAME is required}"
: "${MINIO_APP_ACCESS_KEY:?MINIO_APP_ACCESS_KEY is required}"
: "${MINIO_APP_SECRET_KEY:?MINIO_APP_SECRET_KEY is required}"

until mc alias set local "http://${MINIO_HOST}:${MINIO_PORT}" "${MINIO_ROOT_USER}" "${MINIO_ROOT_PASSWORD}"; do
  sleep 2
done

mc mb --ignore-existing "local/${MINIO_BUCKET_NAME}"
mc anonymous set download "local/${MINIO_BUCKET_NAME}"
mc admin user svcacct add local "${MINIO_ROOT_USER}" --access-key "${MINIO_APP_ACCESS_KEY}" --secret-key "${MINIO_APP_SECRET_KEY}" || true
