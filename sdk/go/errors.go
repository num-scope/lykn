package lykn

import (
	"errors"
	"fmt"
)

var (
	ErrLicenseFile        = errors.New("license file error")
	ErrLicenseSignature   = errors.New("license signature verification failed")
	ErrLicenseExpired     = errors.New("license has expired")
	ErrLicenseNotYetValid = errors.New("license is not yet valid")
	ErrHardwareMismatch   = errors.New("hardware fingerprint mismatch")
	ErrFeatureNotLicensed = errors.New("requested feature is not licensed")
)

func wrapError(base error, format string, args ...any) error {
	if len(args) == 0 {
		return fmt.Errorf("%w: %s", base, format)
	}
	return fmt.Errorf("%w: %s", base, fmt.Sprintf(format, args...))
}
