package lykn

import (
	"net"
	"os"
	"os/exec"
	"runtime"
	"slices"
	"sort"
	"strings"
)

func NormalizeMAC(value string) string {
	text := strings.ToUpper(strings.ReplaceAll(strings.TrimSpace(value), "-", ":"))
	parts := strings.Split(text, ":")
	if len(parts) != 6 {
		return ""
	}
	for i, part := range parts {
		if len(part) == 1 {
			part = "0" + part
		}
		if len(part) != 2 {
			return ""
		}
		for _, ch := range part {
			if !((ch >= '0' && ch <= '9') || (ch >= 'A' && ch <= 'F')) {
				return ""
			}
		}
		parts[i] = part
	}
	return strings.Join(parts, ":")
}

func NormalizeHardware(hw Hardware) Hardware {
	macSet := map[string]struct{}{}
	for _, mac := range hw.MACAddresses {
		if normalized := NormalizeMAC(mac); normalized != "" && normalized != "00:00:00:00:00:00" {
			macSet[normalized] = struct{}{}
		}
	}

	ipSet := map[string]struct{}{}
	for _, ip := range hw.IPAddresses {
		if trimmed := strings.TrimSpace(strings.Split(ip, "%")[0]); trimmed != "" {
			ipSet[trimmed] = struct{}{}
		}
	}

	hw.MACAddresses = sortedKeys(macSet)
	hw.IPAddresses = sortedKeys(ipSet)
	hw.Hostname = strings.TrimSpace(hw.Hostname)
	hw.CPUID = strings.TrimSpace(hw.CPUID)
	hw.DiskSerial = strings.TrimSpace(hw.DiskSerial)
	return hw
}

func sortedKeys(values map[string]struct{}) []string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func CollectHardware() (Hardware, error) {
	hostname, _ := os.Hostname()
	hw := Hardware{
		Hostname:     hostname,
		MACAddresses: collectMACAddresses(),
		IPAddresses:  collectIPAddresses(),
		CPUID:        firstCommandOutput(cpuCommands()),
		DiskSerial:   firstCommandOutput(diskCommands()),
	}
	return NormalizeHardware(hw), nil
}

func collectMACAddresses() []string {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil
	}
	values := make([]string, 0, len(interfaces))
	for _, iface := range interfaces {
		if mac := NormalizeMAC(iface.HardwareAddr.String()); mac != "" {
			values = append(values, mac)
		}
	}
	return values
}

func collectIPAddresses() []string {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil
	}
	var values []string
	for _, iface := range interfaces {
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			ipText := addr.String()
			if host, _, err := net.ParseCIDR(ipText); err == nil {
				ipText = host.String()
			}
			ip := net.ParseIP(strings.Split(ipText, "%")[0])
			if ip == nil || ip.IsLoopback() {
				continue
			}
			values = append(values, ip.String())
		}
	}
	return values
}

func cpuCommands() [][]string {
	switch runtime.GOOS {
	case "linux":
		return [][]string{{"sh", "-c", "awk -F': ' '/Serial/{print $2; exit}' /proc/cpuinfo"}}
	case "darwin":
		return [][]string{{"sysctl", "-n", "machdep.cpu.brand_string"}}
	case "windows":
		return [][]string{{"wmic", "cpu", "get", "ProcessorId"}}
	default:
		return nil
	}
}

func diskCommands() [][]string {
	switch runtime.GOOS {
	case "linux":
		return [][]string{{"sh", "-c", "lsblk -ndo SERIAL | head -n 1"}}
	case "darwin":
		return [][]string{{"sh", "-c", "system_profiler SPNVMeDataType | awk -F': ' '/Serial Number/{print $2; exit}'"}}
	case "windows":
		return [][]string{{"wmic", "diskdrive", "get", "SerialNumber"}}
	default:
		return nil
	}
}

func firstCommandOutput(commands [][]string) string {
	for _, command := range commands {
		if len(command) == 0 {
			continue
		}
		out, err := exec.Command(command[0], command[1:]...).Output()
		if err != nil {
			continue
		}
		for _, line := range strings.Split(string(out), "\n") {
			line = strings.TrimSpace(line)
			lower := strings.ToLower(line)
			if line != "" && !strings.Contains(lower, "serial") && !strings.Contains(lower, "processorid") {
				return line
			}
		}
	}
	return ""
}

func ValidateHardware(required *Hardware, current Hardware) error {
	if required == nil {
		return nil
	}

	req := NormalizeHardware(*required)
	cur := NormalizeHardware(current)
	if req.Hostname != "" && req.Hostname != cur.Hostname {
		return wrapError(ErrHardwareMismatch, "hostname mismatch")
	}
	if req.CPUID != "" && req.CPUID != cur.CPUID {
		return wrapError(ErrHardwareMismatch, "CPU id mismatch")
	}
	if req.DiskSerial != "" && req.DiskSerial != cur.DiskSerial {
		return wrapError(ErrHardwareMismatch, "disk serial mismatch")
	}
	if !isSubset(req.MACAddresses, cur.MACAddresses) {
		return wrapError(ErrHardwareMismatch, "MAC addresses mismatch")
	}
	if !isSubset(req.IPAddresses, cur.IPAddresses) {
		return wrapError(ErrHardwareMismatch, "IP addresses mismatch")
	}
	return nil
}

func isSubset(required []string, current []string) bool {
	for _, item := range required {
		if !slices.Contains(current, item) {
			return false
		}
	}
	return true
}
