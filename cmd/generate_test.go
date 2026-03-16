package cmd

import (
	"strings"
	"testing"
	"text/template"

	"github.com/Maya-Mohamed/emt-migrate/pkg/snapshot"
)

func renderTemplate(td templateData) (string, error) {
	funcMap := template.FuncMap{
		"sanitize": sanitize,
	}
	tmpl, err := template.New("script").Funcs(funcMap).Parse(scriptTemplate)
	if err != nil {
		return "", err
	}
	var buf strings.Builder
	err = tmpl.Execute(&buf, td)
	return buf.String(), err
}

func TestGenerateRHELScript(t *testing.T) {
	td := templateData{
		Hostname:     "emt-edge-node-01",
		Target:       "rhel",
		PkgManager:   "dnf",
		PackageList:  "        docker-ce \\\n        openssl \\\n        curl",
		PackageCount: 3,
		SkippedCount: 1,
		MappedCount:  1,
		HasPackages:  true,
		DockerImages: []snapshot.DockerImage{
			{Repository: "openvino/model_server", Tag: "latest", ID: "sha256:abc", Size: "1.2GB"},
		},
		Services: []snapshot.SystemdService{
			{Name: "sshd.service", State: "enabled"},
		},
	}

	output, err := renderTemplate(td)
	if err != nil {
		t.Fatalf("template render failed: %v", err)
	}

	checks := []string{
		"dnf install -y",
		"docker-ce",
		"openssl",
		"curl",
		"docker pull \"openvino/model_server:latest\"",
		"systemctl enable sshd.service",
		"Target: rhel",
		"emt-edge-node-01",
	}
	for _, check := range checks {
		if !strings.Contains(output, check) {
			t.Errorf("output missing expected string: %q", check)
		}
	}
}

func TestGenerateSUSEScript(t *testing.T) {
	td := templateData{
		Hostname:     "emt-node",
		Target:       "suse",
		PkgManager:   "zypper",
		PackageList:  "        kernel-default \\\n        docker",
		PackageCount: 2,
		HasPackages:  true,
	}

	output, err := renderTemplate(td)
	if err != nil {
		t.Fatalf("template render failed: %v", err)
	}

	if !strings.Contains(output, "zypper install -y") {
		t.Error("SUSE script should use zypper")
	}
	if strings.Contains(output, "dnf install") {
		t.Error("SUSE script should not contain dnf")
	}
}

func TestGenerateEmptyPackageList(t *testing.T) {
	td := templateData{
		Hostname:    "emt-node",
		Target:      "rhel",
		PkgManager:  "dnf",
		HasPackages: false,
	}

	output, err := renderTemplate(td)
	if err != nil {
		t.Fatalf("template render failed: %v", err)
	}

	if strings.Contains(output, "dnf install -y") {
		t.Error("empty package list should not produce dnf install command")
	}
	if !strings.Contains(output, "No packages to install") {
		t.Error("empty package list should say 'No packages to install'")
	}
}

func TestSanitize(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"openvino/model_server", "openvino_model_server"},
		{"image:latest", "image_latest"},
		{"edge-ai-suite/pipeline", "edge_ai_suite_pipeline"},
		{"sshd.service", "sshd_service"},
		{"simple", "simple"},
	}
	for _, tt := range tests {
		got := sanitize(tt.input)
		if got != tt.expected {
			t.Errorf("sanitize(%q) = %q, want %q", tt.input, got, tt.expected)
		}
	}
}
