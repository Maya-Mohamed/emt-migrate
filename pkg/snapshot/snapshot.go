// pkg/snapshot/snapshot.go
package snapshot

import "time"

type Snapshot struct {
	Timestamp    time.Time     `json:"timestamp"`
	Hostname     string        `json:"hostname"`
	Packages     []RPMPackage  `json:"packages"`
	DockerImages []DockerImage `json:"docker_images"`
}

type RPMPackage struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Release     string `json:"release"`
	Arch        string `json:"arch"`
	InstallTime string `json:"install_time,omitempty"`
	Vendor      string `json:"vendor,omitempty"`
	SourceRPM   string `json:"source_rpm,omitempty"`
}

type DockerImage struct {
	Repository string `json:"repository"`
	Tag        string `json:"tag"`
	ID         string `json:"id"`
	Size       string `json:"size"`
}
