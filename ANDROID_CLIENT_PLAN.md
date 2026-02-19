# Eco Android Client — PLAN.md

## 🎯 Goal

Build the Android client for the Eco ecosystem that integrates with the existing EcoCLI via WebSockets.

The Android app is built with:

- **Expo / React Native (TypeScript)** → UI, networking, state, business logic
- **Minimal Kotlin native layer** → Android system listeners only

**Principle:** Native code must remain thin and dumb. All intelligence lives in TypeScript.

---

# 🧠 Architecture Overview

```
Android System
   ↓
Kotlin Listeners / Services
   ↓ (DeviceEventEmitter)
Expo React Native Layer
   ↓ (WebSocket)
EcoCLI (Linux)
```

---

# 📦 Responsibilities

## Kotlin (native layer)

ONLY responsible for:

- Listening to Android system events
- Emitting events to React Native
- Performing small OS actions when requested

**Kotlin must NOT contain business logic.**

---

## Expo / TypeScript layer

Responsible for:

- WebSocket connection to EcoCLI
- Device pairing & identity
- Sync logic & deduplication
- UI
- Retry & offline handling
- Feature flags
- State management

---

## EcoCLI (already implemented)

Agent must read the CLI codebase to understand:

- WebSocket URL and auth
- Event names
- Payload schemas
- Expected acknowledgements

---

# 🚀 MVP Features

## 1. Clipboard Duplex Sync

### Goal

Real-time clipboard sync between:

- Android ⇄ PC

---

## Android Requirements (Kotlin)

Implement:

- ForegroundService
- ClipboardManager.OnPrimaryClipChangedListener

### Behavior

When clipboard changes on Android:

1. Read text content
2. Emit event to React Native:

```
Event name: EcoClipboardChanged
Payload:
{
  text: string,
  timestamp: number,
  source: "android"
}
```

---

## React Native Responsibilities

On receiving `EcoClipboardChanged`:

1. Deduplicate
2. Send to EcoCLI via WebSocket
3. Handle incoming clipboard updates from PC
4. Write to Android clipboard using Expo Clipboard API

---

## Important Rules

- Native layer MUST NOT send directly to WebSocket
- Native layer MUST NOT dedupe
- Native layer MUST NOT contain sync logic

---

# 🔔 2. Notification Singleton Sync (Mobile → PC)

### Goal

Mirror Android notifications to PC.

One-way sync for MVP.

---

## Android Requirements (Kotlin)

Implement:

- NotificationListenerService

### Behavior

On notification posted:

Extract:

- app package
- title
- text
- timestamp

Emit to React Native:

```
Event name: EcoNotificationPosted
Payload:
{
  packageName: string,
  title: string,
  text: string,
  timestamp: number
}
```

---

## React Native Responsibilities

Upon receiving event:

1. Filter (optional future)
2. Forward to EcoCLI via WebSocket
3. Handle reconnection buffering

---

## Explicit Non-Goals (MVP)

- ❌ Notification actions
- ❌ Notification dismissal sync
- ❌ Notification reply

---

# 📞 3. Call Notifications + PC Action Handler

### Goal

When phone receives a call:

- PC gets notified
- PC can trigger actions (answer / reject)

---

## Android Requirements (Kotlin)

Implement:

### Listener

Use:

- TelephonyManager
- PhoneStateListener or modern callback

Detect:

- RINGING
- OFFHOOK
- IDLE

---

### Emit to React Native

```
Event name: EcoCallStateChanged
Payload:
{
  state: "ringing" | "offhook" | "idle",
  phoneNumber?: string,
  timestamp: number
}
```

---

## Action Bridge (React Native → Kotlin)

React Native must be able to call native methods:

Native module methods required:

- `answerCall()`
- `rejectCall()`

⚠️ Implementation may use:

- TelecomManager
- ACTION_ANSWER
- ACTION_DECLINE

Exact implementation left to agent.

---

## React Native Responsibilities

- Forward call events to EcoCLI
- Listen for CLI commands
- Invoke native call actions

---

# 🧱 Required Native Modules

Agent must implement these minimal modules.

---

## Module: EcoClipboardService

Responsibilities:

- Foreground service
- Clipboard listener
- Emit `EcoClipboardChanged`

---

## Module: EcoNotificationListener

Responsibilities:

- NotificationListenerService
- Emit `EcoNotificationPosted`

---

## Module: EcoCallListener

Responsibilities:

- Phone state listener
- Emit `EcoCallStateChanged`

---

## Module: EcoCallActions

Exposed to React Native:

- answerCall()
- rejectCall()

---

# 📡 Event Bridge

All native → JS communication must use:

- DeviceEventEmitter (React Native)

Naming convention:

- Prefix all events with `Eco`

---

# 🔌 WebSocket Integration (Expo side)

Agent must:

1. Read EcoCLI repo

2. Discover:
   - socket URL
   - auth mechanism
   - event names
   - payload formats

3. Implement compatible client

---

# 🧪 Development Phases

## Phase 1 (must complete)

- WebSocket manual test
- Clipboard listener
- Notification listener
- Call state listener
- Basic UI log screen

---

## Phase 2 (later)

- File transfer duplex
- Overlay bubble UX
- Boot persistence
- OEM battery handling
- Encryption

---

# ⚠️ Critical Constraints

## Native Code Rules

- Keep Kotlin minimal
- No business logic in native
- No WebSocket in native
- No state machines in native
- No retries in native

---

## Performance Rules

- Use foreground service for clipboard
- Avoid polling
- Emit events only on change
- Deduplication belongs in TypeScript

---

## Reliability Rules

React Native layer must handle:

- reconnection
- buffering
- dedupe
- device identity
- feature flags

---

# ✅ Definition of Done (MVP)

Android client can:

- Sync clipboard both ways
- Mirror notifications to PC
- Send call state to PC
- Execute call actions from PC
- Maintain stable WebSocket connection

---

# 🧭 Agent Instructions

When implementing:

1. First read EcoCLI protocol
2. Implement WebSocket client in Expo
3. Add native modules incrementally
4. Verify each listener independently
5. Keep native surface minimal
6. Prefer Expo APIs when possible
7. Do NOT over-engineer Kotlin

---

**End of PLAN.md**
