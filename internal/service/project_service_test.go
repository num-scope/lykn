package service

import (
	"context"
	"strings"
	"testing"

	"github.com/wu-clan/lykn/internal/dao"
	"github.com/wu-clan/lykn/internal/dto"
)

func TestCreateProjectEncryptsPrivateKey(t *testing.T) {
	db := newTestDB(t)
	projectDAO := dao.NewProjectDAO(db)
	licenseDAO := dao.NewLicenseDAO(db)
	svc := NewProjectService(projectDAO, licenseDAO, "encrypt-secret")

	resp, err := svc.Create(context.Background(), dto.CreateProjectRequest{Name: "demo", Description: "desc", KeyBits: 2048})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if resp.PublicKey == "" {
		t.Fatal("public key is empty")
	}
	project, err := projectDAO.GetByID(context.Background(), resp.ID)
	if err != nil {
		t.Fatalf("GetByID: %v", err)
	}
	if strings.Contains(project.PrivateKey, "PRIVATE KEY") {
		t.Fatal("private key should be encrypted before storing")
	}
}

func TestDeleteProjectRejectsWhenLicensesExist(t *testing.T) {
	db := newTestDB(t)
	projectDAO := dao.NewProjectDAO(db)
	licenseDAO := dao.NewLicenseDAO(db)
	svc := NewProjectService(projectDAO, licenseDAO, "encrypt-secret")

	projectResp, err := svc.Create(context.Background(), dto.CreateProjectRequest{Name: "demo", KeyBits: 2048})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if err := licenseDAO.Create(context.Background(), newTestLicense(projectResp.ID)); err != nil {
		t.Fatalf("Create license: %v", err)
	}
	if err := svc.Delete(context.Background(), projectResp.ID); err == nil {
		t.Fatal("expected delete conflict when licenses exist")
	}
}
