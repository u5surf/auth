VERSION := $(shell grep -Eo '(v[0-9]+[\.][0-9]+[\.][0-9]+(-dev)?)' main.go)

.PHONY: build docker release

build:
	go fmt ./...
	CGO_ENABLED=1 go build -o bin/auth .

docker:
	docker build -t moov/auth:$(VERSION) -f Dockerfile .
	docker tag moov/auth:$(VERSION) moov/auth:latest

release: docker
	go vet ./...
	go test ./...
	git tag $(VERSION)

release-push:
	echo "$DOCKER_PASSWORD" | docker login -u wadearnold --password-stdin 
	git push origin $(VERSION)
	docker push moov/auth:$(VERSION)
