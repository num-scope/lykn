package main

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func writeTestMaterials(t *testing.T) (string, string) {
	t.Helper()
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	pubBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		t.Fatalf("marshal public key: %v", err)
	}
	publicPEM := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: pubBytes})
	now := time.Now().UTC()
	payload, err := json.Marshal(map[string]any{
		"id":         "lic-cli-001",
		"version":    1,
		"subject":    map[string]any{"name": "CLI Demo"},
		"issued_at":  now.Format(time.RFC3339),
		"not_before": now.Add(-time.Minute).Format(time.RFC3339),
		"not_after":  now.Add(time.Hour).Format(time.RFC3339),
		"features":   []string{"reports"},
	})
	if err != nil {
		t.Fatalf("marshal payload: %v", err)
	}
	hash := sha256.Sum256(payload)
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hash[:])
	if err != nil {
		t.Fatalf("sign payload: %v", err)
	}
	licenseContent, err := json.Marshal(map[string]string{
		"payload":   base64.StdEncoding.EncodeToString(payload),
		"signature": base64.StdEncoding.EncodeToString(signature),
	})
	if err != nil {
		t.Fatalf("marshal license: %v", err)
	}
	dir := t.TempDir()
	publicKeyPath := filepath.Join(dir, "public.pem")
	licensePath := filepath.Join(dir, "license.lic")
	if err := os.WriteFile(publicKeyPath, publicPEM, 0o644); err != nil {
		t.Fatalf("write public key: %v", err)
	}
	if err := os.WriteFile(licensePath, licenseContent, 0o644); err != nil {
		t.Fatalf("write license: %v", err)
	}
	return publicKeyPath, licensePath
}

func TestRunVerify(t *testing.T) {
	publicKeyPath, licensePath := writeTestMaterials(t)
	var stdout, stderr bytes.Buffer
	code := run([]string{"verify", "--public-key", publicKeyPath, "--license", licensePath, "--feature", "reports"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("code = %d stderr = %s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "lic-cli-001") {
		t.Fatalf("verify output missing license id: %s", stdout.String())
	}
}

func TestRunInspect(t *testing.T) {
	_, licensePath := writeTestMaterials(t)
	var stdout, stderr bytes.Buffer
	code := run([]string{"inspect", "--license", licensePath}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("code = %d stderr = %s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "lic-cli-001") {
		t.Fatalf("inspect output missing license id: %s", stdout.String())
	}
}

func TestRunRejectsUnknownCommand(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := run([]string{"unknown"}, &stdout, &stderr)
	if code != 2 {
		t.Fatalf("code = %d, want 2", code)
	}
	if !strings.Contains(stderr.String(), "unknown command") {
		t.Fatalf("stderr = %s", stderr.String())
	}
}
