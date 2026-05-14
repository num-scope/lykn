package lykn

import (
	"bytes"
	"crypto/rsa"
	"encoding/json"
	"os"
	"slices"
	"strings"
	"sync"
	"time"
)

type ValidatorOptions struct {
	PublicKeyPath    string
	PublicKeyPEM     []byte
	LicensePath      string
	LicenseContent   []byte
	CheckInterval    time.Duration
	HardwareProvider func() (Hardware, error)
}

type Validator struct {
	publicKeyPath    string
	publicKeyPEM     []byte
	licensePath      string
	licenseContent   []byte
	checkInterval    time.Duration
	hardwareProvider func() (Hardware, error)
	verifyOnce       func(requiredFeatures ...string) (*LicenseData, error)

	mu        sync.RWMutex
	license   *LicenseData
	lastError error
	invalid   bool
	callbacks []func(error)
	stopCh    chan struct{}
	doneCh    chan struct{}
	running   bool
}

func NewValidator(options ValidatorOptions) (*Validator, error) {
	if (options.PublicKeyPath == "") == (len(options.PublicKeyPEM) == 0) {
		return nil, wrapError(ErrLicenseFile, "exactly one public key source is required")
	}
	if (options.LicensePath == "") == (len(options.LicenseContent) == 0) {
		return nil, wrapError(ErrLicenseFile, "exactly one license source is required")
	}

	provider := options.HardwareProvider
	if provider == nil {
		provider = CollectHardware
	}
	v := &Validator{
		publicKeyPath:    options.PublicKeyPath,
		publicKeyPEM:     bytes.Clone(options.PublicKeyPEM),
		licensePath:      options.LicensePath,
		licenseContent:   bytes.Clone(options.LicenseContent),
		checkInterval:    options.CheckInterval,
		hardwareProvider: provider,
	}
	v.verifyOnce = v.defaultVerifyOnce
	return v, nil
}

func (v *Validator) publicKey() (*rsa.PublicKey, error) {
	if v.publicKeyPath != "" {
		return LoadPublicKey([]byte(v.publicKeyPath))
	}
	return LoadPublicKey(v.publicKeyPEM)
}

func (v *Validator) readLicenseContent() ([]byte, error) {
	if len(v.licenseContent) > 0 {
		return bytes.Clone(v.licenseContent), nil
	}
	data, err := os.ReadFile(v.licensePath)
	if err != nil {
		return nil, wrapError(ErrLicenseFile, "read license file: %v", err)
	}
	return data, nil
}

func (v *Validator) defaultVerifyOnce(requiredFeatures ...string) (*LicenseData, error) {
	publicKey, err := v.publicKey()
	if err != nil {
		return nil, err
	}
	content, err := v.readLicenseContent()
	if err != nil {
		return nil, err
	}
	payload, err := VerifyLicensePayload(content, publicKey)
	if err != nil {
		return nil, err
	}

	var data LicenseData
	if err := json.Unmarshal(payload, &data); err != nil {
		return nil, wrapError(ErrLicenseFile, "parse license payload: %v", err)
	}
	data.applyDefaults()
	if err := validateRequiredFields(data); err != nil {
		return nil, err
	}
	if err := validateTime(data); err != nil {
		return nil, err
	}
	current, err := v.hardwareProvider()
	if err != nil {
		return nil, wrapError(ErrHardwareMismatch, "collect hardware: %v", err)
	}
	if err := ValidateHardware(data.Hardware, current); err != nil {
		return nil, err
	}
	if err := validateFeatures(data, requiredFeatures); err != nil {
		return nil, err
	}
	return &data, nil
}

func validateRequiredFields(data LicenseData) error {
	if data.ID == "" {
		return wrapError(ErrLicenseFile, "missing license id")
	}
	if data.Subject.Name == "" {
		return wrapError(ErrLicenseFile, "missing subject name")
	}
	if data.IssuedAt.IsZero() || data.NotBefore.IsZero() || data.NotAfter.IsZero() {
		return wrapError(ErrLicenseFile, "missing license time window")
	}
	return nil
}

func validateTime(data LicenseData) error {
	now := time.Now().UTC()
	if data.NotBefore.UTC().After(now) {
		return ErrLicenseNotYetValid
	}
	if data.NotAfter.UTC().Before(now) {
		return ErrLicenseExpired
	}
	return nil
}

func validateFeatures(data LicenseData, required []string) error {
	if len(required) == 0 {
		return nil
	}
	var missing []string
	for _, feature := range required {
		if !slices.Contains(data.Features, feature) {
			missing = append(missing, feature)
		}
	}
	if len(missing) > 0 {
		return wrapError(ErrFeatureNotLicensed, "missing licensed features: %s", strings.Join(missing, ", "))
	}
	return nil
}

func (v *Validator) Verify(requiredFeatures ...string) (*LicenseData, error) {
	data, err := v.verifyOnce(requiredFeatures...)
	v.mu.Lock()
	defer v.mu.Unlock()
	if err != nil {
		v.lastError = err
		return nil, err
	}
	v.license = data
	v.lastError = nil
	v.invalid = false
	return data, nil
}

func (v *Validator) HasFeature(feature string) bool {
	v.mu.RLock()
	loaded := v.license != nil
	v.mu.RUnlock()
	if !loaded {
		if _, err := v.Verify(); err != nil {
			return false
		}
	}

	v.mu.RLock()
	defer v.mu.RUnlock()
	if v.license == nil {
		return false
	}
	return slices.Contains(v.license.Features, feature)
}

func (v *Validator) Start() (*LicenseData, error) {
	data, err := v.Verify()
	if err != nil || v.checkInterval <= 0 {
		return data, err
	}

	v.mu.Lock()
	if v.running {
		v.mu.Unlock()
		return data, nil
	}
	v.stopCh = make(chan struct{})
	v.doneCh = make(chan struct{})
	v.running = true
	v.mu.Unlock()

	go v.runLoop()
	return data, nil
}

func (v *Validator) runLoop() {
	defer func() {
		v.mu.Lock()
		v.running = false
		close(v.doneCh)
		v.mu.Unlock()
	}()

	ticker := time.NewTicker(v.checkInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if _, err := v.Verify(); err != nil {
				v.markInvalid(err)
				return
			}
		case <-v.stopCh:
			return
		}
	}
}

func (v *Validator) Stop() {
	v.mu.RLock()
	if !v.running {
		v.mu.RUnlock()
		return
	}
	stopCh := v.stopCh
	doneCh := v.doneCh
	v.mu.RUnlock()

	select {
	case <-stopCh:
	default:
		close(stopCh)
	}
	<-doneCh
}

func (v *Validator) OnInvalid(callback func(error)) {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.callbacks = append(v.callbacks, callback)
}

func (v *Validator) markInvalid(err error) {
	v.mu.Lock()
	v.invalid = true
	v.lastError = err
	callbacks := append([]func(error){}, v.callbacks...)
	v.mu.Unlock()

	for _, callback := range callbacks {
		callback(err)
	}
}

func (v *Validator) License() *LicenseData {
	v.mu.RLock()
	defer v.mu.RUnlock()
	return v.license
}

func (v *Validator) LastError() error {
	v.mu.RLock()
	defer v.mu.RUnlock()
	return v.lastError
}

func (v *Validator) Invalid() bool {
	v.mu.RLock()
	defer v.mu.RUnlock()
	return v.invalid
}
