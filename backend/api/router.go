package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	v1 "github.com/wu-clan/lykn/backend/api/v1"
	"github.com/wu-clan/lykn/backend/config"
	"github.com/wu-clan/lykn/backend/internal/middleware"
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

	return router
}
