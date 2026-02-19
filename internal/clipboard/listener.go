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
	stop        chan struct{}
	cmd         *exec.Cmd
}

// NewListener creates a new clipboard listener
func NewListener(onChange func(content string)) *Listener {
	return &Listener{
		lastContent: "",
		onChange:    onChange,
		running:     false,
		stop:        make(chan struct{}),
		cmd:         nil,
	}
}

// Start begins monitoring the clipboard
func (l *Listener) Start() error {
	_, err := exec.LookPath("wl-paste")
	if err != nil {
		return err
	}

	cmd := exec.Command("wl-paste", "--watch", "cat")
	l.cmd = cmd
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	go func() {
		scanner := bufio.NewScanner(stdout)

		for scanner.Scan() {
			line := scanner.Text()

			if line != l.lastContent {
				l.onChange(line)
				l.lastContent = line
			}
		}
	}()

	l.running = true
	return nil
}

// Stop halts the clipboard monitoring
func (l *Listener) Stop() error {
	if !l.running {
		return nil
	}

	close(l.stop)

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
