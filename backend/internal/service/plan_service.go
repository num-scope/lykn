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

type PlanService struct {
	planDAO    *dao.PlanDAO
	featureDAO *dao.FeatureDAO
}

func NewPlanService(planDAO *dao.PlanDAO, featureDAO *dao.FeatureDAO) *PlanService {
	return &PlanService{planDAO: planDAO, featureDAO: featureDAO}
}

func (s *PlanService) List(ctx context.Context) ([]dto.PlanResponse, error) {
	plans, err := s.planDAO.List(ctx)
	if err != nil {
		return nil, common.NewInternal(common.CodeInternal, "list plans failed")
	}
	items := make([]dto.PlanResponse, 0, len(plans))
	for _, plan := range plans {
		items = append(items, toPlanResponse(&plan))
	}
	return items, nil
}

func (s *PlanService) Create(ctx context.Context, req dto.CreatePlanRequest) (*dto.PlanResponse, error) {
	code, name, description, err := normalizePlanInput(req.Code, req.Name, req.Description, req.MaxUsers, req.MaxDevices)
	if err != nil {
		return nil, err
	}
	if err := s.ensurePlanCodeAvailable(ctx, code, 0); err != nil {
		return nil, err
	}
	features, err := s.resolveEnabledFeatures(ctx, req.FeatureIDs)
	if err != nil {
		return nil, err
	}
	plan := &model.Plan{
		Code:        code,
		Name:        name,
		Description: description,
		MaxUsers:    req.MaxUsers,
		MaxDevices:  req.MaxDevices,
		Enabled:     req.Enabled,
	}
	if err := s.planDAO.Create(ctx, plan, features); err != nil {
		return nil, common.NewInternal(common.CodeInternal, "create plan failed")
	}
	plan.Features = features
	resp := toPlanResponse(plan)
	return &resp, nil
}

func (s *PlanService) Update(ctx context.Context, id uint, req dto.UpdatePlanRequest) (*dto.PlanResponse, error) {
	plan, err := s.getPlanModel(ctx, id)
	if err != nil {
		return nil, err
	}
	code, name, description, err := normalizePlanInput(req.Code, req.Name, req.Description, req.MaxUsers, req.MaxDevices)
	if err != nil {
		return nil, err
	}
	if err := s.ensurePlanCodeAvailable(ctx, code, id); err != nil {
		return nil, err
	}
	features, err := s.resolveEnabledFeatures(ctx, req.FeatureIDs)
	if err != nil {
		return nil, err
	}
	plan.Code = code
	plan.Name = name
	plan.Description = description
	plan.MaxUsers = req.MaxUsers
	plan.MaxDevices = req.MaxDevices
	plan.Enabled = req.Enabled
	if err := s.planDAO.Update(ctx, plan, features); err != nil {
		return nil, common.NewInternal(common.CodeInternal, "update plan failed")
	}
	plan.Features = features
	resp := toPlanResponse(plan)
	return &resp, nil
}

func (s *PlanService) Delete(ctx context.Context, id uint) error {
	plan, err := s.getPlanModel(ctx, id)
	if err != nil {
		return err
	}
	if err := s.planDAO.Delete(ctx, plan); err != nil {
		return common.NewInternal(common.CodeInternal, "delete plan failed")
	}
	return nil
}

func (s *PlanService) getPlanModel(ctx context.Context, id uint) (*model.Plan, error) {
	plan, err := s.planDAO.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, common.ErrNotFound) {
			return nil, common.NewNotFound("plan not found")
		}
		return nil, common.NewInternal(common.CodeInternal, "query plan failed")
	}
	return plan, nil
}

func (s *PlanService) ensurePlanCodeAvailable(ctx context.Context, code string, currentID uint) error {
	existing, err := s.planDAO.GetByCode(ctx, code)
	if err == nil {
		if existing.ID != currentID {
			return common.NewConflict("plan code already exists")
		}
		return nil
	}
	if errors.Is(err, common.ErrNotFound) {
		return nil
	}
	return common.NewInternal(common.CodeInternal, "query plan code failed")
}

func (s *PlanService) resolveEnabledFeatures(ctx context.Context, ids []uint) ([]model.Feature, error) {
	seen := map[uint]bool{}
	features := make([]model.Feature, 0, len(ids))
	for _, id := range ids {
		if id == 0 || seen[id] {
			continue
		}
		seen[id] = true
		feature, err := s.featureDAO.GetByID(ctx, id)
		if err != nil {
			if errors.Is(err, common.ErrNotFound) {
				return nil, common.NewInvalidRequest("feature not found")
			}
			return nil, common.NewInternal(common.CodeInternal, "query feature failed")
		}
		if !feature.Enabled {
			return nil, common.NewInvalidRequest("feature is disabled")
		}
		features = append(features, *feature)
	}
	return features, nil
}

func normalizePlanInput(code string, name string, description string, maxUsers int, maxDevices int) (string, string, string, error) {
	code = strings.TrimSpace(code)
	name = strings.TrimSpace(name)
	description = strings.TrimSpace(description)
	if code == "" {
		return "", "", "", common.NewInvalidRequest("plan code is required")
	}
	if name == "" {
		return "", "", "", common.NewInvalidRequest("plan name is required")
	}
	if maxUsers < 0 {
		return "", "", "", common.NewInvalidRequest("max_users cannot be negative")
	}
	if maxDevices < 1 {
		return "", "", "", common.NewInvalidRequest("max_devices must be at least 1")
	}
	return code, name, description, nil
}

func ListPlans(ctx context.Context) ([]dto.PlanResponse, error) {
	if planRuntime == nil {
		return nil, common.NewInternal(common.CodeInternal, "service not initialized")
	}
	return planRuntime.List(ctx)
}

func CreatePlan(ctx context.Context, req dto.CreatePlanRequest) (*dto.PlanResponse, error) {
	if planRuntime == nil {
		return nil, common.NewInternal(common.CodeInternal, "service not initialized")
	}
	return planRuntime.Create(ctx, req)
}

func UpdatePlan(ctx context.Context, id uint, req dto.UpdatePlanRequest) (*dto.PlanResponse, error) {
	if planRuntime == nil {
		return nil, common.NewInternal(common.CodeInternal, "service not initialized")
	}
	return planRuntime.Update(ctx, id, req)
}

func DeletePlan(ctx context.Context, id uint) error {
	if planRuntime == nil {
		return common.NewInternal(common.CodeInternal, "service not initialized")
	}
	return planRuntime.Delete(ctx, id)
}

func toPlanResponse(plan *model.Plan) dto.PlanResponse {
	features := make([]dto.FeatureResponse, 0, len(plan.Features))
	for _, feature := range plan.Features {
		features = append(features, toFeatureResponse(&feature))
	}
	return dto.PlanResponse{
		ID:          plan.ID,
		Code:        plan.Code,
		Name:        plan.Name,
		Description: plan.Description,
		Features:    features,
		MaxUsers:    plan.MaxUsers,
		MaxDevices:  plan.MaxDevices,
		Enabled:     plan.Enabled,
		CreatedAt:   plan.CreatedAt,
		UpdatedAt:   plan.UpdatedAt,
	}
}
