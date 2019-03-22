BIN := fillin
BUILD_LDFLAGS := "-s -w"
export GO111MODULE=on

.PHONY: all
all: clean build

.PHONY: build
build: deps
	go build -ldflags=$(BUILD_LDFLAGS) -o build/$(BIN) .

.PHONY: install
install: deps
	go install -ldflags=$(BUILD_LDFLAGS) ./...

.PHONY: deps
deps:
	go get -d -v ./...

.PHONY: cross
cross: crossdeps
	goxz -os=linux,darwin,freebsd,netbsd,windows -arch=386,amd64 -n $(BIN) .

.PHONY: crossdeps
crossdeps: deps
	GO111MODULE=off go get github.com/Songmu/goxz/cmd/goxz

.PHONY: test
test: build
	go test -v ./...

.PHONY: lint
lint: lintdeps
	go vet ./...
	golint -set_exit_status ./...

.PHONY: lintdeps
lintdeps:
	GO111MODULE=off go get -u golang.org/x/lint/golint

.PHONY: clean
clean:
	rm -rf build goxz
	go clean
