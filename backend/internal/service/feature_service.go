package service

import (
	"context"
	"errors"
	"strings"

	"github.com/wu-clan/lykn/backend/internal/common"
	"github.com/wu-clan/lykn/backend/internal/dao"
	"github.com/wu-clan/lykn/backend/internal/dto"
	"github.com/wu-clan/lykn/backend/internal/model"
)

type FeatureService struct {
	featureDAO *dao.FeatureDAO
}

func NewFeatureService(featureDAO *dao.FeatureDAO) *FeatureService {
	return &FeatureService{featureDAO: featureDAO}
}

func (s *FeatureService) List(ctx context.Context) ([]dto.FeatureResponse, error) {
	features, err := s.featureDAO.List(ctx)
	if err != nil {
		return nil, common.NewInternal(common.CodeInternal, "list features failed")
	}
	items := make([]dto.FeatureResponse, 0, len(features))
	for _, feature := range features {
		items = append(items, toFeatureResponse(&feature))
	}
	return items, nil
}

func (s *FeatureService) Create(ctx context.Context, req dto.CreateFeatureRequest) (*dto.FeatureResponse, error) {
	code, name, description, err := normalizeFeatureInput(req.Code, req.Name, req.Description)
	if err != nil {
		return nil, err
	}
	if err := s.ensureFeatureCodeAvailable(ctx, code, 0); err != nil {
		return nil, err
	}
	feature := &model.Feature{
		Code:        code,
		Name:        name,
		Description: description,
		Enabled:     req.Enabled,
	}
	if err := s.featureDAO.Create(ctx, feature); err != nil {
		return nil, common.NewInternal(common.CodeInternal, "create feature failed")
	}
	resp := toFeatureResponse(feature)
	return &resp, nil
}

func (s *FeatureService) Update(ctx context.Context, id uint, req dto.UpdateFeatureRequest) (*dto.FeatureResponse, error) {
	feature, err := s.getFeatureModel(ctx, id)
	if err != nil {
		return nil, err
	}
	code, name, description, err := normalizeFeatureInput(req.Code, req.Name, req.Description)
	if err != nil {
		return nil, err
	}
	if err := s.ensureFeatureCodeAvailable(ctx, code, id); err != nil {
		return nil, err
	}
	feature.Code = code
	feature.Name = name
	feature.Description = description
	feature.Enabled = req.Enabled
	if err := s.featureDAO.Update(ctx, feature); err != nil {
		return nil, common.NewInternal(common.CodeInternal, "update feature failed")
	}
	resp := toFeatureResponse(feature)
	return &resp, nil
}

func (s *FeatureService) Delete(ctx context.Context, id uint) error {
	feature, err := s.getFeatureModel(ctx, id)
	if err != nil {
		return err
	}
	count, err := s.featureDAO.CountPlansByFeatureID(ctx, id)
	if err != nil {
		return common.NewInternal(common.CodeInternal, "count feature plans failed")
	}
	if count > 0 {
		return common.NewConflict("feature is used by plans and cannot be deleted")
	}
	if err := s.featureDAO.Delete(ctx, feature); err != nil {
		return common.NewInternal(common.CodeInternal, "delete feature failed")
	}
	return nil
}

func (s *FeatureService) getFeatureModel(ctx context.Context, id uint) (*model.Feature, error) {
	feature, err := s.featureDAO.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, common.ErrNotFound) {
			return nil, common.NewNotFound("feature not found")
		}
		return nil, common.NewInternal(common.CodeInternal, "query feature failed")
	}
	return feature, nil
}

func (s *FeatureService) ensureFeatureCodeAvailable(ctx context.Context, code string, currentID uint) error {
	existing, err := s.featureDAO.GetByCode(ctx, code)
	if err == nil {
		if existing.ID != currentID {
			return common.NewConflict("feature code already exists")
		}
		return nil
	}
	if errors.Is(err, common.ErrNotFound) {
		return nil
	}
	return common.NewInternal(common.CodeInternal, "query feature code failed")
}

func normalizeFeatureInput(code string, name string, description string) (string, string, string, error) {
	code = strings.TrimSpace(code)
	name = strings.TrimSpace(name)
	description = strings.TrimSpace(description)
	if code == "" {
		return "", "", "", common.NewInvalidRequest("feature code is required")
	}
	if name == "" {
		return "", "", "", common.NewInvalidRequest("feature name is required")
	}
	return code, name, description, nil
}

func ListFeatures(ctx context.Context) ([]dto.FeatureResponse, error) {
	if featureRuntime == nil {
		return nil, common.NewInternal(common.CodeInternal, "service not initialized")
	}
	return featureRuntime.List(ctx)
}

func CreateFeature(ctx context.Context, req dto.CreateFeatureRequest) (*dto.FeatureResponse, error) {
	if featureRuntime == nil {
		return nil, common.NewInternal(common.CodeInternal, "service not initialized")
	}
	return featureRuntime.Create(ctx, req)
}

func UpdateFeature(ctx context.Context, id uint, req dto.UpdateFeatureRequest) (*dto.FeatureResponse, error) {
	if featureRuntime == nil {
		return nil, common.NewInternal(common.CodeInternal, "service not initialized")
	}
	return featureRuntime.Update(ctx, id, req)
}

func DeleteFeature(ctx context.Context, id uint) error {
	if featureRuntime == nil {
		return common.NewInternal(common.CodeInternal, "service not initialized")
	}
	return featureRuntime.Delete(ctx, id)
}

func toFeatureResponse(feature *model.Feature) dto.FeatureResponse {
	return dto.FeatureResponse{
		ID:          feature.ID,
		Code:        feature.Code,
		Name:        feature.Name,
		Description: feature.Description,
		Enabled:     feature.Enabled,
		CreatedAt:   feature.CreatedAt,
		UpdatedAt:   feature.UpdatedAt,
	}
}
