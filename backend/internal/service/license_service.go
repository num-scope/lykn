package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/wu-clan/lykn/backend/internal/common"
	"github.com/wu-clan/lykn/backend/internal/dao"
	"github.com/wu-clan/lykn/backend/internal/dto"
	"github.com/wu-clan/lykn/backend/internal/model"
	backendcrypto "github.com/wu-clan/lykn/backend/pkg/crypto"
	"gorm.io/datatypes"
)

type LicenseService struct {
	projectDAO *dao.ProjectDAO
	licenseDAO *dao.LicenseDAO
	planDAO    *dao.PlanDAO
	secret     string
}

type issueSnapshot struct {
	PlanID   *uint
	PlanCode string
	PlanName string
	Features []string
	Limits   dto.LicenseLimitsResponse
	Hardware map[string]any
	Metadata map[string]any
}

func NewLicenseService(projectDAO *dao.ProjectDAO, licenseDAO *dao.LicenseDAO, planDAO *dao.PlanDAO, secret string) *LicenseService {
	return &LicenseService{projectDAO: projectDAO, licenseDAO: licenseDAO, planDAO: planDAO, secret: secret}
}

func (s *LicenseService) Issue(ctx context.Context, projectID uint, req dto.IssueLicenseRequest) (*dto.LicenseResponse, string, error) {
	if req.Subject.Name == "" {
		return nil, "", common.NewInvalidRequest("subject.name is required")
	}
	if !req.NotAfter.After(req.NotBefore) {
		return nil, "", common.NewInvalidRequest("not_after must be later than not_before")
	}
	snapshot, err := s.resolveIssueSnapshot(ctx, req)
	if err != nil {
		return nil, "", err
	}
	project, err := s.projectDAO.GetByID(ctx, projectID)
	if err != nil {
		if errors.Is(err, common.ErrNotFound) {
			return nil, "", common.NewNotFound("project not found")
		}
		return nil, "", common.NewInternal(common.CodeInternal, "query project failed")
	}
	privPEM, err := common.DecryptPrivateKey(project.PrivateKey, s.secret)
	if err != nil {
		return nil, "", common.NewInternal(common.CodeLicenseIssueFailed, "decrypt project private key failed")
	}
	licenseUUID := uuid.NewString()
	issuedAt := time.Now().UTC().Format("2006-01-02T15:04:05Z")
	payload := map[string]any{
		"id":         licenseUUID,
		"version":    1,
		"subject":    map[string]string{"name": req.Subject.Name, "email": req.Subject.Email, "organization": req.Subject.Organization},
		"plan":       snapshot.PlanCode,
		"plan_name":  snapshot.PlanName,
		"issued_at":  issuedAt,
		"not_before": req.NotBefore.UTC().Format("2006-01-02T15:04:05Z"),
		"not_after":  req.NotAfter.UTC().Format("2006-01-02T15:04:05Z"),
		"hardware":   snapshot.Hardware,
		"features":   snapshot.Features,
		"limits":     snapshot.Limits,
		"metadata":   snapshot.Metadata,
	}
	licenseJSON, err := json.Marshal(payload)
	if err != nil {
		return nil, "", common.NewInternal(common.CodeLicenseIssueFailed, "marshal license payload failed")
	}
	licContent, err := backendcrypto.SignLicense(licenseJSON, privPEM)
	if err != nil {
		return nil, "", common.NewInternal(common.CodeLicenseIssueFailed, "sign license failed")
	}
	featuresJSON, _ := json.Marshal(snapshot.Features)
	metadataJSON, _ := json.Marshal(snapshot.Metadata)
	hardwareJSON, _ := json.Marshal(snapshot.Hardware)
	limitsJSON, _ := json.Marshal(snapshot.Limits)
	license := &model.License{
		UUID:         licenseUUID,
		ProjectID:    project.ID,
		SubjectName:  req.Subject.Name,
		SubjectEmail: req.Subject.Email,
		SubjectOrg:   req.Subject.Organization,
		PlanID:       snapshot.PlanID,
		PlanName:     snapshot.PlanName,
		Plan:         snapshot.PlanCode,
		NotBefore:    req.NotBefore,
		NotAfter:     req.NotAfter,
		Hardware:     datatypes.JSON(hardwareJSON),
		Features:     datatypes.JSON(featuresJSON),
		Limits:       datatypes.JSON(limitsJSON),
		Metadata:     datatypes.JSON(metadataJSON),
		LicContent:   string(licContent),
	}
	if err := s.licenseDAO.Create(ctx, license); err != nil {
		return nil, "", common.NewInternal(common.CodeInternal, "save license failed")
	}
	resp := toLicenseResponse(license)
	return &resp, string(licContent), nil
}

func (s *LicenseService) ListByProjectID(ctx context.Context, projectID uint) ([]dto.LicenseResponse, error) {
	licenses, err := s.licenseDAO.ListByProjectID(ctx, projectID)
	if err != nil {
		return nil, common.NewInternal(common.CodeInternal, "list licenses failed")
	}
	items := make([]dto.LicenseResponse, 0, len(licenses))
	for _, license := range licenses {
		items = append(items, toLicenseResponse(&license))
	}
	return items, nil
}

func (s *LicenseService) Get(ctx context.Context, id uint) (*dto.LicenseResponse, error) {
	license, err := s.getLicenseModel(ctx, id)
	if err != nil {
		return nil, err
	}
	resp := toLicenseResponse(license)
	return &resp, nil
}

func (s *LicenseService) getLicenseModel(ctx context.Context, id uint) (*model.License, error) {
	license, err := s.licenseDAO.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, common.ErrNotFound) {
			return nil, common.NewNotFound("license not found")
		}
		return nil, common.NewInternal(common.CodeInternal, "query license failed")
	}
	return license, nil
}

func (s *LicenseService) Download(ctx context.Context, id uint) (string, string, error) {
	license, err := s.getLicenseModel(ctx, id)
	if err != nil {
		return "", "", err
	}
	return fmt.Sprintf("%s.lic", license.UUID), license.LicContent, nil
}

func (s *LicenseService) resolveIssueSnapshot(ctx context.Context, req dto.IssueLicenseRequest) (*issueSnapshot, error) {
	hardware := buildHardwarePayload(req.Hardware)
	if req.PlanID == 0 {
		metadata := req.Metadata
		if metadata == nil {
			metadata = map[string]any{}
		}
		return &issueSnapshot{
			PlanCode: strings.TrimSpace(req.Plan),
			Features: normalizeFeatureCodes(req.Features),
			Limits:   dto.LicenseLimitsResponse{},
			Hardware: hardware,
			Metadata: metadata,
		}, nil
	}
	if s.planDAO == nil {
		return nil, common.NewInternal(common.CodeInternal, "plan service not initialized")
	}
	plan, err := s.planDAO.GetByID(ctx, req.PlanID)
	if err != nil {
		if errors.Is(err, common.ErrNotFound) {
			return nil, common.NewInvalidRequest("plan not found")
		}
		return nil, common.NewInternal(common.CodeInternal, "query plan failed")
	}
	if !plan.Enabled {
		return nil, common.NewInvalidRequest("plan is disabled")
	}
	featureCodes := make([]string, 0, len(plan.Features))
	for _, feature := range plan.Features {
		if feature.Enabled {
			featureCodes = append(featureCodes, feature.Code)
		}
	}
	planID := plan.ID
	return &issueSnapshot{
		PlanID:   &planID,
		PlanCode: plan.Code,
		PlanName: plan.Name,
		Features: normalizeFeatureCodes(featureCodes),
		Limits: dto.LicenseLimitsResponse{
			MaxUsers:   plan.MaxUsers,
			MaxDevices: plan.MaxDevices,
		},
		Hardware: hardware,
		Metadata: map[string]any{},
	}, nil
}

func ListProjectLicenses(ctx context.Context, projectID uint) ([]dto.LicenseResponse, error) {
	if licenseRuntime == nil {
		return nil, common.NewInternal(common.CodeInternal, "service not initialized")
	}
	return licenseRuntime.ListByProjectID(ctx, projectID)
}

func IssueLicense(ctx context.Context, projectID uint, req dto.IssueLicenseRequest) (*dto.LicenseResponse, string, error) {
	if licenseRuntime == nil {
		return nil, "", common.NewInternal(common.CodeInternal, "service not initialized")
	}
	return licenseRuntime.Issue(ctx, projectID, req)
}

func GetLicense(ctx context.Context, id uint) (*dto.LicenseResponse, error) {
	if licenseRuntime == nil {
		return nil, common.NewInternal(common.CodeInternal, "service not initialized")
	}
	return licenseRuntime.Get(ctx, id)
}

func DownloadLicense(ctx context.Context, id uint) (string, string, error) {
	if licenseRuntime == nil {
		return "", "", common.NewInternal(common.CodeInternal, "service not initialized")
	}
	return licenseRuntime.Download(ctx, id)
}

func toLicenseResponse(license *model.License) dto.LicenseResponse {
	return dto.LicenseResponse{
		ID:           license.ID,
		UUID:         license.UUID,
		ProjectID:    license.ProjectID,
		SubjectName:  license.SubjectName,
		SubjectEmail: license.SubjectEmail,
		SubjectOrg:   license.SubjectOrg,
		PlanID:       license.PlanID,
		PlanName:     license.PlanName,
		Plan:         license.Plan,
		NotBefore:    license.NotBefore,
		NotAfter:     license.NotAfter,
		Features:     parseFeatures(license.Features),
		Limits:       parseLimits(license.Limits),
		Metadata:     parseMetadata(license.Metadata),
		CreatedAt:    license.CreatedAt,
	}
}

func normalizeFeatureCodes(features []string) []string {
	seen := map[string]bool{}
	items := make([]string, 0, len(features))
	for _, feature := range features {
		value := strings.TrimSpace(feature)
		if value == "" || seen[value] {
			continue
		}
		seen[value] = true
		items = append(items, value)
	}
	return items
}

func buildHardwarePayload(req dto.LicenseHardwareRequest) map[string]any {
	payload := map[string]any{}
	if value := strings.TrimSpace(req.Hostname); value != "" {
		payload["hostname"] = value
	}
	if value := strings.TrimSpace(req.CPUID); value != "" {
		payload["cpu_id"] = value
	}
	if value := strings.TrimSpace(req.DiskSerial); value != "" {
		payload["disk_serial"] = value
	}
	if values := normalizeStringList(req.MACAddresses); len(values) > 0 {
		payload["mac_addresses"] = values
	}
	if values := normalizeStringList(req.IPAddresses); len(values) > 0 {
		payload["ip_addresses"] = values
	}
	if len(payload) == 0 {
		return nil
	}
	return payload
}

func normalizeStringList(values []string) []string {
	seen := map[string]bool{}
	items := make([]string, 0, len(values))
	for _, value := range values {
		item := strings.TrimSpace(value)
		if item == "" || seen[item] {
			continue
		}
		seen[item] = true
		items = append(items, item)
	}
	return items
}

func parseFeatures(raw []byte) []string {
	if len(raw) == 0 {
		return []string{}
	}
	var features []string
	if err := json.Unmarshal(raw, &features); err != nil {
		return []string{}
	}
	return features
}

func parseMetadata(raw []byte) map[string]any {
	if len(raw) == 0 {
		return map[string]any{}
	}
	var metadata map[string]any
	if err := json.Unmarshal(raw, &metadata); err != nil {
		return map[string]any{}
	}
	if metadata == nil {
		return map[string]any{}
	}
	return metadata
}

func parseLimits(raw []byte) dto.LicenseLimitsResponse {
	if len(raw) == 0 {
		return dto.LicenseLimitsResponse{}
	}
	var limits dto.LicenseLimitsResponse
	if err := json.Unmarshal(raw, &limits); err != nil {
		return dto.LicenseLimitsResponse{}
	}
	return limits
}
