# Phase 2 – Extended Device Capabilities

This phase builds on the MVP (clipboard + notifications + calls) and expands the ecosystem into a true **Android ↔ Linux peripheral bridge**.

Target audience: power users / personal experimentation / open source hackers.

No Play Store constraints assumed.
Sideloading, Accessibility Service, foreground services, and adb are acceptable.

---

## Goals

Turn the Android device into:

- A wireless storage device
- A camera + microphone peripheral
- An input device (keyboard / mouse / touchpad)
- A remotely controllable system target

All over the existing WebSocket control plane.

No cloud.
LAN-first.
CLI-driven on Linux.

---

## Architecture Recap

Linux remains the authority node.

```

Linux Daemon (Go)
↕ WebSocket (control plane)
Android Agent (Kotlin)

````

Phase 2 introduces **data planes** in addition to WebSocket:

- HTTP (file transfer)
- WebRTC (camera + mic)

WebSocket is used only for:

- signaling
- commands
- metadata
- auth

---

# Feature Set

---

## 1. File Transfer

### Description

Bidirectional file transfer between Linux and Android.

Supports:

- single files
- folders (recursive)
- progress reporting

---

### Implementation

### Control Plane (WebSocket)

#### Linux → Mobile

```json
{
  "type": "file.offer",
  "name": "photo.jpg",
  "size": 213123,
  "id": "file-1"
}
````

Mobile replies:

```json
{
  "type": "file.accept",
  "id": "file-1"
}
```

---

### Data Plane (HTTP)

Linux exposes:

```
POST /upload
GET /download/:id
```

Android uploads / downloads using HTTP streams.

Chunked transfer.

---

### Linux

* Go HTTP server
* resumable uploads (optional)
* save location configurable

---

### Android

* Storage Access Framework picker
* Foreground service during transfer
* persistent notification

---

### Status

Phase 2.1 – Required

---

## 2. Mobile → PC Input Control

### Description

Use Android as:

* keyboard
* mouse / touchpad
* scroll wheel

---

### Flow

Android sends input events:

```json
{
  "type": "input.key",
  "key": "A"
}
```

```json
{
  "type": "input.mouse",
  "dx": 12,
  "dy": -4
}
```

---

Linux injects using:

* uinput
* ydotool
* evdev

---

### Linux

Creates virtual input device.

Maps received events to:

* key presses
* mouse movement
* scroll

---

### Android

Custom UI:

* touchpad area
* keyboard field
* gesture support

---

### Status

Phase 2.1 – Required

---

## 3. Open Apps on Mobile (Intent Control)

### Description

Launch Android apps from Linux CLI.

---

### Example

Linux:

```bash
eco mobile open com.spotify.music
```

Sends:

```json
{
  "type": "mobile.open_app",
  "package": "com.spotify.music"
}
```

---

### Android

Uses:

```kotlin
startActivity(packageManager.getLaunchIntentForPackage(...))
```

---

Limitations:

* app must be launchable
* foreground service required on modern Android
* OEM restrictions may apply

---

### Status

Phase 2.2 – Optional

---

## 4. External Camera Feed (Android → Linux)

### Description

Use phone camera as Linux webcam.

---

### Implementation

Uses WebRTC.

---

### Flow

1. Linux sends SDP offer via WebSocket
2. Android replies SDP answer
3. ICE candidates exchanged over WS
4. Media flows directly over RTP

---

### Android

* CameraX
* WebRTC native SDK
* foreground service

---

### Linux

* GStreamer + webrtcbin OR libwebrtc
* v4l2loopback

Creates virtual webcam:

```
/dev/videoX
```

Now usable by:

* OBS
* Zoom
* browsers

---

### Status

Phase 2.3 – Advanced

---

## 5. External Microphone (Android → Linux)

Same WebRTC connection as camera.

---

Android:

* AudioRecord
* Opus encoding

Linux:

* PipeWire / PulseAudio sink
* RTP → virtual mic

---

Phone becomes wireless microphone.

---

### Status

Phase 2.3 – Advanced

---

# Explicit Non-Goals (Phase 2)

* PC → Android raw input injection (without adb/root)
* Internet relay
* multi-device mesh
* encryption beyond shared secret
* iOS support

---

# Permissions Required (Android)

* FOREGROUND_SERVICE
* CAMERA
* RECORD_AUDIO
* READ_EXTERNAL_STORAGE / SAF
* POST_NOTIFICATIONS
* BIND_ACCESSIBILITY_SERVICE (input control)
* INTERNET

Optional (calls already covered in MVP):

* READ_PHONE_STATE
* CALL_PHONE

---

# Development Order

Recommended:

## Phase 2.1

* File transfer
* Mobile → PC input

## Phase 2.2

* Open apps intent

## Phase 2.3

* Camera stream
* Microphone stream

---

# Philosophy

Android is treated as:

* peripheral
* sensor
* input device

Linux is the control plane.

This is a personal distributed operating environment.

Not a consumer app.

---

End of Phase 2.

```

---

If you want, next we can write:

- `ANDROID_AGENT.md`
- `PROTOCOL.md`
- `WEBRTC.md`
- or a `ROADMAP.md`

Just tell me 👍
```
