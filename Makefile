GOFLAGS ?= $(GOFLAGS:)

VERSION=v0.1.0

ifdef BUILD_NUMBER
VERSION:=$(VERSION)+$(BUILD_NUMBER)
endif

ifdef RELEASE_VERSION
ifneq ($(RELEASE_VERSION),none)
VERSION=$(RELEASE_VERSION)
endif
endif

GO_LDFLAGS=-ldflags "-X github.com/achuchev/pump_monitor/v4.versionString=$(VERSION) -X github.com/achuchev/pump_monitor/v4.versionBuildTimeStamp=`date -u +%Y%m%d.%H%M%S` -s -w"
version:
	echo "$(VERSION)"

get: gofmt
	go get $(GOFLAGS) ./...

build_quick: get
	env GOOS=darwin  GOARCH=arm64 go build $(GO_LDFLAGS) -o bin/darwin/pump-monitor        ./cmd

build: get
	env GOOS=linux   GOARCH=arm GOARM=7 go build $(GO_LDFLAGS) -o bin/rpi/pump-monitor         ./cmd
	env GOOS=darwin  GOARCH=arm64 go build $(GO_LDFLAGS) -o bin/darwin/pump-monitor	   ./cmd
	env GOOS=windows GOARCH=amd64 go build $(GO_LDFLAGS) -o bin/windows/pump-monitor.exe   ./cmd

gofmt:
	! gofmt -l . | grep -v ^vendor/ | grep .

linter:
	@golangci-lint --version || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b /go/bin
	golangci-lint run
