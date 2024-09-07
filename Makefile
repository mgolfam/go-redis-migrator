# Define output directory
BINDIR = ./bin
SRC = ./src
BINARY = $(BINDIR)/redis-migrator

# Ensure bin directory exists
$(BINDIR): build
	mkdir -p $(BINDIR)

# Build the binary
build: 
	go build -o $(BINARY) $(SRC)/main.go

# Clean the build
clean:
	rm -rf $(BINDIR)

# Run the binary
run: build
	$(BINARY)

# Install dependencies
deps:
	go mod tidy
	go mod download

.PHONY: build clean run deps
