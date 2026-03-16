// pkg/mapping/mapping.go
package mapping

import "github.com/Maya-Mohamed/emt-migrate/pkg/snapshot"

type PackageMapping struct {
	EMTName    string
	RHELName   string
	FedoraName string
	CentOSName string
	SUSEName   string
	Note       string
}

type MappedPackage struct {
	Original   snapshot.RPMPackage
	TargetName string
	Skip       bool
	Note       string
	Changed    bool
}

var knownMappings = []PackageMapping{
	{EMTName: "moby-engine", RHELName: "docker-ce", FedoraName: "docker-ce", CentOSName: "docker-ce", SUSEName: "docker", Note: "EMT uses Moby (upstream Docker); RHEL/Fedora use docker-ce from Docker repo"},
	{EMTName: "moby-cli", RHELName: "docker-ce-cli", FedoraName: "docker-ce-cli", CentOSName: "docker-ce-cli", SUSEName: "docker", Note: "EMT moby-cli was renamed to docker-cli; RHEL uses docker-ce-cli"},
	{EMTName: "containerd", RHELName: "containerd.io", FedoraName: "containerd.io", CentOSName: "containerd.io", SUSEName: "containerd", Note: ""},
	{EMTName: "docker-init", RHELName: "docker-ce-rootless-extras", FedoraName: "docker-ce-rootless-extras", CentOSName: "docker-ce-rootless-extras", SUSEName: "-", Note: "docker-init (tini) bundled differently per distro"},
	{EMTName: "tdnf", RHELName: "-", FedoraName: "-", CentOSName: "-", SUSEName: "-", Note: "EMT-only package manager (tdnf); target distros use dnf/zypper natively"},
	{EMTName: "edge-release", RHELName: "-", FedoraName: "-", CentOSName: "-", SUSEName: "-", Note: "EMT-specific release package"},
	{EMTName: "azure-linux-repos", RHELName: "-", FedoraName: "-", CentOSName: "-", SUSEName: "-", Note: "EMT/Azure Linux repository config"},
	{EMTName: "msft-golang", RHELName: "golang", FedoraName: "golang", CentOSName: "golang", SUSEName: "go", Note: "Microsoft Go fork maps to standard Go"},
	{EMTName: "vim", RHELName: "vim-enhanced", FedoraName: "vim-enhanced", CentOSName: "vim-enhanced", SUSEName: "vim", Note: ""},
	{EMTName: "python3", RHELName: "python3", FedoraName: "python3", CentOSName: "python3", SUSEName: "python3", Note: ""},
	{EMTName: "kernel", RHELName: "kernel", FedoraName: "kernel", CentOSName: "kernel", SUSEName: "kernel-default", Note: "SUSE uses kernel-default"},
	{EMTName: "glibc-devel", RHELName: "glibc-devel", FedoraName: "glibc-devel", CentOSName: "glibc-devel", SUSEName: "glibc-devel", Note: ""},
	{EMTName: "openssl", RHELName: "openssl", FedoraName: "openssl", CentOSName: "openssl", SUSEName: "openssl", Note: ""},
	{EMTName: "curl", RHELName: "curl", FedoraName: "curl", CentOSName: "curl", SUSEName: "curl", Note: ""},
	{EMTName: "systemd", RHELName: "systemd", FedoraName: "systemd", CentOSName: "systemd", SUSEName: "systemd", Note: ""},
	{EMTName: "qemu-img", RHELName: "qemu-img", FedoraName: "qemu-img", CentOSName: "qemu-img", SUSEName: "qemu-tools", Note: "SUSE packages qemu tools differently"},
}

func findMapping(emtName string) (PackageMapping, bool) {
	for _, m := range knownMappings {
		if m.EMTName == emtName {
			return m, true
		}
	}
	return PackageMapping{}, false
}

func MapPackage(emtName string, targetDistro string) (string, bool, string) {
	m, exists := findMapping(emtName)
	if !exists {
		return emtName, false, ""
	}

	var targetName string
	switch targetDistro {
	case "rhel":
		targetName = m.RHELName
	case "fedora":
		targetName = m.FedoraName
	case "centos":
		targetName = m.CentOSName
	case "suse":
		targetName = m.SUSEName
	default:
		targetName = ""
	}

	if targetName == "-" {
		return "", true, m.Note
	}
	if targetName == "" {
		return emtName, false, m.Note
	}
	return targetName, false, m.Note
}

func MapAll(packages []snapshot.RPMPackage, targetDistro string) []MappedPackage {
	var mapped []MappedPackage
	for _, pkg := range packages {
		targetName, skip, note := MapPackage(pkg.Name, targetDistro)
		mapped = append(mapped, MappedPackage{
			Original:   pkg,
			TargetName: targetName,
			Skip:       skip,
			Note:       note,
			Changed:    targetName != pkg.Name && !skip,
		})
	}
	return mapped
}
