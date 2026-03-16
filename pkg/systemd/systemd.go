package systemd

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/Maya-Mohamed/emt-migrate/pkg/snapshot"
)

func GetEnabledServices() ([]snapshot.SystemdService, error) {
	cmd := exec.Command("systemctl", "list-unit-files", "--type=service", "--state=enabled", "--no-pager", "--no-legend")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to query systemd services: %w", err)
	}

	var services []snapshot.SystemdService
	for _, line := range strings.Split(strings.TrimSpace(string(output)), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		services = append(services, snapshot.SystemdService{
			Name:  fields[0],
			State: fields[1],
		})
	}
	return services, nil
}
