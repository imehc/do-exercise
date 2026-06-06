package internal

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/imehc/do-exercise/server/global"
)

func TestInitConfigLoadsServiceConfigFromEnv(t *testing.T) {
	originalConfig := global.Config
	t.Cleanup(func() {
		global.Config = originalConfig
	})

	configPath := filepath.Join(t.TempDir(), "config.yaml")
	configData := []byte(`
database:
  pool:
    max_connections: 100
minio:
  expires: 3600
`)
	if err := os.WriteFile(configPath, configData, 0o644); err != nil {
		t.Fatal(err)
	}

	t.Setenv("POSTGRES_HOST", "env-postgres")
	t.Setenv("POSTGRES_PORT", "15432")
	t.Setenv("POSTGRES_DB", "env-db")
	t.Setenv("POSTGRES_USER", "env-user")
	t.Setenv("POSTGRES_PASSWORD", "env-password")
	t.Setenv("MINIO_HOST", "env-minio")
	t.Setenv("MINIO_PORT", "19000")
	t.Setenv("MINIO_BUCKET_NAME", "env-bucket")
	t.Setenv("MINIO_APP_ACCESS_KEY", "env-access")
	t.Setenv("MINIO_APP_SECRET_KEY", "env-secret")
	t.Setenv("MINIO_PRESIGNED_HOST", "/env-oss")

	InitConfig(configPath)

	if global.Config.Database.Host != "env-postgres" {
		t.Fatalf("expected database host from env, got %q", global.Config.Database.Host)
	}
	if global.Config.Database.Port != 15432 {
		t.Fatalf("expected database port from env, got %d", global.Config.Database.Port)
	}
	if global.Config.Database.Database != "env-db" {
		t.Fatalf("expected database name from env, got %q", global.Config.Database.Database)
	}
	if global.Config.Database.Username != "env-user" {
		t.Fatalf("expected database username from env, got %q", global.Config.Database.Username)
	}
	if global.Config.Database.Password != "env-password" {
		t.Fatalf("expected database password from env, got %q", global.Config.Database.Password)
	}
	if global.Config.Minio.Host != "env-minio" {
		t.Fatalf("expected minio host from env, got %q", global.Config.Minio.Host)
	}
	if global.Config.Minio.Port != 19000 {
		t.Fatalf("expected minio port from env, got %d", global.Config.Minio.Port)
	}
	if global.Config.Minio.BucketName != "env-bucket" {
		t.Fatalf("expected bucket from env, got %q", global.Config.Minio.BucketName)
	}
	if global.Config.Minio.AccessKey != "env-access" {
		t.Fatalf("expected access key from env, got %q", global.Config.Minio.AccessKey)
	}
	if global.Config.Minio.SecretKey != "env-secret" {
		t.Fatalf("expected secret key from env, got %q", global.Config.Minio.SecretKey)
	}
	if global.Config.Minio.PresignedHost != "/env-oss" {
		t.Fatalf("expected presigned host from env, got %q", global.Config.Minio.PresignedHost)
	}
	if global.Config.Minio.Expires != 3600 {
		t.Fatalf("expected minio expires from config, got %d", global.Config.Minio.Expires)
	}
}
