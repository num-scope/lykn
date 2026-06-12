package dto

import "time"

type UserResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
}

type LoginResponse struct {
	AccessToken string       `json:"access_token"`
	TokenType   string       `json:"token_type"`
	ExpiresAt   time.Time    `json:"expires_at"`
	User        UserResponse `json:"user"`
}

type DashboardSummaryResponse struct {
	ProjectCount        int64 `json:"project_count"`
	LicenseCount        int64 `json:"license_count"`
	ActiveLicenseCount  int64 `json:"active_license_count"`
	ExpiredLicenseCount int64 `json:"expired_license_count"`
}

type ProjectResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	PublicKey   string    `json:"public_key,omitempty"`
	KeyBits     int       `json:"key_bits"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type FeatureResponse struct {
	ID          uint      `json:"id"`
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Enabled     bool      `json:"enabled"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type PlanResponse struct {
	ID          uint              `json:"id"`
	Code        string            `json:"code"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Features    []FeatureResponse `json:"features"`
	MaxUsers    int               `json:"max_users"`
	MaxDevices  int               `json:"max_devices"`
	Enabled     bool              `json:"enabled"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

type LicenseLimitsResponse struct {
	MaxUsers   int `json:"max_users"`
	MaxDevices int `json:"max_devices"`
}

type LicenseResponse struct {
	ID           uint                  `json:"id"`
	UUID         string                `json:"uuid"`
	ProjectID    uint                  `json:"project_id"`
	SubjectName  string                `json:"subject_name"`
	SubjectEmail string                `json:"subject_email"`
	SubjectOrg   string                `json:"subject_org"`
	PlanID       *uint                 `json:"plan_id,omitempty"`
	PlanName     string                `json:"plan_name"`
	Plan         string                `json:"plan"`
	NotBefore    time.Time             `json:"not_before"`
	NotAfter     time.Time             `json:"not_after"`
	Features     []string              `json:"features"`
	Limits       LicenseLimitsResponse `json:"limits"`
	Metadata     map[string]any        `json:"metadata"`
	CreatedAt    time.Time             `json:"created_at"`
}
