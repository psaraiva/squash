# 🎮 Squash Game

[![project](https://img.shields.io/badge/github-psaraiva%2Fsquash-blue)](https://img.shields.io/badge/github-psaraiva%2F-squash-blue)
[![License](https://img.shields.io/badge/license-MIT-%233DA639.svg)](https://opensource.org/licenses/MIT)

[![Go Report Card](https://goreportcard.com/badge/github.com/psaraiva/squash)](https://goreportcard.com/report/github.com/psaraiva/squash)
![Codecov](https://img.shields.io/codecov/c/github/psaraiva/squash)
[![Language: Português](https://img.shields.io/badge/Language-Portugu%C3%AAs-green?style=flat-square)](./README_pt_br.md)

A classic Squash game developed in **Go** (TinyGO) and compiled to **WebAssembly (WASM)**, running directly in the browser without any plugins.

---

## 📑 Quick Navigation

**👾 For Players:**
- [🕹️ How to Play](#️-how-to-play)
- [🎛️ Game Settings](#️-url-configuration)

**👨‍💻 For Developers:**
- [⚙️ Installation and Execution](#️-installation-and-execution)
- [🧪 Tests](#-tests)

**🏗️ For Software Engineers:**
- [📊 Technical Data](#️-technical-data) - Architecture, Design Patterns and Stack

---

## 📋 About the Game

Squash is an arcade game where you control a paddle and must bounce the ball against the wall, preventing it from escaping through the left side. With each successful hit, you earn points and the game becomes progressively more challenging.

### 🎬 Game Preview

![Squash Gameplay](asset/squash-gameplay.png)

## 🕹️ How to Play

### Objective
- Bounce the ball with the paddle
- Prevent the ball from escaping through the left side
- Each hit earns you points (**10 points**)
- Every **100 points**, you advance to the next level and the ball gets faster
- You start with **3 lives** (default)

### Controls

> **Note**: The game requires a **mouse** or **stylus** (pen) for tablets. Touch controls are not supported.

- **Move paddle**: Move the mouse (or stylus) vertically
- **Start game**: Left click
- **Pause**: Right click
- **Restart**: Left click (on Game Over screen)

## 🎛️ URL Configuration

You can customize the game through query parameters:

```
http://localhost:8080?debug=true&lives=5&level=10&boost=0.8&ballsize=0.7&fps=60
```

### Available parameters:

| Parameter  | Type      | Range       | Description                              |
|------------|-----------|-------------|------------------------------------------|
| `debug`    | boolean   | true/false  | Enables debug mode with information      |
| `lives`    | int       | 1 - 99      | Initial number of lives                  |
| `level`    | int       | 0 - 50      | Starting game level                      |
| `boost`    | float     | 0.0 - 1.0   | Speed increment per level                |
| `ballsize` | float     | 0.0 - 1.0   | Ball size scale                          |
| `fps`      | int       | 30 or 60    | Frames per second (update rate)          |

---

## ⚙️ Installation and Execution

<details>
<summary><b>📦 Prerequisites and Installation</b> (click to expand)</summary>

### Prerequisites

**Option 1: Local Execution**
- Go 1.23+
- TinyGo (to compile to WASM)

**Option 2: Docker Execution** 🐳
- Docker installed

### Installation

```bash
# Clone the repository
git clone https://github.com/psaraiva/squash.git
cd squash
```

### Build and Execution

#### 🐳 **Option 1: Docker (Recommended)**

Simpler! No need to install Go or TinyGo.

```bash
# Build and run in a single command
make docker-deploy
```

Access: `http://localhost:8080`

**Available Docker commands:**

```bash
make docker-build    # Build Docker image
make docker-run      # Run container
make docker-stop     # Stop and remove container
make docker-clean    # Remove container and image
```

#### 💻 **Option 2: Local Execution**

Requires Go 1.23+ and TinyGo installed.

```bash
# Install dependencies
go mod download

# Build and start local server
make web-deploy-local

# Or run commands separately:
make web-build        # Compile to WASM
make web-serve-start  # Start HTTP server
```

Access: `http://localhost:8080`

**Cleanup:**

```bash
make web-clean  # Remove compiled files (local)
```

</details>

<details>
<summary><b>🧪 Tests</b> (click to expand)</summary>

```bash
# Run all tests with coverage
make go-test-all

# Run only unit tests
make go-test

# Run only WASM tests
make go-test-wasm

# Generate interface mocks
make go-mock
```

**Coverage:** 100% of statements tested

</details>

---

## 🏗️ Technical Data

> **For software engineers**: This project demonstrates **Clean Architecture** and **Hexagonal Architecture (Ports & Adapters)** in Go with WebAssembly, 100% testable and extensible.

<details open>
<summary>🎯 <b>Technical Features</b> (summary)</summary>

- 🌐 Runs in the browser via WebAssembly
- 🎮 Control via **mouse** or **stylus** (does not support touch)
- 🎚️ Progressive level system with increasing difficulty
- 🎨 Clean and responsive interface  
- 🐛 Debug mode for developers
- ⚙️ Customizable settings via query string
- ✅ **100% test coverage**

</details>

<details>
<summary>🚀 <b>Technology Stack</b> (click to expand)</summary>

### Core
- **Go 1.23** - Main programming language
- **TinyGo** - Optimized compiler for WebAssembly
- **WebAssembly (WASM)** - Technology to run Go code in the browser
- **JavaScript** - Integration with browser APIs via `syscall/js`

### Architecture
- **Clean Architecture** - Clear separation of layers (domain, ports, adapters)
- **Hexagonal Architecture** - Ports & Adapters pattern
- **Dependency Injection** - Interfaces for decoupling
- **Strategy Pattern** - Mouse input strategy

### Testing
- **Go Testing** - Native testing framework
- **Custom Mocks** - Own implementation without external dependencies
- **Table-Driven Tests** - Go-recommended testing pattern
- **TinyGo Test** - WASM target compatible tests

### Development
- **Make** - Build and deployment automation
- **Docker** - Containerization with multi-stage build
- **Go Modules** - Dependency management

</details>

<details>
<summary>🏗️ <b>Project Structure</b> (click to expand)</summary>

```
squash/
├── cmd/                  # Entry points (delivery interfaces)
│   └── wasm/             # WebAssembly implementation
│       ├── main.go       # Wire-up and initialization
│       └── index.html    # HTML interface
│
├── internal/             # Domain core (business logic)
│   ├── app/              # Game engine and business rules
│   │   ├── config.go     # Configuration and default values
│   │   ├── engine.go     # Physics and game mechanics
│   │   └── game.go       # State and game entities
│   │
│   └── ports/            # Contracts/Interfaces
│       ├── config.go     # ConfigProvider interface
│       ├── renderer.go   # Renderer interface
│       └── mocks/        # Generated mocks
│
├── pkg/                  # Reusable code (infrastructure)
│   └── adapters/         # Port implementations
│       ├── input/        # Input adapters
│       │   ├── wasm/     # WASM config loader
│       │   └── web/      # UI and rendering
│       └── output/       # Output adapters  
│           └── web/      # Canvas renderer
│
└── bin/                  # Compiled artifacts
    └── web/              # WASM assets
```

</details>

<details>
<summary>📊 <b>Architecture and Design Patterns</b> (click to expand)</summary>

This project was developed following the principles of **Clean Architecture** and **Hexagonal Architecture (Ports & Adapters)**, making the code highly testable, maintainable, and extensible for different platforms.

### 🎯 Architecture Layers

#### 1. **Core Domain** (`internal/app/`)
- **Responsibility**: Pure business logic, game rules, physics
- **Files**: `engine.go` (physics and mechanics), `game.go` (state), `config.go`
- **Independent**: Doesn't know infrastructure details (Web, CLI, etc)
- **Testable**: 100% testable without external dependencies

#### 2. **Ports** (`internal/ports/`)
- **Responsibility**: Contracts/interfaces that the domain expects
- **Interfaces**: `ConfigProvider`, `Renderer`
- **Dependency Inversion**: Domain defines, adapters implement

#### 3. **Adapters** (`pkg/adapters/`)
- **Responsibility**: Concrete implementations of ports
- **Input Adapters**:
  - `input/wasm/config_loader.go` - Reads config from query string
  - `input/wasm/handler.go` - Captures mouse events
  - `input/web/ui.go` - UI rendering logic
- **Output Adapters**:
  - `output/web/canvas.go` - Canvas 2D Renderer
  - `output/web/jscontext.go` - Wrapper for syscall/js
- **Interchangeable**: Easy to swap implementations without affecting the core

#### 4. **Entry Points** (`cmd/`)
- **Responsibility**: Composition (wire-up) and initialization
- **Minimal logic**: Only instantiates and connects components

### 🔌 Extensibility: New Implementations

The architecture allows easily creating new versions of the game for different platforms:

#### 🖥️ **Example: New Version**
```
cmd/new/
  └── main.go                    # entry point

pkg/adapters/
  ├── input/new/
  │   ├── config_loader.go       # Reads config from flags/env
  │   └── keyboard.go            # Captures input
  └── output/new/
      └── renderer.go            # rendering
```

**Usage example:**
```bash
go run cmd/new/main.go
```

**The core (`internal/app`) remains 100% unchanged!**

### 🧩 SOLID Principles Applied

| Principle | Application in Project |
|-----------|------------------------|
| **S**RP | Each package has a single responsibility |
| **O**CP | Extensible via new adapters without modifying the core |
| **L**SP | `Renderer`, `ConfigProvider` interfaces are substitutable |
| **I**SP | Small and focused interfaces |
| **D**IP | `internal/app` depends on abstractions (`ports`), not implementations |

### 🎨 Design Patterns Used

- **Hexagonal/Ports & Adapters**: Isolated core, adapters connect infrastructure
- **Dependency Injection**: Components receive dependencies via constructor
- **Strategy Pattern**: Mouse input strategy
- **Factory Pattern**: `NewSquash()`, `NewRenderer()`, `NewConfigLoader()`
- **Template Method**: `Renderer.Render()` with specific implementations

### ✨ Architecture Benefits

✅ **Testability**: Core testable without complex mocks (100% coverage)  
✅ **Maintainability**: Changes isolated to specific layers  
✅ **Reusability**: Game logic reusable on any platform  
✅ **Evolution**: Easy to add features without breaking existing code  
✅ **Independence**: Core doesn't depend on external frameworks  
✅ **Portability**: Same core for Web, CLI, Mobile, Desktop

</details>

---

##  License

This project is open source and available under the MIT License.

## 👨‍💻 Author

Developed by [@psaraiva](https://github.com/psaraiva)

---

**Have fun playing! 🎉**
