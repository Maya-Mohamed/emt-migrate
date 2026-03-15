// pkg/docker/docker.go
package docker

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"github.com/Maya-Mohamed/emt-migrate/pkg/snapshot"
)

type dockerImageJSON struct {
	Repository string `json:"Repository"`
	Tag        string `json:"Tag"`
	ID         string `json:"ID"`
	Size       string `json:"Size"`
}

func GetImages() ([]snapshot.DockerImage, error) {
	cmd := exec.Command("docker", "images", "--format", "{{json .}}")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list docker images: %w", err)
	}

	var images []snapshot.DockerImage
	for _, line := range strings.Split(strings.TrimSpace(string(output)), "\n") {
		if line == "" {
			continue
		}
		var img dockerImageJSON
		if err := json.Unmarshal([]byte(line), &img); err != nil {
			continue
		}
		images = append(images, snapshot.DockerImage{
			Repository: img.Repository,
			Tag:        img.Tag,
			ID:         img.ID,
			Size:       img.Size,
		})
	}
	return images, nil
}
