#flags=-ldflags="-s -w"
flags=-ldflags="-s -w -extldflags -static"
TAG := $(shell git tag | sed -e "s,v,,g" | sort -r | head -n 1)

all: build

vet:
	go vet .

build:
	go clean; rm -rf pkg; CGO_ENABLED=0 go build -o chapaas ${flags}

build_debug:
	go clean; rm -rf pkg; CGO_ENABLED=0 go build -o chapaas ${flags} -gcflags="-m -m"

build_amd64: build_linux

build_darwin:
	go clean; rm -rf pkg chapaas; GOOS=darwin CGO_ENABLED=0 go build -o chapaas ${flags}

build_linux:
	go clean; rm -rf pkg chapaas; GOOS=linux CGO_ENABLED=0 go build -o chapaas ${flags}
	mkdir -p /tmp/chapaas/amd64
	cp chapaas /tmp/chapaas/amd64

build_power8:
	go clean; rm -rf pkg chapaas; GOARCH=ppc64le GOOS=linux CGO_ENABLED=0 go build -o chapaas ${flags}
	mkdir -p /tmp/chapaas/power8
	cp chapaas /tmp/chapaas/power8

build_arm64:
	go clean; rm -rf pkg chapaas; GOARCH=arm64 GOOS=linux CGO_ENABLED=0 go build -o chapaas ${flags}
	mkdir -p /tmp/chapaas/arm64
	cp chapaas /tmp/chapaas/arm64

build_windows:
	go clean; rm -rf pkg chapaas; GOARCH=amd64 GOOS=windows CGO_ENABLED=0 go build -o chapaas ${flags}
	mkdir -p /tmp/chapaas/windows
	cp chapaas /tmp/chapaas/windows

install:
	go install

clean:
	go clean; rm -rf pkg; rm -rf chapaas

test : test1

test1:
	go test -v -bench=.

tarball:
	cp -r /tmp/chapaas .
	tar cfz chapaas.tar.gz chapaas
	rm -rf /tmp/chapaas

release: clean build_amd64 build_arm64 build_windows build_power8 build_darwin tarball
