# Go Programming Basics

## Introduction

Go (or Golang) is an open-source, statically typed language designed at Google. It's great for building fast, concurrent services and command-line tools.

## Key Features

- **Static typing** with type inference
- **Garbage collection** for automatic memory management
- **Built-in concurrency** with goroutines and channels
- **Fast compilation** and execution
- **Cross-platform** support

## Hello World

```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
}
```

## Variables

Go has several ways to declare variables:

```go
var name string = "John"
age := 25
const pi = 3.14159
```

## Functions

Functions are first-class citizens in Go:

```go
func add(a, b int) int {
    return a + b
}

func swap(a, b string) (string, string) {
    return b, a
}
```

## Structs

```go
type Person struct {
    Name string
    Age  int
}

func (p Person) Greet() string {
    return fmt.Sprintf("Hello, I'm %s and I'm %d years old", p.Name, p.Age)
}
```
