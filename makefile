NAME = harmonia
VERSION = 0.0.1

# Included as relative path because of compile step, if not Go will look in GOROOT
SRC_DIR = ./src
BIN_DIR = bin
ENV := $(if $(ENV),$(ENV),dev)

####### GO TARGETS ########
.PHONY: swag compile run godoc test tidy

# Constructs bin/ directory for holding compiled binaries
$(BIN_DIR):
	mkdir $(BIN_DIR)

# Builds out swagger documentation and outputs the specs to the docs/ directory
swag:
	go install github.com/swaggo/swag/cmd/swag@v1.8.0
# --parseDependency, --parseInternal and --parseDepth are being used here to parse definitions outside of main package
	swag init -d $(SRC_DIR)/main --parseDependency --parseDepth 1 -g server.go -o $(SRC_DIR)/main/docs

# Compiles build of the src/main Go application and outputs binary to the bin/ directory
# Compiles with different flags depending on environment
# This is dependent on the swagger documentation and bin/ directory being present
compile: swag $(BIN_DIR)
ifneq ($(ENV), prod)
	@echo "compiling non-release build"
	go build -gcflags=all="-N -l" -ldflags "-X main.harmoniaVersion=$(VERSION)" -o $(BIN_DIR)/$(NAME) $(SRC_DIR)/main
else
	@echo "compiling release build"
	go build -ldflags "-X main.harmoniaVersion=$(VERSION)" -o $(BIN_DIR)/$(NAME) $(SRC_DIR)/main
endif

# Runs the compiled version of the Go application located in the bin/ directory
run:
	./$(BIN_DIR)/$(NAME)

# Serves Go documentation for the Go application based on docs
godoc:
	godoc -http=:6060

# Runs Go application tests
test: swag
	go test ./...

# Cleans up Go application dependencies
tidy:
	go mod tidy

# Lints source code
lint: swag
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.46.2
	golangci-lint run

####### MISCELLANEOUS TARGETS ########
.PHONY: local get-version get-tag get-name clean

# Entrypoint to setting up local environment configuration
local:
	@sh local.sh

# Outputs application version - used by Jenkins primarily for building and tagging image 
get-version:
	@echo -n $(VERSION)

# Outputs app name
get-name:
	@echo -n $(NAME)

# Cleans up artifacts
clean:
	rm -rf $(BIN_DIR)
	rm -rf $(SRC_DIR)/main/docs
