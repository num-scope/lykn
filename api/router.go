package api

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	v1 "github.com/wu-clan/lykn/api/v1"
	"github.com/wu-clan/lykn/config"
	"github.com/wu-clan/lykn/internal/middleware"
	"github.com/wu-clan/lykn/pkg/response"
)

func RegisterRoutes(cfg *config.Config) *gin.Engine {
	middleware.SetAuthSecret(cfg.Auth.SecretKey)

	router := gin.Default()
	router.Use(middleware.CORS(cfg.Server.AllowOrigin))

	apiV1 := router.Group("/api/v1")

	router.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	// Auth routes.
	apiV1.POST("/auth/login", v1.Login)
	apiV1.GET("/auth/me", middleware.RequireLogin(), v1.GetCurrentUser)

	// Protected routes.
	secure := apiV1.Group("")
	secure.Use(middleware.RequireLogin())

	secure.GET("/dashboard/summary", v1.GetDashboardSummary)

	secure.GET("/features", v1.ListFeatures)
	secure.POST("/features", v1.CreateFeature)
	secure.PUT("/features/:id", v1.UpdateFeature)
	secure.DELETE("/features/:id", v1.DeleteFeature)

	secure.GET("/plans", v1.ListPlans)
	secure.POST("/plans", v1.CreatePlan)
	secure.PUT("/plans/:id", v1.UpdatePlan)
	secure.DELETE("/plans/:id", v1.DeletePlan)

	secure.GET("/projects", v1.ListProjects)
	secure.POST("/projects", v1.CreateProject)
	secure.GET("/projects/:id", v1.GetProject)
	secure.PUT("/projects/:id", v1.UpdateProject)
	secure.DELETE("/projects/:id", v1.DeleteProject)
	secure.GET("/projects/:id/public-key", v1.DownloadProjectPublicKey)
	secure.GET("/projects/:id/licenses", v1.ListProjectLicenses)
	secure.POST("/projects/:id/licenses", v1.IssueLicense)
	secure.GET("/licenses/:id", v1.GetLicense)
	secure.GET("/licenses/:id/download", v1.DownloadLicense)

	registerWebRoutes(router, webDistDir())

	return router
}

func webDistDir() string {
	if dir := os.Getenv("LYKN_WEB_DIST"); dir != "" {
		return dir
	}
	return filepath.Join("frontend", "dist")
}

func registerWebRoutes(router *gin.Engine, distDir string) {
	indexPath := filepath.Join(distDir, "index.html")
	if _, err := os.Stat(indexPath); err != nil {
		return
	}

	router.StaticFS("/assets", http.Dir(filepath.Join(distDir, "assets")))
	router.StaticFile("/favicon.ico", filepath.Join(distDir, "favicon.ico"))
	router.StaticFile("/favicon.svg", filepath.Join(distDir, "favicon.svg"))

	router.NoRoute(func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/api/") {
			response.Error(c, http.StatusNotFound, "not found")
			return
		}
		if c.Request.Method != http.MethodGet && c.Request.Method != http.MethodHead {
			response.Error(c, http.StatusNotFound, "not found")
			return
		}
		c.File(indexPath)
	})
}
