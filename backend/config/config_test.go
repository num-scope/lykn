package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

const validFileConfig = `
server:
  port: 9090
  allow_origin: http://localhost:5173
database:
  host: localhost
  port: 5432
  dbname: testdb
  username: testuser
  password: testpass
auth:
  secret_key: file-auth-secret
  expire: 12h
encryption:
  secret_key: file-encryption-secret
`

func writeConfig(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(cfgPath, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return cfgPath
}

func TestLoad_ReadsAllConfigFromYAML(t *testing.T) {
	cfg, err := Load(writeConfig(t, validFileConfig))
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.Server.Port != 9090 || cfg.Server.AllowOrigin != "http://localhost:5173" {
		t.Fatalf("unexpected server config: %+v", cfg.Server)
	}
	if cfg.Database.Host != "localhost" || cfg.Database.Port != 5432 {
		t.Fatalf("unexpected database config: %+v", cfg.Database)
	}
	if cfg.Auth.SecretKey != "file-auth-secret" || cfg.Encryption.SecretKey != "file-encryption-secret" {
		t.Fatalf("unexpected secrets: auth=%q encryption=%q", cfg.Auth.SecretKey, cfg.Encryption.SecretKey)
	}
	ttl, err := cfg.Auth.ExpireDuration()
	if err != nil {
		t.Fatalf("ExpireDuration: %v", err)
	}
	if ttl != 12*time.Hour {
		t.Fatalf("ttl = %s, want 12h", ttl)
	}
}

func TestLoad_RejectsMissingOrInvalidFileConfig(t *testing.T) {
	cases := []struct {
		name    string
		path    string
		content string
		want    string
	}{
		{name: "missing file", path: filepath.Join(t.TempDir(), "missing.yaml"), want: "read config"},
		{name: "missing required fields", content: `
server:
  port: 8080
database:
  host: ""
auth:
  expire: 24h
`, want: "server.allow_origin"},
		{name: "invalid auth expire", content: `
server:
  port: 8080
  allow_origin: "*"
database:
  host: localhost
  port: 5432
  dbname: testdb
  username: testuser
  password: testpass
auth:
  secret_key: file-auth-secret
  expire: never
encryption:
  secret_key: file-encryption-secret
`, want: "auth.expire is invalid"},
		{name: "missing auth secret", content: `
server:
  port: 8080
  allow_origin: "*"
database:
  host: localhost
  port: 5432
  dbname: testdb
  username: testuser
  password: testpass
auth:
  expire: 24h
encryption:
  secret_key: file-encryption-secret
`, want: "auth.secret_key"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			path := tc.path
			if path == "" {
				path = writeConfig(t, tc.content)
			}
			_, err := Load(path)
			if err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("error = %v, want containing %q", err, tc.want)
			}
		})
	}
}

func TestLoad_DoesNotReadEnvironmentOverrides(t *testing.T) {
	t.Setenv("IGNORED_DATABASE_HOST", "env-db-host")
	t.Setenv("IGNORED_AUTH_SECRET_KEY", "env-auth-secret")
	t.Setenv("IGNORED_ENCRYPTION_SECRET_KEY", "env-encryption-secret")

	cfg, err := Load(writeConfig(t, validFileConfig))
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.Database.Host != "localhost" || cfg.Auth.SecretKey != "file-auth-secret" || cfg.Encryption.SecretKey != "file-encryption-secret" {
		t.Fatalf("environment should be ignored: db=%q auth=%q encryption=%q", cfg.Database.Host, cfg.Auth.SecretKey, cfg.Encryption.SecretKey)
	}
}
