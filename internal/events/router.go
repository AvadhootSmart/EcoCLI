package events

import (
	"eco/internal/clipboard"
	"eco/internal/device"
	"eco/internal/notifications"
	"eco/internal/protocol"
	"encoding/json"
	"fmt"
	"log"
)

// Router handles routing system events to the connected device
type Router struct {
	deviceConn      *device.Connection
	eventChan       chan Event
	stop            chan struct{}
	running         bool
	clipboardSetter *clipboard.Setter
}

// Event represents a system event to be sent to the device
type Event struct {
	Type    protocol.MessageType
	Payload interface{}
}

// NewRouter creates a new event router
func NewRouter() *Router {
	return &Router{
		deviceConn:      nil,
		eventChan:       make(chan Event, 256),
		stop:            make(chan struct{}),
		running:         false,
		clipboardSetter: clipboard.NewSetter(),
	}
}

// SetDeviceConnection sets the current device connection
// This should be called when a device connects/disconnects
func (r *Router) SetDeviceConnection(conn *device.Connection) {
	r.deviceConn = conn
}

// Start begins processing events
func (r *Router) Start() error {
	r.running = true

	go func() {
		for {
			select {
			case event := <-r.eventChan:
				if r.deviceConn != nil && r.deviceConn.IsConnected() {
					log.Printf("Router: Routing event %s to device", event.Type)
					newMsg, err := protocol.NewMessage(event.Type, r.deviceConn.GetDeviceID(), "", event.Payload)
					if err == nil {
						if err := r.deviceConn.Send(newMsg); err != nil {
							log.Printf("Router: Failed to send message: %v", err)
						}
					} else {
						log.Printf("Router: Failed to create message: %v", err)
					}

				} else {
					log.Printf("Router: Dropping event %s (no device connected)", event.Type)
				}
			case <-r.stop:
				return
			}
		}
	}()

	return nil
}

// Stop halts the event router
func (r *Router) Stop() error {
	close(r.stop)
	r.running = false

	return nil
}

// IsRunning returns whether the router is active
func (r *Router) IsRunning() bool {
	return r.running
}

// RouteEvent queues an event to be sent to the device
func (r *Router) RouteEvent(eventType protocol.MessageType, payload any) error {
	newEvent := Event{
		Type:    eventType,
		Payload: payload,
	}

	select {
	case r.eventChan <- newEvent:
		return nil
	default:
		return fmt.Errorf("event queue full")
	}

}

// RouteClipboardChange sends clipboard change event to device
func (r *Router) RouteClipboardChange(content string) error {
	err := r.RouteEvent(protocol.MessageTypeClipboardChanged, content)
	if err != nil {
		return err
	}

	return nil
}

// RouteNotification sends notification event to device
func (r *Router) RouteNotification(app, title, body string) error {
	err := r.RouteEvent(protocol.MessageTypeNotificationPush, &protocol.NotificationPayload{
		App:   app,
		Title: title,
		Body:  body,
	})
	if err != nil {
		return err
	}
	return nil
}

// RouteCallIncoming sends incoming call event to device
func (r *Router) RouteCallIncoming(number string) error {
	err := r.RouteEvent(protocol.MessageTypeCallIncoming, &protocol.CallPayload{
		Number: number,
	})
	if err != nil {
		return err
	}
	return nil
}

// CreateMessageHandler creates a handler function for incoming messages
func (r *Router) CreateMessageHandler() func(*protocol.Message) {
	return func(msg *protocol.Message) {
		r.handleIncomingMessage(msg)
	}
}

// handleIncomingMessage processes incoming messages from the device
func (r *Router) handleIncomingMessage(msg *protocol.Message) {
	log.Printf("Router: Handling incoming message of type: %s", msg.Type)
	switch msg.Type {
	case protocol.MessageTypeClipboardSet:
		var payload protocol.ClipboardPayload
		if err := json.Unmarshal(msg.Payload, &payload); err != nil {
			fmt.Printf("Error parsing clipboard payload: %v\n", err)
			return
		}
		if err := r.clipboardSetter.SetText(payload.Data); err != nil {
			fmt.Printf("Error setting clipboard: %v\n", err)
		}

	case protocol.MessageTypeNotificationPush:
		var payload protocol.NotificationPayload
		if err := json.Unmarshal(msg.Payload, &payload); err != nil {
			fmt.Printf("Error parsing notification payload: %v\n", err)
			return
		}
		if err := notifications.Send(payload.Title, payload.Body); err != nil {
			fmt.Printf("Error sending notification: %v\n", err)
		}

	case protocol.MessageTypeCallAnswer:
		var payload protocol.CallPayload
		if err := json.Unmarshal(msg.Payload, &payload); err != nil {
			fmt.Printf("Error parsing call answer payload: %v\n", err)
			return
		}
		fmt.Printf("Call answered: %s\n", payload.Number)

	case protocol.MessageTypeCallHangup:
		fmt.Println("Call hung up")

	case protocol.MessageTypeDevicePing:

	case protocol.MessageTypeDeviceDisconnect:
		if r.deviceConn != nil {
			r.deviceConn.Stop()
		}

	default:
		fmt.Printf("Unknown message type: %s\n", msg.Type)
	}
}
