package service

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"testing"
	"time"

	"github.com/wu-clan/lykn/backend/internal/dao"
	"github.com/wu-clan/lykn/backend/internal/dto"
)

func TestIssueLicenseRejectsInvalidTimeRange(t *testing.T) {
	db := newTestDB(t)
	projectDAO := dao.NewProjectDAO(db)
	licenseDAO := dao.NewLicenseDAO(db)
	planDAO := dao.NewPlanDAO(db)
	projectSvc := NewProjectService(projectDAO, licenseDAO, "encrypt-secret")
	licenseSvc := NewLicenseService(projectDAO, licenseDAO, planDAO, "encrypt-secret")
	projectResp, err := projectSvc.Create(context.Background(), dto.CreateProjectRequest{Name: "demo", KeyBits: 2048})
	if err != nil {
		t.Fatalf("Create project: %v", err)
	}
	_, _, err = licenseSvc.Issue(context.Background(), projectResp.ID, dto.IssueLicenseRequest{
		Subject:   dto.SubjectRequest{Name: "Acme"},
		NotBefore: time.Now().Add(time.Hour),
		NotAfter:  time.Now(),
	})
	if err == nil {
		t.Fatal("expected invalid request error")
	}
}

func TestIssueLicenseUsesPlanSnapshotAndStructuredHardware(t *testing.T) {
	db := newTestDB(t)
	projectDAO := dao.NewProjectDAO(db)
	licenseDAO := dao.NewLicenseDAO(db)
	featureDAO := dao.NewFeatureDAO(db)
	planDAO := dao.NewPlanDAO(db)
	projectSvc := NewProjectService(projectDAO, licenseDAO, "encrypt-secret")
	featureSvc := NewFeatureService(featureDAO)
	planSvc := NewPlanService(planDAO, featureDAO)
	licenseSvc := NewLicenseService(projectDAO, licenseDAO, planDAO, "encrypt-secret")

	projectResp, err := projectSvc.Create(context.Background(), dto.CreateProjectRequest{Name: "demo", KeyBits: 2048})
	if err != nil {
		t.Fatalf("Create project: %v", err)
	}
	feature, err := featureSvc.Create(context.Background(), dto.CreateFeatureRequest{
		Code:    "reports",
		Name:    "Reports",
		Enabled: true,
	})
	if err != nil {
		t.Fatalf("Create feature: %v", err)
	}
	plan, err := planSvc.Create(context.Background(), dto.CreatePlanRequest{
		Code:       "pro",
		Name:       "Pro Plan",
		FeatureIDs: []uint{feature.ID},
		MaxUsers:   20,
		MaxDevices: 2,
		Enabled:    true,
	})
	if err != nil {
		t.Fatalf("Create plan: %v", err)
	}

	resp, _, err := licenseSvc.Issue(context.Background(), projectResp.ID, dto.IssueLicenseRequest{
		Subject:   dto.SubjectRequest{Name: "Acme", Email: "ops@acme.test", Organization: "Acme"},
		PlanID:    plan.ID,
		NotBefore: time.Date(2026, 4, 23, 0, 0, 0, 0, time.UTC),
		NotAfter:  time.Date(2027, 4, 23, 0, 0, 0, 0, time.UTC),
		Hardware: dto.LicenseHardwareRequest{
			CPUID:        "CPU-1",
			DiskSerial:   "DISK-1",
			MACAddresses: []string{"AA:BB:CC:DD:EE:FF", "AA:BB:CC:DD:EE:FF"},
		},
	})
	if err != nil {
		t.Fatalf("Issue: %v", err)
	}
	if resp.Plan != "pro" || resp.PlanName != "Pro Plan" {
		t.Fatalf("unexpected plan snapshot: %+v", resp)
	}
	if len(resp.Features) != 1 || resp.Features[0] != "reports" {
		t.Fatalf("unexpected feature snapshot: %+v", resp.Features)
	}
	if resp.Limits.MaxUsers != 20 || resp.Limits.MaxDevices != 2 {
		t.Fatalf("unexpected limits: %+v", resp.Limits)
	}

	stored, err := licenseDAO.GetByID(context.Background(), resp.ID)
	if err != nil {
		t.Fatalf("GetByID: %v", err)
	}
	payload := decodeLicensePayload(t, stored.LicContent)
	limits := payload["limits"].(map[string]any)
	if limits["max_users"].(float64) != 20 || limits["max_devices"].(float64) != 2 {
		t.Fatalf("unexpected payload limits: %+v", limits)
	}
	hardware := payload["hardware"].(map[string]any)
	if hardware["cpu_id"] != "CPU-1" || hardware["disk_serial"] != "DISK-1" {
		t.Fatalf("unexpected hardware payload: %+v", hardware)
	}
	macs := hardware["mac_addresses"].([]any)
	if len(macs) != 1 {
		t.Fatalf("mac count = %d, want 1", len(macs))
	}
}

func decodeLicensePayload(t *testing.T, content string) map[string]any {
	t.Helper()
	var lic map[string]string
	if err := json.Unmarshal([]byte(content), &lic); err != nil {
		t.Fatalf("unmarshal lic: %v", err)
	}
	raw, err := base64.StdEncoding.DecodeString(lic["payload"])
	if err != nil {
		t.Fatalf("decode payload: %v", err)
	}
	var payload map[string]any
	if err := json.Unmarshal(raw, &payload); err != nil {
		t.Fatalf("unmarshal payload: %v", err)
	}
	return payload
}
