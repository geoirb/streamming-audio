FROM golang:1.15-alpine AS build

WORKDIR /go/src/player
COPY . ./

RUN apk update \
    && apk add build-base alsa-lib-dev

RUN go build \
    -o /out/service \
    ./cmd/player/main.go

FROM alpine
WORKDIR /app
COPY --from=build /out/service /app/service
RUN  apk add alsa-lib-dev
CMD ["/app/service"]
