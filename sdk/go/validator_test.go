package lykn

import (
	"encoding/json"
	"errors"
	"testing"
	"time"
)

func licensePayloadForTest(t *testing.T, overrides map[string]any) []byte {
	t.Helper()
	now := time.Now().UTC()
	payload := map[string]any{
		"id":         "lic-001",
		"version":    1,
		"subject":    map[string]any{"name": "Demo", "email": "demo@example.com", "organization": "Acme"},
		"plan":       "pro",
		"plan_name":  "Pro Plan",
		"issued_at":  now.Format(time.RFC3339),
		"not_before": now.Add(-5 * time.Minute).Format(time.RFC3339),
		"not_after":  now.Add(30 * 24 * time.Hour).Format(time.RFC3339),
		"features":   []string{"reports", "exports"},
		"limits":     map[string]any{"max_users": 20, "max_devices": 2},
		"metadata":   map[string]any{"env": "test"},
	}
	for key, value := range overrides {
		payload[key] = value
	}
	data, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal payload: %v", err)
	}
	return data
}

func testValidator(t *testing.T, payload []byte, hw Hardware) *Validator {
	t.Helper()
	privateKey, publicPEM := testKeyPair(t)
	lic := signedLicenseForTest(t, privateKey, payload)
	validator, err := NewValidator(ValidatorOptions{
		PublicKeyPEM:   publicPEM,
		LicenseContent: lic,
		HardwareProvider: func() (Hardware, error) {
			return hw, nil
		},
	})
	if err != nil {
		t.Fatalf("NewValidator: %v", err)
	}
	return validator
}

func TestValidatorVerify(t *testing.T) {
	validator := testValidator(t, licensePayloadForTest(t, nil), Hardware{})
	license, err := validator.Verify("reports")
	if err != nil {
		t.Fatalf("Verify: %v", err)
	}
	if license.ID != "lic-001" || license.Plan != "pro" || license.Limits.MaxUsers != 20 {
		t.Fatalf("unexpected license: %+v", license)
	}
	if !validator.HasFeature("exports") {
		t.Fatalf("expected exports feature")
	}
}

func TestValidatorVerifyErrors(t *testing.T) {
	now := time.Now().UTC()
	cases := []struct {
		name      string
		overrides map[string]any
		hardware  Hardware
		features  []string
		wantErr   error
	}{
		{name: "expired", overrides: map[string]any{"not_after": now.Add(-time.Hour).Format(time.RFC3339)}, wantErr: ErrLicenseExpired},
		{name: "not yet valid", overrides: map[string]any{"not_before": now.Add(time.Hour).Format(time.RFC3339)}, wantErr: ErrLicenseNotYetValid},
		{name: "missing feature", features: []string{"billing"}, wantErr: ErrFeatureNotLicensed},
		{name: "hardware mismatch", overrides: map[string]any{"hardware": map[string]any{"hostname": "prod"}}, hardware: Hardware{Hostname: "dev"}, wantErr: ErrHardwareMismatch},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			validator := testValidator(t, licensePayloadForTest(t, tt.overrides), tt.hardware)
			_, err := validator.Verify(tt.features...)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidatorBackgroundInvalidCallback(t *testing.T) {
	validator := testValidator(t, licensePayloadForTest(t, nil), Hardware{})
	validator.checkInterval = 10 * time.Millisecond
	calls := make(chan error, 1)
	validator.OnInvalid(func(err error) { calls <- err })

	count := 0
	validator.verifyOnce = func(requiredFeatures ...string) (*LicenseData, error) {
		count++
		if count == 1 {
			return &LicenseData{ID: "lic-001", Features: []string{}}, nil
		}
		return nil, ErrLicenseExpired
	}

	if _, err := validator.Start(); err != nil {
		t.Fatalf("Start: %v", err)
	}
	defer validator.Stop()
	select {
	case err := <-calls:
		if !errors.Is(err, ErrLicenseExpired) {
			t.Fatalf("callback error = %v", err)
		}
	case <-time.After(time.Second):
		t.Fatalf("timeout waiting for invalid callback")
	}
	if !validator.Invalid() {
		t.Fatalf("validator should be invalid")
	}
}
