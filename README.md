# Lingo - TypeScript-like Meta-Language for Go

Lingo is a meta-language compiler that brings TypeScript-like static typing and null-safety features to Go, while maintaining Go's familiar syntax and idioms. [Dingo](https://github.com/MadAppGang/dingo) has already proved this can happen,So did [Borgo](https://github.com/borgo-lang/borgo)

## Features

- **Static Type Checking**: Comprehensive type inference and checking before compilation
- **Null Safety**: Built-in null-safety with optional types and null coalescing operators
- **Go-like Syntax**: Looks and feels like normal Go code - no new syntax to learn
- **Full Go Compatibility**: Compiles to idiomatic Go that works with the standard toolchain
- **Advanced Type System**: Generics support, union types, and type constraints
- **Error Prevention**: Catch type errors at compile-time, not runtime

## Installation

```bash
git clone https://github.com/MistyPigeon/lingo.git
cd lingo
make deps
make build
```
Binaries will be created in ./bin/

# Quick Start

Create a Lingo file (. lingo)
```bash
package main

func main() {
    var name: ? string = null
    var greeting: string = name ?: "World"
    
    fmt. Println("Hello, " + greeting)
}
```
Compile to Go
```bash
./bin/lingo -file hello.lingo -out hello.go
go run hello.go
```
Syntax Guide
Type Annotations
Variables require type annotations for static checking:

```bash
var x: int = 42
var name: string = "Alice"
var items: []string = []string{"a", "b"}
var config: map[string]int = map[string]int{"count": 5}
Nullable Types
Use ? to mark types as nullable:
```
```bash
var email: ? string = null
var count: ?int = 0
Null Coalescing
```
Use ? : to provide default values for nullable types:

```bash
var result: string = email ?: "no-email@example.com"
```
Functions with Type Safety
```bash
func add(a: int, b: int) int {
    return a + b
}

func greet(name: ? string) string {
    return "Hello, " + (name ?: "Guest")
}
```
Structs
```bash
type User struct {
    id: int
    name: string
    email: ?string
}
```
Methods
```bash
func (u: *User) GetEmail() ? string {
    return u.email
}
```
Generics (Basic)
```bash
func first(items: []interface{}) interface{} {
    return items[0]
}
```
CLI Commands
Compile Lingo to Go
```bash
./bin/lingo -file input.lingo -out output.go [-check] [-v]
```
Options:

-file: Input . lingo file (required)
-out: Output . go file (default: same name as input with .go extension)
-check: Only perform type checking without generating code
-v: Verbose output (show tokens and AST)
Lexical Analysis
```bash
./bin/lingoctl -cmd lex -file input.lingo
```
Parse to AST
```bash
./bin/lingoctl -cmd parse -file input.lingo
```
Examples
See the examples/ directory for more detailed examples:

basic.lingo - Basic types and functions
nullsafe.lingo - Null safety features
generics.lingo - Generic-like patterns

# How It Works

Lexical Analysis: Source code is tokenized
Parsing: Tokens are parsed into an Abstract Syntax Tree (AST)
Type Checking: The AST is analyzed for type errors
Code Generation: Type-safe AST is compiled to Go code
Execution: Generated Go code is compiled and run with go build or go run
Type System
Basic Types
int, int8, int16, int32, int64
uint, uint8, uint16, uint32, uint64
float32, float64
string
bool
byte, rune
Composite Types
Arrays: []T
Maps: map[K]V
Pointers: *T
Slices: []T
Nullable Types
Any type can be marked nullable with ?:
```bash
? string
?int
?[]string
?*User
```
Null Safety
Nullable Variables
```bash
var user: ? User = nil

// Error at compile time - cannot use potentially nil value
// name := user.name

// Correct - using null coalescing
name := user ?: defaultUser
Null Checks

if user != null {
    // user is guaranteed non-nil here
    name := user.name
}
```
Null Coalescing Operator
```bash
var email: string = userEmail ?: "noemail@example.com"
```
Development
Build
bash
make build
Run Tests
```bash
make test
```
Format Code
```bash
make fmt
```
Run Full Pipeline
```bash
make all
```
# Contributing
Contributions are welcome! Please:

Fork the repository
Create a feature branch
Make your changes
Add tests
Run make fmt and make test
Submit a pull request
License
MIT License - see LICENSE file for details

Roadmap
 Generic types support
 Interface embedding
 Error handling patterns (Result types)
 Pattern matching
 Async/await patterns
 LSP (Language Server Protocol) support
 IDE extensions (VS Code)
 Standard library bindings
 Package manager integration
Performance
Lingo adds no runtime overhead. The generated Go code is optimized and performs identically to hand-written Go.

FAQ
**Q: Do I have to rewrite my Go code? ** A: No! You can mix Lingo and Go files. Lingo compiles to standard Go that works with any Go code.

Q: What about existing Go packages? A: The generated Go code can import and use any Go package directly.

Q: Is this production-ready? A: Currently in beta. The core features are stable, but use with caution in production.

Q: How does it compare to TypeScript? A: Like TypeScript for JavaScript, Lingo adds types and safety to Go while preserving its essence. However, Go already has static typing, so Lingo focuses on null-safety and advanced type patterns.

Contact
For issues, questions, or suggestions: GitHub Issues

Made with ❤️ for Go developers who want TypeScript-like safety
