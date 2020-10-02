FROM golang:1.15-alpine AS build

WORKDIR /go/src/server
COPY . ./

RUN go build \
    -o /out/service \
    ./cmd/server/main.go

FROM alpine
WORKDIR /app
COPY --from=build /out/service /app/service
CMD ["/app/service"]
