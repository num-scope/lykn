package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	lykn "github.com/wu-clan/lykn/sdk/go"
)

type stringList []string

func (s *stringList) String() string { return strings.Join(*s, ",") }

func (s *stringList) Set(value string) error {
	*s = append(*s, value)
	return nil
}

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout io.Writer, stderr io.Writer) int {
	if len(args) == 0 {
		printUsage(stderr)
		return 2
	}
	switch args[0] {
	case "verify":
		return runVerify(args[1:], stdout, stderr)
	case "inspect":
		return runInspect(args[1:], stdout, stderr)
	case "hardware-info":
		return runHardwareInfo(args[1:], stdout, stderr)
	default:
		fmt.Fprintf(stderr, "unknown command: %s\n", args[0])
		printUsage(stderr)
		return 2
	}
}

func printUsage(w io.Writer) {
	fmt.Fprintln(w, "usage: lykn <verify|inspect|hardware-info> [options]")
}

func runVerify(args []string, stdout io.Writer, stderr io.Writer) int {
	fs := flag.NewFlagSet("verify", flag.ContinueOnError)
	fs.SetOutput(stderr)
	publicKey := fs.String("public-key", "", "path to public key PEM")
	licensePath := fs.String("license", "", "path to license file")
	var features stringList
	fs.Var(&features, "feature", "required feature; may be repeated")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if *publicKey == "" || *licensePath == "" {
		fmt.Fprintln(stderr, "verify requires --public-key and --license")
		return 2
	}

	validator, err := lykn.NewValidator(lykn.ValidatorOptions{PublicKeyPath: *publicKey, LicensePath: *licensePath})
	if err != nil {
		fmt.Fprintf(stderr, "error: %v\n", err)
		return 1
	}
	license, err := validator.Verify(features...)
	if err != nil {
		fmt.Fprintf(stderr, "error: %v\n", err)
		return 1
	}
	fmt.Fprintf(stdout, "License %s is valid\n", license.ID)
	return 0
}

func runInspect(args []string, stdout io.Writer, stderr io.Writer) int {
	fs := flag.NewFlagSet("inspect", flag.ContinueOnError)
	fs.SetOutput(stderr)
	licensePath := fs.String("license", "", "path to license file")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if *licensePath == "" {
		fmt.Fprintln(stderr, "inspect requires --license")
		return 2
	}
	content, err := os.ReadFile(*licensePath)
	if err != nil {
		fmt.Fprintf(stderr, "error: %v\n", err)
		return 1
	}
	payload, err := lykn.InspectLicensePayload(content)
	if err != nil {
		fmt.Fprintf(stderr, "error: %v\n", err)
		return 1
	}
	var out any
	if err := json.Unmarshal(payload, &out); err != nil {
		fmt.Fprintf(stderr, "error: %v\n", err)
		return 1
	}
	encoded, _ := json.MarshalIndent(out, "", "  ")
	fmt.Fprintln(stdout, string(encoded))
	return 0
}

func runHardwareInfo(args []string, stdout io.Writer, stderr io.Writer) int {
	fs := flag.NewFlagSet("hardware-info", flag.ContinueOnError)
	fs.SetOutput(stderr)
	format := fs.String("format", "json", "output format: json or table")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	hardware, err := lykn.CollectHardware()
	if err != nil {
		fmt.Fprintf(stderr, "error: %v\n", err)
		return 1
	}
	switch *format {
	case "json":
		encoded, _ := json.MarshalIndent(hardware, "", "  ")
		fmt.Fprintln(stdout, string(encoded))
		return 0
	case "table":
		fmt.Fprintf(stdout, "hostname\t%s\n", hardware.Hostname)
		fmt.Fprintf(stdout, "cpu_id\t%s\n", hardware.CPUID)
		fmt.Fprintf(stdout, "disk_serial\t%s\n", hardware.DiskSerial)
		fmt.Fprintf(stdout, "mac_addresses\t%s\n", strings.Join(hardware.MACAddresses, ", "))
		fmt.Fprintf(stdout, "ip_addresses\t%s\n", strings.Join(hardware.IPAddresses, ", "))
		return 0
	default:
		fmt.Fprintf(stderr, "unsupported format: %s\n", *format)
		return 2
	}
}
