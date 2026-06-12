package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/wu-clan/lykn/pkg/crypto"
)

type subjectFixture struct {
	Name         string `json:"name"`
	Email        string `json:"email"`
	Organization string `json:"organization"`
}

type hardwareFixture struct {
	Hostname     string   `json:"hostname"`
	CPUID        string   `json:"cpu_id"`
	DiskSerial   string   `json:"disk_serial"`
	MACAddresses []string `json:"mac_addresses"`
}

type limitsFixture struct {
	MaxUsers   int `json:"max_users"`
	MaxDevices int `json:"max_devices"`
}

type licenseFixture struct {
	ID        string          `json:"id"`
	Version   int             `json:"version"`
	Subject   subjectFixture  `json:"subject"`
	Plan      string          `json:"plan"`
	PlanName  string          `json:"plan_name"`
	IssuedAt  string          `json:"issued_at"`
	NotBefore string          `json:"not_before"`
	NotAfter  string          `json:"not_after"`
	Hardware  hardwareFixture `json:"hardware"`
	Features  []string        `json:"features"`
	Limits    limitsFixture   `json:"limits"`
	Metadata  map[string]any  `json:"metadata"`
}

func main() {
	// Demo fixture files are written to repo-level tests/fixtures by default.
	outDir := filepath.Join("tests", "fixtures")
	if len(os.Args) > 1 {
		outDir = os.Args[1]
	}

	if err := os.MkdirAll(outDir, 0o755); err != nil {
		log.Fatalf("create output dir: %v", err)
	}

	privPEM, pubPEM, err := crypto.GenerateKeyPair(2048)
	if err != nil {
		log.Fatalf("generate key pair: %v", err)
	}

	license := licenseFixture{
		ID:      "fixture-uuid-001",
		Version: 1,
		Subject: subjectFixture{
			Name:         "Fixture Test Corp",
			Email:        "admin@fixture.test",
			Organization: "Fixture Test Corp",
		},
		Plan:      "enterprise",
		PlanName:  "Enterprise Plan",
		IssuedAt:  "2026-04-01T00:00:00Z",
		NotBefore: "2026-04-01T00:00:00Z",
		NotAfter:  "2027-04-01T00:00:00Z",
		Hardware: hardwareFixture{
			Hostname:     "fixture-host",
			CPUID:        "CPU-FIXTURE-001",
			DiskSerial:   "DISK-FIXTURE-001",
			MACAddresses: []string{"AA:BB:CC:DD:EE:FF"},
		},
		Features: []string{"reports", "exports"},
		Limits: limitsFixture{
			MaxUsers:   20,
			MaxDevices: 1,
		},
		Metadata: map[string]any{
			"customer_id": "cust_fixture_001",
			"environment": "demo",
			"issued_by":   "cmd/demo",
		},
	}

	licenseJSON, err := json.MarshalIndent(license, "", "  ")
	if err != nil {
		log.Fatalf("marshal license: %v", err)
	}
	licenseJSON = append(licenseJSON, '\n')

	licContent, err := crypto.SignLicense(licenseJSON, privPEM)
	if err != nil {
		log.Fatalf("sign license: %v", err)
	}

	files := map[string][]byte{
		"public.pem":   pubPEM,
		"license.lic":  licContent,
		"license.json": licenseJSON,
	}

	for name, data := range files {
		path := filepath.Join(outDir, name)
		if err := os.WriteFile(path, data, 0o644); err != nil {
			log.Fatalf("write %s: %v", name, err)
		}
	}

	fmt.Printf("Demo fixtures generated in %s\n", outDir)
	fmt.Println("  - public.pem")
	fmt.Println("  - license.lic")
	fmt.Println("  - license.json")
}
