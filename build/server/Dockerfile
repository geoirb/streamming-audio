FROM golang:1.14-alpine AS build

WORKDIR /go/src/server
COPY . ./

RUN go build \
    -o /out/service \
    ./cmd/server/main.go

# copy to alpine image
FROM alpine:3.8
WORKDIR /app
COPY --from=build /out/service /app/service
CMD ["/app/service"]
