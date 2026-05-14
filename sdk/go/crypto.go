package lykn

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"os"
	"strings"
)

type licenseFile struct {
	Payload   string `json:"payload"`
	Signature string `json:"signature"`
}

func LoadPublicKey(source []byte) (*rsa.PublicKey, error) {
	pemData := source
	trimmed := strings.TrimSpace(string(source))
	if trimmed != "" && !strings.HasPrefix(trimmed, "-----") {
		data, err := os.ReadFile(trimmed)
		if err == nil {
			pemData = data
		}
	}

	block, _ := pem.Decode(pemData)
	if block == nil {
		return nil, wrapError(ErrLicenseFile, "failed to decode public key PEM")
	}

	key, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err == nil {
		publicKey, ok := key.(*rsa.PublicKey)
		if !ok {
			return nil, wrapError(ErrLicenseFile, "public key is not RSA")
		}
		return publicKey, nil
	}

	pkcs1Key, pkcs1Err := x509.ParsePKCS1PublicKey(block.Bytes)
	if pkcs1Err == nil {
		return pkcs1Key, nil
	}
	return nil, wrapError(ErrLicenseFile, "parse public key: %v", err)
}

func parseLicenseFile(content []byte) (licenseFile, []byte, []byte, error) {
	var lic licenseFile
	if err := json.Unmarshal(content, &lic); err != nil {
		return lic, nil, nil, wrapError(ErrLicenseFile, "invalid .lic JSON format: %v", err)
	}
	if lic.Payload == "" || lic.Signature == "" {
		return lic, nil, nil, wrapError(ErrLicenseFile, "missing payload or signature field")
	}

	payload, err := base64.StdEncoding.DecodeString(lic.Payload)
	if err != nil {
		return lic, nil, nil, wrapError(ErrLicenseFile, "decode payload: %v", err)
	}
	signature, err := base64.StdEncoding.DecodeString(lic.Signature)
	if err != nil {
		return lic, nil, nil, wrapError(ErrLicenseFile, "decode signature: %v", err)
	}
	return lic, payload, signature, nil
}

func VerifyLicensePayload(content []byte, publicKey *rsa.PublicKey) ([]byte, error) {
	_, payload, signature, err := parseLicenseFile(content)
	if err != nil {
		return nil, err
	}
	if !json.Valid(payload) {
		return nil, wrapError(ErrLicenseFile, "invalid license payload JSON")
	}

	hash := sha256.Sum256(payload)
	if err := rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hash[:], signature); err != nil {
		return nil, wrapError(ErrLicenseSignature, "%v", err)
	}
	return payload, nil
}

func InspectLicensePayload(content []byte) ([]byte, error) {
	_, payload, _, err := parseLicenseFile(content)
	if err != nil {
		return nil, err
	}
	if !json.Valid(payload) {
		return nil, wrapError(ErrLicenseFile, "invalid license payload JSON")
	}
	return payload, nil
}
