# Chapter 1: Getting Started

## Introduction

Go (or Golang) is an open-source, statically typed language designed at Google. It's great for building fast, concurrent services and command-line tools.

## Installation

Download the binary from [golang.org/dl](https://golang.org/dl) and follow the installer for your OS. After installing, verify with `go version`.

## Your First Program

Create a file called `hello.go`:

```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
}
```

Run it with:

```bash
go run hello.go
```

## Go Modules

Initialize a new module:

```bash
go mod init myproject
```

Add dependencies:

```bash
go get github.com/gorilla/mux
```

## Project Structure

A typical Go project structure:

```
myproject/
├── cmd/
│   └── main.go
├── internal/
│   └── handlers/
├── pkg/
├── go.mod
└── README.md
```
