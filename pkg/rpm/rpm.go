// pkg/rpm/rpm.go
package rpm

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/Maya-Mohamed/emt-migrate/pkg/snapshot"
)

const queryFormat = `%{NAME}\t%{VERSION}\t%{RELEASE}\t%{ARCH}\t%{INSTALLTIME}\t%{VENDOR}\t%{SOURCERPM}\n`

func GetInstalledPackages() ([]snapshot.RPMPackage, error) {
	cmd := exec.Command("rpm", "-qa", "--queryformat", queryFormat)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to query rpm: %w", err)
	}

	var packages []snapshot.RPMPackage
	for _, line := range strings.Split(strings.TrimSpace(string(output)), "\n") {
		if line == "" {
			continue
		}
		fields := strings.Split(line, "\t")
		if len(fields) < 7 {
			continue
		}
		packages = append(packages, snapshot.RPMPackage{
			Name:        fields[0],
			Version:     fields[1],
			Release:     fields[2],
			Arch:        fields[3],
			InstallTime: fields[4],
			Vendor:      fields[5],
			SourceRPM:   fields[6],
		})
	}
	return packages, nil
}
