# Chapter 2: Core Concepts

## Variables & Types

Declare variables with `var x int` or use the shorthand `x := 42`. Go has built-in types like `int`, `string`, `bool`, plus slices, maps, structs, and more.

### Basic Types

```go
var (
    name   string = "Alice"
    age    int    = 30
    height float64 = 1.75
    isStudent bool = true
)
```

### Type Inference

```go
// Go can infer types
message := "Hello"        // string
count := 42              // int
pi := 3.14159           // float64
active := true           // bool
```

## Functions

Functions are first-class. Declare with `func add(a, b int) int { return a + b }`. Multiple return values are supported: `func swap(a, b string) (string, string)`.

### Function Examples

```go
// Simple function
func greet(name string) string {
    return "Hello, " + name + "!"
}

// Multiple return values
func divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, errors.New("division by zero")
    }
    return a / b, nil
}

// Named return values
func getCoordinates() (x, y int) {
    x = 10
    y = 20
    return // naked return
}
```

## Slices and Maps

### Slices

```go
// Create slices
numbers := []int{1, 2, 3, 4, 5}
names := make([]string, 0, 10)

// Append to slices
names = append(names, "Alice", "Bob")

// Slice operations
firstTwo := numbers[:2]
lastThree := numbers[2:]
```

### Maps

```go
// Create maps
ages := map[string]int{
    "Alice": 25,
    "Bob":   30,
}

// Access and modify
ages["Charlie"] = 35
delete(ages, "Bob")

// Check if key exists
if age, exists := ages["Alice"]; exists {
    fmt.Printf("Alice is %d years old\n", age)
}
```
