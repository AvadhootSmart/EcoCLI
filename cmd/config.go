package cmd

import (
	"fmt"

	"eco/internal/config"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configDeleteCmd)
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage eco configuration",
	Long:  `Manage eco configuration settings and files.`,
}

var configDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete eco configuration",
	Long: `Delete the eco configuration file and all associated settings.
	
This will remove:
  - ~/.config/eco/config.json

You will need to run 'eco init' again to use eco.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load()
		if err != nil {
			fmt.Printf("Error loading config: %s\n", err)
			return
		}

		if !cfg.IsInitialized() {
			fmt.Println("No configuration found to delete.")
			return
		}

		err = cfg.DeleteConfig()
		if err != nil {
			fmt.Printf("Error deleting config: %s\n", err)
			return
		}

		fmt.Println("✓ Eco configuration deleted successfully.")
		fmt.Println("Run 'eco init' to set up eco again.")
	},
}
