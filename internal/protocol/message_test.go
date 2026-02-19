package protocol

import (
	"encoding/json"
	"testing"
)

func TestNewMessage(t *testing.T) {
	tests := []struct {
		name     string
		msgType  MessageType
		deviceID string
		secret   string
		payload  interface{}
		wantErr  bool
	}{
		{
			name:     "Valid clipboard message",
			msgType:  MessageTypeClipboardSet,
			deviceID: "test-device",
			secret:   "test-secret",
			payload:  &ClipboardPayload{Data: "Hello"},
			wantErr:  false,
		},
		{
			name:     "Valid notification message",
			msgType:  MessageTypeNotificationPush,
			deviceID: "test-device",
			secret:   "test-secret",
			payload:  &NotificationPayload{App: "Test", Title: "Title", Body: "Body"},
			wantErr:  false,
		},
		{
			name:     "Message with nil payload",
			msgType:  MessageTypeDevicePing,
			deviceID: "test-device",
			secret:   "test-secret",
			payload:  nil,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg, err := NewMessage(tt.msgType, tt.deviceID, tt.secret, tt.payload)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if msg == nil {
				t.Error("NewMessage() returned nil message")
				return
			}
			if msg.Type != tt.msgType {
				t.Errorf("NewMessage() Type = %v, want %v", msg.Type, tt.msgType)
			}
			if msg.DeviceID != tt.deviceID {
				t.Errorf("NewMessage() DeviceID = %v, want %v", msg.DeviceID, tt.deviceID)
			}
			if msg.Secret != tt.secret {
				t.Errorf("NewMessage() Secret = %v, want %v", msg.Secret, tt.secret)
			}
		})
	}
}

func TestParseMessage(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		wantErr bool
	}{
		{
			name:    "Valid JSON message",
			data:    []byte(`{"type":"clipboard.set","device_id":"test","secret":"secret","payload":{"data":"Hello"}}`),
			wantErr: false,
		},
		{
			name:    "Invalid JSON",
			data:    []byte(`{"type":invalid}`),
			wantErr: true,
		},
		{
			name:    "Empty data",
			data:    []byte{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg, err := ParseMessage(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && msg == nil {
				t.Error("ParseMessage() returned nil message for valid input")
			}
		})
	}
}

func TestMessageJSONTags(t *testing.T) {
	// Verify JSON serialization produces lowercase field names (Android compatible)
	msg := &Message{
		Type:     MessageTypeClipboardSet,
		DeviceID: "test-device",
		Secret:   "test-secret",
	}
	payload := &ClipboardPayload{Data: "Hello World"}
	payloadBytes, _ := json.Marshal(payload)
	msg.Payload = payloadBytes

	data, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("Failed to marshal message: %v", err)
	}

	// Verify lowercase field names
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if _, ok := raw["device_id"]; !ok {
		t.Error("Missing 'device_id' field (expected lowercase)")
	}
	if _, ok := raw["DeviceID"]; ok {
		t.Error("Found 'DeviceID' field (should be lowercase 'device_id')")
	}

	// Verify payload has lowercase fields
	if payloadRaw, ok := raw["payload"].(map[string]interface{}); ok {
		if _, ok := payloadRaw["data"]; !ok {
			t.Error("Payload missing 'data' field (expected lowercase)")
		}
	}
}

func TestPayloadStructures(t *testing.T) {
	tests := []struct {
		name     string
		payload  interface{}
		expected map[string]interface{}
	}{
		{
			name:    "ClipboardPayload lowercase",
			payload: &ClipboardPayload{Data: "test"},
			expected: map[string]interface{}{
				"data": "test",
			},
		},
		{
			name:    "NotificationPayload lowercase",
			payload: &NotificationPayload{App: "TestApp", Title: "Title", Body: "Body"},
			expected: map[string]interface{}{
				"app":   "TestApp",
				"title": "Title",
				"body":  "Body",
			},
		},
		{
			name:    "CallPayload lowercase",
			payload: &CallPayload{Number: "+1234567890"},
			expected: map[string]interface{}{
				"number": "+1234567890",
			},
		},
		{
			name:    "DevicePayload lowercase",
			payload: &DevicePayload{DeviceName: "android"},
			expected: map[string]interface{}{
				"device_name": "android",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.payload)
			if err != nil {
				t.Fatalf("Failed to marshal payload: %v", err)
			}

			var result map[string]interface{}
			if err := json.Unmarshal(data, &result); err != nil {
				t.Fatalf("Failed to unmarshal: %v", err)
			}

			for key, expectedVal := range tt.expected {
				if val, ok := result[key]; !ok {
					t.Errorf("Missing field '%s'", key)
				} else if val != expectedVal {
					t.Errorf("Field '%s' = %v, want %v", key, val, expectedVal)
				}
			}
		})
	}
}
