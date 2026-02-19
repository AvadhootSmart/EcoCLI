package cmd

import (
	"fmt"
	"os"
	"strconv"
	"syscall"
	"time"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(stopCmd)
}

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop eco daemon",
	Long: `Stop the running eco daemon gracefully.
	
This command sends a termination signal to the daemon process
causing it to shutdown cleanly and close all connections.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Read PID from file
		pidBytes, err := os.ReadFile("/tmp/eco.pid")
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Println("Daemon is not running (PID file not found)")
				return
			}
			fmt.Printf("Error reading daemon PID: %s\n", err)
			return
		}

		// Parse PID
		pidStr := string(pidBytes)
		if pidStr == "" {
			fmt.Println("Daemon is not running (PID file is empty)")
			return
		}

		pid, err := strconv.Atoi(pidStr)
		if err != nil {
			fmt.Printf("Invalid PID in file: %s\n", err)
			return
		}

		// Find process
		process, err := os.FindProcess(pid)
		if err != nil {
			fmt.Printf("Error finding daemon process: %s\n", err)
			return
		}

		// Send SIGTERM signal
		fmt.Printf("Stopping eco daemon (PID: %d)...\n", pid)
		err = process.Signal(syscall.SIGTERM)
		if err != nil {
			fmt.Printf("Error sending signal to daemon: %s\n", err)
			return
		}

		// Wait briefly for graceful shutdown
		time.Sleep(2 * time.Second)

		// Check if process still exists (it might have already exited)
		err = process.Signal(syscall.Signal(0))
		if err == nil {
			// Process still running, try SIGKILL as last resort
			fmt.Println("Process still running, sending SIGKILL...")
			process.Signal(syscall.SIGKILL)
			time.Sleep(1 * time.Second)
		}

		// Clean up PID file
		err = os.Remove("/tmp/eco.pid")
		if err != nil && !os.IsNotExist(err) {
			fmt.Printf("Warning: Could not remove PID file: %s\n", err)
		}

		fmt.Println("Eco daemon stopped")
	},
}
