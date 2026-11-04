#!/usr/bin/env bash
set -euo pipefail
umask 077
if [[ $# != 2 || $1 != --empty-targets ]]; then
  echo 'Usage: PG_CONTAINER=... REDIS_CONTAINER=... OSS_CONTAINER=... DB_USER=... DB_NAME=... ./restore.sh --empty-targets /absolute/backup-directory' >&2
  exit 2
fi
: "${PG_CONTAINER:?}" "${REDIS_CONTAINER:?}" "${OSS_CONTAINER:?}" "${DB_USER:?}" "${DB_NAME:?}"
for data_container in "$REDIS_CONTAINER" "$OSS_CONTAINER"; do
  data_mount=$(docker inspect -f '{{range .Mounts}}{{if eq .Destination "/data"}}{{.Destination}}{{end}}{{end}}' "$data_container")
  [[ $data_mount = /data ]] || { echo 'Redis and object storage require a mounted /data volume' >&2; exit 2; }
done
backup_dir=$2
[[ $backup_dir = /* && -f $backup_dir/SHA256SUMS ]] || { echo 'Absolute backup directory with SHA256SUMS required' >&2; exit 2; }
(cd "$backup_dir" && shasum -a 256 -c SHA256SUMS)
mapfile_compatible() { while IFS= read -r line; do backup_images+=("$line"); done < "$backup_dir/images.txt"; }
backup_images=();mapfile_compatible
[[ ${#backup_images[@]} = 3 ]] || { echo 'Invalid image manifest' >&2; exit 2; }
containers=("$PG_CONTAINER" "$REDIS_CONTAINER" "$OSS_CONTAINER")
for i in 0 1 2; do
  [[ $(docker inspect -f '{{.Image}}' "${containers[$i]}") = "${backup_images[$i]}" ]] || { echo 'Target image differs from snapshot image' >&2; exit 2; }
done
for app_container in ${APP_CONTAINERS:-}; do
  [[ $(docker inspect -f '{{.State.Running}}' "$app_container") = false ]] || { echo "Application still running: $app_container" >&2; exit 2; }
done
for i in 1 2; do
  [[ $(docker inspect -f '{{.State.Running}}' "${containers[$i]}") = false ]] || { echo 'Redis and object storage targets must be stopped' >&2; exit 2; }
  docker run --rm --user 0 --volumes-from "${containers[$i]}" --entrypoint sh "${backup_images[$i]}" -c 'test -z "$(ls -A /data)"' || { echo 'Target data volume is not empty' >&2; exit 2; }
done
# PostgreSQL must be freshly initialized and contain no application objects.
objects=$(docker exec "$PG_CONTAINER" psql -U "$DB_USER" -d "$DB_NAME" -Atc "SELECT count(*) FROM pg_class c JOIN pg_namespace n ON n.oid=c.relnamespace WHERE n.nspname NOT IN ('pg_catalog','information_schema') AND n.nspname NOT LIKE 'pg_toast%' AND c.relkind IN ('r','p','v','m','S','f')")
[[ $objects = 0 ]] || { echo 'Target database is not empty' >&2; exit 2; }
docker exec -i "$PG_CONTAINER" pg_restore -U "$DB_USER" -d "$DB_NAME" --no-owner --no-privileges --exit-on-error --single-transaction < "$backup_dir/postgres.dump"
docker run --rm -i --user 0 --volumes-from "$REDIS_CONTAINER" --entrypoint tar "${backup_images[1]}" -C /data -xf - < "$backup_dir/redis.tar"
docker run --rm -i --user 0 --volumes-from "$OSS_CONTAINER" --entrypoint tar "${backup_images[2]}" -C /data -xf - < "$backup_dir/oss.tar"
docker start "$REDIS_CONTAINER" "$OSS_CONTAINER" >/dev/null
echo 'Data restored. Keep application traffic fenced until operational checks and session policy are complete.'
