package config

import (
	"os"
	"testing"
)

func TestConfigPath(t *testing.T) {
	path, err := ConfigPath()
	if err != nil {
		t.Fatalf("ConfigPath() error = %v", err)
	}

	if path == "" {
		t.Error("ConfigPath() returned empty path")
	}

	// Should contain config directory and file
	if !contains(path, ConfigDir) {
		t.Errorf("ConfigPath() %v does not contain %v", path, ConfigDir)
	}
	if !contains(path, ConfigFile) {
		t.Errorf("ConfigPath() %v does not contain %v", path, ConfigFile)
	}
}

func TestLoadAndSave(t *testing.T) {
	// Create temporary directory for test
	tempDir, err := os.MkdirTemp("", "eco-config-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Temporarily override config path
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	t.Run("Load non-existent config", func(t *testing.T) {
		cfg, err := Load()
		if err != nil {
			t.Errorf("Load() error = %v", err)
		}
		if cfg == nil {
			t.Error("Load() returned nil config")
		}
		if cfg.IsInitialized() {
			t.Error("Load() returned initialized config for non-existent file")
		}
	})

	t.Run("Save and load config", func(t *testing.T) {
		cfg := &Config{
			DeviceID:     "test-device-123",
			SharedSecret: "test-secret-abc",
			Port:         4949,
		}

		if err := cfg.Save(); err != nil {
			t.Fatalf("Save() error = %v", err)
		}

		// Verify file exists
		cfgPath, _ := ConfigPath()
		if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
			t.Error("Config file was not created")
		}

		// Load and verify
		loaded, err := Load()
		if err != nil {
			t.Fatalf("Load() error = %v", err)
		}

		if loaded.DeviceID != cfg.DeviceID {
			t.Errorf("Load() DeviceID = %v, want %v", loaded.DeviceID, cfg.DeviceID)
		}
		if loaded.SharedSecret != cfg.SharedSecret {
			t.Errorf("Load() SharedSecret = %v, want %v", loaded.SharedSecret, cfg.SharedSecret)
		}
		if loaded.Port != cfg.Port {
			t.Errorf("Load() Port = %v, want %v", loaded.Port, cfg.Port)
		}
	})
}

func TestIsInitialized(t *testing.T) {
	tests := []struct {
		name     string
		config   Config
		expected bool
	}{
		{
			name:     "Fully initialized",
			config:   Config{DeviceID: "device-123", SharedSecret: "secret-abc"},
			expected: true,
		},
		{
			name:     "Missing DeviceID",
			config:   Config{DeviceID: "", SharedSecret: "secret-abc"},
			expected: false,
		},
		{
			name:     "Missing SharedSecret",
			config:   Config{DeviceID: "device-123", SharedSecret: ""},
			expected: false,
		},
		{
			name:     "Empty config",
			config:   Config{},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.config.IsInitialized(); got != tt.expected {
				t.Errorf("IsInitialized() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestGetDeviceCredentials(t *testing.T) {
	cfg := &Config{
		DeviceID:     "test-device",
		SharedSecret: "test-secret",
	}

	deviceID, secret := cfg.GetDeviceCredentials()

	if deviceID != cfg.DeviceID {
		t.Errorf("GetDeviceCredentials() deviceID = %v, want %v", deviceID, cfg.DeviceID)
	}
	if secret != cfg.SharedSecret {
		t.Errorf("GetDeviceCredentials() secret = %v, want %v", secret, cfg.SharedSecret)
	}
}

func TestDeleteConfig(t *testing.T) {
	// Create temporary directory for test
	tempDir, err := os.MkdirTemp("", "eco-config-test-delete")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Temporarily override config path
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	t.Run("Delete non-existent config", func(t *testing.T) {
		cfg := &Config{}
		err := cfg.DeleteConfig()
		if err == nil {
			t.Error("DeleteConfig() should error for non-initialized config")
		}
	})

	t.Run("Delete existing config", func(t *testing.T) {
		cfg := &Config{
			DeviceID:     "test-device",
			SharedSecret: "test-secret",
		}

		// Save first
		if err := cfg.Save(); err != nil {
			t.Fatalf("Save() error = %v", err)
		}

		// Verify file exists
		cfgPath, _ := ConfigPath()
		if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
			t.Fatal("Config file was not created")
		}

		// Delete
		if err := cfg.DeleteConfig(); err != nil {
			t.Errorf("DeleteConfig() error = %v", err)
		}

		// Verify file is deleted
		if _, err := os.Stat(cfgPath); !os.IsNotExist(err) {
			t.Error("Config file was not deleted")
		}
	})
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
