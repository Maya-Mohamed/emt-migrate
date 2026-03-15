# emt-migrate  
  
A CLI tool that captures Edge Microvisor Toolkit (EMT) system state and generates portable migration scripts for enterprise Linux distributions.  
  
## Overview  
  
`emt-migrate` automates the "works on EMT, deploy anywhere" workflow by:  
1. **Capturing** installed RPM packages and Docker images from an EMT system  
2. **Generating** idempotent shell scripts for target distributions (RHEL, Fedora, CentOS, SUSE)  
  
This is a proof-of-concept for the [GSoC 2025 System Migration Tool project](https://github.com/intel/edge-microvisor-toolkit).  
  
## Design  
  
The tool's design is informed by EMT's internal RPM tooling patterns, specifically:  
- `toolkit/tools/internal/rpm/rpm.go` — RPM query and FQN parsing  
- `toolkit/tools/imagegen/installutils/installutils.go` — Container manifest generation  
- `toolkit/scripts/spec_source_attributions.py` — Cross-distro package origin tracking  
  
## Installation  
  
```bash  
go build -o emt-migrate .