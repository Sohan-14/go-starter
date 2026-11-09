# Configuration Guide

This configuration package uses zero third-party dependencies, relying entirely on Go's standard library (**os.Getenv**).

Because Go does not natively parse .env files, environment variables must be loaded into the process environment before starting the application (for example, via your **shell**, **Docker --env-file**, or **system configuration**).

Overview of how the **`internal/platform/config`** package works and how to modify it.

---

## File Overview

| File              | Purpose             | Edit When...                                      |
| :---------------- | :------------------ | :------------------------------------------------ |
| **`config.go`**   | Constants & Structs | Adding or changing config fields                  |
| **`load.go`**     | Wiring and reading  | Connecting fields to environment keys             |
| **`validate.go`** | Business rules      | Adding validation checks (ports, values)          |
| **`parse.go`**    | Type parsers        | Creating new parsing rules (e.g., bool, duration) |
| **`reader.go`**   | Core engine         | Don't need                                        |

## How the Configuration Pipeline Works

##### The configuration engine uses a two-phase load strategy:

- **Phase 1** (Reading & Parsing): All environment variables are read using required or optional helpers via the reader. Any missing variables or parsing failures are accumulated in reader.errs.

- **Phase 2** (Validation): If reading succeeds without syntax or missing key errors, business logic rules are evaluated across all configuration structs.

---

## How to Add a New Config (e.g. `Log`)

### 1. `config.go` – Add Key and Struct

```go
const envLogLevel = "LOG_LEVEL"

type Config struct {
    // ... existings
    Log *LOG // Add section
}

type LOG struct {
    Level slog.Level
}
```

### 2. `load.go` – Read the Value

```go
func Load(getenv getEnv) (*Config, error) {
  cfg := &Config{
		// ... existings
    Log: loadLog(r), // Wire it here
	}
}

func loadLog(r *reader) *LOG {
    return &LOG{
        Level: optional(r, envLogLevel, slog.LevelInfo, asLevel),
    }
}
```

### 3. `validate.go` – Validate (Optional)

```go
// wire log validate into Load validate
func (c *Config) validate() error {
  if err := c.Log.validate(); err != nil { // Wire validation
    errs = append(errs, err)
  }
}

func (l *LOG) validate() error {
  // Add custom validation rules if necessary
  return nil
}
```

---

## Modifying Existing Configs

- **Change Env Key:** Edit the `const` in `config.go`.
- **Add Parser:** Add a `func(string) (T, error)` signature inside `parse.go`.

## Call from the main

```go
func main() {
	cfg, err := config.Load(os.Getenv)
}
```
