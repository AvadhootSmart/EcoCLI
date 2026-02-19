package cmd

import (
	"fmt"

	"eco/internal/config"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(devicesCmd)
}

var devicesCmd = &cobra.Command{
	Use:   "devices",
	Short: "Show registered device information",
	Long: `Display information about the registered mobile device.
	
This command shows:
  - Registered device ID
  - Connection status (if daemon is running)
  - Device credentials (optional, with --show-secret flag)

For security, the shared secret is not displayed by default.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load()
		if err != nil {
			fmt.Println(err)
			return
		}

		if !cfg.IsInitialized() {
			fmt.Println("No device registered. Run 'eco init' to register a device.")
			return
		}

		fmt.Println("Registered Device")
		fmt.Println("=================")
		fmt.Println("Device ID: " + cfg.DeviceID)

		if cmd.Flags().Changed("show-secret") {
			fmt.Println("Secret: " + cfg.SharedSecret)
		} else {
			fmt.Println("Secret:    ******** (use --show-secret to reveal)")
		}

		//   5. Display connection status (optional)
		//      - Check if daemon is running
		//      - Print: "Status: connected / disconnected / unknown"
		//      - Note: Actual connection status requires IPC with daemon

		// fmt.Println("eco devices")
	},
}

func init() {
	devicesCmd.Flags().Bool("show-secret", false, "Show the shared secret")
}
