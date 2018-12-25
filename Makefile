BIN = fillin

.PHONY: all
all: clean build

.PHONY: build
build: deps
	go build -o build/$(BIN) .

.PHONY: install
install: deps
	go install

.PHONY: deps
deps:
	go get -d -v .

.PHONY: cross
cross: crossdeps
	goxz -os=linux,darwin,freebsd,netbsd,windows -arch=386,amd64 -n $(BIN) -d snapshot .

.PHONY: crossdeps
crossdeps: deps
	go get github.com/Songmu/goxz/cmd/goxz

.PHONY: test
test: testdeps build clean-test-tmp
	go test -v ./...

.PHONY: testdeps
testdeps:
	go get -d -v -t .

.PHONY: lint
lint: lintdeps build
	go vet
	golint -set_exit_status ./...

.PHONY: lintdeps
lintdeps:
	go get -d -v -t .
	command -v golint >/dev/null || go get -u golang.org/x/lint/golint

.PHONY: clean
clean: clean-test-tmp
	rm -rf build  snapshot debian
	go clean

.PHONY: clean-test-tmp
clean-test-tmp:
	rm -rf .test
