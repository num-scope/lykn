package gin

import (
	"net/http"

	githubgin "github.com/gin-gonic/gin"
	lykn "github.com/wu-clan/lykn/sdk/go"
)

const licenseContextKey = "lykn_license"

type verifier interface {
	Verify(requiredFeatures ...string) (*lykn.LicenseData, error)
}

func RequireLicense(v verifier) githubgin.HandlerFunc {
	return RequireFeatures(v)
}

func RequireFeatures(v verifier, features ...string) githubgin.HandlerFunc {
	return func(c *githubgin.Context) {
		license, err := v.Verify(features...)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusForbidden, githubgin.H{"error": "license validation failed"})
			return
		}
		c.Set(licenseContextKey, license)
		c.Next()
	}
}

func License(c *githubgin.Context) (*lykn.LicenseData, bool) {
	value, ok := c.Get(licenseContextKey)
	if !ok {
		return nil, false
	}
	license, ok := value.(*lykn.LicenseData)
	return license, ok
}

func MustLicense(c *githubgin.Context) *lykn.LicenseData {
	license, ok := License(c)
	if !ok {
		panic("lykn license missing from gin context")
	}
	return license
}
