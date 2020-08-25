lint:
	go fmt ./...
	go vet ./...
	go get -u golang.org/x/lint/golint	
	golint -set_exit_status $(go list ./...)
	go get -u github.com/golangci/golangci-lint/cmd/golangci-lint
	golangci-lint run -E gofmt -E golint -E vet

build-media:
	docker build -t $(tag) -f build/media/Dockerfile .

run-media:
	docker run -p 0.0.0.0:$(port):$(port) -p 0.0.0.0:$(port):$(port)/udp  --device /dev/snd --rm  $(image)

build-server:
	docker build -t $(tag) -f build/server/Dockerfile .

run-server:
	docker run -p $(port):$(port) --rm -e FILE=/app/test.wav server