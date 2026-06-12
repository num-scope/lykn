package service

import (
	"context"
	"testing"

	"github.com/wu-clan/lykn/internal/dao"
	"github.com/wu-clan/lykn/internal/dto"
)

func TestFeatureRejectsDuplicateCode(t *testing.T) {
	db := newTestDB(t)
	featureSvc := NewFeatureService(dao.NewFeatureDAO(db))

	_, err := featureSvc.Create(context.Background(), dto.CreateFeatureRequest{
		Code:    "reports",
		Name:    "Reports",
		Enabled: true,
	})
	if err != nil {
		t.Fatalf("Create feature: %v", err)
	}

	_, err = featureSvc.Create(context.Background(), dto.CreateFeatureRequest{
		Code:    "reports",
		Name:    "Reports Copy",
		Enabled: true,
	})
	if err == nil {
		t.Fatal("expected duplicate feature code error")
	}
}

func TestFeatureDeleteRejectsPlanUsage(t *testing.T) {
	db := newTestDB(t)
	featureDAO := dao.NewFeatureDAO(db)
	planDAO := dao.NewPlanDAO(db)
	featureSvc := NewFeatureService(featureDAO)
	planSvc := NewPlanService(planDAO, featureDAO)

	feature, err := featureSvc.Create(context.Background(), dto.CreateFeatureRequest{
		Code:    "reports",
		Name:    "Reports",
		Enabled: true,
	})
	if err != nil {
		t.Fatalf("Create feature: %v", err)
	}
	_, err = planSvc.Create(context.Background(), dto.CreatePlanRequest{
		Code:       "pro",
		Name:       "Pro",
		FeatureIDs: []uint{feature.ID},
		MaxDevices: 1,
		Enabled:    true,
	})
	if err != nil {
		t.Fatalf("Create plan: %v", err)
	}

	err = featureSvc.Delete(context.Background(), feature.ID)
	if err == nil {
		t.Fatal("expected feature-in-use delete error")
	}
}

func TestPlanCreateWithFeatures(t *testing.T) {
	db := newTestDB(t)
	featureDAO := dao.NewFeatureDAO(db)
	planSvc := NewPlanService(dao.NewPlanDAO(db), featureDAO)
	featureSvc := NewFeatureService(featureDAO)

	reports, err := featureSvc.Create(context.Background(), dto.CreateFeatureRequest{
		Code:    "reports",
		Name:    "Reports",
		Enabled: true,
	})
	if err != nil {
		t.Fatalf("Create reports feature: %v", err)
	}
	export, err := featureSvc.Create(context.Background(), dto.CreateFeatureRequest{
		Code:    "export",
		Name:    "Export",
		Enabled: true,
	})
	if err != nil {
		t.Fatalf("Create export feature: %v", err)
	}

	plan, err := planSvc.Create(context.Background(), dto.CreatePlanRequest{
		Code:       "enterprise",
		Name:       "Enterprise",
		FeatureIDs: []uint{reports.ID, export.ID},
		MaxUsers:   100,
		MaxDevices: 3,
		Enabled:    true,
	})
	if err != nil {
		t.Fatalf("Create plan: %v", err)
	}
	if len(plan.Features) != 2 {
		t.Fatalf("feature count = %d, want 2", len(plan.Features))
	}
	if plan.MaxUsers != 100 || plan.MaxDevices != 3 {
		t.Fatalf("unexpected limits: %+v", plan)
	}
}
