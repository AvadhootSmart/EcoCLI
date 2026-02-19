package cmd

import (
	"fmt"

	"eco/internal/config"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(statusCmd)
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show eco daemon status",
	Long: `Display the current status of the eco system including:
  - Initialization status
  - Device connection status
  - Active listeners (clipboard, notifications)
  - Last sync time (optional, for future enhancement)`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Implement status command
		// Steps:
		//   1. Load configuration
		//      - cfg, err := config.Load()

		cfg, err := config.Load()
		if err != nil {
			fmt.Println(err)
			return
		}

		if !cfg.IsInitialized() {
			fmt.Println("Eco is not initialized. Run 'eco init' to set up.")
			return
		}

		fmt.Println("Eco Status")
		fmt.Println("==========")
		fmt.Println("Initialized: yes")
		fmt.Println("Device ID:   " + cfg.DeviceID)

		cfgPath, err := config.ConfigPath()
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("Config path", cfgPath)

		//   3. Display configuration info
		//      - Device ID
		//      - Config file path
		//   4. Check daemon status (optional for MVP)
		//      - Check if daemon is running (could check PID file or just state)
		//      - Print: "Daemon: running / not running"
		//
		//   5. Display device connection status (optional for MVP)
		//      - This would require inter-process communication
		//      - Could be implemented later with a status file or socket
		//      - For now, just print: "Run 'eco daemon start' to begin"
		//
		//   6. Format output nicely
		//      Example:
		//      "Eco Status"
		//      "=========="
		//      "Initialized: yes"
		//      "Device ID:   mobile-xxxxx"
		//      "Daemon:      not running"
		//      ""
		//      "Run 'eco daemon start' to start the daemon."

		// fmt.Println("eco status")
	},
}
