BINARY := adam-younes

# Default target
all: build

# Build the binary with version info (if using git tags)
build:
	@echo "Building $(BINARY)..."
	go build -o $(BINARY) main.go

# Run the binary locally
run: build
	@nohup air & disown
	@echo "Spawning air live build process"
	@echo "Starting tailwindcss watch"
	@npx @tailwindcss/cli -i ./static/css/style.css -o ./static/css/output.css --watch

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
