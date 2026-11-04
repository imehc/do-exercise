#!/usr/bin/env bash
set -euo pipefail
umask 077

# All writers and object-upload ingress must be fenced before invoking this tool.
if [[ $# != 2 || $1 != --quiesced ]]; then
  echo 'Usage: PG_CONTAINER=... REDIS_CONTAINER=... OSS_CONTAINER=... DB_USER=... DB_NAME=... ./backup.sh --quiesced /absolute/backup-directory' >&2
  exit 2
fi
: "${PG_CONTAINER:?}" "${REDIS_CONTAINER:?}" "${OSS_CONTAINER:?}" "${DB_USER:?}" "${DB_NAME:?}"
for data_container in "$REDIS_CONTAINER" "$OSS_CONTAINER"; do
  data_mount=$(docker inspect -f '{{range .Mounts}}{{if eq .Destination "/data"}}{{.Destination}}{{end}}{{end}}' "$data_container")
  [[ $data_mount = /data ]] || { echo 'Redis and object storage require a mounted /data volume' >&2; exit 2; }
done
backup_dir=$2
[[ $backup_dir = /* && ! -e $backup_dir ]] || { echo 'Destination must be a new absolute directory' >&2; exit 2; }
mkdir -m 700 "$backup_dir"
redis_stopped=0
oss_stopped=0
resume() {
  if [[ $redis_stopped = 1 ]]; then docker start "$REDIS_CONTAINER" >/dev/null; fi
  if [[ $oss_stopped = 1 ]]; then docker start "$OSS_CONTAINER" >/dev/null; fi
}
trap resume EXIT
for app_container in ${APP_CONTAINERS:-}; do
  [[ $(docker inspect -f '{{.State.Running}}' "$app_container") = false ]] || { echo "Application still running: $app_container" >&2; exit 2; }
done
pg_image=$(docker inspect -f '{{.Image}}' "$PG_CONTAINER")
redis_image=$(docker inspect -f '{{.Image}}' "$REDIS_CONTAINER")
oss_image=$(docker inspect -f '{{.Image}}' "$OSS_CONTAINER")
# Stopping S3 closes outstanding POST requests before the database snapshot.
docker stop "$OSS_CONTAINER" >/dev/null
oss_stopped=1
if [[ -n ${REDIS_PASSWORD:-} ]]; then
  docker exec -e REDISCLI_AUTH="$REDIS_PASSWORD" "$REDIS_CONTAINER" redis-cli SAVE >/dev/null
else
  docker exec "$REDIS_CONTAINER" redis-cli SAVE >/dev/null
fi
docker stop "$REDIS_CONTAINER" >/dev/null
redis_stopped=1
docker exec "$PG_CONTAINER" pg_dump -U "$DB_USER" -d "$DB_NAME" -Fc > "$backup_dir/postgres.dump"
docker run --rm --user 0 --volumes-from "$REDIS_CONTAINER":ro --entrypoint tar "$redis_image" -C /data -cf - . > "$backup_dir/redis.tar"
docker run --rm --user 0 --volumes-from "$OSS_CONTAINER":ro --entrypoint tar "$oss_image" -C /data -cf - . > "$backup_dir/oss.tar"
printf '%s\n' "$pg_image" "$redis_image" "$oss_image" > "$backup_dir/images.txt"
printf 'database=%s\ncreated_utc=%s\n' "$DB_NAME" "$(date -u +%FT%TZ)" > "$backup_dir/metadata.txt"
(cd "$backup_dir" && shasum -a 256 postgres.dump redis.tar oss.tar images.txt metadata.txt > SHA256SUMS)
echo "Backup created: $backup_dir"
