package gin

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	githubgin "github.com/gin-gonic/gin"
	lykn "github.com/wu-clan/lykn/sdk/go"
)

type stubVerifier struct {
	license *lykn.LicenseData
	err     error
	seen    []string
}

func (s *stubVerifier) Verify(requiredFeatures ...string) (*lykn.LicenseData, error) {
	s.seen = append([]string{}, requiredFeatures...)
	if s.err != nil {
		return nil, s.err
	}
	return s.license, nil
}

func testLicense() *lykn.LicenseData {
	now := time.Now().UTC()
	return &lykn.LicenseData{
		ID:        "lic-001",
		Subject:   lykn.Subject{Name: "Demo"},
		IssuedAt:  now,
		NotBefore: now.Add(-time.Minute),
		NotAfter:  now.Add(time.Hour),
		Features:  []string{"reports"},
	}
}

func TestRequireLicenseSuccess(t *testing.T) {
	githubgin.SetMode(githubgin.TestMode)
	verifier := &stubVerifier{license: testLicense()}
	router := githubgin.New()
	router.GET("/protected", RequireLicense(verifier), func(c *githubgin.Context) {
		license := MustLicense(c)
		c.JSON(http.StatusOK, githubgin.H{"license": license.ID})
	})

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/protected", nil))
	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d body = %s", recorder.Code, recorder.Body.String())
	}
}

func TestRequireFeaturesFailure(t *testing.T) {
	githubgin.SetMode(githubgin.TestMode)
	verifier := &stubVerifier{err: lykn.ErrFeatureNotLicensed}
	router := githubgin.New()
	router.GET("/reports", RequireFeatures(verifier, "reports"), func(c *githubgin.Context) {
		c.Status(http.StatusOK)
	})

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/reports", nil))
	if recorder.Code != http.StatusForbidden {
		t.Fatalf("status = %d body = %s", recorder.Code, recorder.Body.String())
	}
	if len(verifier.seen) != 1 || verifier.seen[0] != "reports" {
		t.Fatalf("features = %+v", verifier.seen)
	}
}

func TestLicenseMissing(t *testing.T) {
	githubgin.SetMode(githubgin.TestMode)
	ctx, _ := githubgin.CreateTestContext(httptest.NewRecorder())
	if license, ok := License(ctx); ok || license != nil {
		t.Fatalf("expected no license")
	}
	defer func() {
		if recovered := recover(); recovered == nil {
			t.Fatalf("expected panic")
		}
	}()
	_ = MustLicense(ctx)
}
