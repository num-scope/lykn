package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/wu-clan/lykn/config"
	"github.com/wu-clan/lykn/internal/model"
	"github.com/wu-clan/lykn/internal/service"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func newRouterTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	name := strings.NewReplacer("/", "_", " ", "_").Replace(t.Name())
	db, err := gorm.Open(sqlite.Open(fmt.Sprintf("file:%s?mode=memory&cache=shared", name)), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := db.AutoMigrate(&model.User{}, &model.Project{}, &model.Feature{}, &model.Plan{}, &model.PlanFeature{}, &model.License{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	return db
}

func newRouterTestConfig() *config.Config {
	return &config.Config{
		Server: config.ServerConfig{Port: 8080, AllowOrigin: "*"},
		Auth: config.AuthConfig{
			SecretKey: "jwt-secret",
			Expire:    "24h",
		},
		Encryption: config.EncryptionConfig{SecretKey: "encrypt-secret"},
	}
}

func TestLoginAndProtectedRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := newRouterTestDB(t)
	cfg := newRouterTestConfig()
	if err := service.Init(db, cfg); err != nil {
		t.Fatalf("Init: %v", err)
	}
	r := RegisterRoutes(cfg)

	loginBody := []byte(`{"username":"admin","password":"admin123"}`)
	loginReq := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(loginBody))
	loginReq.Header.Set("Content-Type", "application/json")
	loginW := httptest.NewRecorder()
	r.ServeHTTP(loginW, loginReq)
	if loginW.Code != http.StatusOK {
		t.Fatalf("login status = %d, want 200", loginW.Code)
	}
	var loginResp struct {
		Data struct {
			AccessToken string `json:"access_token"`
		} `json:"data"`
	}
	if err := json.NewDecoder(loginW.Body).Decode(&loginResp); err != nil {
		t.Fatalf("decode login response: %v", err)
	}
	if loginResp.Data.AccessToken == "" {
		t.Fatal("access token is empty")
	}

	projectReq := httptest.NewRequest(http.MethodGet, "/api/v1/projects", nil)
	projectW := httptest.NewRecorder()
	r.ServeHTTP(projectW, projectReq)
	if projectW.Code != http.StatusUnauthorized {
		t.Fatalf("project status = %d, want 401", projectW.Code)
	}

	meReq := httptest.NewRequest(http.MethodGet, "/api/v1/auth/me", nil)
	meReq.Header.Set("Authorization", "Bearer "+loginResp.Data.AccessToken)
	meW := httptest.NewRecorder()
	r.ServeHTTP(meW, meReq)
	if meW.Code != http.StatusOK {
		t.Fatalf("me status = %d, want 200", meW.Code)
	}

	summaryReq := httptest.NewRequest(http.MethodGet, "/api/v1/dashboard/summary", nil)
	summaryReq.Header.Set("Authorization", "Bearer "+loginResp.Data.AccessToken)
	summaryW := httptest.NewRecorder()
	r.ServeHTTP(summaryW, summaryReq)
	if summaryW.Code != http.StatusOK {
		t.Fatalf("summary status = %d, want 200", summaryW.Code)
	}
}

func TestHealthAndCORS(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cfg := newRouterTestConfig()
	cfg.Server.AllowOrigin = "http://localhost:5173"
	r := RegisterRoutes(cfg)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("health status = %d, want 200", w.Code)
	}
	if w.Body.String() != "OK" {
		t.Fatalf("health body = %q, want OK", w.Body.String())
	}
	if got := w.Header().Get("Access-Control-Allow-Origin"); got != "http://localhost:5173" {
		t.Fatalf("allow origin = %q, want http://localhost:5173", got)
	}
}
