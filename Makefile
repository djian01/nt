# Define the binary name and Docker image
BINARY_NAME = yourtool
DOCKER_IMAGE = yourtool-image
OUTPUT_DIR = ./output

# Output directories for different platforms
OUTPUT_DIR_LINUX = $(OUTPUT_DIR)/linux
OUTPUT_DIR_WINDOWS = $(OUTPUT_DIR)/windows
# OUTPUT_DIR_MACOS = $(OUTPUT_DIR)/macos

# Check if Docker is installed
.PHONY: check-docker
check-docker:
	@command -v docker >/dev/null 2>&1 || { echo >&2 "Docker is not installed. Please install Docker and try again."; exit 1; }

# Build for Linux
.PHONY: build-linux
build-linux: check-docker
	@echo "Building for Linux inside Docker..."
	mkdir -p $(OUTPUT_DIR_LINUX)
	docker build --rm -t $(DOCKER_IMAGE) .
	docker run --rm -v $(PWD)/$(OUTPUT_DIR_LINUX):/output $(DOCKER_IMAGE) \
		/bin/sh -c "GOOS=linux GOARCH=amd64 go build -o /output/$(BINARY_NAME)"
	@echo "Linux binary built: $(OUTPUT_DIR_LINUX)/$(BINARY_NAME)"

# Build for Windows
.PHONY: build-windows
build-windows: check-docker
	@echo "Building for Windows inside Docker..."
	mkdir -p $(OUTPUT_DIR_WINDOWS)
	docker build --rm -t $(DOCKER_IMAGE) .
	docker run --rm -v $(PWD)/$(OUTPUT_DIR_WINDOWS):/output $(DOCKER_IMAGE) \
		/bin/sh -c "GOOS=windows GOARCH=amd64 go build -o /output/$(BINARY_NAME).exe"
	@echo "Windows binary built: $(OUTPUT_DIR_WINDOWS)/$(BINARY_NAME).exe"

# Build for macOS
# .PHONY: build-macos
# build-macos: check-docker
# 	@echo "Building for macOS inside Docker..."
# 	mkdir -p $(OUTPUT_DIR_MACOS)
# 	docker build --rm -t $(DOCKER_IMAGE) .
# 	docker run --rm -v $(PWD)/$(OUTPUT_DIR_MACOS):/output $(DOCKER_IMAGE) \
# 		/bin/sh -c "GOOS=darwin GOARCH=amd64 go build -o /output/$(BINARY_NAME)"
# 	@echo "macOS binary built: $(OUTPUT_DIR_MACOS)/$(BINARY_NAME)"

# Build all platforms
.PHONY: build-all
build-all: build-linux build-windows build-macos
	@echo "Built binaries for Linux and Windows."

# Clean up the build artifacts and Docker images
.PHONY: clean
clean:
	@echo "Cleaning up build artifacts..."
	rm -rf $(OUTPUT_DIR)
	docker rmi $(DOCKER_IMAGE) || true
