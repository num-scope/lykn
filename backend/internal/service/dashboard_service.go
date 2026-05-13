package service

import (
	"context"
	"time"

	"github.com/wu-clan/lykn/backend/internal/common"
	"github.com/wu-clan/lykn/backend/internal/dao"
	"github.com/wu-clan/lykn/backend/internal/dto"
)

type DashboardService struct {
	projectDAO *dao.ProjectDAO
	licenseDAO *dao.LicenseDAO
}

func NewDashboardService(projectDAO *dao.ProjectDAO, licenseDAO *dao.LicenseDAO) *DashboardService {
	return &DashboardService{projectDAO: projectDAO, licenseDAO: licenseDAO}
}

func (s *DashboardService) GetSummary(ctx context.Context) (*dto.DashboardSummaryResponse, error) {
	if _, ok := common.UserIDFromContext(ctx); !ok {
		return nil, common.NewUnauthorized("missing user context")
	}
	projectCount, err := s.projectDAO.Count(ctx)
	if err != nil {
		return nil, common.NewInternal(common.CodeInternal, "count projects failed")
	}
	licenseCount, err := s.licenseDAO.Count(ctx)
	if err != nil {
		return nil, common.NewInternal(common.CodeInternal, "count licenses failed")
	}
	now := time.Now()
	activeLicenseCount, err := s.licenseDAO.CountActive(ctx, now)
	if err != nil {
		return nil, common.NewInternal(common.CodeInternal, "count active licenses failed")
	}
	expiredLicenseCount, err := s.licenseDAO.CountExpired(ctx, now)
	if err != nil {
		return nil, common.NewInternal(common.CodeInternal, "count expired licenses failed")
	}
	return &dto.DashboardSummaryResponse{
		ProjectCount:        projectCount,
		LicenseCount:        licenseCount,
		ActiveLicenseCount:  activeLicenseCount,
		ExpiredLicenseCount: expiredLicenseCount,
	}, nil
}

func GetDashboardSummary(ctx context.Context) (*dto.DashboardSummaryResponse, error) {
	if dashboardRuntime == nil {
		return nil, common.NewInternal(common.CodeInternal, "service not initialized")
	}
	return dashboardRuntime.GetSummary(ctx)
}
