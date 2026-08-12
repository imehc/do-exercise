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
auth:
  access_expire_time: 2h
  refresh_expire_time: 7d
database:
  pool:
    max_connections: 25
oss:
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
	t.Setenv("OSS_HOST", "env-oss")
	t.Setenv("OSS_PORT", "19000")
	t.Setenv("OSS_BUCKET_NAME", "env-bucket")
	t.Setenv("OSS_APP_ACCESS_KEY", "env-access")
	t.Setenv("OSS_APP_SECRET_KEY", "env-secret")
	t.Setenv("OSS_PRESIGNED_HOST", "/env-oss")
	t.Setenv("EMAIL_HOST", "env-smtp")
	t.Setenv("EMAIL_PORT", "1587")
	t.Setenv("EMAIL_USER", "env-mail-user")
	t.Setenv("EMAIL_PASS", "env-mail-pass")

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
	if global.Config.Oss.Host != "env-oss" {
		t.Fatalf("expected oss host from env, got %q", global.Config.Oss.Host)
	}
	if global.Config.Oss.Port != 19000 {
		t.Fatalf("expected oss port from env, got %d", global.Config.Oss.Port)
	}
	if global.Config.Oss.BucketName != "env-bucket" {
		t.Fatalf("expected bucket from env, got %q", global.Config.Oss.BucketName)
	}
	if global.Config.Oss.AccessKey != "env-access" {
		t.Fatalf("expected access key from env, got %q", global.Config.Oss.AccessKey)
	}
	if global.Config.Oss.SecretKey != "env-secret" {
		t.Fatalf("expected secret key from env, got %q", global.Config.Oss.SecretKey)
	}
	if global.Config.Oss.PresignedHost != "/env-oss" {
		t.Fatalf("expected presigned host from env, got %q", global.Config.Oss.PresignedHost)
	}
	if global.Config.Oss.Expires != 3600 {
		t.Fatalf("expected oss expires from config, got %d", global.Config.Oss.Expires)
	}
	// 邮箱凭据必须可由环境变量注入，否则只能写进会被打进镜像层的 config.yaml
	if global.Config.Email.Host != "env-smtp" {
		t.Fatalf("expected email host from env, got %q", global.Config.Email.Host)
	}
	if global.Config.Email.Port != 1587 {
		t.Fatalf("expected email port from env, got %d", global.Config.Email.Port)
	}
	if global.Config.Email.User != "env-mail-user" {
		t.Fatalf("expected email user from env, got %q", global.Config.Email.User)
	}
	if global.Config.Email.Pass != "env-mail-pass" {
		t.Fatalf("expected email pass from env, got %q", global.Config.Email.Pass)
	}
}
