# Logi

**A language for building languages.**

Logi is a Go-based DSL framework that lets you define custom domain-specific languages declaratively. Write a **macro** to describe your language's grammar, then write **definitions** in that language. Logi parses, validates, compiles to JSON, or executes them at runtime via a built-in Virtual Machine.

## Why Logi?

Every domain has its own vocabulary. Credit scoring, circuit design, access control, chatbots -- each benefits from a tailored language that domain experts can read and write. Logi gives you the tools to create that language without writing a parser from scratch.

- **Declarative macro system** -- define your DSL grammar in a `.lgm` file
- **Structured parsing** -- definitions in `.lg` files are validated against your macros
- **JSON compilation** -- compile definitions to JSON for use in any application
- **Go Virtual Machine** -- load and execute definitions at runtime with custom semantics
- **LSP support** -- editor integration with semantic highlighting and diagnostics
- **AI generation** -- experimental Gemini/DeepSeek integration for generating definitions from macros

## Quick Start

### Installation

Download a prebuilt binary from the [releases page](https://github.com/tislib/logi/releases), or build from source:

```bash
go install github.com/tislib/logi/cmd/logi@latest
```

### 1. Define a macro

Create `credit-rule.lgm`:

```logi
macro creditRule {
    kind Syntax

    syntax {
        creditScore <min int> <max int>
        income <min int> <max int>
        age <min int> <max int>
    }
}
```

### 2. Write definitions

Create `credit-rule.lg`:

```logi
creditRule GoldTier {
    creditScore 700 850
    income 50000 200000
    age 25 65
}

creditRule SilverTier {
    creditScore 600 700
    income 30000 50000
    age 18 65
}
```

### 3. Compile

```bash
logi compile -i credit-rule.lg -m .
```

Output is structured JSON ready for your application. See the [full result](examples/credit-rule/credit-rule.json).

## CLI Usage

```
logi compile -i <file.lg> -m <macro-dir> [-o <output-dir>] [-k normal|plain]
logi serve -m <macro-dir>          # HTTP server on :7051
logi lsp                           # Language Server Protocol server
```

| Flag | Description |
|------|-------------|
| `-i` | Input `.lg` file |
| `-m` | Directory containing `.lgm` macro files (default: `.`) |
| `-o` | Output directory for JSON (default: stdout) |
| `-k` | Parse mode: `normal` (macro-matched) or `plain` (raw) |
| `-v` | Verbose logging |

## Using the Virtual Machine

For runtime execution in Go, use the VM with a custom `Implementer`:

```go
package main

import (
    "fmt"
    "github.com/tislib/logi/pkg/vm"
    logiAst "github.com/tislib/logi/pkg/ast/logi"
)

func main() {
    v := vm.New()

    // Load macro and definition files
    v.LoadMacroFromFile("chatbot.lgm")
    v.LoadLogiFromFile("chatbot.lg")

    // Execute with custom semantics
    v.Execute("MyChatbot", vm.NewImplementerFunc(
        func(vm vm.VirtualMachine, stmt logiAst.Statement) error {
            fmt.Printf("Command: %s, Args: %v\n", stmt.Command, stmt.Parameters)
            return nil
        },
    ))
}
```

See [examples/chat-bot/main.go](examples/chat-bot/main.go) for a complete working example.

## Examples

| Example | Description |
|---------|-------------|
| [credit-rule](examples/credit-rule/) | Credit scoring rules DSL |
| [chat-bot](examples/chat-bot/) | Chatbot with intents and responses |
| [circuit](examples/circuit/) | Arduino-style circuit control with nested scopes |
| [user-access-control](examples/user-access-control/) | RBAC policy definitions |
| [arrow-ops](examples/arrow-ops/) | Operator-based DSL |

## A More Advanced Example: Circuit DSL

Logi supports nested scopes, conditionals, and event handlers:

```logi
circuit simple1 {
    components {
        Led     yellowLed 5
        Led     redLed 6
        Button  button1 17
    }

    actions {
        on(yellowLed)
        on(redLed)

        on_click(button1) {
            if (status(button2) == 'on') {
                on(yellowLed)
            } else {
                off(yellowLed)
            }
        }
    }
}
```

The macro that defines this DSL uses **scopes** for nested blocks. See the [full circuit example](examples/circuit/) and the [syntax reference](docs/syntax-reference.md#4-scope) for details.

## Syntax Overview

Logi has two file types:

| File | Extension | Purpose |
|------|-----------|---------|
| Macro | `.lgm` | Defines your DSL's grammar |
| Definition | `.lg` | Instances written in your DSL |

### Syntax Elements

Macros are built from these element types:

| Element | Syntax | Purpose |
|---------|--------|---------|
| Keyword | `name` | Structural token for matching |
| VariableKeyword | `<name string>` | Captures a typed value |
| ParameterList | `(<a string>, <b int>)` | Comma-separated value list |
| ArgumentList | `(...[<args Type<string>>])` | Property list (name + type) |
| AttributeList | `[required boolean, default string]` | Optional attributes in `[]` |
| Scope | `{ scopeName }` | Nested block of statements |
| TypeReference | `<TypeName>` | Inline a type from `types` |
| Combination | `(<a string> \| <a Name>)` | OR between alternatives |

### Supported Value Types

`string`, `int`, `float`, `bool`, `date`, `time`, `datetime`, `duration`, `money`, `unit`, `Name`, `Type`

For the complete syntax specification, see the [Syntax Reference](docs/syntax-reference.md).

## Project Structure

```
cmd/
  logi/          # CLI binary (compile, serve, lsp)
  logi-lsp/      # Standalone LSP server
pkg/
  ast/           # AST types (common, logi, macro, plain)
  parser/        # Parsers (YACC-based) for .lg and .lgm files
  vm/            # Virtual Machine runtime
  lsp/           # Language Server Protocol implementation
  ai/            # AI-powered definition generator
examples/        # Example DSLs with macros and definitions
docs/            # Language documentation
```

## Building from Source

```bash
git clone https://github.com/tislib/logi.git
cd logi
go build ./cmd/logi
go test ./...
```

## License

See [LICENSE](LICENSE) for details.
