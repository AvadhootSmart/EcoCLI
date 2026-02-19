# EcoCLI Android Client

Android client for the EcoCLI ecosystem.

## Building

```bash
cd android
./gradlew assembleDebug
```

APK will be at `app/build/outputs/apk/debug/app-debug.apk`

## Running Tests

```bash
./gradlew test
```

## Project Structure

```
app/src/main/java/dev/eco/
├── core/
│   ├── network/       # WebSocket client
│   └── protocol/      # Message types & handlers
├── feature/
│   ├── clipboard/     # Clipboard sync
│   ├── call/          # Call control
│   ├── connection/    # Connection setup UI
│   ├── notification/  # Linux notifications
│   └── settings/      # App settings
├── service/           # Foreground service
└── ui/                # Theme & navigation
```

## Phase 2 Extension Points

Interfaces are stubbed for future features:

- `feature/file/` - File transfer (Phase 2.1)
- `feature/input/` - Input control (Phase 2.1)
- `feature/media/` - Camera/mic streaming (Phase 2.3)
