BIN = fillin

all: clean build

build: deps
	go build -o build/$(BIN) .

install: deps
	go install

deps:
	go get -d -v .

cross: crossdeps
	goxz -os=linux,darwin,freebsd,netbsd,windows -arch=386,amd64 -n $(BIN) -d snapshot .

crossdeps: deps
	go get github.com/Songmu/goxz/cmd/goxz

test: testdeps build clean-test-tmp
	go test -v ./...

testdeps:
	go get -d -v -t .

lint: lintdeps build
	go vet
	golint -set_exit_status ./...

lintdeps:
	go get -d -v -t .
	command -v golint >/dev/null || go get -u golang.org/x/lint/golint

clean: clean-test-tmp
	rm -rf build  snapshot debian
	go clean

clean-test-tmp:
	rm -rf .test

.PHONY: build install deps cross crossdeps test clean-test-tmp testdeps lint lintdeps clean
