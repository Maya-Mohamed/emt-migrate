# emt-migrate

A CLI tool that captures an Edge Microvisor Toolkit (EMT) system state and generates portable migration scripts for target RPM-based distributions (RHEL, Fedora, CentOS, SUSE).

**This is a proof-of-concept for the GSoC 2025 project: "System Migration Tool for Enterprise Linux".**

## Features

- **System capture**: Queries installed RPM packages, Docker images, and enabled systemd services
- **Cross-distro package mapping**: Translates EMT package names to target distro equivalents (e.g., `moby-engine` → `docker-ce` on RHEL)
- **Migration script generation**: Produces idempotent bash scripts with checkpoint/resume logic
- **Docker offline support**: Generated scripts check for local `.tar` archives before falling back to `docker pull`
- **Multi-distro targets**: RHEL, Fedora, CentOS, SUSE

## Design

Package mappings are informed by EMT's own spec files:
- [`SPECS/moby-engine/moby-engine.spec`](https://github.com/open-edge-platform/edge-microvisor-toolkit/blob/3.0/SPECS/moby-engine/moby-engine.spec) — `Obsoletes: docker-ce`
- [`SPECS/docker-cli/docker-cli.spec`](https://github.com/open-edge-platform/edge-microvisor-toolkit/blob/3.0/SPECS/docker-cli/docker-cli.spec) — `Provides: moby-cli`
- [`SPECS/edge-release/90-default.preset`](https://github.com/open-edge-platform/edge-microvisor-toolkit/blob/3.0/SPECS/edge-release/90-default.preset) — default systemd service presets

## Usage  
  
```bash  
# Capture current system state  
./emt-migrate capture -o snapshot.json  
  
# Show package mapping for a target distro  
./emt-migrate map -s snapshot.json -t rhel  
  
# Generate migration script  
./emt-migrate generate -s snapshot.json -t rhel -o migrate.sh  
./emt-migrate generate -s snapshot.json -t suse -o migrate-suse.sh