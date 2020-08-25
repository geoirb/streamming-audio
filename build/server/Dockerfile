FROM golang AS build

WORKDIR /go/src/server
COPY . ./

RUN go build \
    -o /out/service \
    ./cmd/server/main.go

FROM alpine
WORKDIR /app
COPY --from=build /out/service /app/service
COPY ./audio /audio
CMD ["/app/service"]
