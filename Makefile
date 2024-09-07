# Variables
GO_BIN=go
APP_NAME=redis-migrator
SRC_DIR=./src

# Default target
all: build

# Build the project
build:
	$(GO_BIN) build -o $(APP_NAME) $(SRC_DIR)/main.go

# Run the project
run:
	$(GO_BIN) run $(SRC_DIR)/main.go

# Clean build artifacts
clean:
	rm -f $(APP_NAME)

.PHONY: all build run clean
