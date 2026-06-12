package common

import "testing"

func TestEncryptDecryptPrivateKey(t *testing.T) {
	secret := "0123456789abcdef0123456789abcdef"
	plain := []byte("-----BEGIN PRIVATE KEY-----\nabc\n-----END PRIVATE KEY-----")

	ciphertext, err := EncryptPrivateKey(plain, secret)
	if err != nil {
		t.Fatalf("EncryptPrivateKey: %v", err)
	}
	if ciphertext == string(plain) {
		t.Fatal("ciphertext should not equal plaintext")
	}

	decoded, err := DecryptPrivateKey(ciphertext, secret)
	if err != nil {
		t.Fatalf("DecryptPrivateKey: %v", err)
	}
	if string(decoded) != string(plain) {
		t.Fatalf("decoded = %q, want %q", decoded, plain)
	}
}
