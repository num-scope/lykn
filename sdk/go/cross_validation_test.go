package lykn

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestVerifyBackendSignedFixture(t *testing.T) {
	fixtureDir := filepath.Join("..", "..", "tests", "fixtures")
	publicPEM, err := os.ReadFile(filepath.Join(fixtureDir, "public.pem"))
	if err != nil {
		t.Fatalf("read public.pem: %v", err)
	}
	licContent, err := os.ReadFile(filepath.Join(fixtureDir, "license.lic"))
	if err != nil {
		t.Fatalf("read license.lic: %v", err)
	}
	expectedBytes, err := os.ReadFile(filepath.Join(fixtureDir, "license.json"))
	if err != nil {
		t.Fatalf("read license.json: %v", err)
	}
	var expected LicenseData
	if err := json.Unmarshal(expectedBytes, &expected); err != nil {
		t.Fatalf("parse expected license: %v", err)
	}
	expected.applyDefaults()

	validator, err := NewValidator(ValidatorOptions{
		PublicKeyPEM:   publicPEM,
		LicenseContent: licContent,
		HardwareProvider: func() (Hardware, error) {
			return *expected.Hardware, nil
		},
	})
	if err != nil {
		t.Fatalf("NewValidator: %v", err)
	}
	actual, err := validator.Verify("reports", "exports")
	if err != nil {
		t.Fatalf("Verify: %v", err)
	}
	if actual.ID != expected.ID || actual.Plan != expected.Plan || actual.PlanName != expected.PlanName {
		t.Fatalf("license mismatch: actual=%+v expected=%+v", actual, expected)
	}
	if actual.Subject.Name != expected.Subject.Name || actual.Subject.Email != expected.Subject.Email {
		t.Fatalf("subject mismatch: actual=%+v expected=%+v", actual.Subject, expected.Subject)
	}
	if actual.Limits.MaxUsers != expected.Limits.MaxUsers || actual.Limits.MaxDevices != expected.Limits.MaxDevices {
		t.Fatalf("limits mismatch: actual=%+v expected=%+v", actual.Limits, expected.Limits)
	}
	if actual.Hardware == nil || actual.Hardware.Hostname != expected.Hardware.Hostname {
		t.Fatalf("hardware mismatch: actual=%+v expected=%+v", actual.Hardware, expected.Hardware)
	}
}
