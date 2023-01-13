build:
	go build -ldflags "-X main.commit=$(shell git rev-parse HEAD) -X main.version=$(shell cat VERSION)" -o accicalc main.go

install: build
	go install

lint:
	golangci-lint run

release:
	git checkout main
	git pull
	git tag -a $(shell cat VERSION) -m "Release $(shell cat VERSION)"
	git push origin $(shell cat VERSION)

clean:
	go clean

.PHONY: build install lint release
