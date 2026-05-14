package lykn

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"testing"
)

func testKeyPair(t *testing.T) (*rsa.PrivateKey, []byte) {
	t.Helper()
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	pubBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		t.Fatalf("marshal public key: %v", err)
	}
	pubPEM := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: pubBytes})
	return privateKey, pubPEM
}

func signedLicenseForTest(t *testing.T, privateKey *rsa.PrivateKey, payload []byte) []byte {
	t.Helper()
	hash := sha256.Sum256(payload)
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hash[:])
	if err != nil {
		t.Fatalf("sign payload: %v", err)
	}
	lic, err := json.Marshal(map[string]string{
		"payload":   base64.StdEncoding.EncodeToString(payload),
		"signature": base64.StdEncoding.EncodeToString(signature),
	})
	if err != nil {
		t.Fatalf("marshal license: %v", err)
	}
	return lic
}

func TestVerifyLicensePayload(t *testing.T) {
	privateKey, publicPEM := testKeyPair(t)
	payload := []byte(`{"id":"lic-001","version":1}`)
	lic := signedLicenseForTest(t, privateKey, payload)

	publicKey, err := LoadPublicKey(publicPEM)
	if err != nil {
		t.Fatalf("LoadPublicKey: %v", err)
	}
	got, err := VerifyLicensePayload(lic, publicKey)
	if err != nil {
		t.Fatalf("VerifyLicensePayload: %v", err)
	}
	if string(got) != string(payload) {
		t.Fatalf("payload mismatch: got %s want %s", got, payload)
	}
}

func TestVerifyLicensePayloadErrors(t *testing.T) {
	privateKey, publicPEM := testKeyPair(t)
	publicKey, err := LoadPublicKey(publicPEM)
	if err != nil {
		t.Fatalf("LoadPublicKey: %v", err)
	}
	otherKey, _ := testKeyPair(t)

	cases := []struct {
		name    string
		lic     []byte
		wantErr error
	}{
		{name: "invalid json", lic: []byte(`{`), wantErr: ErrLicenseFile},
		{name: "missing fields", lic: []byte(`{"payload":"abc"}`), wantErr: ErrLicenseFile},
		{name: "bad payload base64", lic: []byte(`{"payload":"***","signature":"abc"}`), wantErr: ErrLicenseFile},
		{name: "bad signature base64", lic: []byte(`{"payload":"YWJj","signature":"***"}`), wantErr: ErrLicenseFile},
		{name: "invalid payload json", lic: signedLicenseForTest(t, privateKey, []byte(`{`)), wantErr: ErrLicenseFile},
		{name: "bad signature", lic: signedLicenseForTest(t, otherKey, []byte(`{"id":"lic-001"}`)), wantErr: ErrLicenseSignature},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			_, err := VerifyLicensePayload(tt.lic, publicKey)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestInspectLicensePayload(t *testing.T) {
	privateKey, _ := testKeyPair(t)
	payload := []byte(`{"id":"lic-001"}`)
	lic := signedLicenseForTest(t, privateKey, payload)

	got, err := InspectLicensePayload(lic)
	if err != nil {
		t.Fatalf("InspectLicensePayload: %v", err)
	}
	if string(got) != string(payload) {
		t.Fatalf("payload mismatch: got %s want %s", got, payload)
	}
}
