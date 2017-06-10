BIN = fillin

all: clean build

build: deps
	go build -o build/$(BIN) .

install: deps
	go install

deps:
	go get -d -v .

test: testdeps build
	go test -v ./...

testdeps:
	go get -d -v -t .

lint: lintdeps build
	go vet
	golint -set_exit_status ./...

lintdeps:
	go get -d -v -t .
	go get -u github.com/golang/lint/golint

clean:
	rm -rf build .test
	go clean

.PHONY: build deps test testdeps lint lintdeps clean
