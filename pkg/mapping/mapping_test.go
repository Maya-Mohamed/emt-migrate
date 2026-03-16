package mapping

import "testing"

func TestMapPackageKnownRHEL(t *testing.T) {
	name, skip, _ := MapPackage("moby-engine", "rhel")
	if skip {
		t.Fatal("moby-engine should not be skipped")
	}
	if name != "docker-ce" {
		t.Errorf("expected docker-ce, got %s", name)
	}
}

func TestMapPackageKnownSUSE(t *testing.T) {
	name, skip, _ := MapPackage("moby-engine", "suse")
	if skip {
		t.Fatal("moby-engine should not be skipped")
	}
	if name != "docker" {
		t.Errorf("expected docker, got %s", name)
	}
}

func TestMapPackageUnknown(t *testing.T) {
	name, skip, _ := MapPackage("some-unknown-pkg", "rhel")
	if skip {
		t.Fatal("unknown package should not be skipped")
	}
	if name != "some-unknown-pkg" {
		t.Errorf("expected some-unknown-pkg, got %s", name)
	}
}

func TestMapPackageSkipped(t *testing.T) {
	_, skip, _ := MapPackage("tdnf", "rhel")
	if !skip {
		t.Fatal("tdnf should be skipped")
	}
}

func TestValidateTargetValid(t *testing.T) {
	for _, d := range SupportedDistros {
		if err := ValidateTarget(d); err != nil {
			t.Errorf("expected %s to be valid, got error: %v", d, err)
		}
	}
}

func TestValidateTargetInvalid(t *testing.T) {
	if err := ValidateTarget("ubuntu"); err == nil {
		t.Fatal("expected error for unsupported distro ubuntu")
	}
}

func TestMappingIndexContainsAllEntries(t *testing.T) {
	for _, m := range knownMappings {
		if _, ok := mappingIndex[m.EMTName]; !ok {
			t.Errorf("mappingIndex missing entry for %s", m.EMTName)
		}
	}
}
