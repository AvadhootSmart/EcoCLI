package cmd

import (
	"fmt"

	"eco/internal/config"
	"eco/internal/crypto"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize eco (generate secret + config)",
	Long: `Initialize the eco daemon by generating device credentials and saving configuration.
	
This command will:
  1. Generate a unique device ID
  2. Generate a 32-byte shared secret for authentication
  3. Save configuration to ~/.config/eco/config.json
  4. Display the device ID and secret for mobile pairing

The secret must be entered manually on the mobile device to establish a secure connection.`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Implement init command
		// Steps:
		//   1. Check if already initialized
		//      - Load existing config
		//      - If IsInitialized(), warn user and ask to overwrite (or provide --force flag)

		cfg, err := config.Load()
		if err != nil {
			fmt.Println(err)
			return
		}
		if cfg.IsInitialized() {
			fmt.Println("eco is already initialized. Use --force to overwrite.")
			return
		}

		//   2. Generate device credentials
		//      - deviceID := crypto.GenerateDeviceID()
		//      - secret, err := crypto.GenerateSecret()

		deviceID := crypto.GenerateDeviceID()
		secret, err := crypto.GenerateSecret()
		if err != nil {
			fmt.Println(err)
			return
		}

		//   3. Create and save config
		//      - cfg := &config.Config{...}
		//      - err := cfg.Save()

		newConfig := &config.Config{
			DeviceID:     deviceID,
			SharedSecret: secret,
		}

		cfgSaveErr := newConfig.Save()
		if cfgSaveErr != nil {
			fmt.Println(cfgSaveErr)
			return
		}

		//   4. Display results to user
		//      - Print device ID
		//      - Print secret (make it prominent, maybe with formatting)
		//      - Instructions for mobile pairing
		//      - Example output:
		//        "✓ Eco initialized successfully!"
		//        ""
		//        "Device ID: mobile-xxxxx"
		//        "Secret:    xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
		//        ""
		//        "Enter these credentials in the Eco mobile app to connect."
		//        "Config saved to: ~/.config/eco/config.json"

		fmt.Println("✓ Eco initialized successfully!")
		fmt.Println("")
		fmt.Printf("Device ID: %s\n", deviceID)
		fmt.Printf("Secret:    %s\n", secret)
		fmt.Println("")
		fmt.Println("Enter these credentials in the Eco mobile app to connect.")

		cfgPath, err := config.ConfigPath()
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("Config saved to: %s\n", cfgPath)
		// println("eco init")
	},
}
