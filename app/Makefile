GOCMD = go
GOBUILD = $(GOCMD) build
GOCLEAN = $(GOCMD) clean
GOTEST = $(GOCMD) test
GOGET = $(GOCMD) get

BINARY_NAME = app

SRC_DIR = ./cmd
BIN_DIR = ./bin

all: build

build:
	@echo "Building $(BINARY_NAME)..."
	@$(GOBUILD) -o $(BIN_DIR)/$(BINARY_NAME) $(SRC_DIR)/main.go
	@echo "Build complete. Binary: $(BIN_DIR)/$(BINARY_NAME)"

# test:
# 	@echo "Running tests..."
# 	@$(GOTEST) ./...
# 	@echo "Tests passed!"

clean:
	@echo "Cleaning up..."
	@$(GOCLEAN)
	@rm -f $(BIN_DIR)/$(BINARY_NAME)
	@echo "Cleaned up."

run:
	@echo "Running $(BINARY_NAME)..."
	@$(GOBUILD) -o $(BIN_DIR)/$(BINARY_NAME) $(SRC_DIR)/main.go
	@$(BIN_DIR)/$(BINARY_NAME)
