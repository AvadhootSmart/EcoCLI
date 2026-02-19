# Cobra Cheatsheet (for `eco`)

> Goal: minimal Cobra setup without generators, focused on real systems work.

You’ll use Cobra as a **command router**, nothing more.

---

# 0. Install Cobra

```bash
go get github.com/spf13/cobra@latest
```

That’s it.

No generators.

---

# 1. Minimal Project Layout

Recommended:

```
eco/
  go.mod
  main.go
  cmd/
    root.go
    init.go
    daemon.go
    file.go
```

Later you’ll add:

```
clipboard/
ws/
input/
android/
```

But start small.

---

# 2. main.go

`main.go` should be tiny:

```go
package main

import "eco/cmd"

func main() {
	cmd.Execute()
}
```

---

# 3. Root Command

`cmd/root.go`

```go
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "eco",
	Short: "eco - Linux ↔ Android ecosystem CLI",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
```

This defines:

```
eco
```

Nothing else yet.

---

# 4. Add a Command (`eco init`)

Create `cmd/init.go`

```go
package cmd

import "github.com/spf13/cobra"

func init() {
	rootCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize eco (generate secret + config)",
	Run: func(cmd *cobra.Command, args []string) {
		println("eco init")
	},
}
```

Run:

```bash
go run . init
```

You now have:

```
eco init
```

---

# 5. Nested Commands (`eco daemon start`)

`cmd/daemon.go`

```go
package cmd

import "github.com/spf13/cobra"

func init() {
	rootCmd.AddCommand(daemonCmd)
	daemonCmd.AddCommand(daemonStartCmd)
}

var daemonCmd = &cobra.Command{
	Use:   "daemon",
	Short: "Daemon control",
}

var daemonStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start eco daemon",
	Run: func(cmd *cobra.Command, args []string) {
		println("daemon start")
	},
}
```

Now:

```
eco daemon start
```

---

# 6. Flags

## Local flag

```go
daemonStartCmd.Flags().BoolP("debug", "d", false, "enable debug")
```

Access inside Run:

```go
debug, _ := cmd.Flags().GetBool("debug")
```

Usage:

```
eco daemon start --debug
```

---

## Persistent flag (shared by children)

```go
daemonCmd.PersistentFlags().String("config", "", "config file")
```

Available to:

```
eco daemon start --config x
```

---

# 7. Arguments

```go
var fileSendCmd = &cobra.Command{
	Use:  "send [path]",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		path := args[0]
	},
}
```

Validators:

```go
cobra.MinimumNArgs(1)
cobra.MaximumNArgs(2)
cobra.ExactArgs(1)
```

---

# 8. Required Flags

```go
cmd.MarkFlagRequired("device")
```

---

# 9. Completion (free)

You get shell completion automatically:

```bash
eco completion bash
eco completion zsh
```

Huge win for zero effort.

---

# 10. Help (free)

```
eco --help
eco daemon --help
```

Automatically generated.

---

# 11. Typical `eco` Command Tree

You’re aiming for:

```
eco init
eco daemon start
eco daemon stop
eco status

eco file send
eco file receive

eco mobile open

eco devices list
```

Structure maps directly to Cobra.

---

# 12. Where Your Real Logic Goes

IMPORTANT:

Don’t write logic inside Cobra commands.

Instead:

```go
Run: func(cmd *cobra.Command, args []string) {
	daemon.Start()
}
```

Actual implementation:

```
daemon/start.go
```

Cobra should only route.

---

# 13. Signals + Daemon Mode

Inside your daemon:

```go
ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
defer stop()

go startServer()

<-ctx.Done()
```

Cobra doesn’t manage lifecycle — you do.

---

# 14. Config Loading

Later:

```go
eco init -> writes ~/.config/eco/config.json
eco daemon start -> reads it
```

Don’t use Viper unless you really need it.

Plain JSON is fine.

---

# 15. Common Mistakes (avoid)

❌ cobra generators
❌ viper on day 1
❌ putting logic in Run blocks
❌ global variables everywhere
❌ giant root.go

---

# Recommended Pattern

Cobra only handles:

* args
* flags
* routing

Everything else lives elsewhere.

Think:

Cobra = HTTP router
Your code = services

---

# Mental Model

Cobra is just:

```
switch os.Args[1] {
  case "init":
  case "daemon":
}
```

But done properly.

---

# TLDR

Minimal Cobra:

1. root.go
2. add commands manually
3. route to real code
4. no generators
5. no Viper

You’ll be productive in ~30 minutes.

---

If you want next, I can give:

👉 starter repo skeleton
👉 eco command tree diagram
👉 daemon bootstrap template

Just tell me 👍
