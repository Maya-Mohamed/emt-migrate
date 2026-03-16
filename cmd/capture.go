package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/Maya-Mohamed/emt-migrate/pkg/docker"
	"github.com/Maya-Mohamed/emt-migrate/pkg/rpm"
	"github.com/Maya-Mohamed/emt-migrate/pkg/snapshot"
	"github.com/Maya-Mohamed/emt-migrate/pkg/systemd"
	"github.com/spf13/cobra"
)

var captureOutput string

var captureCmd = &cobra.Command{
	Use:   "capture",
	Short: "Capture current system state (RPM packages, Docker images, systemd services)",
	RunE: func(cmd *cobra.Command, args []string) error {
		hostname, _ := os.Hostname()
		snap := snapshot.Snapshot{
			Timestamp: time.Now(),
			Hostname:  hostname,
		}

		fmt.Println("Querying installed RPM packages...")
		pkgs, err := rpm.GetInstalledPackages()
		if err != nil {
			return fmt.Errorf("rpm capture failed: %w", err)
		}
		snap.Packages = pkgs
		fmt.Printf("Found %d packages\n", len(pkgs))

		fmt.Println("Querying Docker images...")
		imgs, err := docker.GetImages()
		if err != nil {
			fmt.Printf("Warning: Docker query failed: %s (continuing without Docker images)\n", err)
		} else {
			snap.DockerImages = imgs
			fmt.Printf("Found %d Docker images\n", len(imgs))
		}

		fmt.Println("Querying systemd services...")
		svcs, err := systemd.GetEnabledServices()
		if err != nil {
			fmt.Printf("Warning: systemd query failed: %s (continuing without services)\n", err)
		} else {
			snap.Services = svcs
			fmt.Printf("Found %d enabled services\n", len(svcs))
		}

		data, err := json.MarshalIndent(snap, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal snapshot: %w", err)
		}

		if err := os.WriteFile(captureOutput, data, 0644); err != nil {
			return fmt.Errorf("failed to write snapshot: %w", err)
		}

		fmt.Printf("Snapshot saved to %s\n", captureOutput)
		return nil
	},
}

func init() {
	captureCmd.Flags().StringVarP(&captureOutput, "output", "o", "snapshot.json", "Output file path")
	rootCmd.AddCommand(captureCmd)
}
