package server

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"eco/internal/auth"
	"eco/internal/config"
	"eco/internal/device"
	"eco/internal/events"
	"eco/internal/protocol"

	"github.com/gorilla/websocket"
)

// Server manages the WebSocket server and device connection
type Server struct {
	config      *config.Config
	deviceConn  *device.Connection
	upgrader    websocket.Upgrader
	eventRouter *events.Router
	httpServer  *http.Server
	staticPath  string
	pwaBaseURL  string
}

// NewServer creates a new WebSocket server
func NewServer(cfg *config.Config) *Server {
	return &Server{
		config: cfg,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		eventRouter: events.NewRouter(),
		httpServer:  &http.Server{},
	}
}

// SetStaticPath sets the path to serve static PWA files from
func (s *Server) SetStaticPath(path string) {
	s.staticPath = path
}

// SetPWABaseURL sets the base URL for the PWA (used in QR code)
func (s *Server) SetPWABaseURL(url string) {
	s.pwaBaseURL = url
}

// Start begins listening for WebSocket connections
func (s *Server) Start() error {
	// Serve PWA static files FIRST (lower priority)
	if s.staticPath != "" {
		fs := http.FileServer(http.Dir(s.staticPath))
		// Add index.html as default
		http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/" || r.URL.Path == "" {
				http.ServeFile(w, r, s.staticPath+"/index.html")
				return
			}
			fs.ServeHTTP(w, r)
		}))
	}

	// API endpoints (higher priority)
	http.HandleFunc("/ws", s.handleWebSocket)
	http.HandleFunc("/qr", s.handleQRCode)

	s.httpServer.Addr = ":4949"

	err := s.eventRouter.Start()
	if err != nil {
		return err
	}

	s.httpServer.ListenAndServe()
	fmt.Println("HTTP Server started on port 4949")
	return nil
}

// Stop gracefully shuts down the server
func (s *Server) Stop() error {
	err := s.eventRouter.Stop()
	if err != nil {
		return err
	}
	if s.deviceConn != nil && s.deviceConn.IsConnected() {
		s.deviceConn.Stop()
	}
	err = s.httpServer.Shutdown(context.TODO())
	if err != nil {
		return err
	}
	return nil
}

// handleWebSocket upgrades HTTP to WebSocket and handles the connection
func (s *Server) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	if s.deviceConn != nil && s.deviceConn.IsConnected() {
		http.Error(w, "Device already connected", http.StatusServiceUnavailable)
		return
	}

	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WS: Upgrade failed from %s: %v", r.RemoteAddr, err)
		return
	}
	log.Printf("WS: Connection upgraded from %s", r.RemoteAddr)

	// Read the first message for authentication
	_, data, err := conn.ReadMessage()
	if err != nil {
		log.Println("Failed to read auth message:", err)
		conn.Close()
		return
	}

	msg, err := protocol.ParseMessage(data)
	if err != nil {
		log.Println("Failed to parse message:", err)
		conn.Close()
		return
	}

	if auth.NewAuthenticator(s.config).ValidateCredentials(msg.DeviceID, msg.Secret) {
		log.Printf("WS: Authentication successful for device: %s", msg.DeviceID)
		s.deviceConn = device.NewConnection(msg.DeviceID, conn)
		s.deviceConn.SetHandler(s.eventRouter.CreateMessageHandler())
		s.eventRouter.SetDeviceConnection(s.deviceConn)
		go s.deviceConn.Start()
	} else {
		log.Printf("WS: Authentication failed for device: %s (ID: %s, Secret: [REDACTED])", msg.DeviceID, msg.DeviceID)
		conn.Close()
	}
}

// GetDeviceConnection returns the current device connection (if any)
func (s *Server) GetDeviceConnection() *device.Connection {
	return s.deviceConn
}

// IsDeviceConnected returns true if a device is currently connected
func (s *Server) IsDeviceConnected() bool {
	return s.deviceConn != nil && s.deviceConn.IsConnected()
}

// BroadcastEvent sends an event to the connected device
func (s *Server) BroadcastEvent(eventType protocol.MessageType, payload any) error {
	if !s.IsDeviceConnected() {
		return fmt.Errorf("no device connected")
	}

	//check appropriate secret later
	_, err := protocol.NewMessage(eventType, s.deviceConn.GetDeviceID(), s.config.SharedSecret, payload)
	if err != nil {
		return err
	}
	return nil
}

// handleQRCode generates a QR code for easy PWA connection
func (s *Server) handleQRCode(w http.ResponseWriter, r *http.Request) {
	// Determine the base URL for the PWA
	baseURL := s.pwaBaseURL
	if baseURL == "" {
		// Try to determine from request
		host := r.Host
		scheme := "http"
		if r.TLS != nil {
			scheme = "https"
		}
		baseURL = fmt.Sprintf("%s://%s", scheme, host)
	}

	// Create QR data with server info
	qrData := fmt.Sprintf("eco://connect?server=%s&secret=%s", baseURL, s.config.SharedSecret)

	// Generate QR code as PNG using an external API (qrserver.com)
	qrURL := fmt.Sprintf("https://api.qrserver.com/v1/create-qr-code/?size=300x300&data=%s", qrData)

	// Redirect to QR code image
	http.Redirect(w, r, qrURL, http.StatusFound)
}

// GetConnectionURL returns the WebSocket URL for clients
func (s *Server) GetConnectionURL() string {
	return fmt.Sprintf("ws://localhost:4949/ws")
}
