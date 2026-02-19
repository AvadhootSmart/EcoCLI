# Mobile Client (PWA) Integration Guide

This document provides technical details for integrating the Eco PWA mobile client with a signaling server (e.g., a Linux client written in Go).

## Quick Start

1. Host the PWA files on a web server (or serve from the CLI directly)
2. Open the URL in a mobile browser
3. Enter the server WebSocket URL and shared secret
4. Tap "Connect" to pair with the CLI

## Network Architecture

The PWA uses a persistent WebSocket connection to communicate with the CLI server. It features:

- **Auto-reconnect**: Exponential backoff (starting at 1s, max 60s)
- **Message Queuing**: Outgoing messages are queued while offline and flushed upon reconnection
- **Keep-Alive**: Ping messages every 30 seconds
- **Offline Support**: Service worker caches assets for offline viewing

## WebSocket Protocol

All communication happens over a single WebSocket endpoint using JSON-encoded messages.

### Message Envelope

Every message sent or received is wrapped in an envelope:

```json
{
  "device_id": "unique-uuid-or-id",
  "secret": "pre-shared-authentication-secret",
  "type": "message.category.action",
  "payload": { ... }
}
```

- **`device_id`**: Identifies the mobile device
- **`secret`**: Used for simple authentication/pairing
- **`type`**: Dot-notated string (e.g., `clipboard.changed`)
- **`payload`**: Type-specific JSON object (optional)

## Connection Flow

1. **Connect**: User enters server URL and secret, clicks Connect
2. **Handshake**: Client sends a `device.hello` message
3. **Active State**: The connection is ready for bi-directional message exchange
4. **Heartbeat**: Client sends `device.ping` every 30 seconds

## Event Catalog

### Device Management

| Type                | Description       | Payload Schema               |
| :------------------ | :---------------- | :--------------------------- |
| `device.hello`      | Initial handshake | `{"device_name": "PWA"}`    |
| `device.ping`       | Heartbeat         | None                         |
| `device.disconnect` | Graceful shutdown | None                         |

### Clipboard

| Type                | Description                     | Payload Schema             |
| :------------------ | :------------------------------ | :------------------------- |
| `clipboard.changed` | Device clipboard updated        | `{"data": "text content"}` |
| `clipboard.set`     | Request to set device clipboard | `{"data": "text content"}` |

### Notifications

| Type                | Description               | Payload Schema                                   |
| :------------------ | :------------------------ | :----------------------------------------------- |
| `notification.push` | Notification from server | `{"app": "name", "title": "...", "body": "..."}` |

### Telephony (Phase 1)

| Type            | Description                | Payload Schema             |
| :-------------- | :------------------------- | :------------------------- |
| `call.incoming` | Incoming call alert        | `{"number": "1234567890"}` |
| `call.answer`   | Remote request to answer   | None                       |
| `call.hangup`   | Remote request to end call | None                       |

## Go Server Integration

1. **JSON Handling**: Use `encoding/json` with structs having `json:"device_id"` tags
2. **Concurrency**: Handle multiple messages in parallel using goroutines
3. **Authentication**: Validate the `secret` in every envelope
4. **Graceful Disconnects**: Handle WebSocket close codes correctly

## PWA Features

- **Installable**: Add to home screen on iOS and Android
- **Offline Mode**: View cached content when disconnected
- **Push Notifications**: Receive notifications from server
- **Clipboard Sync**: Share clipboard between device and CLI

## Files

- `index.html` - Main application UI
- `app.js` - Application logic and UI handling
- `client.js` - WebSocket client with reconnection
- `styles.css` - Mobile-optimized styling
- `sw.js` - Service worker for offline support
- `manifest.json` - PWA manifest for installability
