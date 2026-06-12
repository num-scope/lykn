package crypto

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
)

type LicFile struct {
	Payload   string `json:"payload"`
	Signature string `json:"signature"`
}

func GenerateKeyPair(bits int) (privateKeyPEM, publicKeyPEM []byte, err error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, nil, fmt.Errorf("generate rsa key: %w", err)
	}

	privBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		return nil, nil, fmt.Errorf("marshal private key: %w", err)
	}
	privateKeyPEM = pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: privBytes,
	})

	pubBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return nil, nil, fmt.Errorf("marshal public key: %w", err)
	}
	publicKeyPEM = pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubBytes,
	})

	return privateKeyPEM, publicKeyPEM, nil
}

func SignLicense(licenseJSON []byte, privateKeyPEM []byte) ([]byte, error) {
	block, _ := pem.Decode(privateKeyPEM)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}

	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("parse private key: %w", err)
	}

	rsaKey, ok := key.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("key is not RSA private key")
	}

	hash := sha256.Sum256(licenseJSON)
	signature, err := rsa.SignPKCS1v15(rand.Reader, rsaKey, crypto.SHA256, hash[:])
	if err != nil {
		return nil, fmt.Errorf("sign license: %w", err)
	}

	lic := LicFile{
		Payload:   base64.StdEncoding.EncodeToString(licenseJSON),
		Signature: base64.StdEncoding.EncodeToString(signature),
	}

	return json.MarshalIndent(lic, "", "    ")
}

func VerifySignature(licContent []byte, publicKeyPEM []byte) ([]byte, error) {
	var lic LicFile
	if err := json.Unmarshal(licContent, &lic); err != nil {
		return nil, fmt.Errorf("parse lic file: %w", err)
	}

	payload, err := base64.StdEncoding.DecodeString(lic.Payload)
	if err != nil {
		return nil, fmt.Errorf("decode payload: %w", err)
	}

	signature, err := base64.StdEncoding.DecodeString(lic.Signature)
	if err != nil {
		return nil, fmt.Errorf("decode signature: %w", err)
	}

	block, _ := pem.Decode(publicKeyPEM)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}

	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("parse public key: %w", err)
	}

	rsaPubKey, ok := pubKey.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("key is not RSA public key")
	}

	hash := sha256.Sum256(payload)
	if err := rsa.VerifyPKCS1v15(rsaPubKey, crypto.SHA256, hash[:], signature); err != nil {
		return nil, fmt.Errorf("signature verification failed: %w", err)
	}

	return payload, nil
}
