package service

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/wu-clan/lykn/backend/internal/model"
	"gorm.io/datatypes"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func newTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	name := strings.NewReplacer("/", "_", " ", "_").Replace(t.Name())
	db, err := gorm.Open(sqlite.Open(fmt.Sprintf("file:%s?mode=memory&cache=shared", name)), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := db.AutoMigrate(&model.User{}, &model.Project{}, &model.Feature{}, &model.Plan{}, &model.PlanFeature{}, &model.License{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	return db
}

func newTestLicense(projectID uint) *model.License {
	return &model.License{
		UUID:        "fixture-license",
		ProjectID:   projectID,
		SubjectName: "Fixture",
		Features:    datatypes.JSON([]byte(`[]`)),
		Limits:      datatypes.JSON([]byte(`{"max_users":0,"max_devices":0}`)),
		Metadata:    datatypes.JSON([]byte(`{}`)),
		LicContent:  `{"payload":"e30=","signature":"e30="}`,
		NotBefore:   time.Now().Add(-time.Hour),
		NotAfter:    time.Now().Add(time.Hour),
	}
}
