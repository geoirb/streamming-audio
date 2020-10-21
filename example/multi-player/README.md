# Test Server and Player

1. Start server
    * `go run cmd/server/main.go`
2. Build and start players on different machines
   * `make build-player tag=player`
   * `docker run --rm -p 8080:8080 -p 8081:8081 --device /dev/snd --name player1 player`
3. Configurate `example.go`:
   * input all players on in map `player`   
4. Start example
   * `go run example/server-player/example.go`