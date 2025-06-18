BINARY := adam-younes

# Default target
all: build

# Build the binary with version info (if using git tags)
build:
	@echo "Building $(BINARY)..."
	go build -o $(BINARY) main.go

# Run the binary locally
run: build
	@echo "Running $(BINARY)..."
	./$(BINARY)

# Clean up the binary
clean:
	@echo "Cleaning..."
	rm -f $(BINARY)

# Format Go code
fmt:
	go fmt ./...

# Run Go tests (if any)
test:
	go test ./...
