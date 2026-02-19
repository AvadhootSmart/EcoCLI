package auth

import (
	"eco/internal/config"
)

// Authenticator validates device connections
type Authenticator struct {
	config *config.Config
}

// NewAuthenticator creates a new authenticator
func NewAuthenticator(cfg *config.Config) *Authenticator {
	return &Authenticator{
		config: cfg,
	}
}

// ValidateCredentials checks if device_id and secret match the config
func (a *Authenticator) ValidateCredentials(deviceID, secret string) bool {
	return a.config.DeviceID == deviceID && a.config.SharedSecret == secret
}

// AuthenticateMessage validates a protocol.Message
// TODO: Implement AuthenticateMessage method (optional)
// This could validate the DeviceID and Secret fields of a Message
// Return error if invalid, nil if valid
// func (a *Authenticator) AuthenticateMessage(deviceID, secret string) error {
// 	return nil
// }
