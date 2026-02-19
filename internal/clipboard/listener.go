package clipboard

import (
	"bufio"
	"os/exec"
)

// Listener monitors the system clipboard for changes
type Listener struct {
	lastContent string
	onChange    func(string)
	running     bool
	cmd         *exec.Cmd
}

// NewListener creates a new clipboard listener
func NewListener(onChange func(content string)) *Listener {
	return &Listener{
		lastContent: "",
		onChange:    onChange,
		running:     false,
		cmd:         nil,
	}
}

// Start begins monitoring the clipboard
func (l *Listener) Start() error {
	_, err := exec.LookPath("wl-paste")
	if err != nil {
		return err
	}

	// We use wl-paste --watch to execute a command whenever the clipboard changes.
	// We use 'printf "!\n"' as a lightweight trigger that we can scan for.
	cmd := exec.Command("wl-paste", "--watch", "sh", "-c", "printf '!\n'")
	l.cmd = cmd
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	l.running = true
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			// Trigger received! Fetch full content.
			content, err := l.getContent()
			if err != nil {
				continue
			}

			if content != "" && content != l.lastContent {
				l.onChange(content)
				l.lastContent = content
			}
		}
	}()

	return nil
}

func (l *Listener) getContent() (string, error) {
	cmd := exec.Command("wl-paste", "--no-newline")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

// Stop halts the clipboard monitoring
func (l *Listener) Stop() error {
	if !l.running {
		return nil
	}

	if l.cmd != nil && l.cmd.Process != nil {
		_ = l.cmd.Process.Kill()
	}

	l.running = false
	return nil
}

// IsRunning returns whether the listener is active
func (l *Listener) IsRunning() bool {
	return l.running
}
