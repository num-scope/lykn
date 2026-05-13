package dto

import "time"

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type CreateProjectRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	KeyBits     int    `json:"key_bits"`
}

type UpdateProjectRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type CreateFeatureRequest struct {
	Code        string `json:"code" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Enabled     bool   `json:"enabled"`
}

type UpdateFeatureRequest struct {
	Code        string `json:"code" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Enabled     bool   `json:"enabled"`
}

type CreatePlanRequest struct {
	Code        string `json:"code" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	FeatureIDs  []uint `json:"feature_ids"`
	MaxUsers    int    `json:"max_users"`
	MaxDevices  int    `json:"max_devices"`
	Enabled     bool   `json:"enabled"`
}

type UpdatePlanRequest struct {
	Code        string `json:"code" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	FeatureIDs  []uint `json:"feature_ids"`
	MaxUsers    int    `json:"max_users"`
	MaxDevices  int    `json:"max_devices"`
	Enabled     bool   `json:"enabled"`
}

type SubjectRequest struct {
	Name         string `json:"name" binding:"required"`
	Email        string `json:"email"`
	Organization string `json:"organization"`
}

type IssueLicenseRequest struct {
	Subject   SubjectRequest         `json:"subject" binding:"required"`
	PlanID    uint                   `json:"plan_id"`
	Plan      string                 `json:"plan"`
	NotBefore time.Time              `json:"not_before" binding:"required"`
	NotAfter  time.Time              `json:"not_after" binding:"required"`
	Hardware  LicenseHardwareRequest `json:"hardware"`
	Features  []string               `json:"features"`
	Metadata  map[string]any         `json:"metadata"`
}

type LicenseHardwareRequest struct {
	Hostname     string   `json:"hostname"`
	CPUID        string   `json:"cpu_id"`
	DiskSerial   string   `json:"disk_serial"`
	MACAddresses []string `json:"mac_addresses"`
	IPAddresses  []string `json:"ip_addresses"`
}
