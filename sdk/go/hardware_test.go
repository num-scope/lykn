package lykn

import (
	"errors"
	"testing"
)

func TestNormalizeHardware(t *testing.T) {
	hw := NormalizeHardware(Hardware{
		MACAddresses: []string{"aa-bb-cc-dd-ee-ff", "AA:BB:CC:DD:EE:FF", "bad"},
		IPAddresses:  []string{" 10.0.0.5 ", "", "10.0.0.5"},
		Hostname:     " host ",
		CPUID:        " cpu ",
		DiskSerial:   " disk ",
	})
	if got, want := len(hw.MACAddresses), 1; got != want {
		t.Fatalf("MAC count = %d, want %d", got, want)
	}
	if hw.MACAddresses[0] != "AA:BB:CC:DD:EE:FF" {
		t.Fatalf("MAC = %q", hw.MACAddresses[0])
	}
	if got, want := len(hw.IPAddresses), 1; got != want {
		t.Fatalf("IP count = %d, want %d", got, want)
	}
	if hw.Hostname != "host" || hw.CPUID != "cpu" || hw.DiskSerial != "disk" {
		t.Fatalf("trimmed fields mismatch: %+v", hw)
	}
}

func TestValidateHardwareMatch(t *testing.T) {
	current := Hardware{Hostname: "prod", CPUID: "CPU-1", DiskSerial: "DISK-1", MACAddresses: []string{"AA:BB:CC:DD:EE:FF"}}
	required := &Hardware{Hostname: "prod", MACAddresses: []string{"aa-bb-cc-dd-ee-ff"}}
	if err := ValidateHardware(required, current); err != nil {
		t.Fatalf("ValidateHardware: %v", err)
	}
}

func TestValidateHardwareMismatch(t *testing.T) {
	current := Hardware{Hostname: "dev", MACAddresses: []string{"11:22:33:44:55:66"}}
	required := &Hardware{Hostname: "prod", MACAddresses: []string{"AA:BB:CC:DD:EE:FF"}}
	if err := ValidateHardware(required, current); !errors.Is(err, ErrHardwareMismatch) {
		t.Fatalf("error = %v, want ErrHardwareMismatch", err)
	}
}

func TestValidateHardwareNilRequired(t *testing.T) {
	if err := ValidateHardware(nil, Hardware{}); err != nil {
		t.Fatalf("ValidateHardware nil required: %v", err)
	}
}
