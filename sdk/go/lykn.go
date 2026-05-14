package lykn

import "time"

type Subject struct {
	Name         string `json:"name"`
	Email        string `json:"email"`
	Organization string `json:"organization"`
}

type Hardware struct {
	MACAddresses []string `json:"mac_addresses"`
	IPAddresses  []string `json:"ip_addresses"`
	Hostname     string   `json:"hostname"`
	CPUID        string   `json:"cpu_id"`
	DiskSerial   string   `json:"disk_serial"`
}

type LicenseLimits struct {
	MaxUsers   int `json:"max_users"`
	MaxDevices int `json:"max_devices"`
}

type LicenseData struct {
	ID        string         `json:"id"`
	Version   int            `json:"version"`
	Subject   Subject        `json:"subject"`
	Plan      string         `json:"plan"`
	PlanName  string         `json:"plan_name"`
	IssuedAt  time.Time      `json:"issued_at"`
	NotBefore time.Time      `json:"not_before"`
	NotAfter  time.Time      `json:"not_after"`
	Hardware  *Hardware      `json:"hardware"`
	Features  []string       `json:"features"`
	Limits    LicenseLimits  `json:"limits"`
	Metadata  map[string]any `json:"metadata"`
}

func (l *LicenseData) applyDefaults() {
	if l.Version == 0 {
		l.Version = 1
	}
	if l.Features == nil {
		l.Features = []string{}
	}
	if l.Metadata == nil {
		l.Metadata = map[string]any{}
	}
	if l.Hardware != nil {
		*l.Hardware = NormalizeHardware(*l.Hardware)
	}
}
