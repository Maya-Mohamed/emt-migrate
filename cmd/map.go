package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/Maya-Mohamed/emt-migrate/pkg/mapping"
	"github.com/Maya-Mohamed/emt-migrate/pkg/snapshot"
	"github.com/spf13/cobra"
)

var mapSnapshot string
var mapTarget string

var mapCmd = &cobra.Command{
	Use:   "map",
	Short: "Show how EMT packages map to a target distribution",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := mapping.ValidateTarget(mapTarget); err != nil {
			return err
		}
		data, err := os.ReadFile(mapSnapshot)
		if err != nil {
			return fmt.Errorf("failed to read snapshot: %w", err)
		}

		var snap snapshot.Snapshot
		if err := json.Unmarshal(data, &snap); err != nil {
			return fmt.Errorf("failed to parse snapshot: %w", err)
		}

		fmt.Printf("Package mapping: EMT -> %s\n", mapTarget)
		fmt.Println("-------------------------------------------")

		skipped := 0
		mapped := 0
		direct := 0

		for _, pkg := range snap.Packages {
			targetName, skip, note := mapping.MapPackage(pkg.Name, mapTarget)
			if skip {
				fmt.Printf("  %-25s -> SKIP (%s)\n", pkg.Name, note)
				skipped++
			} else if targetName != pkg.Name {
				fmt.Printf("  %-25s -> %s (%s)\n", pkg.Name, targetName, note)
				mapped++
			} else {
				fmt.Printf("  %-25s -> %s (direct)\n", pkg.Name, targetName)
				direct++
			}
		}

		fmt.Println("-------------------------------------------")
		fmt.Printf("Total: %d packages | %d mapped | %d direct | %d skipped\n",
			len(snap.Packages), mapped, direct, skipped)

		return nil
	},
}

func init() {
	mapCmd.Flags().StringVarP(&mapSnapshot, "snapshot", "s", "snapshot.json", "Input snapshot file")
	mapCmd.Flags().StringVarP(&mapTarget, "target", "t", "rhel", "Target distro (rhel, fedora, centos, suse)")
	rootCmd.AddCommand(mapCmd)
}
