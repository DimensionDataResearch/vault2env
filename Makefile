VERSION = 0.0.2
VERSION_INFO_FILE = ./$(PROVIDER_NAME)/version-info.go

BIN_DIRECTORY   = _bin
EXECUTABLE_NAME = vault2env
DIST_ZIP_PREFIX = $(EXECUTABLE_NAME).v$(VERSION)

REPO_BASE     = github.com/DimensionDataResearch
REPO_ROOT     = $(REPO_BASE)/vault2env

default: dev

fmt:
	go fmt ./...

clean:
	rm -rf ./_bin
	go clean $(REPO_ROOT)/...

# Peform a development (current-platform-only) build.
dev: fmt
	go build -o _bin/${EXECUTABLE_NAME}

build: fmt build_linux_amd64 build_darwin_amd64 build_windows_amd64

build_linux_amd64:
	GOOS=linux GOARCH=amd64 go build -o _bin/linux-amd64/${EXECUTABLE_NAME}

build_darwin_amd64:
	GOOS=darwin GOARCH=amd64 go build -o _bin/darwin-amd64/${EXECUTABLE_NAME}

build_windows_amd64:
	GOOS=windows GOARCH=amd64 go build -o _bin/windows-amd64/${EXECUTABLE_NAME}.exe

# Produce archives for a GitHub release.
dist: build
	cd $(BIN_DIRECTORY)/windows-amd64 && \
		zip -9 ../$(DIST_ZIP_PREFIX).windows-amd64.zip $(EXECUTABLE_NAME).exe
	cd $(BIN_DIRECTORY)/linux-amd64 && \
		zip -9 ../$(DIST_ZIP_PREFIX).linux-amd64.zip $(EXECUTABLE_NAME)
	cd $(BIN_DIRECTORY)/darwin-amd64 && \
		zip -9 ../$(DIST_ZIP_PREFIX)-darwin-amd64.zip $(EXECUTABLE_NAME)
