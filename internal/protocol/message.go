package protocol

import (
	"encoding/json"
)

// MessageType represents the type of message being sent
type MessageType string

const (
	// Example: MessageTypeClipboardChanged MessageType = "clipboard.changed"
	MessageTypeClipboardChanged MessageType = "clipboard.changed"
	MessageTypeClipboardSet     MessageType = "clipboard.set"
	MessageTypeNotificationPush MessageType = "notification.push"
	MessageTypeCallIncoming     MessageType = "call.incoming"
	MessageTypeCallAnswer       MessageType = "call.answer"
	MessageTypeCallHangup       MessageType = "call.hangup"
	MessageTypeDeviceHello      MessageType = "device.hello"
	MessageTypeDevicePing       MessageType = "device.ping"
	MessageTypeDeviceDisconnect MessageType = "device.disconnect"
)

// Message is the base structure for all WebSocket messages
type Message struct {
	Type     MessageType     `json:"type"`
	DeviceID string          `json:"device_id"`
	Secret   string          `json:"secret"`
	Payload  json.RawMessage `json:"payload"`
}

// ClipboardPayload represents clipboard content
type ClipboardPayload struct {
	Data string `json:"data"`
}

// NotificationPayload represents a system notification
type NotificationPayload struct {
	App   string `json:"app"`
	Title string `json:"title"`
	Body  string `json:"body"`
}

// CallPayload represents call-related events
type CallPayload struct {
	Number string `json:"number"`
}

// DevicePayload represents device handshake info
type DevicePayload struct {
	DeviceName string `json:"device_name"`
}

// NewMessage creates a new Message with the given type
func NewMessage(msgType MessageType, deviceID, secret string, payload any) (*Message, error) {
	var raw json.RawMessage

	if payload != nil {
		b, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}
		raw = b
	}

	return &Message{
		Type:     msgType,
		DeviceID: deviceID,
		Secret:   secret,
		Payload:  raw,
	}, nil
}

// ParseMessage parses a JSON byte slice into a Message
func ParseMessage(data []byte) (*Message, error) {
	var parsedMsg Message

	err := json.Unmarshal(data, &parsedMsg)
	if err != nil {
		return nil, err
	}

	return &parsedMsg, nil
}

// GetPayload unmarshals the Payload field into the target struct
// Usage: msg.GetPayload(&clipboardPayload)
func (m *Message) GetPayload(target any) error {
	return json.Unmarshal(m.Payload, target)
}

// ToJSON marshals the message to JSON bytes
func (m *Message) ToJSON() ([]byte, error) {
	return json.Marshal(m)
}
