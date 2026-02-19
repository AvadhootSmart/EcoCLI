# 🧩 Project: Linux ↔ Mobile Ecosystem (MVP)

## Goal

Create a **local-network ecosystem app** that connects Linux and Mobile to:

* Sync clipboard both ways
* Mirror notifications (Linux → Mobile)
* Control phone calls (answer / hang up from Linux)

No cloud.
No accounts.
No P2P.

LAN only.

---

# 🏗 Architecture

### Topology

```
Linux (Server + CLI + Daemon)
        ↑ WebSocket
        ↓
Mobile (Client + Widget + Background Service)
```

Linux is the **authority node**.

Mobile connects to Linux over LAN via WebSocket.

---

## Transport

* Protocol: WebSocket
* Encoding: JSON
* One persistent connection per device

Example:

```
ws://<linux-ip>:4949
```

---

## Security (MVP)

### Authentication Model

* Linux generates:

  * `device_id`
  * `shared_secret` (32 random bytes)

Mobile receives secret via QR / manual entry.

Every message includes:

```json
{
  "device_id": "mobile-1",
  "secret": "shared_secret",
  "type": "...",
  "payload": {}
}
```

Linux validates:

* device_id exists
* secret matches

Otherwise → drop connection.

---

### Stored Config

Linux:

```
~/.config/eco/config.json
```

Mobile:

Secure storage / keychain.

---

# 🖥 Linux App

## Form

CLI + background daemon.

Binary: `eco`

---

## Commands

```bash
eco init        # generate secret + show QR
eco start       # start websocket server + listeners
eco status
eco devices
eco stop
```

---

## Responsibilities

### 1. WebSocket Server

* Accept connections
* Authenticate clients
* Dispatch events

---

### 2. Clipboard Listener

Linux:

* Wayland: `wl-paste --watch`
* X11: `xclip`

On clipboard change:

Emit:

```json
{
  "type": "clipboard.changed",
  "data": "text"
}
```

---

### 3. Clipboard Setter

On receiving:

```json
{
  "type": "clipboard.set",
  "data": "text"
}
```

Linux updates clipboard.

---

### 4. Notification Listener

Via DBus:

* `org.freedesktop.Notifications`

On notification:

```json
{
  "type": "notification",
  "app": "Slack",
  "title": "...",
  "body": "..."
}
```

Send to mobile.

---

### 5. Call Control Receiver

Accept:

```json
{
  "type": "call.answer"
}
```

```json
{
  "type": "call.hangup"
}
```

Forward to mobile.

Linux does NOT manage calls — only relays intent.

---

# 📱 Mobile App

## Components

* Background WebSocket client
* Home widget
* Notification service
* Telephony integration

---

## Responsibilities

### 1. Connect to Linux

* IP entered or QR scanned
* Persistent WS connection
* Auto reconnect

---

### 2. Clipboard Sync

Widget allows:

* Paste to Linux
* Copy from Linux

Events:

```json
clipboard.set
clipboard.changed
```

---

### 3. Receive Linux Notifications

Display as native mobile notifications.

---

### 4. Phone Call Bridge

Mobile listens for phone state.

On incoming call:

Send to Linux:

```json
{
  "type": "call.incoming",
  "number": "unknown"
}
```

Linux UI shows call popup.

Linux can reply:

```json
call.answer
call.hangup
```

Mobile executes action.

---

# 📡 Message Types

## Clipboard

```
clipboard.changed
clipboard.set
```

---

## Notifications

```
notification.push
```

---

## Calls

```
call.incoming
call.answer
call.hangup
```

---

## System

```
device.hello
device.ping
device.disconnect
```

---

# 🛣 MVP Milestones

## Phase 1

* WS connection
* Secret auth
* Manual IP pairing

---

## Phase 2

* Clipboard both ways

---

## Phase 3

* Linux → Mobile notifications

---

## Phase 4

* Incoming call mirror
* Answer / hangup

---

## Done when:

✅ Copy on phone → Linux updates
✅ Copy on Linux → phone updates
✅ Linux notifications appear on phone
✅ Incoming call visible on Linux
✅ Answer/hangup from Linux works

That’s MVP.

---

# Explicit Non-Goals (v0)

❌ Internet sync
❌ Multi-device mesh
❌ File transfer
❌ Images
❌ Encryption beyond shared secret
❌ User accounts

---

# Philosophy

* LAN first
* Ship fast
* Minimal crypto
* CLI driven
* Power-user friendly

Inspired by:

* KDE Connect
* LocalSend
* Syncthing

But lighter and hacker-centric.

---

If you’re ready, next logical step is:

👉 define exact JSON schemas
👉 Go server skeleton
👉 clipboard implementation (Wayland + X11)
👉 Android permissions for call control

Tell me which you want first.
