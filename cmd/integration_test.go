//go:build integration
// +build integration

package cmd

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// TestCLIIntegration performs end-to-end testing of the CLI
func TestCLIIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Build the binary
	binaryPath := filepath.Join(t.TempDir(), "eco")
	buildCmd := exec.Command("go", "build", "-o", binaryPath, ".")
	buildCmd.Dir = "../" // Project root
	if output, err := buildCmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed to build binary: %v\nOutput: %s", err, output)
	}

	// Set up temporary home directory for isolated config
	tempHome := t.TempDir()
	origHome := os.Getenv("HOME")
	os.Setenv("HOME", tempHome)
	defer os.Setenv("HOME", origHome)

	t.Run("InitCommand", func(t *testing.T) {
		cmd := exec.Command(binaryPath, "init")
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("init command failed: %v\nOutput: %s", err, output)
		}

		outputStr := string(output)
		if !strings.Contains(outputStr, "Eco initialized successfully") {
			t.Errorf("init output missing success message: %s", outputStr)
		}
		if !strings.Contains(outputStr, "Device ID:") {
			t.Error("init output missing Device ID")
		}
		if !strings.Contains(outputStr, "Secret:") {
			t.Error("init output missing Secret")
		}

		// Verify config file exists
		configPath := filepath.Join(tempHome, ".config", "eco", "config.json")
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			t.Error("Config file was not created")
		}
	})

	t.Run("StatusCommand", func(t *testing.T) {
		cmd := exec.Command(binaryPath, "status")
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("status command failed: %v\nOutput: %s", err, output)
		}

		outputStr := string(output)
		if !strings.Contains(outputStr, "Initialized: yes") {
			t.Errorf("status shows not initialized: %s", outputStr)
		}
		if !strings.Contains(outputStr, "Device ID:") {
			t.Error("status output missing Device ID")
		}
	})

	t.Run("DevicesCommand", func(t *testing.T) {
		cmd := exec.Command(binaryPath, "devices")
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("devices command failed: %v\nOutput: %s", err, output)
		}

		outputStr := string(output)
		if !strings.Contains(outputStr, "Registered Device") {
			t.Errorf("devices output missing header: %s", outputStr)
		}
		if !strings.Contains(outputStr, "Device ID:") {
			t.Error("devices output missing Device ID")
		}
	})

	t.Run("DevicesCommandWithSecret", func(t *testing.T) {
		cmd := exec.Command(binaryPath, "devices", "--show-secret")
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("devices --show-secret command failed: %v\nOutput: %s", err, output)
		}

		outputStr := string(output)
		// Should show secret (not masked)
		if strings.Contains(outputStr, "Secret: ********") {
			t.Error("devices --show-secret should show actual secret, not asterisks")
		}
	})

	t.Run("DaemonStartAndStop", func(t *testing.T) {
		// Start daemon
		startCmd := exec.Command(binaryPath, "daemon", "start")
		if err := startCmd.Start(); err != nil {
			t.Fatalf("Failed to start daemon: %v", err)
		}

		// Wait for daemon to start
		time.Sleep(2 * time.Second)

		// Check PID file exists
		pidFile := "/tmp/eco.pid"
		if _, err := os.Stat(pidFile); os.IsNotExist(err) {
			startCmd.Process.Kill()
			t.Fatal("Daemon PID file not created")
		}

		// Stop daemon
		stopCmd := exec.Command(binaryPath, "stop")
		stopOutput, err := stopCmd.CombinedOutput()
		if err != nil {
			startCmd.Process.Kill()
			t.Fatalf("stop command failed: %v\nOutput: %s", err, stopOutput)
		}

		if !strings.Contains(string(stopOutput), "Eco daemon stopped") {
			t.Errorf("stop output missing success message: %s", stopOutput)
		}

		// Wait for daemon to stop
		time.Sleep(1 * time.Second)

		// Verify PID file is removed
		if _, err := os.Stat(pidFile); !os.IsNotExist(err) {
			t.Error("PID file was not removed after stop")
		}
	})

	t.Run("StopWhenNotRunning", func(t *testing.T) {
		// Ensure no PID file
		os.Remove("/tmp/eco.pid")

		cmd := exec.Command(binaryPath, "stop")
		output, err := cmd.CombinedOutput()
		// This should not error, just inform user
		if err != nil {
			t.Logf("stop command returned error (may be expected): %v", err)
		}

		outputStr := string(output)
		if !strings.Contains(outputStr, "not running") {
			t.Errorf("Expected 'not running' message, got: %s", outputStr)
		}
	})

	t.Run("ConfigDelete", func(t *testing.T) {
		cmd := exec.Command(binaryPath, "config", "delete")
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("config delete command failed: %v\nOutput: %s", err, output)
		}

		if !strings.Contains(string(output), "deleted successfully") {
			t.Errorf("config delete output missing success message: %s", output)
		}

		// Verify config is deleted
		configPath := filepath.Join(tempHome, ".config", "eco", "config.json")
		if _, err := os.Stat(configPath); !os.IsNotExist(err) {
			t.Error("Config file was not deleted")
		}
	})

	t.Run("StatusAfterDelete", func(t *testing.T) {
		cmd := exec.Command(binaryPath, "status")
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("status command failed: %v\nOutput: %s", err, output)
		}

		if !strings.Contains(string(output), "not initialized") {
			t.Errorf("Expected 'not initialized' message after delete: %s", output)
		}
	})

	t.Run("HelpCommands", func(t *testing.T) {
		commands := [][]string{
			{"--help"},
			{"init", "--help"},
			{"daemon", "--help"},
			{"stop", "--help"},
			{"status", "--help"},
			{"devices", "--help"},
			{"config", "--help"},
			{"config", "delete", "--help"},
		}

		for _, args := range commands {
			cmd := exec.Command(binaryPath, args...)
			output, err := cmd.CombinedOutput()
			if err != nil {
				t.Errorf("%v --help failed: %v\nOutput: %s", args, err, output)
			}

			outputStr := string(output)
			if !strings.Contains(outputStr, "Usage:") {
				t.Errorf("%v help missing Usage section", args)
			}
		}
	})
}
