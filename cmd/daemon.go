package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"eco/internal/clipboard"
	"eco/internal/config"
	"eco/internal/events"
	"eco/internal/notifications"
	"eco/internal/server"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(daemonCmd)
	daemonCmd.AddCommand(daemonStartCmd)
}

var daemonCmd = &cobra.Command{
	Use:   "daemon",
	Short: "Daemon control",
}

var daemonStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start eco daemon",
	Long: `Start the eco daemon which runs the WebSocket server and system listeners.
	
The daemon will:
  1. Load configuration from ~/.config/eco/config.json
  2. Start WebSocket server on port 4949
  3. Listen for clipboard changes (Wayland)
  4. Accept connections from authorized mobile devices
  5. Display notifications from mobile (using notify-send)

Run 'eco init' first if you haven't initialized the system.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Steps:
		//   1. Load configuration
		//      - cfg, err := config.Load()
		//      - Check if initialized (cfg.IsInitialized())
		//      - If not initialized, print error and exit with instructions to run 'eco init'

		cfg, err := config.Load()
		if err != nil {
			fmt.Printf("Error loading configuration: %s\n", err)
			return
		}

		if !cfg.IsInitialized() {
			fmt.Println("eco is not initialized. Run 'eco init' first.")
			return
		}

		//   2. Create and start event router
		//      - eventRouter := events.NewRouter()
		//      - eventRouter.Start()
		//      - This routes system events to the connected device

		eventRouter := events.NewRouter()
		err = eventRouter.Start()
		if err != nil {
			fmt.Printf("Error starting event router: %s\n", err)
			return
		}

		//   3. Create and start WebSocket server
		//      - server := server.NewServer(cfg)
		//      - server.Start() (run in goroutine since it blocks)

		srv := server.NewServer(cfg)

		// Find and set PWA static path
		pwaPath := findPWAPath()
		if pwaPath != "" {
			srv.SetStaticPath(pwaPath)
			fmt.Printf("Serving PWA from: %s\n", pwaPath)
		} else {
			fmt.Println("WARNING: PWA not found, serving only API endpoints")
		}

		// Print connection info
		fmt.Println("\n=== Eco Daemon Started ===")
		fmt.Println("WebSocket: ws://localhost:4949/ws")
		fmt.Printf("QR Code: http://localhost:4949/qr\n")
		fmt.Println("PWA: http://localhost:4949/")
		fmt.Println("============================\n")

		go srv.Start()

		//   4. Create and start clipboard listener
		//      - clipboardListener := clipboard.NewListener(func(content string) {
		//          - eventRouter.RouteClipboardChange(content)
		//        })
		//      - clipboardListener.Start()
		//      - Handle Wayland/X11 detection errors

		clipboardListener := clipboard.NewListener(func(content string) {
			eventRouter.RouteClipboardChange(content)
		})
		err = clipboardListener.Start()
		if err != nil {
			fmt.Printf("Error starting clipboard listener: %s\n", err)
			return
		}

		if !notifications.IsAvailable() {
			fmt.Println("notify-send command not found. Notifications will not be available.")
		}

		notifications.Send("Eco", "Eco daemon started")

		//   5. Handle incoming messages
		//      - clipboard.set messages → update local clipboard
		//      - notification messages → display via notify-send
		//      - call messages → handle call control

		// Write PID file for stop command
		err = os.WriteFile("/tmp/eco.pid", []byte(fmt.Sprintf("%d", os.Getpid())), 0644)
		if err != nil {
			fmt.Printf("Error writing PID file: %s\n", err)
		}

		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		defer func() {
			fmt.Println("Shutting down...")
			clipboardListener.Stop()
			srv.Stop()
			eventRouter.Stop()
		}()

		// println("daemon start")
	},
}

func findPWAPath() string {
	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		cwd = ""
	}

	// Check various possible locations for the PWA
	paths := []string{
		"mobile/pwa",
		"../mobile/pwa",
		"../../mobile/pwa",
		filepath.Join(cwd, "mobile/pwa"),
		filepath.Join(filepath.Dir(os.Args[0]), "mobile/pwa"),
		filepath.Join(filepath.Dir(os.Args[0]), "..", "mobile/pwa"),
		"/usr/share/eco/pwa",
		"/usr/local/share/eco/pwa",
	}

	for _, p := range paths {
		absPath, err := filepath.Abs(p)
		if err == nil {
			info, err := os.Stat(absPath)
			if err == nil && info.IsDir() {
				return absPath
			}
		}
	}

	return ""
}
