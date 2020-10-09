lint:
	go fmt ./...
	go vet ./...
	go get golang.org/x/lint/golint	
	golint -set_exit_status $(go list ./...)
	go get github.com/golangci/golangci-lint/cmd/golangci-lint
	golangci-lint run -E gofmt -E golint -E vet 

build-recorder:
	docker build -t $(tag) -f build/recorder/Dockerfile .

build-player:
	docker build -t $(tag) -f build/player/Dockerfile .

build-server:
	docker build -t $(tag) -f build/server/Dockerfile .
