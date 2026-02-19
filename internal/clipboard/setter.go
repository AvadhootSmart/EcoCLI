package clipboard

import (
	// "fmt"
	"os/exec"
)

// Setter writes text to the system clipboard
type Setter struct{}

// NewSetter creates a new clipboard setter
func NewSetter() *Setter {
	return &Setter{}
}

// SetText writes text to the clipboard
func (s *Setter) SetText(text string) error {
	cmd := exec.Command("wl-copy")

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	if _, err := stdin.Write([]byte(text)); err != nil {
		return err
	}

	stdin.Close()

	return cmd.Wait()
}

// SetTextFromMessage extracts text from clipboard.set message and sets clipboard
func (s *Setter) SetTextFromMessage(payload []byte) error {
	err := s.SetText(string(payload))

	if err != nil {
		return err
	}

	return nil
}

// IsAvailable checks if clipboard commands are available
func (s *Setter) IsAvailable() bool {
	_, err := exec.LookPath("wl-copy")
	if err != nil {
		return false
	}
	return true
}
