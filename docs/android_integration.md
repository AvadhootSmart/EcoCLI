# Android Client Integration Guide

This document provides technical details for integrating the Eco Android client with a signaling server (e.g., a Linux client written in Go).

## Network Architecture

The Android client uses a persistent WebSocket connection to communicate with the signaling server. It features:

- **Auto-reconnect**: Exponential backoff (starting at 1s, max 60s).
- **Message Queuing**: Outgoing messages are queued while offline and flushed upon reconnection.
- **Keep-Alive**: TCP Pings every 30 seconds.

## WebSocket Protocol

All communication happens over a single WebSocket endpoint using JSON-encoded messages.

### Message Envelope (`EcoEnvelope`)

Every message sent or received is wrapped in an envelope:

```json
{
  "device_id": "unique-uuid-or-id",
  "secret": "pre-shared-authentication-secret",
  "type": "message.category.action",
  "payload": { ... }
}
```

- **`device_id`**: Identifies the Android device.
- **`secret`**: Used for simple authentication/pairing.
- **`type`**: Dot-notated string (e.g., `clipboard.changed`).
- **`payload`**: Type-specific JSON object (optional).

## Connection Flow

1.  **Handshake**: Upon connection, the client sends a `device.hello` message.
2.  **Queue Flush**: Any pending messages in the client's local queue are sent immediately after the handshake.
3.  **Active State**: The connection is now ready for bi-directional message exchange.
4.  **Heartbeat**: The server should respond to `device.ping` or maintain the connection via standard WebSocket pings.

## Event Catalog

### Device Management

| Type                | Description       | Payload Schema               |
| :------------------ | :---------------- | :--------------------------- |
| `device.hello`      | Initial handshake | `{"device_name": "android"}` |
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
| `notification.push` | Notification from Android | `{"app": "name", "title": "...", "body": "..."}` |

### Telephony (Phase 1)

| Type            | Description                | Payload Schema             |
| :-------------- | :------------------------- | :------------------------- |
| `call.incoming` | Incoming call alert        | `{"number": "1234567890"}` |
| `call.answer`   | Remote request to answer   | None                       |
| `call.hangup`   | Remote request to end call | None                       |

### Advanced Features (Phase 2 Stubs)

- **File Transfer**: `file.offer`, `file.accept`
- **Input Control**: `input.key`, `input.mouse`
- **WebRTC Signaling**: `media.sdp_offer`, `media.sdp_answer`, `media.ice_candidate`

## Go Integration Tips (Signaling Server)

1.  **JSON Handling**: Use `encoding/json`. Use a `struct` with `json:"device_id"` tags for the envelope.
2.  **Concurrency**: The Android client is highly asynchronous. Your Go server should handle multiple messages in parallel using goroutines.
3.  **Authentication**: Validate the `secret` in every envelope until a more robust session-based auth is implemented.
4.  **Graceful Disconnects**: Handle WebSocket close codes correctly to avoid unnecessary reconnect loops on the Android side.
