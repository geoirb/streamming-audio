check-code:
	go get -u golang.org/x/lint/golint
	go get -u github.com/golangci/golangci-lint/cmd/golangci-lint
	go fmt ./...
	go vet ./...
	golint -set_exit_status $(go list ./...)
	golangci-lint run -E gofmt -E golint -E vet
	out=$(go fmt ./...) && if [[ -n "$out" ]]; then echo "$out"; exit 1; fi

media:
	docker build -t $(tag) -f build/media/Dockerfile .

server:
	docker build -t $(tag) -f build/server/Dockerfile .