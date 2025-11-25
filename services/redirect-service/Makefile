# Variables
BINARY_NAME := redirect-service
TARGET := x86_64-unknown-linux-musl
# Default App Name 
APP_NAME ?= us-func-p6ndmuotrzo5a

.PHONY: all build clean deploy local

# Default command: Just build
all: build

# 1. Cross-compile and move binary to root
build:
	@echo "--- Building Static Binary (Musl) ---"
	cargo build --release --target $(TARGET)
	@echo "--- Moving Binary to Root ---"
	cp target/$(TARGET)/release/$(BINARY_NAME) .

# 2. Run locally (Builds first, then starts func)
local: build
	@echo "--- Starting Local Function ---"
	func start

# 3. Deploy to Azure (Builds first, then publishes)
deploy: build
	@echo "--- Deploying to Azure: $(APP_NAME) ---"
	func azure functionapp publish $(APP_NAME)

# Cleanup
clean:
	cargo clean
	rm -f $(BINARY_NAME)
