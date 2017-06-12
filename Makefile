BIN = fillin

all: clean build

build: deps
	go build -o build/$(BIN) .

install: deps
	go install

cross: deps clean-test-tmp
	goxc -max-processors=8 -build-ldflags="" \
		-os="linux darwin freebsd netbsd windows" -arch="386 amd64 arm" -d . \
		-resources-include='README*' -n $(BIN)

deps:
	go get -d -v .

test: testdeps build clean-test-tmp
	go test -v ./...

testdeps:
	go get -d -v -t .

lint: lintdeps build
	go vet
	golint -set_exit_status ./...

lintdeps:
	go get -d -v -t .
	go get -u github.com/golang/lint/golint

clean: clean-test-tmp
	rm -rf build  snapshot debian
	go clean

clean-test-tmp:
	rm -rf .test

.PHONY: build install cross deps test clean-test-tmp testdeps lint lintdeps clean
