package service

import (
	"context"
	"errors"

	"github.com/wu-clan/lykn/backend/internal/common"
	"github.com/wu-clan/lykn/backend/internal/dao"
	"github.com/wu-clan/lykn/backend/internal/dto"
	"github.com/wu-clan/lykn/backend/internal/model"
	backendcrypto "github.com/wu-clan/lykn/backend/pkg/crypto"
)

type ProjectService struct {
	projectDAO *dao.ProjectDAO
	licenseDAO *dao.LicenseDAO
	secret     string
}

func NewProjectService(projectDAO *dao.ProjectDAO, licenseDAO *dao.LicenseDAO, secret string) *ProjectService {
	return &ProjectService{projectDAO: projectDAO, licenseDAO: licenseDAO, secret: secret}
}

func (s *ProjectService) List(ctx context.Context) ([]dto.ProjectResponse, error) {
	projects, err := s.projectDAO.List(ctx)
	if err != nil {
		return nil, common.NewInternal(common.CodeInternal, "list projects failed")
	}
	items := make([]dto.ProjectResponse, 0, len(projects))
	for _, project := range projects {
		items = append(items, toProjectResponse(&project, false))
	}
	return items, nil
}

func (s *ProjectService) Get(ctx context.Context, id uint) (*dto.ProjectResponse, error) {
	project, err := s.getProjectModel(ctx, id)
	if err != nil {
		return nil, err
	}
	resp := toProjectResponse(project, true)
	return &resp, nil
}

func (s *ProjectService) getProjectModel(ctx context.Context, id uint) (*model.Project, error) {
	project, err := s.projectDAO.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, common.ErrNotFound) {
			return nil, common.NewNotFound("project not found")
		}
		return nil, common.NewInternal(common.CodeInternal, "query project failed")
	}
	return project, nil
}

func (s *ProjectService) Create(ctx context.Context, req dto.CreateProjectRequest) (*dto.ProjectResponse, error) {
	bits := req.KeyBits
	if bits == 0 {
		bits = 2048
	}
	privPEM, pubPEM, err := backendcrypto.GenerateKeyPair(bits)
	if err != nil {
		return nil, common.NewInternal(common.CodeProjectKeygenFailed, "generate project key pair failed")
	}
	ciphertext, err := common.EncryptPrivateKey(privPEM, s.secret)
	if err != nil {
		return nil, common.NewInternal(common.CodeProjectKeygenFailed, "encrypt private key failed")
	}
	project := &model.Project{
		Name:        req.Name,
		Description: req.Description,
		PrivateKey:  ciphertext,
		PublicKey:   string(pubPEM),
		KeyBits:     bits,
	}
	if err := s.projectDAO.Create(ctx, project); err != nil {
		return nil, common.NewInternal(common.CodeInternal, "create project failed")
	}
	resp := toProjectResponse(project, true)
	return &resp, nil
}

func (s *ProjectService) Update(ctx context.Context, id uint, req dto.UpdateProjectRequest) (*dto.ProjectResponse, error) {
	project, err := s.getProjectModel(ctx, id)
	if err != nil {
		return nil, err
	}
	project.Name = req.Name
	project.Description = req.Description
	if err := s.projectDAO.Update(ctx, project); err != nil {
		return nil, common.NewInternal(common.CodeInternal, "update project failed")
	}
	resp := toProjectResponse(project, true)
	return &resp, nil
}

func (s *ProjectService) Delete(ctx context.Context, id uint) error {
	project, err := s.getProjectModel(ctx, id)
	if err != nil {
		return err
	}
	count, err := s.licenseDAO.CountByProjectID(ctx, id)
	if err != nil {
		return common.NewInternal(common.CodeInternal, "count project licenses failed")
	}
	if count > 0 {
		return common.NewConflict("project has licenses and cannot be deleted")
	}
	if err := s.projectDAO.Delete(ctx, project); err != nil {
		return common.NewInternal(common.CodeInternal, "delete project failed")
	}
	return nil
}

func ListProjects(ctx context.Context) ([]dto.ProjectResponse, error) {
	if projectRuntime == nil {
		return nil, common.NewInternal(common.CodeInternal, "service not initialized")
	}
	return projectRuntime.List(ctx)
}

func GetProject(ctx context.Context, id uint) (*dto.ProjectResponse, error) {
	if projectRuntime == nil {
		return nil, common.NewInternal(common.CodeInternal, "service not initialized")
	}
	return projectRuntime.Get(ctx, id)
}

func CreateProject(ctx context.Context, req dto.CreateProjectRequest) (*dto.ProjectResponse, error) {
	if projectRuntime == nil {
		return nil, common.NewInternal(common.CodeInternal, "service not initialized")
	}
	return projectRuntime.Create(ctx, req)
}

func UpdateProject(ctx context.Context, id uint, req dto.UpdateProjectRequest) (*dto.ProjectResponse, error) {
	if projectRuntime == nil {
		return nil, common.NewInternal(common.CodeInternal, "service not initialized")
	}
	return projectRuntime.Update(ctx, id, req)
}

func DeleteProject(ctx context.Context, id uint) error {
	if projectRuntime == nil {
		return common.NewInternal(common.CodeInternal, "service not initialized")
	}
	return projectRuntime.Delete(ctx, id)
}

func toProjectResponse(project *model.Project, includePublicKey bool) dto.ProjectResponse {
	resp := dto.ProjectResponse{
		ID:          project.ID,
		Name:        project.Name,
		Description: project.Description,
		KeyBits:     project.KeyBits,
		CreatedAt:   project.CreatedAt,
		UpdatedAt:   project.UpdatedAt,
	}
	if includePublicKey {
		resp.PublicKey = project.PublicKey
	}
	return resp
}
