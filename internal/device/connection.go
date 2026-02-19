package device

import (
	"eco/internal/protocol"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
)

// Connection represents a single WebSocket connection to the mobile device
type Connection struct {
	deviceID  string
	conn      *websocket.Conn
	send      chan *protocol.Message
	recv      chan *protocol.Message
	stop      chan struct{}
	connected bool
	handler   func(*protocol.Message)
}

// NewConnection creates a new device connection
func NewConnection(deviceID string, conn *websocket.Conn) *Connection {
	return &Connection{
		deviceID: deviceID,
		conn:     conn,
		send:     make(chan *protocol.Message, 256),
		recv:     make(chan *protocol.Message),
		stop:     make(chan struct{}),
	}
}

func (c *Connection) IsConnected() bool {
	return c.connected
}

// Start begins the read and write pumps
func (c *Connection) Start() {
	c.connected = true

	go c.readPump()
	go c.writePump()

	c.Send(&protocol.Message{
		Type: protocol.MessageTypeDeviceHello,
	})
}

// Stop gracefully closes the connection
func (c *Connection) Stop() {
	// c.stop <- struct{}{}
	close(c.stop)
	c.connected = false
	c.conn.Close()
}

// Send queues a message to be sent to the device
func (c *Connection) Send(msg *protocol.Message) error {
	if !c.connected {
		return fmt.Errorf("not connected")
	}

	select {
	case c.send <- msg:
		return nil
	default:
		return fmt.Errorf("channel full")
	}
}

// readPump reads messages from the WebSocket connection
func (c *Connection) readPump() {
	defer func() {
		close(c.stop)
		close(c.recv)
		c.conn.Close()
	}()

	c.conn.SetReadLimit(5120)
	c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, data, err := c.conn.ReadMessage()
		if err != nil {
			return
		}

		parsedMsg, err := protocol.ParseMessage(data)
		if err != nil {
			return
		}
		if c.handler != nil {
			c.handler(parsedMsg)
		}

		msg, err := protocol.ParseMessage(data)
		if err != nil {
			continue
		}

		c.recv <- msg
	}

}

// writePump writes messages to the WebSocket connection
func (c *Connection) writePump() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {

		case msg := <-c.send:
			//outgoing msg
			if err := c.conn.WriteJSON(msg); err != nil {
				c.conn.Close()
				return
			}

		case <-ticker.C:
			//ping to keep conn alive
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				c.conn.Close()
				return
			}

		case <-c.stop:
			//exit
			c.conn.Close()
			return
		}
	}
}

// GetDeviceID returns the connected device's ID
func (c *Connection) GetDeviceID() string {
	return c.deviceID
}

// SetHandler sets the message handler for incoming messages
func (c *Connection) SetHandler(handler func(*protocol.Message)) {
	c.handler = handler
}
