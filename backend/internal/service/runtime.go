package service

import (
	"context"

	"github.com/wu-clan/lykn/backend/config"
	"github.com/wu-clan/lykn/backend/internal/dao"
	"gorm.io/gorm"
)

var (
	userDAO          *dao.UserDAO
	projectDAO       *dao.ProjectDAO
	licenseDAO       *dao.LicenseDAO
	featureDAO       *dao.FeatureDAO
	planDAO          *dao.PlanDAO
	authRuntime      *AuthService
	projectRuntime   *ProjectService
	licenseRuntime   *LicenseService
	dashboardRuntime *DashboardService
	featureRuntime   *FeatureService
	planRuntime      *PlanService
)

func Init(db *gorm.DB, cfg *config.Config) error {
	userDAO = dao.NewUserDAO(db)
	projectDAO = dao.NewProjectDAO(db)
	licenseDAO = dao.NewLicenseDAO(db)
	featureDAO = dao.NewFeatureDAO(db)
	planDAO = dao.NewPlanDAO(db)
	authTTL, err := cfg.Auth.ExpireDuration()
	if err != nil {
		return err
	}
	authRuntime = NewAuthService(userDAO, cfg.Auth.SecretKey, authTTL)
	projectRuntime = NewProjectService(projectDAO, licenseDAO, cfg.Encryption.SecretKey)
	licenseRuntime = NewLicenseService(projectDAO, licenseDAO, planDAO, cfg.Encryption.SecretKey)
	dashboardRuntime = NewDashboardService(projectDAO, licenseDAO)
	featureRuntime = NewFeatureService(featureDAO)
	planRuntime = NewPlanService(planDAO, featureDAO)
	return authRuntime.EnsureDefaultUser(context.Background())
}
