FROM golang:alpine AS build

WORKDIR /go/src/media
COPY . ./

RUN apk update \
    && apk add build-base alsa-lib-dev

RUN go build \
    -o /out/service \
    ./cmd/media/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=build /out/service /app/service
RUN  apk add alsa-lib-dev
CMD ["/app/service"]
