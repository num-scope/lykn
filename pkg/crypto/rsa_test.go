package crypto

import (
	"encoding/json"
	"testing"
)

func TestSignAndVerify(t *testing.T) {
	privPEM, pubPEM, err := GenerateKeyPair(2048)
	if err != nil {
		t.Fatalf("GenerateKeyPair: %v", err)
	}

	license := map[string]any{
		"id":         "test-uuid",
		"version":    1,
		"subject":    map[string]string{"name": "Test User", "email": "test@example.com"},
		"plan":       "enterprise",
		"issued_at":  "2026-01-01T00:00:00Z",
		"not_before": "2026-01-01T00:00:00Z",
		"not_after":  "2027-01-01T00:00:00Z",
		"features":   []string{"feature_a"},
		"metadata":   map[string]any{},
	}

	licenseJSON, err := json.Marshal(license)
	if err != nil {
		t.Fatalf("marshal license: %v", err)
	}

	licContent, err := SignLicense(licenseJSON, privPEM)
	if err != nil {
		t.Fatalf("SignLicense: %v", err)
	}

	var lic LicFile
	if err := json.Unmarshal(licContent, &lic); err != nil {
		t.Fatalf("unmarshal lic: %v", err)
	}
	if lic.Payload == "" {
		t.Error("payload is empty")
	}
	if lic.Signature == "" {
		t.Error("signature is empty")
	}

	payload, err := VerifySignature(licContent, pubPEM)
	if err != nil {
		t.Fatalf("VerifySignature: %v", err)
	}

	var decoded map[string]any
	if err := json.Unmarshal(payload, &decoded); err != nil {
		t.Fatalf("unmarshal payload: %v", err)
	}
	if decoded["id"] != "test-uuid" {
		t.Errorf("id = %v, want test-uuid", decoded["id"])
	}
}

func TestVerifySignature_WrongKey(t *testing.T) {
	privPEM1, _, err := GenerateKeyPair(2048)
	if err != nil {
		t.Fatalf("GenerateKeyPair 1: %v", err)
	}
	_, pubPEM2, err := GenerateKeyPair(2048)
	if err != nil {
		t.Fatalf("GenerateKeyPair 2: %v", err)
	}

	licenseJSON := []byte(`{"id":"test"}`)
	licContent, err := SignLicense(licenseJSON, privPEM1)
	if err != nil {
		t.Fatalf("SignLicense: %v", err)
	}

	_, err = VerifySignature(licContent, pubPEM2)
	if err == nil {
		t.Fatal("expected verification error for wrong key")
	}
}
