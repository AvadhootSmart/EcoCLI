# 🧰 Go CLI Application Cheat Sheet

## 1. Create a New CLI Project

```bash
mkdir mycli
cd mycli
go mod init github.com/you/mycli
touch main.go
```

Basic skeleton:

```go
package main

import "fmt"

func main() {
	fmt.Println("Hello CLI!")
}
```

Run it:

```bash
go run .
```

Build binary:

```bash
go build -o mycli
./mycli
```

---

## 2. Access Command-Line Arguments (`os.Args`)

```go
import (
	"fmt"
	"os"
)

func main() {
	fmt.Println(os.Args)
}
```

Example:

```bash
./mycli foo bar
```

Output:

```
[./mycli foo bar]
```

Common pattern:

```go
args := os.Args[1:] // skip program name
```

---

## 3. Parse Flags (Standard Library)

Use `flag` for simple CLIs.

### Basic Flags

```go
import (
	"flag"
	"fmt"
)

func main() {
	name := flag.String("name", "world", "name to greet")
	age := flag.Int("age", 0, "your age")

	flag.Parse()

	fmt.Printf("Hello %s (%d)\n", *name, *age)
}
```

Run:

```bash
./mycli --name Alice --age 30
```

Short flags:

```go
flag.StringVar(&name, "n", "world", "short name flag")
```

---

## 4. Remaining Args After Flags

```go
flag.Parse()
args := flag.Args()

fmt.Println("Extra args:", args)
```

---

## 5. Simple Subcommands (Manual)

Go doesn’t have built-in subcommands, but you can switch on `os.Args`.

```go
func main() {
	if len(os.Args) < 2 {
		fmt.Println("expected 'add' or 'remove'")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "add":
		fmt.Println("adding...")
	case "remove":
		fmt.Println("removing...")
	default:
		fmt.Println("unknown command")
	}
}
```

With flags per subcommand:

```go
addCmd := flag.NewFlagSet("add", flag.ExitOnError)
name := addCmd.String("name", "", "item name")

switch os.Args[1] {
case "add":
	addCmd.Parse(os.Args[2:])
	fmt.Println("adding", *name)
}
```

---

## 6. Print Output (stdout / stderr)

```go
fmt.Println("normal output")

fmt.Fprintln(os.Stderr, "error output")
```

---

## 7. Exit Codes

```go
os.Exit(0) // success
os.Exit(1) // generic error
```

Pattern:

```go
if err != nil {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
```

---

## 8. Read From stdin

Useful for pipes:

```go
import "io"

data, _ := io.ReadAll(os.Stdin)
fmt.Println(string(data))
```

Example:

```bash
echo "hello" | ./mycli
```

---

## 9. Read / Write Files

Read:

```go
b, err := os.ReadFile("file.txt")
if err != nil {
	panic(err)
}
fmt.Println(string(b))
```

Write:

```go
err := os.WriteFile("out.txt", []byte("hello"), 0644)
```

---

## 10. Prompt User for Input

```go
reader := bufio.NewReader(os.Stdin)
fmt.Print("Enter name: ")
text, _ := reader.ReadString('\n')
```

---

## 11. Organizing a Real CLI Project

Typical layout:

```
mycli/
  main.go
  cmd/
    add.go
    remove.go
  internal/
```

`main.go`:

```go
func main() {
	cmd.Execute()
}
```

Put command logic in `cmd/`.

---

## 12. Helpful Patterns

### Central error handler

```go
func fatal(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
```

---

### Usage Message

```go
flag.Usage = func() {
	fmt.Println("Usage: mycli [options]")
	flag.PrintDefaults()
}
```

---

## 13. Building for Distribution

Cross compile:

```bash
GOOS=linux GOARCH=amd64 go build
GOOS=darwin GOARCH=arm64 go build
GOOS=windows GOARCH=amd64 go build
```

---

## 14. Popular CLI Libraries (Optional)

Once your CLI grows, these help a lot:

* **cobra** – full-featured subcommands, help, completions
* **urfave/cli** – simple and expressive
* **spf13/pflag** – POSIX-style flags

If you’re just learning: start with **standard library first**.

---

# 🚀 Minimal Example: Real CLI

```go
package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	msg := flag.String("msg", "hello", "message")
	flag.Parse()

	if len(flag.Args()) == 0 {
		fmt.Println(*msg)
		return
	}

	for _, a := range flag.Args() {
		fmt.Println(*msg, a)
	}

	os.Exit(0)
}
```

Run:

```bash
./mycli --msg hi Alice Bob
```

---

If you’d like next, I can also give you:

✅ A Cobra-based starter template
✅ A real-world CLI example (todo app, file tool, etc.)
✅ Testing CLIs in Go
✅ Argument completion (bash/zsh)
✅ Packaging with Homebrew

Just tell me 👍
