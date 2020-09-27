lint:
	go mod tidy
	go fmt ./...
	go vet ./...
	go get -u golang.org/x/lint/golint	
	golint -set_exit_status $(go list ./...)
	go get -u github.com/golangci/golangci-lint/cmd/golangci-lint
	golangci-lint run -E gofmt -E golint -E vet

build-player:
	docker build -t $(tag) -f build/player/Dockerfile .

build-server:
	docker build -t $(tag) -f build/server/Dockerfile .
