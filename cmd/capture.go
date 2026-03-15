// cmd/capture.go
package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/Maya-Mohamed/emt-migrate/pkg/docker"
	"github.com/Maya-Mohamed/emt-migrate/pkg/rpm"
	"github.com/Maya-Mohamed/emt-migrate/pkg/snapshot"
	"github.com/spf13/cobra"
)

var outputFile string

var captureCmd = &cobra.Command{
	Use:   "capture",
	Short: "Capture current system state (RPM packages, Docker images)",
	RunE: func(cmd *cobra.Command, args []string) error {
		hostname, _ := os.Hostname()

		fmt.Println("Querying installed RPM packages...")
		packages, err := rpm.GetInstalledPackages()
		if err != nil {
			return fmt.Errorf("rpm capture failed: %w", err)
		}
		fmt.Printf("Found %d packages\n", len(packages))

		fmt.Println("Querying Docker images...")
		images, err := docker.GetImages()
		if err != nil {
			fmt.Printf("Warning: Docker query failed: %v (continuing without Docker images)\n", err)
		} else {
			fmt.Printf("Found %d images\n", len(images))
		}

		snap := snapshot.Snapshot{
			Timestamp:    time.Now(),
			Hostname:     hostname,
			Packages:     packages,
			DockerImages: images,
		}

		data, err := json.MarshalIndent(snap, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal snapshot: %w", err)
		}

		if err := os.WriteFile(outputFile, data, 0644); err != nil {
			return fmt.Errorf("failed to write snapshot: %w", err)
		}

		fmt.Printf("Snapshot saved to %s\n", outputFile)
		return nil
	},
}

func init() {
	captureCmd.Flags().StringVarP(&outputFile, "output", "o", "snapshot.json", "Output file path")
	rootCmd.AddCommand(captureCmd)
}
