package notifications

import (
	// "fmt"
	"os/exec"
)

// Send displays a desktop notification using notify-send
func Send(title, body string) error {
	if IsAvailable() {
		cmd := exec.Command("notify-send", title, body)
		return cmd.Run()
	}
	return nil
}

// SendWithIcon displays a notification with an icon
func SendWithIcon(title, body, iconPath string) error {
	if IsAvailable() {
		cmd := exec.Command("notify-send", "--icon", iconPath, title, body)
		return cmd.Run()
	}
	return nil
}

// IsAvailable checks if notify-send command exists
func IsAvailable() bool {
	_, err := exec.LookPath("notify-send")
	if err != nil {
		return false
	}
	return true
}
