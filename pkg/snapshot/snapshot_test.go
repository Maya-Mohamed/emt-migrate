package snapshot

import (
	"encoding/json"
	"testing"
	"time"
)

func TestSnapshotRoundTrip(t *testing.T) {
	original := Snapshot{
		Timestamp: time.Now().Truncate(time.Second),
		Hostname:  "test-host",
		Packages: []RPMPackage{
			{Name: "kernel", Version: "6.6.51", Release: "1.emt3", Arch: "x86_64"},
		},
		DockerImages: []DockerImage{
			{Repository: "nginx", Tag: "latest", ID: "sha256:abc", Size: "100MB"},
		},
		Services: []SystemdService{
			{Name: "sshd.service", State: "enabled"},
		},
	}

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	var decoded Snapshot
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	if decoded.Hostname != original.Hostname {
		t.Errorf("hostname mismatch: got %s, want %s", decoded.Hostname, original.Hostname)
	}
	if len(decoded.Packages) != 1 {
		t.Errorf("expected 1 package, got %d", len(decoded.Packages))
	}
	if len(decoded.DockerImages) != 1 {
		t.Errorf("expected 1 docker image, got %d", len(decoded.DockerImages))
	}
	if len(decoded.Services) != 1 {
		t.Errorf("expected 1 service, got %d", len(decoded.Services))
	}
}

func TestSnapshotEmptyFields(t *testing.T) {
	original := Snapshot{
		Hostname: "empty-host",
	}

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	var decoded Snapshot
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	if decoded.Packages != nil {
		t.Errorf("expected nil packages, got %v", decoded.Packages)
	}
}
