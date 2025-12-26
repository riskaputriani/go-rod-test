# Makefile

# The name of your application
APP_NAME=go-rod-testing-browser-restrict

# The main package of your application
MAIN_PACKAGE=main.go

# OS and Architecture for the build
GOOS=linux
GOARCH=amd64

.PHONY: build clean

build:
	@echo "Building for $(GOOS)/$(GOARCH)..."
	@GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o $(APP_NAME) $(MAIN_PACKAGE)

clean:
	@echo "Cleaning..."
	@powershell -Command "Remove-Item -Path $(APP_NAME) -ErrorAction SilentlyContinue"
