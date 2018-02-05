EXECUTABLE_NAME=vault2env

default: dev

fmt:
	go fmt ./...

clean:
	rm -rf ./_bin

dev: fmt
	go build -o _bin/${EXECUTABLE_NAME}

build: fmt build_linux_amd64 build_darwin_amd64 build_windows_amd64

build_linux_amd64:
	GOOS=linux GOARCH=amd64 go build -o _bin/${EXECUTABLE_NAME}

build_darwin_amd64:
	GOOS=darwin GOARCH=amd64 go build -o _bin/${EXECUTABLE_NAME}

build_windows_amd64:
	GOOS=windows GOARCH=amd64 go build -o _bin/${EXECUTABLE_NAME}.exe
